<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import Board from '$lib/components/Board.svelte';
	import Button from '$lib/components/Button.svelte';
	import PlayerCard from '$lib/components/PlayerCard.svelte';
	import Result from '$lib/components/Result.svelte';
	import Timer from '$lib/components/Timer.svelte';
	import { client, getSocket } from '$lib/services/nakama';
	import { sessionStore } from '$lib/stores/session.svelte';
	import { matchStore } from '$lib/stores/match.svelte';
	import type { MatchStatePayload, EndPayload } from '$lib/types';

	const OP_MOVE = 1;
	const OP_STATE_UPDATE = 2;
	const OP_GAME_START = 3;
	const OP_GAME_ENDED = 4;
	const OP_PLAYER_DISCONNECTED = 5;
	const OP_PLAYER_RECONNECTED = 6;

	const decoder = new TextDecoder();

	let remainingSeconds = $state(0);
	let timerInterval: ReturnType<typeof setInterval> | null = null;

	// Remaining time derived from `matchStore.deadline` (unix timestamp from server), not client-side logic.
	function startTimerFromDeadline(deadline: number) {
		stopTimer();
		if (deadline <= 0) return; // Classic mode — no timer
		const tick = () => {
			remainingSeconds = Math.max(0, deadline - Math.floor(Date.now() / 1000));
		};
		tick();
		timerInterval = setInterval(tick, 1000);
	}

	function stopTimer() {
		if (timerInterval) {
			clearInterval(timerInterval);
			timerInterval = null;
		}
	}

	async function handleMove(pos: number) {
		const matchId = sessionStore.matchId;
		if (!matchId || matchStore.gamePhase !== 'playing') return;
		// Only allow moves on my turn
		if (matchStore.turn !== sessionStore.role) return;

		const socket = await getSocket();
		// sendMatchState(matchId, opCode, data, presences?) — empty presences = broadcast to all
		await socket.sendMatchState(matchId, OP_MOVE, JSON.stringify({ position: pos }));
	}

	// Calling leaveMatch triggers MatchLeave on the server.
	async function handleQuit() {
		const matchId = sessionStore.matchId;
		if (matchId) {
			try {
				const socket = await getSocket();
				await socket.leaveMatch(matchId);
			} catch (err) {
				console.error('leaveMatch error:', err);
			}
		} else {
			console.error('No matchId found');
			sessionStore.clearMatch();
			matchStore.reset();
			stopTimer();
			goto('/');
		}
	}

	// register socket handlers first, THEN join
	// This order prevents the race condition where OpCodeGameStart could fire
	// before this page's onmatchdata handler is registered.
	onMount(async () => {
		const matchId = sessionStore.matchId;
		if (!matchId) {
			goto('/');
			return;
		}

		const socket = await getSocket();

		socket.onmatchdata = async (matchData) => {
			const text = decoder.decode(matchData.data);
			const payload = JSON.parse(text);

			switch (matchData.op_code) {
				case OP_GAME_START: {
					const state = payload as MatchStatePayload;
					matchStore.applyState(state);

					// Fetch display names for both players in one call
					// getUsers(session, ids?, usernames?, facebookIds?)
					const usersResult = await client.getUsers(sessionStore.session!, [
						state.player_x,
						state.player_o
					]);
					const users = usersResult.users ?? [];
					const xUser = users.find((u) => u.id === state.player_x);
					const oUser = users.find((u) => u.id === state.player_o);

					matchStore.setPlayers(
						{ userId: state.player_x, username: xUser?.username ?? 'Player X', wins: 0 },
						{ userId: state.player_o, username: oUser?.username ?? 'Player O', wins: 0 }
					);

					// Determine my role (X or O) based on my user ID
					sessionStore.role = sessionStore.session?.user_id === state.player_x ? 'X' : 'O';

					startTimerFromDeadline(state.deadline);
					break;
				}

				case OP_STATE_UPDATE: {
					const state = payload as MatchStatePayload;
					matchStore.applyState(state);
					startTimerFromDeadline(state.deadline);
					break;
				}

				case OP_GAME_ENDED: {
					matchStore.applyEnd(payload as EndPayload);
					stopTimer();
					break;
				}

				case OP_PLAYER_DISCONNECTED: {
					matchStore.setOpponentDisconnected(true);
					break;
				}

				case OP_PLAYER_RECONNECTED: {
					matchStore.setOpponentDisconnected(false);
					break;
				}
			}
		};

		// Now join — any server messages from here on will hit the handler above
		await socket.joinMatch(matchId);
	});

	onDestroy(() => {
		stopTimer();
	});

	const myPlayer = $derived(sessionStore.role === 'X' ? matchStore.playerX : matchStore.playerO);
	const opponentPlayer = $derived(
		sessionStore.role === 'X' ? matchStore.playerO : matchStore.playerX
	);
</script>

{#if matchStore.gamePhase === 'ended'}
	<!-- Result screen -->
	<Result
		winnerName={matchStore.result.winner === 'draw'
			? 'Nobody'
			: matchStore.result.winner === sessionStore.role
				? (myPlayer?.username ?? 'You')
				: (opponentPlayer?.username ?? 'Opponent')}
		reason={matchStore.result.reason}
		playerXStats={{
			name: matchStore.playerX?.username ?? 'Player X',
			wins: matchStore.playerX?.wins ?? 0,
			losses: 0,
			draws: 0
		}}
		playerOStats={{
			name: matchStore.playerO?.username ?? 'Player O',
			wins: matchStore.playerO?.wins ?? 0,
			losses: 0,
			draws: 0
		}}
	/>
{:else}
	<!-- Active game -->
	{#if matchStore.opponentDisconnected}
		<p class="mb-4 rounded-lg bg-yellow-500/20 px-4 py-2 text-sm text-yellow-300">
			Opponent disconnected — waiting up to 30s for them to reconnect…
		</p>
	{/if}

	<div class="mb-8 flex w-full gap-2">
		<PlayerCard
			name={matchStore.playerX?.username ?? 'Waiting…'}
			class="shrink-0 basis-1/2"
			role="X"
			isActive={matchStore.turn === 'X' && matchStore.gamePhase === 'playing'}
			points={String(matchStore.playerX?.wins ?? 0)}
		/>
		<PlayerCard
			name={matchStore.playerO?.username ?? 'Waiting…'}
			class="shrink-0 basis-1/2"
			role="O"
			isActive={matchStore.turn === 'O' && matchStore.gamePhase === 'playing'}
			points={String(matchStore.playerO?.wins ?? 0)}
		/>
	</div>

	{#if matchStore.gamePhase === 'playing' && matchStore.deadline > 0}
		<Timer {remainingSeconds} class="mb-8" />
	{/if}

	<Board onMove={handleMove} currentTurn={matchStore.turn} inputs={matchStore.board} class="mb-8" />

	<Button variant="outline" class="mb-8" onclick={handleQuit}>Quit Match</Button>
{/if}
