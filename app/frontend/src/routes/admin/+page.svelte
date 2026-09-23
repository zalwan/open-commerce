<script lang="ts">
	import { api } from '$lib/api';
	import { getToken, setToken } from '$lib/session';

	let email = 'admin@shop.test';
	let password = 'admin123';
	let msg = getToken() ? 'Sudah login (token tersimpan).' : '';

	async function login() {
		try {
			const res = await api.login(email, password);
			setToken(res.token);
			msg = `Login OK sebagai ${res.role}.`;
		} catch (e) {
			msg = e instanceof Error ? e.message : 'Gagal';
		}
	}
</script>

<h1>Admin</h1>
<p>Kredensial dummy dev-only: <code>admin@shop.test / admin123</code></p>
<form on:submit|preventDefault={login}>
	<label>Email <input bind:value={email} /></label>
	<label>Password <input type="password" bind:value={password} /></label>
	<button type="submit">Login (stub)</button>
</form>
{#if msg}<p>{msg}</p>{/if}
<ul>
	<li><a href="/admin/products">Kelola produk</a></li>
	<li><a href="/admin/orders">Kelola order</a></li>
</ul>
