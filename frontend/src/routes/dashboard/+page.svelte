<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { api, type Document, type HealthMetric, type Goal } from '$lib/api';
	import DocumentCard from '$lib/components/DocumentCard.svelte';
	import BMITracker from '$lib/components/BMITracker.svelte';
	import GoalCard from '$lib/components/GoalCard.svelte';
import { FilePlus2, CheckCircle, Target, RefreshCw, Plus, Calendar, Flame, Brain, Dumbbell, Utensils } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let documents = $state<Document[]>([]);
	let metrics = $state<HealthMetric[]>([]);
	let goals = $state<Goal[]>([]);
	let isLoading = $state(true);

	// Goal form state
	let newGoalTitle = $state('');
	let newGoalDate = $state('');
	let isAddingGoal = $state(false);

	let poller: any = null;

	onMount(async () => {
		try {
			await loadAllData();
		} catch (err) {
			console.error(err);
		} finally {
			isLoading = false;
		}
	});

	onDestroy(() => {
		stopPolling();
	});

	async function loadAllData() {
		const [docsRes, metricsRes, goalsRes] = await Promise.all([
			api.listDocuments(),
			api.getHealthMetrics(),
			api.getGoals()
		]);

		if (docsRes.success && docsRes.data) {
			documents = docsRes.data;
			checkAndStartPolling();
		}

		if (metricsRes.success && metricsRes.data) {
			metrics = metricsRes.data;
		}

		if (goalsRes.success && goalsRes.data) {
			goals = goalsRes.data;
		}
	}

	function checkAndStartPolling() {
		const hasProcessing = documents.some(
			(d) => d.status === 'pending' || d.status === 'processing'
		);
		if (hasProcessing) {
			startPolling();
		} else {
			stopPolling();
		}
	}

	function startPolling() {
		if (poller) return;
		poller = setInterval(async () => {
			const res = await api.listDocuments();
			if (res.success && res.data) {
				documents = res.data;
				const stillProcessing = res.data.some(
					(d) => d.status === 'pending' || d.status === 'processing'
				);
				if (!stillProcessing) {
					stopPolling();
					toast.success('Document analysis completed! Your files are ready for chat.');
				}
			}
		}, 5000);
	}

	function stopPolling() {
		if (poller) {
			clearInterval(poller);
			poller = null;
		}
	}

	async function handleDeleteDoc(id: string) {
		const res = await api.deleteDocument(id);
		if (res.success) {
			documents = documents.filter((d) => d.id !== id);
			toast.success('Document deleted successfully');
			checkAndStartPolling();
		} else {
			toast.error(res.error || 'Failed to delete document');
		}
	}

	async function handleLogBMI(heightVal: number, weightVal: number) {
		const res = await api.createHealthMetric(heightVal, weightVal);
		if (res.success && res.data) {
			metrics = [...metrics, res.data];
		} else {
			throw new Error(res.error || 'Failed to record health metric');
		}
	}

	async function handleAddGoal(e: Event) {
		e.preventDefault();
		if (!newGoalTitle.trim() || isAddingGoal) return;
		isAddingGoal = true;

		try {
			const res = await api.createGoal(newGoalTitle, newGoalDate);
			if (res.success && res.data) {
				goals = [res.data, ...goals];
				toast.success('Health goal created!');
				newGoalTitle = '';
				newGoalDate = '';
			} else {
				toast.error(res.error || 'Failed to create goal');
			}
		} catch (err: any) {
			toast.error(err.message || 'Error occurred');
		} finally {
			isAddingGoal = false;
		}
	}

	async function handleToggleGoal(id: string, nextStatus: 'active' | 'completed') {
		const res = await api.updateGoalStatus(id, nextStatus);
		if (res.success) {
			goals = goals.map((g) => (g.id === id ? { ...g, status: nextStatus } : g));
		} else {
			throw new Error(res.error || 'Failed to toggle status');
		}
	}

	// Filter computed goals
	let activeGoals = $derived(goals.filter((g) => g.status === 'active'));
	let completedGoals = $derived(goals.filter((g) => g.status === 'completed'));
