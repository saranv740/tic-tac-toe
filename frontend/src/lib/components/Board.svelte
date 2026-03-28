<script lang="ts">
	import type { Turn } from '$lib/types';

	interface Props {
		rows?: number;
		cols?: number;
		onMove: (row: number, col: number) => void;
		currentTurn: Turn;
	}

	const { rows = 3, cols = 3, onMove, currentTurn }: Props = $props();
</script>

<div class="board-container p-12">
	<div class="decoration-top"></div>
	<div class="board" style:--row-count={rows} style:--col-count={cols}>
		{#each { length: rows }, r}
			{#each { length: cols }, c}
				<button
					type="button"
					class="h-16 w-16 bg-[#0C0D18] sm:h-32 sm:w-32"
					onclick={() => onMove(r, c)}
				>
					{currentTurn}
				</button>
			{/each}
		{/each}
	</div>
	<div class="decoration-bottom"></div>
</div>

<style>
	.board {
		display: grid;
		gap: 8px;
		grid-template-columns: repeat(var(--row-count), 1fr);
		grid-template-rows: repeat(var(--col-count), 1fr);
	}
</style>
