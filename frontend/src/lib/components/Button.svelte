<script lang="ts">
	import { cn } from '$lib';
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Variants = 'primary' | 'secondary' | 'outline';
	type Props = HTMLButtonAttributes & {
		children: Snippet;
		variant: Variants;
	};

	const { children, variant = 'primary', ...props }: Props = $props();
</script>

<button
	{...props}
	class={cn(
		{
			'bg-[#aec6ff] text-on-primary': variant === 'primary',
			'bg-surface text-[#aec6ff]': variant === 'secondary',
			'border border-solid border-outline-variant bg-transparent text-on-surface-variant':
				variant === 'outline'
		},
		'inline-block w-full cursor-pointer rounded-sm py-4 text-base font-medium',
		props.class
	)}
>
	{@render children()}
</button>
