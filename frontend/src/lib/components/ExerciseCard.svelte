<script lang="ts">
	import type { Exercise } from '$lib/api';
	import { ChevronDown, ChevronUp, Info, HelpCircle } from 'lucide-svelte';

	let { exercise } = $props<{ exercise: Exercise }>();

	let isExpanded = $state(false);
</script>

<div class="rounded-xl border border-slate-850 bg-slate-950/20 p-4 transition-all duration-300 hover:border-slate-700/60 hover:bg-slate-900/10 shadow-sm">
	<div class="flex items-start justify-between gap-4">
		<div class="space-y-1.5">
			<h5 class="text-sm font-extrabold text-slate-100 tracking-tight">{exercise.name}</h5>
			<!-- Sets, Reps, Rest -->
			<div class="flex items-center gap-2 flex-wrap text-[10px] font-bold text-slate-450 uppercase tracking-wider">
				<span class="bg-slate-900 px-2 py-0.5 rounded border border-slate-850">{exercise.sets} sets</span>
				<span>•</span>
				<span class="bg-slate-900 px-2 py-0.5 rounded border border-slate-850">{exercise.reps} reps</span>
				<span>•</span>
				<span class="text-emerald-450 bg-emerald-950/25 px-2 py-0.5 rounded border border-emerald-500/10">{exercise.rest_seconds}s rest</span>
			</div>
		</div>

		<button
			onclick={() => isExpanded = !isExpanded}
			class="p-1 rounded-lg hover:bg-slate-900/60 text-slate-500 hover:text-slate-350 transition-all cursor-pointer"
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
	<div class="mt-3 flex items-center gap-1.5 flex-wrap">
		{#each exercise.muscle_groups as group}
			<span class="inline-flex items-center px-2 py-0.5 rounded-full text-[9px] font-extrabold bg-slate-900 border border-slate-850 text-slate-450 uppercase tracking-wider">
				{group}
			</span>
		{/each}
	</div>

	{#if isExpanded}
		<!-- Expandable: Form Instructions & Modifications -->
		<div class="mt-4 pt-3.5 border-t border-slate-900 space-y-4 text-xs text-slate-300 leading-relaxed animate-fadeIn">
			<div class="space-y-1.5 bg-slate-900/15 p-3 rounded-xl border border-slate-850">
				<span class="font-extrabold text-[9px] text-slate-400 uppercase tracking-wider flex items-center gap-1.5">
					<Info class="h-3.5 w-3.5 text-emerald-400 shrink-0" />
					<span>Form Cues & Instructions</span>
				</span>
				<p class="font-semibold text-slate-350 leading-relaxed">{exercise.instructions}</p>
			</div>

			<div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-1">
				{#if exercise.modification_easier}
					<div class="p-3 rounded-xl bg-sky-950/15 border border-sky-500/15 space-y-1.5">
						<span class="font-extrabold text-[9px] text-sky-400 uppercase tracking-wider block">🟢 Easier Alternative</span>
						<p class="text-[11px] text-slate-400 font-semibold leading-relaxed">{exercise.modification_easier}</p>
					</div>
				{/if}

				{#if exercise.modification_harder}
					<div class="p-3 rounded-xl bg-rose-950/15 border border-rose-500/15 space-y-1.5">
						<span class="font-extrabold text-[9px] text-rose-450 uppercase tracking-wider block">🔴 Harder Alternative</span>
						<p class="text-[11px] text-slate-400 font-semibold leading-relaxed">{exercise.modification_harder}</p>
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
