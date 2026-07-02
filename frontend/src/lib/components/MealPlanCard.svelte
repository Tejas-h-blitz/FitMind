<script lang="ts">
	import type { MealDay } from '$lib/api';
	import { ChevronDown, ChevronUp, Flame, Check, Info } from 'lucide-svelte';

	let { dayPlan } = $props<{ dayPlan: MealDay }>();

	let isExpanded = $state(false);
	let activeTab = $state<'breakfast' | 'lunch' | 'dinner' | 'snacks'>('breakfast');

	const tabs = [
		{ id: 'breakfast', label: 'Breakfast' },
		{ id: 'lunch', label: 'Lunch' },
		{ id: 'dinner', label: 'Dinner' },
		{ id: 'snacks', label: 'Snacks' }
	] as const;

	let currentMeal = $derived(dayPlan.meals[activeTab]);
</script>

<div class="rounded-xl border border-slate-800 bg-slate-900/30 overflow-hidden shadow-md transition-all duration-300 {isExpanded ? 'border-slate-700/80 bg-slate-900/50' : 'hover:border-slate-800/80 hover:bg-slate-900/40'}">
	<!-- Accordion Header -->
	<button
		onclick={() => isExpanded = !isExpanded}
		class="w-full flex items-center justify-between px-5 py-4 text-slate-100 font-bold text-sm cursor-pointer select-none text-left"
	>
		<div class="flex items-center gap-3">
			<span class="inline-flex items-center justify-center h-7 w-7 rounded-full bg-emerald-500/10 text-emerald-400 text-xs">
				{dayPlan.day[0]}
			</span>
			<div>
				<span class="block text-slate-200 font-bold">{dayPlan.day}</span>
				{#if !isExpanded}
					<span class="block text-[11px] text-slate-500 font-normal mt-0.5 truncate max-w-xs sm:max-w-md">
						Breakfast: {dayPlan.meals.breakfast.name} • {dayPlan.meals.breakfast.calories} kcal
					</span>
				{/if}
			</div>
		</div>

		<div class="flex items-center gap-3">
			<span class="text-xs font-semibold text-slate-400 hidden sm:inline">
				{dayPlan.meals.breakfast.calories + dayPlan.meals.lunch.calories + dayPlan.meals.dinner.calories + dayPlan.meals.snacks.calories} kcal
			</span>
			{#if isExpanded}
				<ChevronUp class="h-4.5 w-4.5 text-slate-400" />
			{:else}
				<ChevronDown class="h-4.5 w-4.5 text-slate-400" />
			{/if}
		</div>
	</button>

	<!-- Collapsible Content -->
	{#if isExpanded}
		<div class="border-t border-slate-900 p-5 bg-slate-950/20 space-y-5 animate-slideDown">
			<!-- Tabs Bar -->
			<div class="flex gap-1.5 border-b border-slate-900 pb-2 overflow-x-auto scrollbar-none">
				{#each tabs as tab}
					<button
						onclick={() => activeTab = tab.id}
						class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all whitespace-nowrap cursor-pointer {activeTab === tab.id ? 'bg-emerald-600 text-white shadow-sm' : 'text-slate-400 hover:text-white hover:bg-slate-900/60'}"
					>
						{tab.label}
					</button>
				{/each}
			</div>

			<!-- Tab Content: Active Meal -->
			{#if currentMeal}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-5 items-start">
					<!-- Left: Name + Ingredients -->
					<div class="space-y-4">
						<div>
							<h4 class="text-base font-bold text-slate-100">{currentMeal.name}</h4>
							<span class="mt-1.5 inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-amber-950/40 border border-amber-500/30 text-amber-400">
								<Flame class="h-3 w-3 fill-amber-950" />
								{currentMeal.calories} Calories
							</span>
						</div>

						<div class="space-y-2">
							<span class="block text-[11px] font-bold text-slate-500 uppercase tracking-wider">Ingredients List</span>
							<ul class="grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs text-slate-300">
								{#each currentMeal.ingredients as ingredient}
									<li class="flex items-center gap-2 bg-slate-900/40 p-2 rounded-lg border border-slate-900">
										<Check class="h-3.5 w-3.5 text-emerald-400 shrink-0" />
										<span class="truncate">{ingredient}</span>
									</li>
								{/each}
							</ul>
						</div>
					</div>

					<!-- Right: Benefits -->
					<div class="p-4 rounded-xl border border-slate-800 bg-slate-950/40 h-full flex flex-col justify-center">
						<span class="block text-[10px] font-bold text-slate-500 uppercase tracking-wider mb-2 flex items-center gap-1">
							<Info class="h-3.5 w-3.5 text-emerald-400" />
							<span>Why it helps</span>
						</span>
						<p class="text-xs text-slate-300 leading-relaxed italic">
							"{currentMeal.benefits}"
						</p>
					</div>
				</div>
			{/if}

			<!-- Daily Tip Box (highlighted in blue) -->
			<div class="p-3.5 rounded-lg bg-blue-950/20 border border-blue-500/20 text-blue-400 text-xs leading-relaxed flex items-start gap-2.5">
				<span class="text-base select-none shrink-0">💡</span>
				<div>
					<span class="font-bold block mb-0.5 text-blue-300">Actionable Daily Tip</span>
					<span>{dayPlan.daily_tip}</span>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	@keyframes slideDown {
		from { opacity: 0; transform: translateY(-5px); }
		to { opacity: 1; transform: translateY(0); }
	}
	.animate-slideDown {
		animation: slideDown 0.25s ease-out forwards;
	}
</style>
