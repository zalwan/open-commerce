<script lang="ts">
	import { page } from '$app/stores';
	import { api, formatIDR } from '$lib/api';
	import { statusClass } from '$lib/ui';

	const id: string = $page.params.id ?? '';

</script>

<div class="panel success-hero">
	<div class="check">🎉</div>
	<h1>Thank you! Your order is in.</h1>
	{#await api.order(id) then o}
		<p>Order <strong>{o.id}</strong> · <span class={statusClass(o.status)}>{o.status}</span></p>
		<ul class="rows" style="text-align:left">
			{#each o.items as it}
				<li class="row">
					<div class="grow">{it.name} <span class="muted">× {it.qty}</span></div>
					<div>{formatIDR(it.priceMinor * it.qty)}</div>
				</li>
			{/each}
		</ul>
		<p class="total">Total: {formatIDR(o.totalMinor)}</p>
		<a class="btn btn-primary" href="/">Shop again</a>
	{:catch e}
		<p class="alert-error">Failed to load order: {e.message}</p>
	{/await}
</div>
