<script lang="ts">
	import { cn } from '$lib';
	import type { Turn } from '$lib/types';

	interface Props {
		name: string;
		role: Turn;
		isActive: boolean;
		points: string;
		class?: string;
	}
	const { name, role, isActive, points, ...props }: Props = $props();
</script>

<div
	class={cn(
		'w-full max-w-100 rounded-lg border-l-4 border-solid bg-surface-container-high p-6 text-on-surface opacity-50',
		{
			'border-primary': isActive && role === 'X',
			'border-on-secondary-container': isActive && role === 'O',
			'opacity-100': isActive
		},
		props.class
	)}
>
	<p
		class={cn('mb-2 text-[10px] font-medium tracking-[1px] uppercase', {
			'text-primary': isActive && role === 'X',
			'text-on-secondary-container': isActive && role === 'O'
		})}
	>
		current turn
	</p>
	<p class="mb-1 text-base font-bold text-on-surface">{name} ({role})</p>
	<p class="tracking-wide text-on-surface-variant">POINTS: {points}</p>
</div>
