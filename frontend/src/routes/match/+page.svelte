<script lang="ts">
	import Board from '$lib/components/Board.svelte';
	import PlayerCard from '$lib/components/PlayerCard.svelte';
	import type { Turn } from '$lib/types';

	let turn = $state<Turn>('X');

	function nextTurn() {
		if (turn === 'X') {
			turn = 'O';
		} else {
			turn = 'X';
		}
	}
</script>

<button class="text-white" type="button" onclick={nextTurn}>click</button>
<div class="flex gap-2">
	<PlayerCard
		name="Saran 1"
		class="shrink-0 basis-1/2"
		role="X"
		isActive={turn === 'X'}
		points="1200"
	/>
	<PlayerCard name="Saran 2" role="O" isActive={turn === 'O'} points="1200" />
</div>
<Board
	onMove={(r, c) => {
		console.log(r, c);
	}}
	currentTurn={turn}
/>
