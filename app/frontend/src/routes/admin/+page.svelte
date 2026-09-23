<script lang="ts">
	import { api } from '$lib/api';
	import { getToken, setToken } from '$lib/session';

	let email = 'admin@shop.test';
	let password = 'admin123';
	let msg = getToken() ? 'Sudah login (token tersimpan).' : '';
	let isErr = false;

	async function login() {
		try {
			const res = await api.login(email, password);
			setToken(res.token);
			msg = `Login OK sebagai ${res.role}.`;
			isErr = false;
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Gagal';
			isErr = true;
		}
	}
</script>

<h1>Admin</h1>
<div class="panel">
	<h2>Masuk</h2>
	<p class="muted" style="margin-top:0">Kredensial dummy dev-only: <code>admin@shop.test / admin123</code></p>
	<form on:submit|preventDefault={login}>
		<label class="field">Email<input class="input" bind:value={email} /></label>
		<label class="field">Password<input class="input" type="password" bind:value={password} /></label>
		<button class="btn btn-primary" type="submit">Login (stub)</button>
	</form>
	{#if msg}<p class={isErr ? 'alert-error' : 'alert-ok'}>{msg}</p>{/if}
</div>

<div class="panel">
	<h2>Kelola</h2>
	<div class="btn-row">
		<a class="btn btn-ghost" href="/admin/products">📦 Produk</a>
		<a class="btn btn-ghost" href="/admin/orders">🧾 Order</a>
	</div>
</div>
