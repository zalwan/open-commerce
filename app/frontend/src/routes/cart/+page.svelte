<script lang="ts">
	import { onMount } from 'svelte';
	import { api, formatIDR } from '$lib/api';
	import { getSessionId } from '$lib/session';
	import type { Cart } from '$lib/api';

	let cart: Cart | null = null;
	let error = '';

	onMount(async () => {
		try {
			cart = await api.cart(getSessionId());
		} catch (e) {
			error = e instanceof Error ? e.message : 'Gagal memuat';
		}
	});
</script>

<h1>Keranjang</h1>
{#if error}<p style="color:red">{error}</p>{/if}
{#if !cart}
	<p>Memuat…</p>
{:else if cart.items.length === 0}
	<p>Keranjang kosong. <a href="/">Belanja dulu</a></p>
{:else}
	<ul>
		{#each cart.items as it}
			<li>{it.name} × {it.qty} — {formatIDR(it.priceMinor * it.qty)}</li>
		{/each}
	</ul>
	<p><strong>Subtotal: {formatIDR(cart.subtotalMinor)}</strong></p>
	<a href="/checkout">Lanjut checkout →</a>
{/if}
