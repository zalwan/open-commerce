package catalog

import "testing"

func TestValidate(t *testing.T) {
	if err := Validate(Product{Name: "A", PriceMinor: 100, Stock: 1}); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
	if err := Validate(Product{PriceMinor: 100}); err != ErrInvalidName {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
	if err := Validate(Product{Name: "A", PriceMinor: 0}); err != ErrInvalidPrice {
		t.Fatalf("expected ErrInvalidPrice, got %v", err)
	}
}

func TestCreateAndList(t *testing.T) {
	svc := NewService(NewMemoryStore())
	before := len(svc.List(""))
	p, err := svc.Create(Product{Name: "Keyboard", PriceMinor: 5000000, Stock: 5})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if p.ID == "" || p.Currency != "IDR" {
		t.Fatalf("unexpected product: %+v", p)
	}
	if got := len(svc.List("keyboard")); got != 1 {
		t.Fatalf("expected 1 search hit, got %d", got)
	}
	if got := len(svc.List("")); got != before+1 {
		t.Fatalf("expected %d products, got %d", before+1, got)
	}
}
