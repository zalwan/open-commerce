<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { getSessionId } from '$lib/session';

	let email = '';
	let error = '';
	let loading = false;

	async function submit() {
		loading = true;
		error = '';
		try {
			const order = await api.checkout(getSessionId(), email);
			await goto(`/orders/${order.id}`);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Checkout gagal';
		} finally {
			loading = false;
		}
	}
</script>

<h1>Checkout</h1>
<p>Pembayaran via <em>stub mock</em> — tanpa uang asli.</p>
<form on:submit|preventDefault={submit}>
	<label>Email <input type="email" required bind:value={email} placeholder="kamu@mail.test" /></label>
	<button type="submit" disabled={loading}>{loading ? 'Memproses…' : 'Bayar (mock)'}</button>
</form>
{#if error}<p style="color:red">{error}</p>{/if}
