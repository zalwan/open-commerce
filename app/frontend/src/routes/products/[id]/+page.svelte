<script lang="ts">
	import { page } from '$app/stores';
	import { api, formatIDR } from '$lib/api';
	import { getSessionId } from '$lib/session';

	const id: string = $page.params.id ?? '';
	let qty = 1;
	let msg = '';

	async function add() {
		try {
			await api.addToCart(getSessionId(), id, qty);
			msg = 'Ditambahkan ke keranjang.';
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Gagal';
		}
	}
</script>

<a href="/">← Kembali</a>
{#await api.product(id) then p}
	<h1>{p.name}</h1>
	<p>{p.description ?? '-'}</p>
	<p><strong>{formatIDR(p.priceMinor)}</strong> · stok {p.stock}</p>
	<label>Qty <input type="number" min="1" max={p.stock} bind:value={qty} /></label>
	<button on:click={add}>Tambah ke keranjang</button>
	{#if msg}<p>{msg}</p>{/if}
{:catch e}
	<p style="color:red">Gagal: {e.message}</p>
{/await}
