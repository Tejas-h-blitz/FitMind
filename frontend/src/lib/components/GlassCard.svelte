<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children, class: className = '', onclick = undefined } = $props<{
		children: Snippet;
		class?: string;
		onclick?: (e: MouseEvent) => void;
	}>();
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	{onclick}
	class="rounded-2xl border border-slate-850 bg-slate-950/45 p-6 backdrop-blur-md shadow-xl transition-all duration-300 {onclick ? 'cursor-pointer hover:border-emerald-500/30 hover:bg-slate-900/20 hover:-translate-y-0.5 hover:shadow-emerald-950/5' : ''} {className}"
>
	{@render children()}
</div>
