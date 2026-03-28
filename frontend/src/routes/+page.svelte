<script lang="ts">
	import Button from '$lib/components/Button.svelte';
	import decorativeBoard from '$lib/assets/decorative-board.svg';
</script>

<!-- <h1 class="text-center text-6xl font-bold text-on-surface uppercase">Game Time?</h1> -->
<div class="mt-20 flex w-full max-w-md flex-col gap-4">
	<Button variant="primary">Quick match</Button>
	<Button variant="secondary">Join room</Button>
	<Button variant="secondary">Create room</Button>
	<a
		href="/"
		class="inline-block w-full cursor-pointer rounded-sm border border-solid border-outline-variant bg-transparent py-4 text-center text-base font-medium text-on-surface-variant underline-offset-1 hover:underline"
	>
		Leaderboard
	</a>
</div>
<img src={decorativeBoard} class="mt-10" alt="decorative board" aria-hidden="true" />
