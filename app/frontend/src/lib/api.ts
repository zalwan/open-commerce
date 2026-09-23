// Typed REST client for the Go API. Base URL via PUBLIC_API_BASE_URL.
import { env } from '$env/dynamic/public';

export const API_BASE = env.PUBLIC_API_BASE_URL ?? 'http://localhost:8080';

export interface Product {
	id: string;
	name: string;
	description?: string;
	priceMinor: number;
	currency: string;
	stock: number;
	imageUrl?: string;
}

export interface CartItem {
	productId: string;
	name: string;
	priceMinor: number;
	qty: number;
}

export interface Cart {
	sessionId: string;
	items: CartItem[];
	subtotalMinor: number;
}

export interface Order {
	id: string;
	email: string;
	items: CartItem[];
	totalMinor: number;
	currency: string;
	status: string;
	paymentTx?: string;
}

async function req<T>(path: string, init?: RequestInit, token?: string): Promise<T> {
	const headers: Record<string, string> = { 'Content-Type': 'application/json' };
	if (init?.headers) {
		for (const [k, v] of Object.entries(init.headers as Record<string, string>)) headers[k] = v;
	}
	if (token) headers['Authorization'] = `Bearer ${token}`;
	const res = await fetch(`${API_BASE}${path}`, { ...init, headers });
	if (!res.ok) {
		const body = await res.text();
		throw new Error(`API ${res.status}: ${body}`);
	}
	if (res.status === 204) return undefined as T;
	return (await res.json()) as T;
}

export const api = {
	products: (q = '') => req<Product[]>(`/api/v1/products${q ? `?q=${encodeURIComponent(q)}` : ''}`),
	product: (id: string) => req<Product>(`/api/v1/products/${id}`),
	cart: (sessionId: string) => req<Cart>(`/api/v1/cart?session_id=${encodeURIComponent(sessionId)}`),
	addToCart: (sessionId: string, productId: string, qty: number) =>
		req<Cart>('/api/v1/cart/items', {
			method: 'POST',
			body: JSON.stringify({ session_id: sessionId, product_id: productId, qty })
		}),
	checkout: (sessionId: string, email: string) =>
		req<Order>('/api/v1/orders/checkout', {
			method: 'POST',
			body: JSON.stringify({ session_id: sessionId, email })
		}),
	order: (id: string) => req<Order>(`/api/v1/orders/${id}`),
	login: (email: string, password: string) =>
		req<{ token: string; role: string }>('/api/v1/auth/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		}),
	// Admin (butuh token stub admin).
	createProduct: (token: string, p: Partial<Product>) =>
		req<Product>('/api/v1/admin/products', { method: 'POST', body: JSON.stringify(p) }, token),
	updateProduct: (token: string, id: string, p: Partial<Product>) =>
		req<Product>(`/api/v1/admin/products/${id}`, { method: 'PUT', body: JSON.stringify(p) }, token),
	deleteProduct: (token: string, id: string) =>
		req<void>(`/api/v1/admin/products/${id}`, { method: 'DELETE' }, token),
	adminOrders: (token: string) => req<Order[]>('/api/v1/admin/orders', {}, token),
	setOrderStatus: (token: string, id: string, status: string) =>
		req<Order>(`/api/v1/admin/orders/${id}/status`, { method: 'POST', body: JSON.stringify({ status }) }, token)
};

export function formatIDR(minor: number): string {
	return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(minor / 100);
}
