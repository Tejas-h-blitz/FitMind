<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		children,
		class: className = '',
		type = 'button',
		variant = 'primary',
		disabled = false,
		onclick = undefined
	} = $props<{
		children: Snippet;
		class?: string;
		type?: 'button' | 'submit' | 'reset';
		variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
		disabled?: boolean;
		onclick?: (e: MouseEvent) => void;
	}>();
</script>

<button
	{type}
	{disabled}
	{onclick}
	class="inline-flex items-center justify-center gap-2 rounded-xl text-xs font-bold transition-all duration-300 active:scale-[0.97] disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer
		{variant === 'primary'
			? 'bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-550 hover:to-teal-550 text-white shadow-md shadow-emerald-950/30 hover:shadow-lg hover:shadow-emerald-950/40 hover:-translate-y-0.5'
			: variant === 'secondary'
				? 'bg-slate-900/60 hover:bg-slate-800/80 text-slate-350 border border-slate-850 hover:border-slate-700 hover:-translate-y-0.5'
				: variant === 'ghost'
					? 'text-slate-450 hover:text-slate-200 hover:bg-slate-900/40'
					: 'bg-rose-650 hover:bg-rose-550 text-white shadow-md shadow-rose-950/30 hover:-translate-y-0.5'
		} {className}"
>
	{@render children()}
</button>
