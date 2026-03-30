<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/Button.svelte';
	import JoinRoom from '$lib/components/JoinRoom.svelte';
	import decorativeBoard from '$lib/assets/decorative-board.svg';
	import { client, getSocket } from '$lib/services/nakama';
	import { sessionStore } from '$lib/stores/session.svelte';
	import type { MatchLabel, MatchType } from '$lib/types';

	let joinRoomDialog: HTMLDialogElement | undefined;
	let isSearching = $state(false);
	let isCreating = $state(false);
	let searchError = $state('');

	/** Navigate to /match after storing the match ID in the session store. */
	function enterMatch(matchId: string) {
		sessionStore.matchId = matchId;
		goto('/match');
	}

	async function handleQuickMatch() {
		if (!sessionStore.session) return;
		isSearching = true;
		searchError = '';

		try {
			const socket = await getSocket();

			// Register the one-time handler BEFORE adding to matchmaker
			socket.onmatchmakermatched = (matched) => {
				// match_id is set for server-authoritative matches; token for relayed
				const matchId = matched.match_id ?? matched.token;
				if (matchId) {
					// Store match ID and navigate — the match page joins after registering handlers
					enterMatch(matchId);
				}
			};

			// query '*' = any opponent, min/max = 2 players
			await socket.addMatchmaker('*', 2, 2, { mode: 'classic' });
		} catch (err) {
			console.error('Matchmaking error:', err);
			searchError = 'Failed to start matchmaking. Is the server running?';
			isSearching = false;
		}
	}

	async function handleCreateRoom() {
		if (!sessionStore.session) return;
		isCreating = true;
		searchError = '';

		try {
			// Call our backend RPC — client.rpc already parses the JSON payload
			const response = await client.rpc(sessionStore.session, 'create_match', { mode: 'classic' });
			const data: Record<string, string> = (response.payload as Record<string, string>) || {};
			if (!data.match_id) {
				throw new Error('No match_id in RPC response');
			}

			const matchId: string = data.match_id;

			// Navigate — match page will join the socket after registering handlers
			enterMatch(matchId);
		} catch (err) {
			console.error('Create room error:', err);
			searchError = 'Failed to create room.';
		} finally {
			isCreating = false;
		}
	}

	async function fetchOpenRooms(): Promise<{ matchId: string; name: string; mode: MatchType }[]> {
		if (!sessionStore.session) return [];

		// Query server for authoritative matches where label.open == 1
		const result = await client.listMatches(
			sessionStore.session,
			20, // limit
			true, // authoritative only
			undefined, // label (exact) — we use query instead
			1, // minSize: at least 1 player inside
			1, // maxSize: only rooms waiting for a second player
			'+label.open:1' // query: open rooms
		);

		const matches = result.matches ?? [];

		return matches
			.map((m) => {
				try {
					const label: MatchLabel = JSON.parse(m.label ?? '{}');
					return {
						matchId: m.match_id,
						name: label.name ?? 'Unknown room',
						mode: label.mode ?? 'classic'
					};
				} catch {
					return null;
				}
			})
			.filter(Boolean) as { matchId: string; name: string; mode: MatchType }[];
	}

	async function handleJoin(matchId: string) {
		if (!matchId) return;
		// Navigate — match page will join after registering handlers
		enterMatch(matchId);
	}
</script>

<h1 class="text-center text-6xl font-bold text-on-surface uppercase">Game Time?</h1>
<div class="mt-20 flex w-full max-w-md flex-col gap-4">
	<!-- Quick match -->
	<Button
		type="button"
		variant="primary"
		onclick={handleQuickMatch}
		disabled={isSearching || isCreating}
	>
		{isSearching ? 'Searching...' : 'Quick match'}
	</Button>

	<!-- Join room -->
	<Button
		type="button"
		variant="secondary"
		disabled={isSearching || isCreating}
		onclick={() => joinRoomDialog?.showModal()}
	>
		Join room
	</Button>

	<!-- Create room -->
	<Button
		type="button"
		variant="secondary"
		onclick={handleCreateRoom}
		disabled={isSearching || isCreating}
	>
		{isCreating ? 'Creating...' : 'Create room'}
	</Button>

	<a
		href="/leaderboard"
		class="inline-block w-full cursor-pointer rounded-sm border border-solid border-outline-variant bg-transparent py-4 text-center text-base font-medium text-on-surface-variant underline-offset-1 hover:underline"
	>
		Leaderboard
	</a>

	{#if searchError}
		<p class="text-center text-sm text-red-400">{searchError}</p>
	{/if}
</div>

<img
	src={decorativeBoard}
	class="absolute top-[33%] right-[10%] -z-1"
	alt="decorative board"
	aria-hidden="true"
/>

<dialog bind:this={joinRoomDialog}>
	<JoinRoom handleClose={() => joinRoomDialog?.close()} {handleJoin} {fetchOpenRooms} />
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
