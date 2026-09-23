<script lang="ts">
	import { api, formatIDR } from '$lib/api';
	import type { Product } from '$lib/api';

	let q = '';
	let promise = api.products();

	function search() {
		promise = api.products(q);
	}

	function stockBadge(p: Product): string {
		if (p.stock <= 0) return 'badge badge-err';
		if (p.stock <= 5) return 'badge badge-warn';
		return 'badge badge-ok';
	}

	function stockLabel(p: Product): string {
		if (p.stock <= 0) return 'Out of stock';
		if (p.stock <= 5) return `Only ${p.stock} left`;
		return `In stock (${p.stock})`;
	}
</script>

<div class="hero">
	<h1>Shop straight from our store</h1>
	<p>Catalog, cart, and checkout in one open source app.</p>
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
			{#each products as p}
				<li class="card">
					<a class="thumb" href={`/products/${p.id}`} aria-label={p.name}>
						{#if p.imageUrl}
							<img src={p.imageUrl} alt={p.name} />
						{:else}
							{p.name.slice(0, 1)}
						{/if}
					</a>
					<div class="card-body">
						<a class="title" href={`/products/${p.id}`}>{p.name}</a>
						<div class="price">{formatIDR(p.priceMinor)}</div>
						<span class={stockBadge(p)}>{stockLabel(p)}</span>
					</div>
				</li>
			{/each}
		</ul>
	{/if}
{:catch e}
	<p class="alert-error">Failed to load: {e.message}</p>
{/await}
