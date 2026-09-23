<script lang="ts">
	import { page } from '$app/stores';
	import { api, formatIDR } from '$lib/api';

	const id: string = $page.params.id ?? '';
</script>

<h1>Order sukses 🎉</h1>
{#await api.order(id) then o}
	<p>Order <strong>{o.id}</strong> — status <strong>{o.status}</strong></p>
	<ul>
		{#each o.items as it}
			<li>{it.name} × {it.qty}</li>
		{/each}
	</ul>
	<p>Total: <strong>{formatIDR(o.totalMinor)}</strong></p>
	<a href="/">Belanja lagi</a>
{:catch e}
	<p style="color:red">Gagal memuat order: {e.message}</p>
{/await}
