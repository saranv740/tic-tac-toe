<script lang="ts">
	import CloseIcon from '$lib/assets/CloseIcon.svelte';
	import type { MatchType } from '$lib/types';
	import Button from './Button.svelte';

	interface RoomItem {
		matchId: string;
		name: string;
		mode: MatchType;
	}

	interface Props {
		handleClose: () => void;
		handleJoin: (roomID: string) => void;
		fetchOpenRooms: () => Promise<RoomItem[]>;
	}

	const { handleClose, handleJoin, fetchOpenRooms }: Props = $props();

	let rooms = $state<RoomItem[]>([]);
	let isLoading = $state(true);
	let error = $state('');

	// Fetch rooms when component mounts
	$effect(() => {
		fetchOpenRooms()
			.then((r) => {
				rooms = r;
				isLoading = false;
			})
			.catch((err) => {
				console.error('Failed to fetch rooms:', err);
				error = 'Failed to load rooms.';
				isLoading = false;
			});
	});
</script>

{#snippet listItem(room: RoomItem)}
	<li class="flex items-center justify-between rounded-2xl bg-surface p-4">
		<p class="flex flex-col">
			<span class="text-body-md font-medium text-on-surface">{room.name}</span>
			<span class="text-[10px] tracking-widest text-on-surface-variant uppercase">{room.mode}</span>
		</p>
		<Button
			variant="secondary"
			class="max-w-16 bg-surface-container-high py-1.5 text-xs font-bold tracking-[0.6px] uppercase"
			onclick={() => {
				handleJoin(room.matchId);
			}}
		>
			Join
		</Button>
	</li>
{/snippet}

<div class="flex h-full w-full items-center justify-center">
	<div
		class="h-120 w-[90%] max-w-110 overflow-hidden rounded-2xl border border-solid border-outline-variant bg-surface-container-low"
	>
		<div class="flex items-center justify-between bg-surface-container-high px-4 py-4">
			<p class="text-lg font-bold text-on-surface">Available rooms</p>
			<button
				type="button"
				onclick={handleClose}
				class="cursor-pointer bg-transparent p-1"
				aria-label="Close join rooms dialog"
			>
				<CloseIcon />
			</button>
		</div>

		<ul class="flex max-h-105 flex-col gap-4 overflow-scroll p-4">
			{#if isLoading}
				<li class="py-8 text-center text-on-surface-variant">Loading rooms...</li>
			{:else if error}
				<li class="py-8 text-center text-red-400">{error}</li>
			{:else if rooms.length === 0}
				<li class="py-8 text-center text-on-surface-variant">No open rooms found.</li>
			{:else}
				{#each rooms as room}
					{@render listItem(room)}
				{/each}
			{/if}
		</ul>
	</div>
</div>