</script>

<svelte:head>
	<title>FitMind Dashboard - AI Health Advisor</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
	{#if isLoading}
		<div class="h-96 flex flex-col justify-center items-center">
			<div class="h-10 w-10 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
			<span class="text-sm font-semibold text-slate-400 mt-4">Loading your health stats...</span>
		</div>
	{:else}
		<!-- Dashboard Header -->
		<div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-8">
			<div>
				<h1 class="text-2xl sm:text-3xl font-extrabold text-slate-100 flex items-center gap-2">
					<Brain class="text-emerald-400" />
					<span>Health Advisor Workspace</span>
				</h1>
				<p class="text-slate-400 text-sm mt-1">Manage documents, log physical health metrics, and set milestones.</p>
			</div>
			
			<a
				href="/upload"
				class="inline-flex items-center justify-center gap-2 py-2.5 px-4 text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-500 rounded-lg shadow-md transition-all cursor-pointer"
			>
				<Plus class="h-4 w-4" />
				<span>Upload Health Data</span>
			</a>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<!-- Left and Middle Column: Documents + BMI tracker -->
			<div class="lg:col-span-2 space-y-8">
				<!-- Documents Section -->
				<section>
					<div class="flex items-center justify-between mb-4">
						<h2 class="text-lg font-bold text-slate-200 flex items-center gap-2">
							<FilePlus2 class="h-5 w-5 text-emerald-400" />
							<span>Analyzed Documents</span>
						</h2>
						{#if documents.length > 0}
							<span class="text-xs text-slate-500 font-semibold">{documents.length} files total</span>
						{/if}
					</div>

					{#if documents.length === 0}
						<div class="rounded-2xl border border-slate-800 bg-slate-900/10 p-10 text-center flex flex-col items-center">
							<FilePlus2 class="h-10 w-10 text-slate-500 mb-3" />
							<h3 class="font-bold text-slate-200 text-base">No health reports uploaded</h3>
							<p class="text-sm text-slate-400 mt-1 max-w-sm mx-auto leading-relaxed">
								Upload blood work PDFs, prescription files, diet plans or workouts to ask AI questions.
							</p>
							<a
								href="/upload"
								class="mt-5 inline-flex items-center gap-1.5 px-4 py-2 bg-emerald-950/40 hover:bg-emerald-900/40 border border-emerald-500/25 text-emerald-400 rounded-lg text-sm font-bold transition-all"
							>
								<span>Upload First PDF</span>
							</a>
						</div>
					{:else}
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
							{#each documents as document (document.id)}
								<DocumentCard {document} onDelete={handleDeleteDoc} />
							{/each}
						</div>
					{/if}
				</section>

				<!-- BMI Tracker Section -->
				<BMITracker {metrics} onLog={handleLogBMI} />
			</div>

			<!-- Right Column: AI Planners + Goals Checklist -->
			<div class="space-y-8">
				<!-- AI Planners Card -->
				<section class="rounded-xl border border-slate-800/80 bg-slate-900/30 p-6 shadow-xl backdrop-blur-sm">
					<div class="flex items-center gap-2 mb-4">
						<div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400">
							<Brain class="h-5 w-5" />
						</div>
						<h2 class="text-base font-bold text-slate-100">AI Health Planners</h2>
					</div>
					<p class="text-xs text-slate-400 mb-4 leading-relaxed">
						Generate fully personalized programs tailored to your specific health report analyses, BMI tracker log, and fitness milestones.
					</p>
					<div class="grid grid-cols-1 gap-3">
						<a
							href="/meal-plan"
							class="flex items-center justify-between p-3.5 rounded-xl border border-slate-800 bg-slate-950/40 hover:border-emerald-500/40 hover:bg-slate-900/35 transition-all group cursor-pointer"
						>
							<div class="flex items-center gap-3">
								<div class="p-2 rounded-lg bg-emerald-500/10 text-emerald-400 group-hover:bg-emerald-500/20 transition-all">
									<Utensils class="h-5 w-5" />
								</div>
								<div class="text-left">
									<span class="block text-sm font-semibold text-slate-200">Meal Planner</span>
									<span class="block text-[11px] text-slate-500">AI-generated diet based on reports</span>
								</div>
							</div>
							<span class="text-slate-500 group-hover:text-emerald-400 transition-colors">→</span>
						</a>

						<a
							href="/workout"
							class="flex items-center justify-between p-3.5 rounded-xl border border-slate-800 bg-slate-950/40 hover:border-emerald-500/40 hover:bg-slate-900/35 transition-all group cursor-pointer"
						>
							<div class="flex items-center gap-3">
								<div class="p-2 rounded-lg bg-emerald-500/10 text-emerald-400 group-hover:bg-emerald-500/20 transition-all">
									<Dumbbell class="h-5 w-5" />
								</div>
								<div class="text-left">
									<span class="block text-sm font-semibold text-slate-200">Workout Planner</span>
									<span class="block text-[11px] text-slate-500">4-week training program</span>
								</div>
							</div>
							<span class="text-slate-500 group-hover:text-emerald-400 transition-colors">→</span>
						</a>
					</div>
				</section>

				<section class="rounded-xl border border-slate-800/85 bg-slate-900/30 p-6 shadow-xl backdrop-blur-sm">
					<div class="flex items-center gap-2 mb-6">
						<div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400">
							<Target class="h-5 w-5" />
						</div>
						<h2 class="text-base font-bold text-slate-100">Health Milestones</h2>
					</div>

					<!-- Create Goal Form -->
					<form onsubmit={handleAddGoal} class="space-y-3 mb-6 bg-slate-950/40 p-4 rounded-xl border border-slate-900">
						<span class="block text-[11px] font-bold text-slate-500 uppercase tracking-wider">Set New Milestone</span>
						<div>
							<input
								type="text"
								bind:value={newGoalTitle}
								placeholder="e.g. 10k steps daily / Lose 5kg"
								required
								class="w-full px-3 py-2 text-sm rounded-lg border border-slate-800 bg-slate-950 text-slate-100 focus:outline-none focus:border-emerald-500 placeholder-slate-550"
							/>
						</div>
						<div class="relative">
							<input
								type="date"
								bind:value={newGoalDate}
								class="w-full px-3 py-2 text-sm rounded-lg border border-slate-800 bg-slate-950 text-slate-100 focus:outline-none focus:border-emerald-500"
							/>
						</div>
						<button
							type="submit"
							disabled={isAddingGoal || !newGoalTitle.trim()}
							class="w-full flex items-center justify-center gap-1.5 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg text-xs font-bold transition-all cursor-pointer disabled:opacity-50"
						>
							<Plus class="h-3.5 w-3.5" />
							<span>Add Milestone</span>
						</button>
					</form>

					<!-- Goals List -->
					<div class="space-y-4">
						{#if goals.length === 0}
							<div class="text-center py-8 text-xs text-slate-500 italic">
								<p>No goals set yet.</p>
								<p class="mt-1">Log a milestone above to start tracking achievements.</p>
							</div>
						{:else}
							<!-- Active Goals -->
							{#if activeGoals.length > 0}
								<div class="space-y-2">
									<h3 class="text-[10px] uppercase font-bold text-slate-500 tracking-wider">Active Goals ({activeGoals.length})</h3>
									{#each activeGoals as goal (goal.id)}
										<GoalCard {goal} onToggle={handleToggleGoal} />
									{/each}
								</div>
							{/if}

							<!-- Completed Goals -->
							{#if completedGoals.length > 0}
								<div class="space-y-2 pt-2 border-t border-slate-900/60">
									<h3 class="text-[10px] uppercase font-bold text-slate-500 tracking-wider">Completed ({completedGoals.length})</h3>
									{#each completedGoals as goal (goal.id)}
										<GoalCard {goal} onToggle={handleToggleGoal} />
									{/each}
								</div>
							{/if}
						{/if}
					</div>
				</section>
			</div>
		</div>
	{/if}
</div>
