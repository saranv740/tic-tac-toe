<script>
	import { sessionStore } from '$lib/stores/session.svelte';
</script>

<header class="mx-auto flex max-w-5xl items-center justify-between p-2 text-on-surface">
	<a href="/" class="text-lg font-bold">TIC-TAC-TOE</a>
	<p>Welcome, {sessionStore.username}</p>
</header>
