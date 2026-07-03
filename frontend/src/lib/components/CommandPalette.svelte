<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Brain, LayoutDashboard, Upload, Utensils, Dumbbell, Search, Keyboard } from 'lucide-svelte';

	let isOpen = $state(false);
	let search = $state('');
	let selectedIndex = $state(0);

	const actions = [
		{ id: 'dashboard', label: 'Go to Dashboard', category: 'Navigation', icon: LayoutDashboard, action: () => goto('/dashboard') },
		{ id: 'upload', label: 'Go to Upload Files', category: 'Navigation', icon: Upload, action: () => goto('/upload') },
		{ id: 'meal-plan', label: 'Go to Meal Planner', category: 'Navigation', icon: Utensils, action: () => goto('/meal-plan') },
		{ id: 'workout', label: 'Go to Workout Planner', category: 'Navigation', icon: Dumbbell, action: () => goto('/workout') },
		{ id: 'quick-upload', label: 'Quick Action: Upload PDF', category: 'Actions', icon: Upload, action: () => goto('/upload') },
		{ id: 'quick-meal', label: 'Quick Action: Generate Meal Plan', category: 'Actions', icon: Utensils, action: () => goto('/meal-plan') },
		{ id: 'quick-workout', label: 'Quick Action: Generate Workout Plan', category: 'Actions', icon: Dumbbell, action: () => goto('/workout') }
	];

	let filteredActions = $derived(
		actions.filter(a => a.label.toLowerCase().includes(search.toLowerCase()))
	);

	function handleKeydown(e: KeyboardEvent) {
		if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
			e.preventDefault();
			isOpen = !isOpen;
			search = '';
			selectedIndex = 0;
		}

		if (!isOpen) return;

		if (e.key === 'Escape') {
			isOpen = false;
		} else if (e.key === 'ArrowDown') {
			e.preventDefault();
			selectedIndex = filteredActions.length > 0 ? (selectedIndex + 1) % filteredActions.length : 0;
		} else if (e.key === 'ArrowUp') {
			e.preventDefault();
			selectedIndex = filteredActions.length > 0 ? (selectedIndex - 1 + filteredActions.length) % filteredActions.length : 0;
		} else if (e.key === 'Enter') {
			e.preventDefault();
			if (filteredActions[selectedIndex]) {
				filteredActions[selectedIndex].action();
				isOpen = false;
			}
		}
	}

	onMount(() => {
		window.addEventListener('keydown', handleKeydown);
		return () => window.removeEventListener('keydown', handleKeydown);
	});
</script>

{#if isOpen}
	<!-- Backdrop overlay -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		onclick={() => isOpen = false}
		class="fixed inset-0 z-[100] bg-slate-950/70 backdrop-blur-md flex items-start justify-center pt-[15vh] px-4 animate-fade-in"
	>
		<!-- Command Panel -->
		<div
			onclick={(e) => e.stopPropagation()}
			class="w-full max-w-lg bg-slate-955/90 border border-slate-850 rounded-2xl shadow-2xl overflow-hidden backdrop-blur-xl relative animate-scale-up"
		>
			<!-- Search input -->
			<div class="flex items-center px-4 border-b border-slate-900">
				<Search class="h-4 w-4 text-slate-500 shrink-0" />
				<input
					type="text"
					bind:value={search}
					placeholder="Search command palette..."
					class="w-full px-3 py-4 bg-transparent text-slate-205 text-xs font-semibold placeholder-slate-600 focus:outline-none"
					autofocus
				/>
				<span class="text-[9px] font-bold text-slate-500 bg-slate-900 border border-slate-850 px-2 py-0.5 rounded-md shrink-0 uppercase tracking-wider select-none">ESC</span>
			</div>

			<!-- Actions List -->
			<div class="max-h-[300px] overflow-y-auto p-2 space-y-1 scrollbar-none">
				{#if filteredActions.length === 0}
					<div class="text-center py-8 text-xs text-slate-500 font-bold select-none">
						No results found for "{search}"
					</div>
				{:else}
					{#each filteredActions as action, idx}
						<!-- Category Label -->
						{#if idx === 0 || filteredActions[idx - 1].category !== action.category}
							<div class="text-[8px] font-extrabold uppercase tracking-widest text-slate-500 px-3 py-1.5 select-none mt-1">
								{action.category}
							</div>
						{/if}

						<button
							onclick={() => { action.action(); isOpen = false; }}
							class="w-full flex items-center justify-between px-3 py-2.5 rounded-xl text-left text-xs font-bold transition-all cursor-pointer select-none
								{idx === selectedIndex
									? 'bg-emerald-600 text-white shadow-md'
									: 'text-slate-350 hover:bg-slate-900/60 hover:text-slate-100'}"
						>
							<div class="flex items-center gap-3">
								<svelte:component this={action.icon} class="h-4 w-4 shrink-0" />
								<span>{action.label}</span>
							</div>

							{#if idx === selectedIndex}
								<span class="text-[10px] text-emerald-100 select-none">↵</span>
							{/if}
						</button>
					{/each}
				{/if}
			</div>

			<!-- Bottom Hints Bar -->
			<div class="px-4 py-2 bg-slate-950 border-t border-slate-900 flex justify-between items-center text-[9px] text-slate-500 font-bold select-none">
				<span class="flex items-center gap-1.5">
					<Keyboard class="h-3.5 w-3.5" />
					<span>Use arrow keys to navigate, enter to select</span>
				</span>
				<span>FitMind Commands</span>
			</div>
		</div>
	</div>
{/if}

<style>
	@keyframes fadeIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}
	@keyframes scaleUp {
		from { opacity: 0; transform: scale(0.97); }
		to { opacity: 1; transform: scale(1); }
	}
	.animate-fade-in {
		animation: fadeIn 0.15s ease-out forwards;
	}
	.animate-scale-up {
		animation: scaleUp 0.15s cubic-bezier(0.16, 1, 0.3, 1) forwards;
	}
</style>
