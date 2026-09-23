<script lang="ts">
	import { onMount } from 'svelte';
	import { api, formatIDR } from '$lib/api';
	import type { Product } from '$lib/api';
	import { getToken } from '$lib/session';

	let products: Product[] = [];
	let error = '';
	let msg = '';
	let form: { id: string; name: string; priceMinor: number; stock: number; description: string } = {
		id: '',
		name: '',
		priceMinor: 0,
		stock: 0,
		description: ''
	};
	let editing: string | null = null;

	function token(): string {
		const t = getToken();
		if (!t) throw new Error('Belum login — masuk dulu di /admin');
		return t;
	}

	async function reload() {
		products = await api.products();
	}

	onMount(() => {
		reload().catch((e) => (error = e instanceof Error ? e.message : 'Gagal'));
	});

	async function submit() {
		msg = '';
		try {
			if (editing) {
				await api.updateProduct(token(), editing, {
					name: form.name,
					priceMinor: form.priceMinor,
					stock: form.stock,
					description: form.description
				});
				msg = 'Produk diperbarui.';
			} else {
				await api.createProduct(token(), {
					name: form.name,
					priceMinor: form.priceMinor,
					stock: form.stock,
					description: form.description
				});
				msg = 'Produk dibuat.';
			}
			editing = null;
			form = { id: '', name: '', priceMinor: 0, stock: 0, description: '' };
			await reload();
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Gagal menyimpan';
		}
	}

	function edit(p: Product) {
		editing = p.id;
		form = { id: p.id, name: p.name, priceMinor: p.priceMinor, stock: p.stock, description: p.description ?? '' };
	}

	async function remove(id: string) {
		try {
			await api.deleteProduct(token(), id);
			msg = 'Produk dihapus.';
			await reload();
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Gagal hapus';
		}
	}
</script>

<h1>Kelola Produk</h1>
{#if error}<p style="color:red">{error}</p>{/if}

<h2>{editing ? 'Edit' : 'Tambah'} produk</h2>
<form on:submit|preventDefault={submit}>
	<label>Nama <input required bind:value={form.name} /></label>
	<label>Harga (IDR) <input type="number" min="1" required bind:value={form.priceMinor} /></label>
	<label>Stok <input type="number" min="0" required bind:value={form.stock} /></label>
	<label>Deskripsi <input bind:value={form.description} /></label>
	<button type="submit">{editing ? 'Simpan' : 'Buat'}</button>
	{#if editing}<button type="button" on:click={() => (editing = null)}>Batal</button>{/if}
</form>
{#if msg}<p>{msg}</p>{/if}

<h2>Daftar ({products.length})</h2>
<ul style="list-style:none;padding:0;display:grid;gap:0.75rem">
	{#each products as p}
		<li style="border:1px solid #ddd;padding:0.75rem">
			<strong>{p.name}</strong> — {formatIDR(p.priceMinor)} · stok {p.stock}
			<button on:click={() => edit(p)}>Edit</button>
			<button on:click={() => remove(p.id)}>Hapus</button>
		</li>
	{/each}
</ul>
