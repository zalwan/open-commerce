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
			error = e instanceof Error ? e.message : 'Failed to load';
		}
	});
</script>

<h1>Cart</h1>
{#if error}<p class="alert-error">{error}</p>{/if}
{#if !cart}
	<p class="muted">Loading…</p>
{:else if cart.items.length === 0}
	<div class="panel">
		<p>Your cart is empty.</p>
		<a class="btn btn-primary" href="/">Start shopping</a>
	</div>
{:else}
	<ul class="rows">
		{#each cart.items as it}
			<li class="row">
				<div class="grow"><strong>{it.name}</strong><br /><span class="muted">× {it.qty}</span></div>
				<div class="price">{formatIDR(it.priceMinor * it.qty)}</div>
			</li>
		{/each}
	</ul>
	<p class="total">Subtotal: {formatIDR(cart.subtotalMinor)}</p>
	<div class="btn-row" style="justify-content:flex-end">
		<a class="btn btn-ghost" href="/">Continue shopping</a>
		<a class="btn btn-primary" href="/checkout">Proceed to checkout →</a>
	</div>
{/if}
