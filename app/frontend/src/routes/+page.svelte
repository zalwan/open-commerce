<script lang="ts">
	import { api, formatIDR } from '$lib/api';

	let q = '';
	let promise = api.products();

	function search() {
		promise = api.products(q);
	}
</script>

<h1>Katalog</h1>
<form on:submit|preventDefault={search} style="margin-bottom:1rem">
	<input bind:value={q} placeholder="Cari produk…" />
	<button type="submit">Cari</button>
</form>

{#await promise}
	<p>Memuat… (pastikan backend Go jalan di :8080)</p>
{:then products}
	{#if products.length === 0}
		<p>Belum ada produk.</p>
	{:else}
		<ul style="list-style:none;padding:0;display:grid;gap:1rem">
			{#each products as p}
				<li style="border:1px solid #ddd;padding:1rem">
					<a href={`/products/${p.id}`}><strong>{p.name}</strong></a>
					<div>{formatIDR(p.priceMinor)} · stok {p.stock}</div>
				</li>
			{/each}
		</ul>
	{/if}
{:catch e}
	<p style="color:red">Gagal memuat: {e.message}</p>
{/await}
