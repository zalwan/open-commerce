<script lang="ts">
	import { onMount } from 'svelte';
	import { api, formatIDR } from '$lib/api';
	import type { Order } from '$lib/api';
	import { getToken } from '$lib/session';
	import { statusClass } from '$lib/ui';

	let orders: Order[] = [];
	let error = '';
	let msg = '';

	const NEXT: Record<string, string> = { paid: 'shipped', shipped: 'done', payment_failed: 'cancelled' };


	async function reload() {
		const t = getToken();
		if (!t) throw new Error('Not logged in — sign in at /admin first');
		orders = await api.adminOrders(t);
	}

	onMount(() => {
		reload().catch((e) => (error = e instanceof Error ? e.message : 'Failed'));
	});

	async function advance(o: Order) {
		const next = NEXT[o.status];
		if (!next) return;
		try {
			const t = getToken();
			if (!t) throw new Error('Not logged in');
			await api.setOrderStatus(t, o.id, next);
			msg = `Order ${o.id} → ${next}.`;
			await reload();
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Failed to change status';
		}
	}
</script>

<a class="back" href="/admin">← Admin</a>
<h1>Manage Orders</h1>
{#if error}<p class="alert-error">{error}</p>{/if}
{#if msg}<p class="alert-ok">{msg}</p>{/if}
{#if orders.length === 0 && !error}
	<p class="muted">No orders yet.</p>
{:else}
	<ul class="rows">
		{#each orders as o}
			<li class="row">
				<div class="grow">
					<a href={`/admin/orders/${o.id}`}><strong>{o.id}</strong></a> · {o.email}<br />
					<span class="muted">{o.items.length} items · {formatIDR(o.totalMinor)}</span>
					<span class={statusClass(o.status)}>{o.status}</span>
				</div>
				{#if NEXT[o.status]}
					<button class="btn btn-primary" on:click={() => advance(o)}>→ {NEXT[o.status]}</button>
				{/if}
			</li>
		{/each}
	</ul>
{/if}
