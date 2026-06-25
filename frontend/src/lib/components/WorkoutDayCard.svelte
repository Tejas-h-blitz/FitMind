<script lang="ts">
	import type { WorkoutDay } from '$lib/api';
	import ExerciseCard from './ExerciseCard.svelte';
	import { Clock, Dumbbell, ShieldAlert, Award } from 'lucide-svelte';

	let { dayWorkout } = $props<{ dayWorkout: WorkoutDay }>();

	let isExpanded = $state(false);
</script>

{#if dayWorkout.is_rest_day}
	<!-- Rest Day Card: Grey, clean recovery card -->
	<div class="rounded-xl border border-slate-900 bg-slate-950/20 p-6 flex flex-col justify-center items-center text-center space-y-3 h-full min-h-[220px]">
		<div class="p-3 rounded-full bg-slate-900 border border-slate-800 text-slate-500">
			<Award class="h-6 w-6" />
		</div>
		<div>
			<h3 class="font-bold text-slate-300 text-base">{dayWorkout.day}</h3>
			<span class="inline-flex items-center mt-1 px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-slate-900 border border-slate-800 text-slate-400 capitalize">
				{dayWorkout.focus || 'Recovery Day'}
			</span>
		</div>
		<p class="text-xs text-slate-500 max-w-[200px] leading-relaxed italic">
			"Rest is where the muscle grows. Keep active with light walking or mobility stretching."
		</p>
	</div>
{:else}
	<!-- Active Training Day Card -->
	<div class="rounded-xl border border-slate-800 bg-slate-900/30 overflow-hidden shadow-lg transition-all duration-300 {isExpanded ? 'border-emerald-500/30 bg-slate-900/50' : 'hover:border-slate-800/80 hover:bg-slate-900/40'}">
		<!-- Card Header -->
		<button
			onclick={() => isExpanded = !isExpanded}
			class="w-full flex items-center justify-between px-5 py-4 cursor-pointer text-left select-none"
		>
			<div class="flex items-center gap-3">
				<span class="inline-flex items-center justify-center h-7 w-7 rounded-full bg-emerald-500/10 text-emerald-400 text-xs font-bold">
					{dayWorkout.day[0]}
				</span>
				<div>
					<span class="block text-slate-200 font-bold text-sm sm:text-base">{dayWorkout.day}</span>
					<span class="block text-xs text-slate-400 font-semibold mt-0.5">{dayWorkout.focus}</span>
				</div>
			</div>

			<div class="flex items-center gap-4">
				<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold bg-slate-950/60 border border-slate-850 text-slate-400">
					<Clock class="h-3 w-3 text-emerald-400" />
					{dayWorkout.estimated_duration_minutes}m
				</span>
				<span class="text-xs text-emerald-400 font-bold hidden sm:inline">
					{dayWorkout.exercises?.length || 0} Exercises
				</span>
			</div>
		</button>

		<!-- Card Body (Expandable) -->
		{#if isExpanded}
			<div class="border-t border-slate-900 p-5 bg-slate-950/20 space-y-6 animate-slideDown">
				<!-- Warmup Section -->
				{#if dayWorkout.warmup}
					<div class="p-3 bg-emerald-950/10 border border-emerald-500/15 rounded-lg">
						<span class="block text-[10px] font-bold text-emerald-400 uppercase tracking-wider mb-1">🔴 Warmup Routine</span>
						<p class="text-xs text-slate-300 leading-relaxed font-medium">{dayWorkout.warmup}</p>
					</div>
				{/if}

				<!-- Exercises List -->
				<div class="space-y-3">
					<span class="block text-[10px] font-bold text-slate-500 uppercase tracking-wider">Exercise List</span>
					<div class="space-y-3">
						{#each dayWorkout.exercises as exercise}
							<ExerciseCard {exercise} />
						{/each}
					</div>
				</div>

				<!-- Cooldown Section -->
				{#if dayWorkout.cooldown}
					<div class="p-3 bg-sky-950/10 border border-sky-500/15 rounded-lg">
						<span class="block text-[10px] font-bold text-sky-400 uppercase tracking-wider mb-1">🔵 Cooldown Routine</span>
						<p class="text-xs text-slate-300 leading-relaxed font-medium">{dayWorkout.cooldown}</p>
					</div>
				{/if}
			</div>
		{/if}
	</div>
{/if}

<style>
	:global(@keyframes slideDown) {
		from { opacity: 0; transform: translateY(-5px); }
		to { opacity: 1; transform: translateY(0); }
	}
	.animate-slideDown {
		animation: slideDown 0.25s ease-out forwards;
	}
</style>
