// Package media stores product images on local disk and serves them.
// v0.3 scope: single local dir (DATA_DIR), no CDN. Filenames are derived
// from the product ID to avoid traversal and collisions.
package media

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrTooBig        = errors.New("image exceeds 5MB")
	ErrBadType       = errors.New("only jpeg, png, webp or gif images are accepted")
	ErrNoFile        = errors.New("missing image file")
	maxBytes   int64 = 5 << 20
	allowed          = map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
		"image/gif":  ".gif",
	}
)

// Store writes uploads under dir, served publicly at urlPrefix.
type Store struct {
	dir       string
	urlPrefix string // e.g. "/static/"
}

func NewStore(dir, urlPrefix string) *Store {
	return &Store{dir: dir, urlPrefix: urlPrefix}
}

func (s *Store) Dir() string { return s.dir }

// SaveImage validates and stores the uploaded file for productID,
// returning the public URL path.
func (s *Store) SaveImage(productID string, r *http.Request) (string, error) {
	if err := r.ParseMultipartForm(maxBytes + (1 << 20)); err != nil {
		return "", ErrTooBig
	}
	f, _, err := r.FormFile("image")
	if err != nil {
		return "", ErrNoFile
	}
	defer f.Close()
	head := make([]byte, 512)
	n, _ := io.ReadFull(f, head)
	ctype := http.DetectContentType(head[:n])
	ext, ok := allowed[ctype]
	if !ok {
		return "", ErrBadType
	}
	seeked := false
	if seeker, ok := f.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
		seeked = true
	}
	safe := strings.Map(func(c rune) rune {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			return c
		}
		return '-'
	}, productID)
	if safe == "" {
		return "", ErrNoFile
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir uploads: %w", err)
	}
	name := safe + ext
	// Remove stale siblings with other extensions for the same product.
	if entries, err := os.ReadDir(s.dir); err == nil {
		for _, e := range entries {
			base := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
			if base == safe && e.Name() != name {
				_ = os.Remove(filepath.Join(s.dir, e.Name()))
			}
		}
	}
	dst, err := os.Create(filepath.Join(s.dir, name))
	if err != nil {
		return "", fmt.Errorf("save image: %w", err)
	}
	defer dst.Close()
	var src io.Reader = f
	if !seeked {
		// Non-seekable upload: replay the sniffed head before the rest.
		src = io.MultiReader(strings.NewReader(string(head[:n])), f)
	}
	written, err := io.Copy(dst, io.LimitReader(src, maxBytes+1))
	if err != nil {
		_ = os.Remove(dst.Name())
		return "", fmt.Errorf("save image: %w", err)
	}
	if written > maxBytes {
		_ = os.Remove(dst.Name())
		return "", ErrTooBig
	}
	return s.urlPrefix + name, nil
}
