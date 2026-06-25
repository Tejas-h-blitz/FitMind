<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Document, type MealPlan } from '$lib/api';
	import MealPlanCard from '$lib/components/MealPlanCard.svelte';
	import { Brain, Sparkles, AlertCircle, RefreshCw, Printer, Trash2, Check, ArrowRight } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let documents = $state<Document[]>([]);
	let analyzedDocs = $derived(documents.filter((d) => d.status === 'analyzed'));
	
	let selectedDocId = $state('');
	let dietaryPreference = $state<'vegetarian' | 'non-vegetarian' | 'vegan'>('non-vegetarian');
	
	let mealPlan = $state<MealPlan | null>(null);
	let isLoadingDocs = $state(true);
	let isGenerating = $state(false);
	
	// Loading message animation state
	let loadingMsgIndex = $state(0);
	const loadingMessages = [
		'Analyzing your health report...',
		'Checking your goals...',
		'Crafting your personalized plan...'
	];
	let loadingTimer: any = null;

	onMount(async () => {
		try {
			const [docsRes, planRes] = await Promise.all([
				api.listDocuments(),
				api.getLatestMealPlan()
			]);

			if (docsRes.success && docsRes.data) {
				documents = docsRes.data;
				if (analyzedDocs.length > 0) {
					selectedDocId = analyzedDocs[0].id;
				}
			}

			if (planRes.success && planRes.data) {
				mealPlan = planRes.data;
			}
		} catch (err) {
			console.error(err);
			toast.error('Failed to load initial data');
		} finally {
			isLoadingDocs = false;
		}
	});

	function startGeneratingMessages() {
		loadingMsgIndex = 0;
		loadingTimer = setInterval(() => {
			loadingMsgIndex = (loadingMsgIndex + 1) % loadingMessages.length;
		}, 3000);
	}

	function stopGeneratingMessages() {
		if (loadingTimer) {
			clearInterval(loadingTimer);
			loadingTimer = null;
		}
	}

	async function handleGenerate() {
		if (!selectedDocId) {
			toast.error('Please select an analyzed health report.');
			return;
		}

		isGenerating = true;
		startGeneratingMessages();

		try {
			const res = await api.generateMealPlan(selectedDocId, dietaryPreference);
			if (res.success && res.data) {
				mealPlan = res.data;
				toast.success('Your personalized meal plan has been generated!');
			} else {
				toast.error(res.error || 'Failed to generate meal plan');
			}
		} catch (err: any) {
			toast.error(err.message || 'An error occurred during generation.');
		} finally {
			stopGeneratingMessages();
			isGenerating = false;
		}
	}

	async function handleDeletePlan() {
		if (!mealPlan) return;
		if (confirm('Are you sure you want to delete this meal plan?')) {
			const res = await api.deleteMealPlan(mealPlan.id);
			if (res.success) {
				mealPlan = null;
				toast.success('Meal plan deleted');
			} else {
				toast.error(res.error || 'Failed to delete meal plan');
			}
		}
	}

	function handlePrint() {
		window.print();
	}
</script>

<svelte:head>
	<title>AI Personalized Meal Planner - FitMind</title>
</svelte:head>

