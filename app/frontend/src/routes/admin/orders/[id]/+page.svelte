<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api, formatIDR } from '$lib/api';
	import type { Order } from '$lib/api';
	import { statusClass } from '$lib/ui';
	import { getToken } from '$lib/session';

	const id: string = $page.params.id ?? '';
	let order: Order | null = null;
	let error = '';
	let msg = '';

	const NEXT: Record<string, string[]> = {
		paid: ['shipped', 'cancelled'],
		shipped: ['done'],
		payment_failed: ['cancelled']
	};

	async function reload() {
		const res = await api.order(id);
		order = res;
	}

	onMount(async () => {
		try {
			await reload();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed';
		}
	});

	async function setStatus(next: string) {
		try {
			const t = getToken();
			if (!t) throw new Error('Not logged in — sign in at /admin first');
			order = await api.setOrderStatus(t, id, next);
			msg = `Order → ${next}.`;
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Failed to change status';
		}
	}
</script>

<a class="back" href="/admin/orders">← Orders</a>
<h1>Order detail</h1>
{#if error}<p class="alert-error">{error}</p>{/if}
{#if msg}<p class="alert-ok">{msg}</p>{/if}
{#if !order && !error}
	<p class="muted">Loading…</p>
{:else if order}
	<div class="panel">
		<p><strong>{order.id}</strong> · <span class={statusClass(order.status)}>{order.status}</span></p>
		<p class="muted">{order.email} · {order.items.length} items</p>
		<ul class="rows">
			{#each order.items as it}
				<li class="row">
					<div class="grow">{it.name} <span class="muted">× {it.qty}</span></div>
					<div>{formatIDR(it.priceMinor * it.qty)}</div>
				</li>
			{/each}
		</ul>
		<p class="total">Total: {formatIDR(order.totalMinor)}</p>
		{#if order.paymentTx}<p class="muted">Payment: <code>{order.paymentTx}</code></p>{/if}
		{#if NEXT[order.status]}
			<div class="btn-row">
				{#each NEXT[order.status] as next}
					<button class="btn btn-primary" on:click={() => setStatus(next)}>→ {next}</button>
				{/each}
			</div>
		{/if}
	</div>
{/if}
