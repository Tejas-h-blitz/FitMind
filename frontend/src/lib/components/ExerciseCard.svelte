<script lang="ts">
	import type { Exercise } from '$lib/api';
	import { ChevronDown, ChevronUp, Info, HelpCircle } from 'lucide-svelte';

	let { exercise } = $props<{ exercise: Exercise }>();

	let isExpanded = $state(false);
</script>

<div class="rounded-xl border border-slate-900 bg-slate-950/40 p-4 transition-all duration-300 hover:border-slate-800">
	<div class="flex items-start justify-between gap-4">
		<div class="space-y-1">
			<h5 class="text-sm font-bold text-slate-100">{exercise.name}</h5>
			<!-- Sets, Reps, Rest -->
			<div class="flex items-center gap-2 flex-wrap text-[11px] font-semibold text-slate-400">
				<span>{exercise.sets} sets</span>
				<span>•</span>
				<span>{exercise.reps} reps</span>
				<span>•</span>
				<span class="text-emerald-400">{exercise.rest_seconds}s rest</span>
			</div>
		</div>

		<button
			onclick={() => isExpanded = !isExpanded}
			class="p-1 rounded-lg hover:bg-slate-900 text-slate-500 hover:text-slate-300 transition-all cursor-pointer"
			aria-label="Toggle details"
		>
			{#if isExpanded}
				<ChevronUp class="h-4 w-4" />
			{:else}
				<ChevronDown class="h-4 w-4" />
			{/if}
		</button>
	</div>

	<!-- Muscle Groups tags -->
	<div class="mt-2.5 flex items-center gap-1.5 flex-wrap">
		{#each exercise.muscle_groups as group}
			<span class="inline-flex items-center px-2 py-0.5 rounded-full text-[9px] font-bold bg-slate-900 border border-slate-800 text-slate-400 capitalize">
				{group}
			</span>
		{/each}
	</div>

	{#if isExpanded}
		<!-- Expandable: Form Instructions & Modifications -->
		<div class="mt-3.5 pt-3 border-t border-slate-900/60 space-y-3.5 text-xs text-slate-300 leading-relaxed animate-fadeIn">
			<div class="space-y-1 bg-slate-900/10 p-2.5 rounded-lg border border-slate-900/30">
				<span class="font-bold text-[10px] text-slate-500 uppercase tracking-wider flex items-center gap-1">
					<Info class="h-3.5 w-3.5 text-emerald-400 shrink-0" />
					<span>Form Cues & Instructions</span>
				</span>
				<p>{exercise.instructions}</p>
			</div>

			<div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-1">
				{#if exercise.modification_easier}
					<div class="p-2.5 rounded-lg bg-sky-950/10 border border-sky-500/10 space-y-1">
						<span class="font-bold text-[10px] text-sky-400 uppercase tracking-wider">🟢 Easier Alternative</span>
						<p class="text-[11px] text-slate-400">{exercise.modification_easier}</p>
					</div>
				{/if}

				{#if exercise.modification_harder}
					<div class="p-2.5 rounded-lg bg-rose-950/10 border border-rose-500/10 space-y-1">
						<span class="font-bold text-[10px] text-rose-400 uppercase tracking-wider">🔴 Harder Alternative</span>
						<p class="text-[11px] text-slate-400">{exercise.modification_harder}</p>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	@keyframes fadeIn {
		from { opacity: 0; transform: translateY(-3px); }
		to { opacity: 1; transform: translateY(0); }
	}
	.animate-fadeIn {
		animation: fadeIn 0.2s ease-out forwards;
	}
</style>
