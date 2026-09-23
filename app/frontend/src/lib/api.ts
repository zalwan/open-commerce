// Typed REST client for the Go API. Base URL via PUBLIC_API_BASE_URL.
import { env } from '$env/dynamic/public';

export const API_BASE = env.PUBLIC_API_BASE_URL ?? 'http://localhost:8080';

export interface Product {
	id: string;
	name: string;
	description?: string;
	category?: string;
	priceMinor: number;
	currency: string;
	stock: number;
	imageUrl?: string;
}

export interface ProductList {
	items: Product[];
	total: number;
	page: number;
	perPage: number;
}

export interface ProductQuery {
	q?: string;
	category?: string;
	sort?: string;
	page?: number;
	perPage?: number;
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
	products: (query: ProductQuery = {}) => {
		const params = new URLSearchParams();
		if (query.q) params.set('q', query.q);
		if (query.category) params.set('category', query.category);
		if (query.sort) params.set('sort', query.sort);
		if (query.page) params.set('page', String(query.page));
		if (query.perPage) params.set('per_page', String(query.perPage));
		const qs = params.toString();
		return req<ProductList>(`/api/v1/products${qs ? `?${qs}` : ''}`);
	},
	categories: () => req<string[]>('/api/v1/categories'),
	product: (id: string) => req<Product>(`/api/v1/products/${id}`),
	cart: (sessionId: string) => req<Cart>(`/api/v1/cart?session_id=${encodeURIComponent(sessionId)}`),
	addToCart: (sessionId: string, productId: string, qty: number) =>
		req<Cart>('/api/v1/cart/items', {
			method: 'POST',
			body: JSON.stringify({ session_id: sessionId, product_id: productId, qty })
		}),
	setQty: (sessionId: string, productId: string, qty: number) =>
		req<Cart>('/api/v1/cart/items', {
			method: 'PUT',
			body: JSON.stringify({ session_id: sessionId, product_id: productId, qty })
		}),
	removeFromCart: (sessionId: string, productId: string) =>
		req<Cart>(`/api/v1/cart/items/${productId}?session_id=${encodeURIComponent(sessionId)}`, {
			method: 'DELETE'
		}),
	checkout: (sessionId: string, email: string) =>
		req<Order>('/api/v1/orders/checkout', {
			method: 'POST',
			body: JSON.stringify({ session_id: sessionId, email })
		}),
	order: (id: string) => req<Order>(`/api/v1/orders/${id}`),
	myOrders: (email: string) => req<Order[]>(`/api/v1/orders?email=${encodeURIComponent(email)}`),
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
	uploadImage: async (token: string, id: string, file: File): Promise<Product> => {
		const form = new FormData();
		form.append('image', file);
		const res = await fetch(`${API_BASE}/api/v1/admin/products/${id}/image`, {
			method: 'POST',
			headers: { Authorization: `Bearer ${token}` },
			body: form
		});
		if (!res.ok) throw new Error(`API ${res.status}: ${await res.text()}`);
		return (await res.json()) as Product;
	},
	adminOrders: (token: string) => req<Order[]>('/api/v1/admin/orders', {}, token),
	setOrderStatus: (token: string, id: string, status: string) =>
		req<Order>(`/api/v1/admin/orders/${id}/status`, { method: 'POST', body: JSON.stringify({ status }) }, token),
	adminStats: (token: string) =>
		req<{ products: number; orders: number; paidOrders: number; revenueMinor: number; lowStock: number }>(
			'/api/v1/admin/stats',
			{},
			token
		)
};

export function formatIDR(minor: number): string {
	return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(minor / 100);
}

// imgSrc resolves API-relative image paths (e.g. /static/x.png) against
// the API host; absolute URLs pass through untouched.
export function imgSrc(url: string | undefined): string {
	if (!url) return '';
	return url.startsWith('/') ? `${API_BASE}${url}` : url;
}
