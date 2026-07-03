<script lang="ts">
	import type { SourceChunk } from '$lib/api';
	import { ChevronDown, ChevronUp, FileSearch } from 'lucide-svelte';

	let { source } = $props<{ source: SourceChunk }>();
	let isOpen = $state(false);
</script>

<div class="rounded-xl border border-slate-850 bg-slate-950/20 text-slate-350 text-xs overflow-hidden transition-all duration-200 hover:border-slate-700/60">
	<button
		onclick={() => isOpen = !isOpen}
		class="w-full flex justify-between items-center gap-2 px-3 py-2 text-left hover:bg-slate-900/30 cursor-pointer select-none transition-colors"
	>
		<div class="flex items-center gap-2 text-emerald-400 font-bold">
			<FileSearch class="h-3.5 w-3.5 text-emerald-450" />
			<span>Source (Page {source.page_number})</span>
		</div>
		<div class="text-slate-500">
			{#if isOpen}
				<ChevronUp class="h-3.5 w-3.5" />
			{:else}
				<ChevronDown class="h-3.5 w-3.5" />
			{/if}
		</div>
	</button>

	{#if isOpen}
		<div class="px-3 pb-3 pt-1 border-t border-slate-900 bg-slate-950/40 text-slate-400 leading-relaxed max-h-48 overflow-y-auto font-mono text-[10px] whitespace-pre-wrap">
			{source.text.trim()}
		</div>
	{/if}
</div>
