<script lang="ts">
	import Button from '$lib/components/Button.svelte';
	import decorativeBoard from '$lib/assets/decorative-board.svg';
	import JoinRoom from '$lib/components/JoinRoom.svelte';

	let joinRoomDialog: HTMLDialogElement | undefined;
</script>

<h1 class="text-center text-6xl font-bold text-on-surface uppercase">Game Time?</h1>
<div class="mt-20 flex w-full max-w-md flex-col gap-4">
	<Button type="button" variant="primary">Quick match</Button>
	<Button
		type="button"
		variant="secondary"
		onclick={() => {
			if (joinRoomDialog) {
				joinRoomDialog.showModal();
			}
		}}
	>
		Join room
	</Button>
	<Button type="button" variant="secondary">Create room</Button>
	<a
		href="/"
		class="inline-block w-full cursor-pointer rounded-sm border border-solid border-outline-variant bg-transparent py-4 text-center text-base font-medium text-on-surface-variant underline-offset-1 hover:underline"
	>
		Leaderboard
	</a>
</div>
<img
	src={decorativeBoard}
	class="absolute top-[33%] right-[10%] -z-1"
	alt="decorative board"
	aria-hidden="true"
/>
<dialog bind:this={joinRoomDialog}>
	<JoinRoom
		handleClose={() => {
			joinRoomDialog?.close();
		}}
		handleJoin={(roomID) => {
			console.log(roomID);
		}}
	/>
</dialog>

<style>
	dialog {
		width: 100%;
		height: 100%;
		max-width: 100%;
		max-height: 100%;
		background-color: transparent;
	}

	dialog::backdrop {
		background-color: transparent;
		backdrop-filter: blur(10px);
	}
</style>
