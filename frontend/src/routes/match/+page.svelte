<script lang="ts">
	import Board from '$lib/components/Board.svelte';
	import Button from '$lib/components/Button.svelte';
	import PlayerCard from '$lib/components/PlayerCard.svelte';
	import Result from '$lib/components/Result.svelte';
	import Timer from '$lib/components/Timer.svelte';
	import { TURN_TIME } from '$lib/constants';
	import type { BoardInput, Turn } from '$lib/types';

	let turn = $state<Turn>('O');
	let boardState = $state<BoardInput>(['', 'X', '', 'O', '', 'O', '', '', 'X']);
	let remainingSeconds = $state(TURN_TIME);

	$effect(() => {
		turn;
		const intervalId = setInterval(() => {
			if (remainingSeconds > 0) {
				remainingSeconds--;
			}
			if (remainingSeconds === 0) {
				nextTurn();
			}
		}, 1000);

		return () => {
			clearInterval(intervalId);
		};
	});

	function nextTurn() {
		if (turn === 'X') {
			turn = 'O';
		} else {
			turn = 'X';
		}
		remainingSeconds = TURN_TIME;
	}
</script>

{#if false}
	<Result
		winnerName="Saran 1"
		reason="match session concluded"
		playerOStats={{ name: 'Saran_1', wins: 100, losses: 23, draws: 8 }}
		playerXStats={{ name: 'Saran_2', wins: 100, losses: 23, draws: 8 }}
	/>
{:else}
	<div class="mb-8 flex w-full gap-2">
		<PlayerCard
			name="Saran 1"
			class="shrink-0 basis-1/2"
			role="X"
			isActive={turn === 'X'}
			points="1200"
		/>
		<PlayerCard
			name="Saran 2"
			class="shrink-0 basis-1/2"
			role="O"
			isActive={turn === 'O'}
			points="1200"
		/>
	</div>
	<Timer {remainingSeconds} class="mb-8" />
	<Board
		onMove={(r) => {
			console.log(r);
		}}
		currentTurn={turn}
		inputs={boardState}
		class="mb-8"
	/>
	<Button variant="outline" class="mb-8">Quit Match</Button>
{/if}
