<script lang="ts">
	import { afterNavigate } from '$app/navigation';
	import { onMount } from 'svelte';
	import { cartCount, refreshCartCount } from '$lib/cart-count';
	import '../app.css';

	onMount(refreshCartCount);
	afterNavigate(() => void refreshCartCount());
</script>

<header class="topbar">
	<div class="topbar-inner">
		<a class="brand" href="/">Open<span>Commerce</span></a>
		<nav class="nav">
			<a href="/">Catalog</a>
			<a href="/cart">Cart{#if $cartCount > 0} <span class="cart-count">{$cartCount}</span>{/if}</a>
			<a href="/orders/track">Track</a>
			<a href="/admin">Admin</a>
		</nav>
	</div>
</header>
<main class="container">
	<slot />
</main>
<footer class="footer">Open Commerce · open source single-vendor · stub payments (no real money)</footer>
