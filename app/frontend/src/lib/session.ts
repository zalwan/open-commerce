// Session + token helpers (localStorage, client-only).
import { browser } from '$app/environment';

const SESSION_KEY = 'oc_session_id';
const TOKEN_KEY = 'oc_token';

export function getSessionId(): string {
	if (!browser) return 'ssr-session';
	let id = localStorage.getItem(SESSION_KEY);
	if (!id) {
		id = `sess-${Math.random().toString(36).slice(2)}${Date.now().toString(36)}`;
		localStorage.setItem(SESSION_KEY, id);
	}
	return id;
}

export function getToken(): string | null {
	if (!browser) return null;
	return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
	if (browser) localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
	if (browser) localStorage.removeItem(TOKEN_KEY);
}
