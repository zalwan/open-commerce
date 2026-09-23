<script lang="ts">
	import { onMount } from 'svelte';
	import { api, formatIDR } from '$lib/api';
	import { getSessionId } from '$lib/session';
	import type { Cart } from '$lib/api';

	let cart: Cart | null = null;
	let error = '';

	async function reload() {
		cart = await api.cart(getSessionId());
	}

	onMount(async () => {
		try {
			await reload();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load';
		}
	});

	async function changeQty(productId: string, qty: number) {
		error = '';
		try {
			cart = await api.setQty(getSessionId(), productId, qty);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to update';
		}
	}

	async function remove(productId: string) {
		error = '';
		try {
			cart = await api.removeFromCart(getSessionId(), productId);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to remove';
		}
	}
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
		{#each cart.items as it (it.productId)}
			<li class="row">
				<div class="grow"><strong>{it.name}</strong><br /><span class="muted">{formatIDR(it.priceMinor)} each</span></div>
				<div class="stepper">
					<button class="btn btn-ghost btn-step" on:click={() => changeQty(it.productId, it.qty - 1)} disabled={it.qty <= 1}>−</button>
					<span class="qty-num">{it.qty}</span>
					<button class="btn btn-ghost btn-step" on:click={() => changeQty(it.productId, it.qty + 1)}>+</button>
				</div>
				<div class="price">{formatIDR(it.priceMinor * it.qty)}</div>
				<button class="btn btn-danger btn-step" on:click={() => remove(it.productId)}>Remove</button>
			</li>
		{/each}
	</ul>
	<p class="total">Subtotal: {formatIDR(cart.subtotalMinor)}</p>
	<div class="btn-row" style="justify-content:flex-end">
		<a class="btn btn-ghost" href="/">Continue shopping</a>
		<a class="btn btn-primary" href="/checkout">Proceed to checkout →</a>
	</div>
{/if}

<style>
	.stepper {
		display: flex;
		align-items: center;
		gap: 0.4rem;
	}
	.btn-step {
		padding: 0.3rem 0.7rem;
		font-size: 0.85rem;
	}
	.qty-num {
		font-family: var(--mono);
		min-width: 2ch;
		text-align: center;
	}
</style>