<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 printable-area">
	<!-- Header -->
	<div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-8 non-printable">
		<div>
			<h1 class="text-2xl sm:text-3xl font-extrabold text-slate-100 flex items-center gap-2">
				<Sparkles class="text-emerald-400 h-7 w-7" />
				<span>AI Personalized Meal Planner</span>
			</h1>
			<p class="text-slate-400 text-sm mt-1">Generate a comprehensive 7-day meal plan based on health metrics, BMI, and goals.</p>
		</div>

		{#if mealPlan}
			<div class="flex gap-2 shrink-0">
				<button
					onclick={handlePrint}
					class="inline-flex items-center justify-center gap-1.5 py-2 px-3 text-xs font-bold text-slate-300 bg-slate-900 hover:bg-slate-800 border border-slate-800 hover:border-slate-700 rounded-lg cursor-pointer transition-all"
				>
					<Printer class="h-3.5 w-3.5" />
					<span>Download PDF</span>
				</button>
				<button
					onclick={handleDeletePlan}
					class="inline-flex items-center justify-center p-2 text-slate-400 hover:text-rose-400 bg-slate-900/40 hover:bg-rose-950/20 border border-slate-800 hover:border-rose-500/20 rounded-lg cursor-pointer transition-all"
					aria-label="Delete plan"
				>
					<Trash2 class="h-3.5 w-3.5" />
				</button>
			</div>
		{/if}
	</div>

	{#if isLoadingDocs}
		<div class="h-96 flex flex-col justify-center items-center non-printable">
			<div class="h-10 w-10 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
			<span class="text-sm font-semibold text-slate-400 mt-4">Loading your records...</span>
		</div>
	{:else if isGenerating}
		<!-- Generating loading screen -->
		<div class="h-96 flex flex-col justify-center items-center text-center max-w-sm mx-auto non-printable">
			<div class="relative flex items-center justify-center mb-6">
				<div class="h-16 w-16 border-4 border-emerald-500/20 border-t-emerald-400 rounded-full animate-spin"></div>
				<Brain class="absolute h-6 w-6 text-emerald-400 animate-pulse" />
			</div>
			<h3 class="text-lg font-bold text-slate-200 animate-pulse">{loadingMessages[loadingMsgIndex]}</h3>
			<p class="text-xs text-slate-500 mt-2 leading-relaxed">Please don't close this tab. Our nutritionist engine is analyzing your blood work deficiencies and structuring recipes.</p>
		</div>
	{:else}
		<!-- Main Generator Form / Results -->
		{#if !mealPlan}
			<!-- Generator configuration layout -->
			<div class="space-y-8 bg-slate-900/10 border border-slate-800/60 p-6 sm:p-8 rounded-2xl backdrop-blur-sm non-printable">
				<!-- Step 1: Document select -->
				<div class="space-y-3">
					<label for="doc-select" class="block text-sm font-bold text-slate-350">
						Step 1: Select Health Report
					</label>
					<p class="text-xs text-slate-500">We will scan this report to incorporate missing vitamins and minerals into your recipes.</p>
					
					{#if analyzedDocs.length === 0}
						<div class="p-4 rounded-xl border border-rose-500/10 bg-rose-950/5 text-slate-400 text-xs flex items-start gap-2.5">
							<AlertCircle class="h-4 w-4 text-rose-500 shrink-0 mt-0.5" />
							<div>
								<span class="font-bold text-slate-200 block mb-0.5">No Analyzed Reports Available</span>
								<span>You need to upload and ingest a health PDF report first. Svelte matches will show up here once analyzed.</span>
								<a href="/upload" class="block mt-2 font-bold text-emerald-400 hover:underline">Go upload a PDF →</a>
							</div>
						</div>
					{:else}
						<select
							id="doc-select"
							bind:value={selectedDocId}
							class="w-full px-3.5 py-2.5 rounded-lg border border-slate-800 bg-slate-950 text-slate-100 text-sm focus:outline-none focus:border-emerald-500"
						>
							{#each analyzedDocs as doc}
								<option value={doc.id}>{doc.name} (Analyzed: {new Date(doc.created_at).toLocaleDateString()})</option>
							{/each}
						</select>
					{/if}
				</div>

				<!-- Step 2: Dietary preference -->
				<div class="space-y-3">
					<span class="block text-sm font-bold text-slate-350">
						Step 2: Select Dietary Preference
					</span>
					<p class="text-xs text-slate-500">We will customize the recipes based on your dietary guidelines.</p>
					
					<div class="flex gap-3 flex-wrap">
						{#each ['non-vegetarian', 'vegetarian', 'vegan'] as pref}
							<button
								onclick={() => dietaryPreference = pref as any}
								class="px-4 py-2.5 rounded-xl border text-xs font-bold capitalize transition-all cursor-pointer flex items-center gap-2 {dietaryPreference === pref ? 'bg-emerald-600 border-emerald-500 text-white shadow-md' : 'bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200 hover:border-slate-700'}"
							>
								{#if dietaryPreference === pref}
									<Check class="h-3.5 w-3.5" />
								{/if}
								<span>{pref.replace('-', ' ')}</span>
							</button>
						{/each}
					</div>
				</div>

				<!-- Submit button -->
				<div class="pt-4 border-t border-slate-900 flex justify-end">
					<button
						onclick={handleGenerate}
						disabled={analyzedDocs.length === 0}
						class="w-full sm:w-auto inline-flex items-center justify-center gap-2 py-3 px-6 text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg shadow-md transition-all cursor-pointer"
					>
						<span>Generate My Meal Plan</span>
						<ArrowRight class="h-4 w-4" />
					</button>
				</div>
			</div>
		{:else}
			<!-- Generated Meal Plan Layout -->
			<div class="space-y-8">
				<!-- Reasoning box at top -->
				<div class="p-5 rounded-xl border border-slate-800 bg-slate-900/35 backdrop-blur-sm">
					<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider">AI Nutritionist Reasoning</h3>
					<p class="text-sm text-slate-200 mt-2 leading-relaxed italic">
						"{mealPlan.reasoning}"
					</p>
				</div>

				<!-- Stats Cards: Calories + Protein targets -->
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
					<div class="p-5 rounded-xl border border-slate-800 bg-slate-950/30 flex justify-between items-center">
						<div>
							<span class="block text-xs font-bold text-slate-500 uppercase tracking-wider">Daily Calories</span>
							<span class="block text-3xl font-extrabold text-slate-100 mt-1">{mealPlan.daily_calories_target} <span class="text-sm font-normal text-slate-400">kcal</span></span>
						</div>
						<div class="h-10 w-10 rounded-full bg-emerald-500/10 text-emerald-400 flex items-center justify-center text-lg">🔥</div>
					</div>

					<div class="p-5 rounded-xl border border-slate-800 bg-slate-950/30 flex justify-between items-center">
						<div>
							<span class="block text-xs font-bold text-slate-500 uppercase tracking-wider">Daily Protein Target</span>
							<span class="block text-3xl font-extrabold text-slate-100 mt-1">{mealPlan.protein_target_g} <span class="text-sm font-normal text-slate-400">g</span></span>
						</div>
						<div class="h-10 w-10 rounded-full bg-sky-500/10 text-sky-400 flex items-center justify-center text-lg">🥩</div>
					</div>
				</div>

				<!-- 7-day accordion -->
				<div class="space-y-4">
					<div class="flex items-center justify-between">
						<h2 class="text-base font-bold text-slate-200">7-Day Program Guide</h2>
						<span class="text-xs text-slate-500">Click a day to view meal recipes and benefits</span>
					</div>

					<div class="space-y-3">
						{#each mealPlan.days as dayPlan}
							<MealPlanCard {dayPlan} />
						{/each}
					</div>
				</div>

				<!-- Weekly notes at bottom -->
				<div class="p-5 rounded-xl border border-emerald-500/10 bg-emerald-950/5">
					<h4 class="text-xs font-bold uppercase tracking-wider text-emerald-400">Overall Weekly Notes & Guidance</h4>
					<p class="text-sm text-slate-350 mt-2.5 leading-relaxed">
						{mealPlan.weekly_notes}
					</p>
				</div>

				<!-- Disclaimer -->
				<div class="text-[10px] text-slate-500 italic text-center max-w-lg mx-auto">
					⚠️ Medical Disclaimer: This meal plan is AI-generated and not a substitute for professional nutritional advice or clinical care. Please review the ingredients against personal allergies.
				</div>

				<!-- Regenerate Button (non-printable) -->
				<div class="flex justify-center pt-4 border-t border-slate-900 non-printable">
					<button
						onclick={() => mealPlan = null}
						class="inline-flex items-center justify-center gap-2 py-2.5 px-5 text-sm font-bold text-emerald-400 hover:text-emerald-300 bg-emerald-950/20 hover:bg-emerald-950/35 border border-emerald-500/20 hover:border-emerald-500/40 rounded-lg cursor-pointer transition-all"
					>
						<RefreshCw class="h-4 w-4" />
						<span>Regenerate Meal Plan</span>
					</button>
				</div>
			</div>
		{/if}
	{/if}
</div>

<style>
	/* Print specific styling */
	@media print {
		:global(body) {
			background: white !important;
			color: black !important;
		}
		.non-printable {
			display: none !important;
		}
		.printable-area {
			max-width: 100% !important;
			width: 100% !important;
			padding: 0 !important;
			margin: 0 !important;
		}
		:global(nav) {
			display: none !important;
		}
		:global(header) {
			display: none !important;
		}
		:global(button) {
			display: none !important;
		}
	}
</style>
