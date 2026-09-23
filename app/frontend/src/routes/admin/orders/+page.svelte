<script lang="ts">
	import { onMount } from 'svelte';
	import { api, formatIDR } from '$lib/api';
	import type { Order } from '$lib/api';
	import { getToken } from '$lib/session';

	let orders: Order[] = [];
	let error = '';
	let msg = '';

	const NEXT: Record<string, string> = { paid: 'shipped', shipped: 'done', payment_failed: 'cancelled' };

	async function reload() {
		const t = getToken();
		if (!t) throw new Error('Belum login — masuk dulu di /admin');
		orders = await api.adminOrders(t);
	}

	onMount(() => {
		reload().catch((e) => (error = e instanceof Error ? e.message : 'Gagal'));
	});

	async function advance(o: Order) {
		const next = NEXT[o.status];
		if (!next) return;
		try {
			const t = getToken();
			if (!t) throw new Error('Belum login');
			await api.setOrderStatus(t, o.id, next);
			msg = `Order ${o.id} → ${next}.`;
			await reload();
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Gagal ubah status';
		}
	}
</script>

<h1>Kelola Order</h1>
{#if error}<p style="color:red">{error}</p>{/if}
{#if msg}<p>{msg}</p>{/if}
{#if orders.length === 0}
	<p>Belum ada order.</p>
{:else}
	<ul style="list-style:none;padding:0;display:grid;gap:0.75rem">
		{#each orders as o}
			<li style="border:1px solid #ddd;padding:0.75rem">
				<strong>{o.id}</strong> — {o.email} — {formatIDR(o.totalMinor)} — status <strong>{o.status}</strong>
				{#if NEXT[o.status]}
					<button on:click={() => advance(o)}>→ {NEXT[o.status]}</button>
				{/if}
			</li>
		{/each}
	</ul>
{/if}
