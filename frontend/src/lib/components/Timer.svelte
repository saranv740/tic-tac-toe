<script lang="ts">
	import { cn } from '$lib';

	interface Props {
		remainingSeconds?: number;
		class?: string;
	}

	let { remainingSeconds = 25, ...props }: Props = $props();

	let formattedTime = $derived.by(() => {
		const m = Math.floor(Math.max(0, remainingSeconds) / 60);
		const s = Math.max(0, remainingSeconds) % 60;
		return `${m}:${s.toString().padStart(2, '0')}`;
	});
</script>

<div class={cn('flex flex-col items-center justify-center gap-4', props.class)}>
	<div
		class="flex items-center justify-center rounded-2xl bg-surface-container-lowest px-12 py-5 shadow-inner"
	>
		<span
			class="text-2xl font-bold tracking-widest text-on-secondary-container tabular-nums md:text-4xl"
		>
			{formattedTime}
		</span>
	</div>
	<span class="text-[0.75rem] font-medium tracking-[0.1em] text-on-surface-variant uppercase">
		Timed Match
	</span>
</div>
