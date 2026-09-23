<script lang="ts">
	import { api, formatIDR } from '$lib/api';
	import type { Order } from '$lib/api';
	import { statusClass } from '$lib/ui';

	let email = '';
	let orders: Order[] | null = null;
	let error = '';
	let loading = false;


	async function submit() {
		loading = true;
		error = '';
		try {
			orders = await api.myOrders(email);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed';
			orders = null;
		} finally {
			loading = false;
		}
	}
</script>

<h1>Track orders</h1>
<div class="panel">
	<p class="muted" style="margin-top:0">Enter the email you checked out with to see your orders.</p>
	<form on:submit|preventDefault={submit}>
		<label class="field">Email<input class="input" type="email" required bind:value={email} placeholder="you@mail.test" /></label>
		<button class="btn btn-primary" type="submit" disabled={loading}>{loading ? 'Looking up…' : 'Find my orders'}</button>
	</form>
	{#if error}<p class="alert-error">{error}</p>{/if}
</div>

{#if orders !== null}
	{#if orders.length === 0}
		<p class="muted">No orders for this email.</p>
	{:else}
		<ul class="rows">
			{#each orders as o}
				<li class="row">
					<div class="grow">
						<a href={`/orders/${o.id}`}><strong>{o.id}</strong></a><br />
						<span class="muted">{o.items.length} items · {formatIDR(o.totalMinor)}</span>
						<span class={statusClass(o.status)}>{o.status}</span>
					</div>
				</li>
			{/each}
		</ul>
	{/if}
{/if}
