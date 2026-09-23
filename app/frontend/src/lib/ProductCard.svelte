<script lang="ts">
	import { api, formatIDR, imgSrc } from '$lib/api';
	import type { Product } from '$lib/api';
	import { refreshCartCount } from '$lib/cart-count';
	import { getSessionId } from '$lib/session';

	export let product: Product;

	let state: 'idle' | 'adding' | 'added' | 'error' = 'idle';

	function stockBadge(): string {
		if (product.stock <= 0) return 'badge badge-err';
		if (product.stock <= 5) return 'badge badge-warn';
		return 'badge badge-ok';
	}

	function stockLabel(): string {
		if (product.stock <= 0) return 'Out of stock';
		if (product.stock <= 5) return `Only ${product.stock} left`;
		return `In stock (${product.stock})`;
	}

	async function quickAdd() {
		if (state === 'adding' || product.stock <= 0) return;
		state = 'adding';
		try {
			await api.addToCart(getSessionId(), product.id, 1);
			state = 'added';
			void refreshCartCount();
			setTimeout(() => {
				if (state === 'added') state = 'idle';
			}, 1200);
		} catch {
			state = 'error';
			setTimeout(() => {
				if (state === 'error') state = 'idle';
			}, 1600);
		}
	}
</script>

<li class="card">
	<a class="thumb" href={`/products/${product.id}`} aria-label={product.name}>
		{#if product.imageUrl}
			<img src={imgSrc(product.imageUrl)} alt={product.name} />
		{:else}
			{product.name.slice(0, 1)}
		{/if}
	</a>
	<div class="card-body">
		<a class="title" href={`/products/${product.id}`}>{product.name}</a>
		<div class="price">{formatIDR(product.priceMinor)}</div>
		<span class={stockBadge()}>{stockLabel()}</span>
		<button class="btn btn-ghost btn-quick" on:click={quickAdd} disabled={state === 'adding' || product.stock <= 0}>
			{#if state === 'adding'}
				Adding…
			{:else if state === 'added'}
				✓ Added
			{:else if state === 'error'}
				Failed — retry
			{:else}
				+ Quick add
			{/if}
		</button>
	</div>
</li>

<style>
	.btn-quick {
		margin-top: 0.4rem;
		font-size: 0.8rem;
		padding: 0.4rem 0.8rem;
		width: 100%;
	}
</style>
