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

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 relative z-10">
	{#if isLoading}
		<div class="h-[70vh] flex flex-col justify-center items-center">
			<div class="h-10 w-10 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
			<span class="text-sm font-semibold text-slate-400 mt-4 animate-pulse">Loading your health stats...</span>
		</div>
	{:else}
		<!-- Dashboard Header -->
		<div class="flex flex-col md:flex-row md:justify-between md:items-center gap-4 mb-10 pb-6 border-b border-slate-900">
			<div>
				<h1 class="text-3xl font-black text-slate-100 tracking-tight flex items-center gap-3">
					<span class="p-2 rounded-xl bg-emerald-500/10 text-emerald-400">
						<Brain class="h-6 w-6" />
					</span>
					<span>Health Advisor Workspace</span>
				</h1>
				<p class="text-slate-400 text-sm mt-2 max-w-xl">Manage documents, log physical health metrics, and set milestones.</p>
			</div>
			
			<a
				href="/upload"
				class="inline-flex items-center justify-center gap-2 py-3 px-5 text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-500 rounded-xl shadow-lg shadow-emerald-950/20 hover:shadow-emerald-950/40 hover:-translate-y-0.5 transition-all duration-300 cursor-pointer"
			>
				<Plus class="h-4 w-4" />
				<span>Upload Health Data</span>
			</a>
		</div>

		<!-- Bento Grid Layout -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6 items-start">
			<!-- Bento 1: BMI Tracker (Hero Widget, spans 2 columns) -->
			<div class="md:col-span-2">
				<BMITracker {metrics} onLog={handleLogBMI} />
			</div>

			<!-- Bento 2: AI Planners Card (Satellite Widget, spans 1 column) -->
			<section class="rounded-2xl border border-slate-850 bg-slate-950/45 p-6 shadow-xl backdrop-blur-md">
				<div class="flex items-center gap-3 mb-4 pb-3 border-b border-slate-900">
					<div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400">
						<Brain class="h-5 w-5" />
					</div>
					<h2 class="text-base font-extrabold text-slate-100 tracking-tight">AI Health Planners</h2>
				</div>
				<p class="text-xs text-slate-400 mb-5 leading-relaxed">
					Generate fully personalized programs tailored to your specific health report analyses, BMI tracker log, and fitness milestones.
				</p>
				<div class="grid grid-cols-1 gap-3">
					<a
						href="/meal-plan"
						class="flex items-center justify-between p-4 rounded-xl border border-slate-850 bg-slate-900/25 hover:border-emerald-500/30 hover:bg-slate-900/40 hover:-translate-y-0.5 transition-all duration-300 group cursor-pointer"
					>
						<div class="flex items-center gap-3">
							<div class="p-2 rounded-lg bg-emerald-500/10 text-emerald-400 group-hover:bg-emerald-500/25 transition-all duration-300">
								<Utensils class="h-5 w-5" />
							</div>
							<div class="text-left">
								<span class="block text-sm font-semibold text-slate-200">Meal Planner</span>
								<span class="block text-[11px] text-slate-500 font-semibold mt-0.5">AI diet based on reports</span>
							</div>
						</div>
						<span class="text-slate-500 group-hover:text-emerald-400 group-hover:translate-x-0.5 transition-all duration-300">→</span>
					</a>

					<a
						href="/workout"
						class="flex items-center justify-between p-4 rounded-xl border border-slate-850 bg-slate-900/25 hover:border-emerald-500/30 hover:bg-slate-900/40 hover:-translate-y-0.5 transition-all duration-300 group cursor-pointer"
					>
						<div class="flex items-center gap-3">
							<div class="p-2 rounded-lg bg-emerald-500/10 text-emerald-400 group-hover:bg-emerald-500/25 transition-all duration-300">
								<Dumbbell class="h-5 w-5" />
							</div>
							<div class="text-left">
								<span class="block text-sm font-semibold text-slate-200">Workout Planner</span>
								<span class="block text-[11px] text-slate-500 font-semibold mt-0.5">4-week training program</span>
							</div>
						</div>
						<span class="text-slate-500 group-hover:text-emerald-400 group-hover:translate-x-0.5 transition-all duration-300">→</span>
					</a>
				</div>
			</section>

			<!-- Bento 3: Documents Section (spans 2 columns) -->
			<section class="md:col-span-2 space-y-4 bg-slate-950/20 border border-slate-850 p-6 rounded-2xl shadow-xl backdrop-blur-sm">
				<div class="flex items-center justify-between pb-2 border-b border-slate-900/60">
					<h2 class="text-lg font-bold text-slate-200 flex items-center gap-2.5">
						<FilePlus2 class="h-5 w-5 text-emerald-400" />
						<span>Analyzed Documents</span>
					</h2>
					{#if documents.length > 0}
						<span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-slate-900 border border-slate-800 text-slate-400 uppercase tracking-wider">{documents.length} files</span>
					{/if}
				</div>

				{#if documents.length === 0}
					<div class="rounded-xl border border-dashed border-slate-800 bg-slate-950/10 p-10 text-center flex flex-col items-center">
						<div class="p-4 rounded-full bg-slate-900/40 text-slate-500 mb-3 border border-slate-850">
							<FilePlus2 class="h-8 w-8" />
						</div>
						<h3 class="font-bold text-slate-200 text-base">No health reports uploaded</h3>
						<p class="text-xs text-slate-400 mt-1 max-w-sm mx-auto leading-relaxed font-semibold">
							Upload blood work PDFs, prescription files, diet plans or workouts to ask AI questions.
						</p>
						<a
							href="/upload"
							class="mt-5 inline-flex items-center gap-2 px-4 py-2 bg-emerald-950/30 hover:bg-emerald-900/30 border border-emerald-500/20 text-emerald-400 rounded-xl text-xs font-bold transition-all hover:scale-[1.02]"
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

			<!-- Bento 4: Health Milestones Widget (Satellite Widget, spans 1 column) -->
			<section class="rounded-2xl border border-slate-850 bg-slate-950/45 p-6 shadow-xl backdrop-blur-md">
				<div class="flex items-center gap-3 mb-5 pb-3 border-b border-slate-900">
					<div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400">
						<Target class="h-5 w-5" />
					</div>
					<h2 class="text-base font-extrabold text-slate-100 tracking-tight">Health Milestones</h2>
				</div>

				<!-- Create Goal Form -->
				<form onsubmit={handleAddGoal} class="space-y-3 mb-6 bg-slate-900/20 p-4 rounded-xl border border-slate-850">
					<span class="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1">Set New Milestone</span>
					<div>
						<input
							type="text"
							bind:value={newGoalTitle}
							placeholder="e.g. 10k steps daily / Lose 5kg"
							required
							class="w-full px-3 py-2 text-xs rounded-lg border border-slate-800 bg-slate-950 text-slate-100 focus:outline-none focus:border-emerald-500 placeholder-slate-650 transition-colors"
						/>
					</div>
					<div class="relative">
						<input
							type="date"
							bind:value={newGoalDate}
							class="w-full px-3 py-2 text-xs rounded-lg border border-slate-800 bg-slate-950 text-slate-100 focus:outline-none focus:border-emerald-500 transition-colors"
						/>
					</div>
					<button
						type="submit"
						disabled={isAddingGoal || !newGoalTitle.trim()}
						class="w-full flex items-center justify-center gap-1.5 py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg text-xs font-bold transition-all cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed hover:-translate-y-0.5 active:translate-y-0"
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
							<p class="mt-1 font-semibold">Log a milestone above to start tracking achievements.</p>
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
							<div class="space-y-2 pt-3 border-t border-slate-900">
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
	{/if}
</div>
