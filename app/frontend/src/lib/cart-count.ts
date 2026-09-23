// Reactive cart badge count. Mutating pages call refreshCartCount()
// after every cart change so the nav badge updates immediately;
// the layout also refreshes on mount + after each navigation.
import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { api } from '$lib/api';
import { getSessionId } from '$lib/session';

export const cartCount = writable(0);

export async function refreshCartCount(): Promise<void> {
	if (!browser) return;
	try {
		const cart = await api.cart(getSessionId());
		cartCount.set(cart.items.reduce((n, it) => n + it.qty, 0));
	} catch {
		cartCount.set(0);
	}
}
