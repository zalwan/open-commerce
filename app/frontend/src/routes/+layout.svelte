<script lang="ts">
	import '../app.css';
</script>

<header class="topbar">
	<div class="topbar-inner">
		<a class="brand" href="/">Open<span>Commerce</span></a>
		<nav class="nav">
			<a href="/">Catalog</a>
			<a href="/cart">Cart</a>
			<a href="/admin">Admin</a>
		</nav>
	</div>
</header>
<main class="container">
	<slot />
</main>
<footer class="footer">Open Commerce · open source single-vendor · stub payments (no real money)</footer>
