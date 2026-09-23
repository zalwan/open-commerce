<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api, formatIDR } from '$lib/api';
	import type { Cart } from '$lib/api';
	import { getSessionId } from '$lib/session';

	let email = '';
	let error = '';
	let loading = false;
	let cart: Cart | null = null;

	onMount(async () => {
		try {
			cart = await api.cart(getSessionId());
		} catch {
			cart = null;
		}
	});

	async function submit() {
		loading = true;
		error = '';
		try {
			const order = await api.checkout(getSessionId(), email);
			await goto(`/orders/${order.id}`);
		} catch (e) {
			const msg = e instanceof Error ? e.message : 'Checkout failed';
			error = msg.includes('409')
				? 'Not enough stock for some items — reduce the qty or come back later.'
				: msg;
		} finally {
			loading = false;
		}
	}
</script>

<h1>Checkout</h1>
<div class="panel">
	<h2>Summary</h2>
	{#if !cart}
		<p class="muted">Loading cart…</p>
	{:else if cart.items.length === 0}
		<p>Cart is empty. <a href="/">Shop first</a>.</p>
	{:else}
		<ul class="rows">
			{#each cart.items as it}
				<li class="row">
					<div class="grow">{it.name} <span class="muted">× {it.qty}</span></div>
					<div>{formatIDR(it.priceMinor * it.qty)}</div>
				</li>
			{/each}
		</ul>
		<p class="total">Total: {formatIDR(cart.subtotalMinor)}</p>
	{/if}
</div>

<div class="panel">
	<h2>Payment</h2>
	<p class="muted" style="margin-top:0">Via <em>mock stub</em> — no real money.</p>
	<form on:submit|preventDefault={submit}>
		<label class="field">Email<input class="input" type="email" required bind:value={email} placeholder="you@mail.test" /></label>
		<button class="btn btn-primary" type="submit" disabled={loading}>{loading ? 'Processing…' : 'Pay (mock)'}</button>
	</form>
	{#if error}<p class="alert-error">{error}</p>{/if}
</div>
