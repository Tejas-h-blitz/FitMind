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

<div class="rounded-xl border border-slate-850 bg-slate-900/10 overflow-hidden shadow-md transition-all duration-300 {isExpanded ? 'border-slate-700/80 bg-slate-950/40 shadow-xl' : 'hover:border-slate-800/80 hover:bg-slate-900/20 hover:-translate-y-0.5'}">
	<!-- Accordion Header -->
	<button
		onclick={() => isExpanded = !isExpanded}
		class="w-full flex items-center justify-between px-5 py-4 text-slate-100 font-bold text-sm cursor-pointer select-none text-left"
	>
		<div class="flex items-center gap-3">
			<span class="inline-flex items-center justify-center h-7 w-7 rounded-lg bg-emerald-500/10 border border-emerald-500/10 text-emerald-400 text-xs font-bold uppercase">
				{dayPlan.day[0] + dayPlan.day[1]}
			</span>
			<div>
				<span class="block text-slate-200 font-extrabold tracking-tight">{dayPlan.day}</span>
				{#if !isExpanded}
					<span class="block text-[10px] text-slate-500 font-bold mt-0.5 truncate max-w-xs sm:max-w-md">
						Breakfast: {dayPlan.meals.breakfast.name} • {dayPlan.meals.breakfast.calories} kcal
					</span>
				{/if}
			</div>
		</div>

		<div class="flex items-center gap-3">
			<span class="text-xs font-bold text-slate-400 bg-slate-950/60 px-2.5 py-1 rounded-lg border border-slate-900">
				{dayPlan.meals.breakfast.calories + dayPlan.meals.lunch.calories + dayPlan.meals.dinner.calories + dayPlan.meals.snacks.calories} kcal
			</span>
			{#if isExpanded}
				<ChevronUp class="h-4.5 w-4.5 text-slate-450" />
			{:else}
				<ChevronDown class="h-4.5 w-4.5 text-slate-450" />
			{/if}
		</div>
	</button>

	<!-- Collapsible Content -->
	{#if isExpanded}
		<div class="border-t border-slate-900 p-5 bg-slate-950/30 space-y-5 animate-slideDown">
			<!-- Tabs Bar -->
			<div class="flex gap-1.5 border-b border-slate-900/60 pb-2.5 overflow-x-auto scrollbar-none">
				{#each tabs as tab}
					<button
						onclick={() => activeTab = tab.id}
						class="px-4 py-2 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer hover:-translate-y-0.5 {activeTab === tab.id ? 'bg-emerald-600 text-white shadow-md' : 'text-slate-450 hover:text-slate-200 hover:bg-slate-900/40 border border-transparent'}"
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
							<h4 class="text-base font-extrabold text-slate-200 tracking-tight">{currentMeal.name}</h4>
							<span class="mt-2 inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[9px] font-extrabold bg-amber-950/45 border border-amber-500/20 text-amber-400 uppercase tracking-wider">
								<Flame class="h-3 w-3 fill-amber-950 text-amber-500" />
								{currentMeal.calories} Calories
							</span>
						</div>

						<div class="space-y-2">
							<span class="block text-[10px] font-bold text-slate-500 uppercase tracking-wider">Ingredients List</span>
							<ul class="grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs text-slate-300">
								{#each currentMeal.ingredients as ingredient}
									<li class="flex items-center gap-2 bg-slate-950/45 p-2 rounded-lg border border-slate-900/60">
										<Check class="h-3.5 w-3.5 text-emerald-400 shrink-0" />
										<span class="truncate font-semibold text-slate-350">{ingredient}</span>
									</li>
								{/each}
							</ul>
						</div>
					</div>

					<!-- Right: Benefits -->
					<div class="p-4 rounded-xl border border-slate-850 bg-slate-900/10 min-h-[140px] flex flex-col justify-center shadow-inner">
						<span class="block text-[9px] font-extrabold text-slate-500 uppercase tracking-wider mb-2 flex items-center gap-1.5">
							<Info class="h-3.5 w-3.5 text-emerald-400" />
							<span>Why it helps</span>
						</span>
						<p class="text-xs text-slate-300 leading-relaxed italic font-medium">
							"{currentMeal.benefits}"
						</p>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	@keyframes slideDown {
		from { opacity: 0; transform: translateY(-3px); }
		to { opacity: 1; transform: translateY(0); }
	}
	.animate-slideDown {
		animation: slideDown 0.2s ease-out forwards;
	}
</style>
