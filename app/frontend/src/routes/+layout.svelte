<script lang="ts">
	import { browser } from '$app/environment';
	import { afterNavigate } from '$app/navigation';
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { getSessionId } from '$lib/session';
	import '../app.css';

	let count = 0;

	async function refresh() {
		if (!browser) return;
		try {
			const cart = await api.cart(getSessionId());
			count = cart.items.reduce((n, it) => n + it.qty, 0);
		} catch {
			count = 0;
		}
	}

	onMount(refresh);
	afterNavigate(refresh);
</script>

<header class="topbar">
	<div class="topbar-inner">
		<a class="brand" href="/">Open<span>Commerce</span></a>
		<nav class="nav">
			<a href="/">Catalog</a>
			<a href="/cart">Cart{#if count > 0} <span class="cart-count">{count}</span>{/if}</a>
			<a href="/orders/track">Track</a>
			<a href="/admin">Admin</a>
		</nav>
	</div>
</header>
<main class="container">
	<slot />
</main>
<footer class="footer">Open Commerce · open source single-vendor · stub payments (no real money)</footer>
