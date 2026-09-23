<script lang="ts">
	import { api, formatIDR } from '$lib/api';
	import { getToken, setToken } from '$lib/session';

	let email = 'admin@shop.test';
	let password = 'admin123';
	let msg = getToken() ? 'Already logged in (token stored).' : '';
	let isErr = false;
	let stats: { products: number; orders: number; paidOrders: number; revenueMinor: number; lowStock: number } | null = null;

	async function loadStats() {
		const t = getToken();
		if (!t) {
			stats = null;
			return;
		}
		try {
			stats = await api.adminStats(t);
		} catch {
			stats = null;
		}
	}

	async function login() {
		try {
			const res = await api.login(email, password);
			setToken(res.token);
			msg = `Logged in as ${res.role}.`;
			isErr = false;
			await loadStats();
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Failed';
			isErr = true;
		}
	}

	if (getToken()) {
		void loadStats();
	}
</script>

<h1>Admin</h1>
<div class="panel">
	<h2>Sign in</h2>
	<p class="muted" style="margin-top:0">Dev-only dummy credentials: <code>admin@shop.test / admin123</code></p>
	<form on:submit|preventDefault={login}>
		<label class="field">Email<input class="input" bind:value={email} /></label>
		<label class="field">Password<input class="input" type="password" bind:value={password} /></label>
		<button class="btn btn-primary" type="submit">Login (stub)</button>
	</form>
	{#if msg}<p class={isErr ? 'alert-error' : 'alert-ok'}>{msg}</p>{/if}
</div>

<div class="panel">
	<h2>Manage</h2>
	{#if stats}
		<div class="statgrid">
			<div class="stat"><span class="stat-num">{stats.products}</span><span class="stat-label">Products</span></div>
			<div class="stat"><span class="stat-num">{stats.orders}</span><span class="stat-label">Orders</span></div>
			<div class="stat"><span class="stat-num">{formatIDR(stats.revenueMinor)}</span><span class="stat-label">Revenue</span></div>
			<div class="stat"><span class="stat-num">{stats.lowStock}</span><span class="stat-label">Low stock</span></div>
		</div>
	{/if}
	<div class="btn-row">
		<a class="btn btn-ghost" href="/admin/products">📦 Products</a>
		<a class="btn btn-ghost" href="/admin/orders">🧾 Orders</a>
	</div>
</div>

<style>
	.statgrid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
		gap: 0.75rem;
		margin-bottom: 1rem;
	}
	.stat {
		border: 1px solid var(--border);
		border-radius: 8px;
		padding: 0.7rem;
		display: flex;
		flex-direction: column;
		background: rgb(255 255 255 / 0.02);
	}
	.stat-num {
		font-family: var(--mono);
		font-weight: 700;
		font-size: 1.2rem;
		color: var(--brand-dim);
	}
	.stat-label {
		font-size: 0.8rem;
		color: var(--muted);
	}
</style>
