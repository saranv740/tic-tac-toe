<script lang="ts">
	import { goto } from '$app/navigation';
	import { cn } from '$lib';
	import OIcon from '$lib/assets/OIcon.svelte';
	import XIcon from '$lib/assets/XIcon.svelte';
	import type { Turn, UserStats } from '$lib/types';
	import Button from './Button.svelte';

	interface Props {
		playerXStats: UserStats;
		playerOStats: UserStats;
		reason: string;
		winnerName: string;
		class?: string;
	}
	const { winnerName, reason, playerOStats, playerXStats, ...props }: Props = $props();
</script>

{#snippet statLine(key: string, value: number | string)}
	<p class="statLine mb-4">
		<span class="text-sm text-on-surface-variant">{key}</span>
		<span class="text-base font-bold text-on-surface">{value}</span>
	</p>
{/snippet}

{#snippet statCard(stat: UserStats, type: Turn)}
	<div
		class={cn(
			'w-full max-w-100 rounded-lg border-l-4 border-solid bg-surface-container-high p-6 text-on-surface',
			{
				'border-primary': type === 'X',
				'border-on-secondary-container': type === 'O'
			},
			props.class
		)}
	>
		<p
			class={cn('statLine mb-8 text-xs leading-7.5 font-medium tracking-[1px] uppercase', {
				'text-primary': type === 'X',
				'text-on-secondary-container': type === 'O'
			})}
		>
			<span>{playerXStats.name}</span>
			<span class="text-body-md">
				{#if type === 'X'}
					<XIcon />
				{:else}
					<OIcon />
				{/if}
			</span>
		</p>
		{@render statLine('Wins', stat.wins)}
		{@render statLine('Losses', stat.losses)}
		{@render statLine('Draws', stat.draws)}
	</div>
{/snippet}

<h1 class="my-4 text-center text-6xl font-bold text-on-surface">WINNER: {winnerName}</h1>
<p class="mb-10 text-sm text-on-surface-variant uppercase md:mb-20">{reason}</p>
<div class="mb-15 flex w-full flex-col items-center gap-4 md:mb-25 md:flex-row">
	<div class="w-full shrink-0 basis-1 md:basis-1/2">
		{@render statCard(playerXStats, 'X')}
	</div>
	<div class="w-full shrink-0 basis-1 md:basis-1/2">
		{@render statCard(playerOStats, 'O')}
	</div>
</div>
<Button
	onclick={() => {
		goto('/');
	}}
	type="button"
	variant="primary"
	class="mb-8 max-w-80 uppercase">Back to home</Button
>

<style>
	.statLine {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
</style>
