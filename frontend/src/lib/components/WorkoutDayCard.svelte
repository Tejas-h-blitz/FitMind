<script lang="ts">
	import type { WorkoutDay } from '$lib/api';
	import ExerciseCard from './ExerciseCard.svelte';
	import { Clock, Dumbbell, ShieldAlert, Award, ChevronDown, ChevronUp } from 'lucide-svelte';

	let { dayWorkout } = $props<{ dayWorkout: WorkoutDay }>();

	let isExpanded = $state(false);
</script>

{#if dayWorkout.is_rest_day}
	<!-- Rest Day Card: Grey, clean recovery card -->
	<div class="rounded-2xl border border-slate-850 bg-slate-950/20 p-6 flex flex-col justify-center items-center text-center space-y-4 h-full min-h-[220px] shadow-lg backdrop-blur-sm">
		<div class="p-3.5 rounded-xl bg-slate-900/60 border border-slate-800 text-slate-500 shadow-inner">
			<Award class="h-6 w-6" />
		</div>
		<div>
			<h3 class="font-extrabold text-slate-200 text-base tracking-tight">{dayWorkout.day}</h3>
			<span class="inline-flex items-center mt-1.5 px-2.5 py-0.5 rounded-full text-[9px] font-extrabold bg-slate-950 border border-slate-900 text-slate-400 uppercase tracking-wider">
				{dayWorkout.focus || 'Recovery Day'}
			</span>
		</div>
		<p class="text-xs text-slate-500 max-w-[220px] leading-relaxed italic font-medium">
			"Rest is where the muscle grows. Keep active with light walking or mobility stretching."
		</p>
	</div>
{:else}
	<!-- Active Training Day Card -->
	<div class="rounded-2xl border border-slate-850 bg-slate-950/30 overflow-hidden shadow-lg transition-all duration-300 {isExpanded ? 'border-emerald-500/25 bg-slate-950/50 shadow-xl' : 'hover:border-slate-800/80 hover:bg-slate-950/45 hover:-translate-y-0.5'}">
		<!-- Card Header -->
		<button
			onclick={() => isExpanded = !isExpanded}
			class="w-full flex items-center justify-between px-5 py-4 cursor-pointer text-left select-none"
		>
			<div class="flex items-center gap-3">
				<span class="inline-flex items-center justify-center h-8 w-8 rounded-lg bg-emerald-500/10 border border-emerald-500/10 text-emerald-400 text-xs font-black uppercase">
					{dayWorkout.day[0] + dayWorkout.day[1]}
				</span>
				<div>
					<span class="block text-slate-200 font-extrabold text-sm sm:text-base tracking-tight">{dayWorkout.day}</span>
					<span class="block text-[11px] text-slate-450 font-bold mt-0.5">{dayWorkout.focus}</span>
				</div>
			</div>

			<div class="flex items-center gap-3">
				<span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-slate-950/60 border border-slate-900 text-slate-450 text-[10px] font-bold">
					<Clock class="h-3 w-3 text-emerald-400" />
					<span>{dayWorkout.estimated_duration_minutes}m</span>
				</span>
				<span class="text-xs font-bold text-slate-400 bg-slate-950/60 px-2 py-1 rounded-lg border border-slate-900 hidden sm:inline">
					{dayWorkout.exercises?.length || 0} Exercises
				</span>
				{#if isExpanded}
					<ChevronUp class="h-4.5 w-4.5 text-slate-450" />
				{:else}
					<ChevronDown class="h-4.5 w-4.5 text-slate-450" />
				{/if}
			</div>
		</button>

		<!-- Card Body (Expandable) -->
		{#if isExpanded}
			<div class="border-t border-slate-900 p-5 bg-slate-950/30 space-y-6 animate-slideDown">
				<!-- Warmup Section -->
				{#if dayWorkout.warmup}
					<div class="p-4 bg-emerald-950/10 border border-emerald-500/15 rounded-xl shadow-inner">
						<span class="block text-[9px] font-extrabold text-emerald-400 uppercase tracking-wider mb-1.5">🔴 Warmup Routine</span>
						<p class="text-xs text-slate-300 leading-relaxed font-semibold">{dayWorkout.warmup}</p>
					</div>
				{/if}

				<!-- Exercises List -->
				<div class="space-y-3">
					<span class="block text-[9px] font-extrabold text-slate-500 uppercase tracking-wider pl-1">Exercise List</span>
					<div class="space-y-3">
						{#each dayWorkout.exercises as exercise}
							<ExerciseCard {exercise} />
						{/each}
					</div>
				</div>

				<!-- Cooldown Section -->
				{#if dayWorkout.cooldown}
					<div class="p-4 bg-sky-950/10 border border-sky-500/15 rounded-xl shadow-inner">
						<span class="block text-[9px] font-extrabold text-sky-400 uppercase tracking-wider mb-1.5">🔵 Cooldown Routine</span>
						<p class="text-xs text-slate-300 leading-relaxed font-semibold">{dayWorkout.cooldown}</p>
					</div>
				{/if}
			</div>
		{/if}
	</div>
{/if}

<style>
	@keyframes slideDown {
		from { opacity: 0; transform: translateY(-3px); }
		to { opacity: 1; transform: translateY(0); }
	}
	.animate-slideDown {
		animation: slideDown 0.2s ease-out forwards;
	}
</style>
