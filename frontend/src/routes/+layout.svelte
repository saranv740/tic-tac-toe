<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/Header.svelte';

	let { children } = $props();
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
