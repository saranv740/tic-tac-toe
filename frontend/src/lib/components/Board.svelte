<script lang="ts">
	import { cn } from '$lib';
	import type { BoardInput, Turn } from '$lib/types';

	interface Props {
		rows?: number;
		cols?: number;
		inputs: BoardInput;
		onMove: (row: number) => void;
		currentTurn: Turn;
		class?: string;
		isDisabled?: boolean;
	}

	const { rows = 3, cols = 3, onMove, currentTurn, inputs, isDisabled, ...props }: Props = $props();
</script>

<div class={cn('board-container p-6 sm:p-12', props.class)}>
	<div class="decoration-top" aria-hidden="true"></div>
	<div class="board" style:--row-count={rows} style:--col-count={cols}>
		{#each inputs as input, pos}
			<button
				disabled={isDisabled}
				type="button"
				aria-label={`Board position ${pos + 1}. ${input === '' ? 'Not filled' : `Filled by ${input}`}`}
				class={cn(
					'group flex h-16 w-16 cursor-default items-center justify-center bg-surface-container-lowest text-[3.5rem] font-bold transition-all duration-200 ease-neon sm:h-32 sm:w-32',
					{
						'text-primary': input === 'X',
						'text-on-secondary-container': input === 'O',
						'cursor-pointer hover:bg-surface-bright': input === '',
						'pointer-events-none': isDisabled
					}
				)}
				onclick={() => {
					if (input === '') onMove(pos);
				}}
			>
				{#if input !== ''}
					{input}
				{:else}
					<span
						class={[
							'opacity-0 transition-opacity duration-200 ease-neon group-hover:opacity-20',
							currentTurn === 'X' ? 'text-primary' : 'text-on-secondary-container'
						]}
					>
						{currentTurn}
					</span>
				{/if}
			</button>
		{/each}
	</div>
	<div class="decoration-bottom"></div>
</div>

<style>
	.board-container {
		position: relative;
	}

	.board {
		display: grid;
		gap: 8px;
		grid-template-columns: repeat(var(--row-count), 1fr);
		grid-template-rows: repeat(var(--col-count), 1fr);
	}

	.decoration-top,
	.decoration-bottom {
		width: calc(var(--spacing) * 16);
		height: calc(var(--spacing) * 16);
		border: 3px solid #43475133;
		position: absolute;
		z-index: -1;

		@media (width >= 40rem /* 640px */) {
			width: calc(var(--spacing) * 32);
			height: calc(var(--spacing) * 32);
		}
	}

	.decoration-top {
		border-right: transparent;
		border-bottom: transparent;
		border-top-left-radius: 24px;
		top: 0;
		left: 0;
	}

	.decoration-bottom {
		border-left: transparent;
		border-top: transparent;
		border-bottom-right-radius: 24px;
		bottom: 0;
		right: 0;
	}
</style>
