<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import Board from '$lib/components/Board.svelte';
	import Button from '$lib/components/Button.svelte';
	import PlayerCard from '$lib/components/PlayerCard.svelte';
	import Result from '$lib/components/Result.svelte';
	import Timer from '$lib/components/Timer.svelte';
	import { client, getSocket, resetSocket, authenticate } from '$lib/services/nakama';
	import { sessionStore } from '$lib/stores/session.svelte';
	import { matchStore, type PlayerInfo } from '$lib/stores/match.svelte';
	import type { MatchStatePayload, EndPayload } from '$lib/types';

	const OP_MOVE = 1;
	const OP_STATE_UPDATE = 2;
	const OP_GAME_START = 3;
	const OP_GAME_ENDED = 4;
	const OP_PLAYER_DISCONNECTED = 5;
	const OP_PLAYER_RECONNECTED = 6;
	const MATCH_ID_KEY = 'ttt_match_id';

	const decoder = new TextDecoder();

	let remainingSeconds = $state(0);
	let timerInterval: ReturnType<typeof setInterval> | null = null;
	let selfDisconnected = $state(false); // true while WE are reconnecting
	let reconnecting = false; // guard against concurrent reconnect attempts

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
		if (!matchId || matchStore.gamePhase !== 'playing' || matchStore.opponentDisconnected) return;
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
		}

		sessionStorage.removeItem(MATCH_ID_KEY); // clear persisted match on voluntary quit

		sessionStore.clearMatch();
		matchStore.reset();

		stopTimer();
		goto('/');
	}

	// Registers all socket message handlers on the given socket.
	// Extracted so it can be called both on initial mount and after a reconnect.
	function registerHandlers(socket: Awaited<ReturnType<typeof getSocket>>) {
		socket.onmatchdata = async (matchData) => {
			const text = decoder.decode(matchData.data);
			const payload = JSON.parse(text);

			switch (matchData.op_code) {
				case OP_GAME_START: {
					const state = payload as MatchStatePayload;
					matchStore.applyState(state);
					await fetchAndSetPlayers(state);
					startTimerFromDeadline(state.deadline);
					break;
				}

				case OP_STATE_UPDATE: {
					const state = payload as MatchStatePayload;
					matchStore.applyState(state);
					// On a tab-reopen rejoin the server sends OP_STATE_UPDATE, not OP_GAME_START.
					// If player info was lost, restore it here.
					// We check both players — they are always set/cleared together via setPlayers(),
					// so either being null means neither has been populated yet.
					if (!matchStore.playerX || !matchStore.playerO) {
						await fetchAndSetPlayers(state);
					}
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

		// When our own socket drops, attempt to reconnect automatically.
		// The server has already started the 30s forfeit timer via MatchLeave.
		socket.ondisconnect = () => {
			if (matchStore.gamePhase === 'playing') {
				console.warn('Socket disconnected — attempting reconnect…');
				attemptReconnect();
			}
		};
	}

	/**
	 * Fetches user profiles + stored stats for both players and writes them into
	 * matchStore. Also sets sessionStore.role from the state payload.
	 * Called on OP_GAME_START and on OP_STATE_UPDATE when rejoining after a tab close.
	 */
	async function fetchAndSetPlayers(state: MatchStatePayload) {
		const [usersResult, statsResult] = await Promise.all([
			client.getUsers(sessionStore.session!, [state.player_x, state.player_o]),
			client.readStorageObjects(sessionStore.session!, {
				object_ids: [
					{ collection: 'stats', key: 'tic-tac-toe', user_id: state.player_x },
					{ collection: 'stats', key: 'tic-tac-toe', user_id: state.player_o }
				]
			})
		]);

		const users = usersResult.users ?? [];
		const xUser = users.find((u) => u.id === state.player_x);
		const oUser = users.find((u) => u.id === state.player_o);

		const storageObjects = statsResult.objects ?? [];
		const xStats =
			(storageObjects.find((o) => o.user_id === state.player_x)?.value as PlayerInfo) ?? {};
		const oStats =
			(storageObjects.find((o) => o.user_id === state.player_o)?.value as PlayerInfo) ?? {};

		matchStore.setPlayers(
			{
				userId: state.player_x,
				username: xUser?.username ?? 'Player X',
				wins: xStats.wins ?? 0,
				losses: xStats.losses ?? 0,
				draws: xStats.draws ?? 0,
				current_streak: xStats.current_streak ?? 0
			},
			{
				userId: state.player_o,
				username: oUser?.username ?? 'Player O',
				wins: oStats.wins ?? 0,
				losses: oStats.losses ?? 0,
				draws: oStats.draws ?? 0,
				current_streak: oStats.current_streak ?? 0
			}
		);

		// Determine my role (X or O) based on my user ID
		sessionStore.role = sessionStore.session?.user_id === state.player_x ? 'X' : 'O';
	}

	/**
	 * Tries to re-establish the Nakama socket and re-join the current match.
	 * Called by ondisconnect and the window 'online' event.
	 * Retries up to 5 times with exponential back-off (1s, 2s, 4s, 8s, 16s).
	 */
	async function attemptReconnect() {
		if (reconnecting) return; // already retrying

		const matchId = sessionStore.matchId;
		if (!matchId || matchStore.gamePhase !== 'playing') return;

		reconnecting = true;
		selfDisconnected = true;

		const MAX_ATTEMPTS = 5;
		for (let attempt = 1; attempt <= MAX_ATTEMPTS; attempt++) {
			try {
				// Clear the dead cached socket so getSocket() creates a fresh one
				resetSocket();
				// Re-authenticate in case the session token also expired
				sessionStore.session = await authenticate();
				// Get a fresh socket connection
				const socket = await getSocket();
				// Re-register all message + disconnect handlers on the new socket
				registerHandlers(socket);
				// Re-join the match — server will cancel the 30s timer
				await socket.joinMatch(matchId);
				selfDisconnected = false;
				reconnecting = false;
				console.info('Reconnected and rejoined match.');
				return;
			} catch (err) {
				console.warn(`Reconnect attempt ${attempt}/${MAX_ATTEMPTS} failed:`, err);
				if (attempt < MAX_ATTEMPTS) {
					await new Promise((r) => setTimeout(r, 1000 * 2 ** (attempt - 1))); // 1s, 2s, 4s, 8s
				}
			}
		}

		// All attempts exhausted — give up and go home
		console.error('Could not reconnect after 5 attempts.');

		reconnecting = false;
		selfDisconnected = false;

		sessionStorage.removeItem(MATCH_ID_KEY); // clear stale match so next open goes to home

		sessionStore.clearMatch();
		matchStore.reset();

		stopTimer();
		goto('/');
	}

	// register socket handlers first, THEN join.
	// onMount only ever runs client-side, so window access here is safe.
	onMount(async () => {
		// Restore matchId from sessionStorage in case the in-memory store was
		// wiped by a tab close/reopen.
		if (!sessionStore.matchId) {
			sessionStore.session = await authenticate();
			const stored = sessionStorage.getItem(MATCH_ID_KEY);
			if (stored) sessionStore.matchId = stored;
		}

		const matchId = sessionStore.matchId;
		if (!matchId) {
			goto('/');
			return;
		}

		const socket = await getSocket();
		registerHandlers(socket);

		// 'online' fires when the browser regains network access.
		// Use it as a backup reconnect trigger (e.g. laptop lid re-open).
		window.addEventListener('online', attemptReconnect);

		// Now join — any server messages from here on will hit the handler above
		await socket.joinMatch(matchId);
	});

	onDestroy(() => {
		stopTimer();
		if (typeof window !== 'undefined') {
			window.removeEventListener('online', attemptReconnect);
		}
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
			losses: matchStore.playerX?.losses ?? 0,
			draws: matchStore.playerX?.draws ?? 0,
			current_streak: matchStore.playerX?.current_streak ?? 0
		}}
		playerOStats={{
			name: matchStore.playerO?.username ?? 'Player O',
			wins: matchStore.playerO?.wins ?? 0,
			losses: matchStore.playerO?.losses ?? 0,
			draws: matchStore.playerO?.draws ?? 0,
			current_streak: matchStore.playerO?.current_streak ?? 0
		}}
	/>
{:else}
	<!-- Active game -->
	{#if selfDisconnected}
		<p class="mb-4 rounded-lg bg-red-500/20 px-4 py-2 text-sm text-red-300">
			Connection lost — reconnecting… (the server will wait up to 30s)
		</p>
	{/if}

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

	<Board
		onMove={handleMove}
		currentTurn={matchStore.turn}
		inputs={matchStore.board}
		class="mb-8"
		isDisabled={matchStore.opponentDisconnected}
	/>

	<Button variant="outline" class="mb-8" onclick={handleQuit}>Quit Match</Button>
{/if}
