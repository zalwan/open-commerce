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
		if (p.stock <= 0) return 'Stok habis';
		if (p.stock <= 5) return `Sisa ${p.stock}`;
		return `Stok ${p.stock}`;
	}
</script>

<div class="hero">
	<h1>Belanja langsung dari toko kami</h1>
	<p>Katalog, keranjang, dan checkout dalam satu aplikasi open source.</p>
</div>

<form class="searchbar" on:submit|preventDefault={search}>
	<input class="input" bind:value={q} placeholder="Cari produk…" />
	<button class="btn btn-primary" type="submit">Cari</button>
</form>

{#await promise}
	<p class="muted">Memuat… (pastikan backend Go jalan di :8080)</p>
{:then products}
	{#if products.length === 0}
		<p>Belum ada produk yang cocok.</p>
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
	<p class="alert-error">Gagal memuat: {e.message}</p>
{/await}
