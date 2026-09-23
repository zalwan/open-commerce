<script lang="ts">
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import type { HeroScene } from '$lib/hero-scene';
	import ProductCard from '$lib/ProductCard.svelte';

	let q = '';
	let category = '';
	let sort = 'name_asc';
	let page = 1;
	const perPage = 12;
	let cats: string[] = [];
	let promise = load();

	function load() {
		return api.products({ q, category, sort, page, perPage });
	}

	function search() {
		page = 1;
		promise = load();
	}

	function setCategory(c: string) {
		category = c;
		page = 1;
		promise = load();
	}

	function setSort(s: string) {
		sort = s;
		page = 1;
		promise = load();
	}

	function gotoPage(p: number) {
		page = p;
		promise = load();
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}

	let canvas: HTMLCanvasElement | null = null;
	let scene: HeroScene | null = null;
	let webgl = true;

	// Three.js hero is client-only: lazy-load behind a browser guard so
	// SSR/prerender never touches WebGL. Any failure falls back to CSS.
	onMount(() => {
		api.categories().then((c) => (cats = c)).catch(() => (cats = []));
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
	<select class="input" style="max-width:11rem" bind:value={sort} on:change={() => setSort(sort)}>
		<option value="name_asc">Name A–Z</option>
		<option value="price_asc">Price ↑</option>
		<option value="price_desc">Price ↓</option>
		<option value="newest">Newest</option>
	</select>
	<button class="btn btn-primary" type="submit">Search</button>
</form>

{#if cats.length > 0}
	<div class="catrow">
		<button class="chipbtn" class:active={category === ''} on:click={() => setCategory('')}>All</button>
		{#each cats as c}
			<button class="chipbtn" class:active={category === c} on:click={() => setCategory(c)}>{c}</button>
		{/each}
	</div>
{/if}

{#await promise}
	<p class="muted">Loading… (make sure the Go backend runs on :8080)</p>
{:then res}
	{#if res.items.length === 0}
		<p>No matching products yet.</p>
	{:else}
		<p class="muted">{res.total} product{res.total === 1 ? '' : 's'} · page {res.page} of {Math.max(1, Math.ceil(res.total / res.perPage))}</p>
		<ul class="grid">
			{#each res.items as p (p.id)}
				<ProductCard product={p} />
			{/each}
		</ul>
		{#if res.total > res.perPage}
			<div class="btn-row pager">
				<button class="btn btn-ghost" disabled={res.page <= 1} on:click={() => gotoPage(res.page - 1)}>← Prev</button>
				<button class="btn btn-ghost" disabled={res.page * res.perPage >= res.total} on:click={() => gotoPage(res.page + 1)}>Next →</button>
			</div>
		{/if}
		<div class="strip">
			<span class="strip-label">{res.total} products live · powered by</span>
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

<style>
	.catrow {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
		margin-bottom: 1rem;
	}
	.chipbtn {
		font-family: var(--mono);
		font-size: 0.78rem;
		background: transparent;
		color: var(--muted);
		border: 1px solid var(--border);
		border-radius: 999px;
		padding: 0.3rem 0.85rem;
		cursor: pointer;
	}
	.chipbtn.active {
		color: #06281d;
		background: var(--brand);
		border-color: var(--brand);
		font-weight: 700;
	}
	.pager {
		justify-content: center;
		margin-top: 1.5rem;
	}
</style>
