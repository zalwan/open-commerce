<script lang="ts">
	import { page } from '$app/stores';
	import { api, formatIDR, imgSrc } from '$lib/api';
	import { getSessionId } from '$lib/session';

	const id: string = $page.params.id ?? '';
	let qty = 1;
	let msg = '';
	let isErr = false;

	async function add() {
		try {
			await api.addToCart(getSessionId(), id, qty);
			msg = 'Added to cart.';
			isErr = false;
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Failed';
			isErr = true;
		}
	}
</script>

<a class="back" href="/">← Back to catalog</a>
{#await api.product(id) then p}
	<div class="detail">
		<div class="thumb">
			{#if p.imageUrl}
				<img src={imgSrc(p.imageUrl)} alt={p.name} />
			{:else}
				{p.name.slice(0, 1)}
			{/if}
		</div>
		<div>
			<h1 style="margin-top:0">{p.name}</h1>
			<p class="muted">{p.description ?? 'No description.'}</p>
			<p class="price" style="font-size:1.4rem">{formatIDR(p.priceMinor)}</p>
			<p>
				{#if p.stock <= 0}
					<span class="badge badge-err">Out of stock</span>
				{:else if p.stock <= 5}
					<span class="badge badge-warn">Only {p.stock} left</span>
				{:else}
					<span class="badge badge-ok">In stock ({p.stock})</span>
				{/if}
			</p>
			<div class="qty">
				<label for="qty">Qty</label>
				<input id="qty" class="input" type="number" min="1" max={Math.max(p.stock, 1)} bind:value={qty} disabled={p.stock <= 0} />
			</div>
			<div class="btn-row">
				<button class="btn btn-primary" on:click={add} disabled={p.stock <= 0}>Add to cart</button>
				<a class="btn btn-ghost" href="/cart">View cart</a>
			</div>
			{#if msg}<p class={isErr ? 'alert-error' : 'alert-ok'}>{msg}</p>{/if}
		</div>
	</div>
{:catch e}
	<p class="alert-error">Failed: {e.message}</p>
{/await}
