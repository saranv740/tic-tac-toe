<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/Header.svelte';
	import { onMount } from 'svelte';
	import { authenticate, client } from '$lib/services/nakama';
	import { sessionStore } from '$lib/stores/session.svelte';

	let { children } = $props();

	onMount(async () => {
		try {
			const session = await authenticate();
			sessionStore.session = session;

			// Fetch the user's own account to get their username
			const account = await client.getAccount(session);
			sessionStore.username = account.user?.username ?? 'Player';
		} catch (err) {
			console.error('Nakama authentication failed:', err);
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Header />
<main
	class="mx-auto flex min-h-[calc(100vh-44px)] max-w-3xl flex-col items-center justify-center px-4"
>
	{@render children()}
</main>
