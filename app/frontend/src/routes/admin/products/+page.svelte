<script lang="ts">
	import { onMount } from 'svelte';
	import { api, formatIDR } from '$lib/api';
	import type { Product } from '$lib/api';
	import { getToken } from '$lib/session';

	let products: Product[] = [];
	let error = '';
	let msg = '';
	let form: { name: string; priceMinor: number; stock: number; description: string } = {
		name: '',
		priceMinor: 0,
		stock: 0,
		description: ''
	};
	let editing: string | null = null;

	function token(): string {
		const t = getToken();
		if (!t) throw new Error('Not logged in — sign in at /admin first');
		return t;
	}

	async function reload() {
		products = await api.products();
	}

	onMount(() => {
		reload().catch((e) => (error = e instanceof Error ? e.message : 'Failed'));
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
				msg = 'Product updated.';
			} else {
				await api.createProduct(token(), {
					name: form.name,
					priceMinor: form.priceMinor,
					stock: form.stock,
					description: form.description
				});
				msg = 'Product created.';
			}
			editing = null;
			form = { name: '', priceMinor: 0, stock: 0, description: '' };
			await reload();
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Failed to save';
		}
	}

	function edit(p: Product) {
		editing = p.id;
		form = { name: p.name, priceMinor: p.priceMinor, stock: p.stock, description: p.description ?? '' };
	}

	function cancel() {
		editing = null;
		form = { name: '', priceMinor: 0, stock: 0, description: '' };
	}

	async function remove(id: string) {
		try {
			await api.deleteProduct(token(), id);
			msg = 'Product deleted.';
			await reload();
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Failed to delete';
		}
	}
</script>

<a class="back" href="/admin">← Admin</a>
<h1>Manage Products</h1>
{#if error}<p class="alert-error">{error}</p>{/if}

<div class="panel">
	<h2>{editing ? 'Edit product' : 'Add product'}</h2>
	<form on:submit|preventDefault={submit}>
		<label class="field">Name<input class="input" required bind:value={form.name} /></label>
		<label class="field">Price (IDR, no separators)<input class="input" type="number" min="1" required bind:value={form.priceMinor} /></label>
		<label class="field">Stock<input class="input" type="number" min="0" required bind:value={form.stock} /></label>
		<label class="field">Description<input class="input" bind:value={form.description} /></label>
		<div class="btn-row">
			<button class="btn btn-primary" type="submit">{editing ? 'Save' : 'Create'}</button>
			{#if editing}<button class="btn btn-ghost" type="button" on:click={cancel}>Cancel</button>{/if}
		</div>
	</form>
	{#if msg}<p class="alert-ok">{msg}</p>{/if}
</div>

<h2>List ({products.length})</h2>
<ul class="rows">
	{#each products as p}
		<li class="row">
			<div class="grow">
				<strong>{p.name}</strong><br />
				<span class="muted">{formatIDR(p.priceMinor)} · stock {p.stock} · <code>{p.id}</code></span>
			</div>
			<div class="btn-row" style="margin-top:0">
				<button class="btn btn-ghost" on:click={() => edit(p)}>Edit</button>
				<button class="btn btn-danger" on:click={() => remove(p.id)}>Delete</button>
			</div>
		</li>
	{/each}
</ul>
