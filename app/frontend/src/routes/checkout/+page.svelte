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
			const msg = e instanceof Error ? e.message : 'Checkout gagal';
			error = msg.includes('409')
				? 'Stok tidak cukup untuk sebagian produk — kurangi qty atau kembali lagi nanti.'
				: msg;
		} finally {
			loading = false;
		}
	}
</script>

<h1>Checkout</h1>
<div class="panel">
	<h2>Ringkasan</h2>
	{#if !cart}
		<p class="muted">Memuat keranjang…</p>
	{:else if cart.items.length === 0}
		<p>Keranjang kosong. <a href="/">Belanja dulu</a>.</p>
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
	<h2>Pembayaran</h2>
	<p class="muted" style="margin-top:0">Via <em>stub mock</em> — tanpa uang asli.</p>
	<form on:submit|preventDefault={submit}>
		<label class="field">Email<input class="input" type="email" required bind:value={email} placeholder="kamu@mail.test" /></label>
		<button class="btn btn-primary" type="submit" disabled={loading}>{loading ? 'Memproses…' : 'Bayar (mock)'}</button>
	</form>
	{#if error}<p class="alert-error">{error}</p>{/if}
</div>
