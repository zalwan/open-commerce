<script lang="ts">
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import type { HeroScene } from '$lib/hero-scene';
	import ProductCard from '$lib/ProductCard.svelte';

	let q = '';
	let promise = api.products();
	let canvas: HTMLCanvasElement | null = null;
	let scene: HeroScene | null = null;
	let webgl = true;

	// Three.js hero is client-only: lazy-load behind a browser guard so
	// SSR/prerender never touches WebGL. Any failure falls back to CSS.
	onMount(() => {
		if (!browser || !canvas) return;
		let cancelled = false;
		import('$lib/hero-scene')
			.then((mod) => {
				if (cancelled) return;
				try {
					scene = mod.createHeroScene(canvas as HTMLCanvasElement);
				} catch {
					webgl = false;
				}
			})
			.catch(() => {
				webgl = false;
			});
		return () => {
			cancelled = true;
			scene?.destroy();
			scene = null;
		};
	});

	function search() {
		promise = api.products(q);
	}
</script>

<div class="hero">
	{#if webgl}
		<canvas class="hero-canvas" bind:this={canvas} aria-hidden="true"></canvas>
	{/if}
	<div class="hero-copy">
		<h1>Shop straight from our store</h1>
		<p>Catalog, cart, and checkout in one open source app.</p>
	</div>
</div>

<form class="searchbar" on:submit|preventDefault={search}>
	<input class="input" bind:value={q} placeholder="Search products…" />
	<button class="btn btn-primary" type="submit">Search</button>
</form>

{#await promise}
	<p class="muted">Loading… (make sure the Go backend runs on :8080)</p>
{:then products}
	{#if products.length === 0}
		<p>No matching products yet.</p>
	{:else}
		<ul class="grid">
			{#each products as p (p.id)}
				<ProductCard product={p} />
			{/each}
		</ul>
		<div class="strip">
			<span class="strip-label">{products.length} products live · powered by</span>
			<span class="chip">Go 1.25 API</span>
			<span class="chip">Postgres 16</span>
			<span class="chip">SvelteKit 2</span>
			<span class="chip">Three.js hero</span>
			<span class="chip">MIT</span>
		</div>
	{/if}
{:catch e}
	<p class="alert-error">Failed to load: {e.message}</p>
{/await}
