<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type WorkoutPlan } from '$lib/api';
	import WorkoutDayCard from '$lib/components/WorkoutDayCard.svelte';
	import { Dumbbell, ShieldAlert, Sparkles, ArrowRight, Check, RefreshCw, Trash2, Heart, Award, Flame } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let fitnessLevel = $state<'beginner' | 'intermediate' | 'advanced'>('beginner');
	let equipment = $state<'none' | 'home' | 'full_gym'>('none');
	let daysPerWeek = $state(4);
	
	let workoutPlan = $state<WorkoutPlan | null>(null);
	let isLoading = $state(true);
	let isGenerating = $state(false);

	// Loading messages cycle state
	let loadingMsgIndex = $state(0);
	const loadingMessages = [
		'Analyzing your fitness profile...',
		'Building your program...',
		'Calculating progressive overload...'
	];
	let loadingTimer: any = null;

	onMount(async () => {
		try {
			const res = await api.getLatestWorkoutPlan();
			if (res.success && res.data) {
				workoutPlan = res.data;
			}
		} catch (err) {
			console.error(err);
		} finally {
			isLoading = false;
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
		isGenerating = true;
		startGeneratingMessages();

		try {
			const res = await api.generateWorkoutPlan(fitnessLevel, equipment, daysPerWeek);
			if (res.success && res.data) {
				workoutPlan = res.data;
				toast.success('Your custom workout program is ready!');
			} else {
				toast.error(res.error || 'Failed to generate workout plan');
			}
		} catch (err: any) {
			toast.error(err.message || 'Error occurred during generation.');
		} finally {
			stopGeneratingMessages();
			isGenerating = false;
		}
	}

	async function handleDeletePlan() {
		if (!workoutPlan) return;
		if (confirm('Are you sure you want to delete this workout program?')) {
			const res = await api.deleteWorkoutPlan(workoutPlan.id);
			if (res.success) {
				workoutPlan = null;
				toast.success('Workout program deleted');
			} else {
				toast.error(res.error || 'Failed to delete workout plan');
			}
		}
	}
</script>

<svelte:head>
	<title>AI Workout Program Generator - FitMind</title>
</svelte:head>

<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
	<!-- Header -->
	<div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-8">
		<div>
			<h1 class="text-2xl sm:text-3xl font-extrabold text-slate-100 flex items-center gap-2">
				<Dumbbell class="text-emerald-400 h-7 w-7" />
				<span>AI Workout Plan Generator</span>
			</h1>
			<p class="text-slate-400 text-sm mt-1">Design a structured progressive overload routine incorporating your BMI and fitness goals.</p>
		</div>

		{#if workoutPlan}
			<button
				onclick={handleDeletePlan}
				class="inline-flex items-center justify-center gap-2 py-2 px-3.5 text-xs font-bold text-slate-400 hover:text-rose-450 bg-slate-900/60 hover:bg-rose-950/20 border border-slate-800 hover:border-rose-500/20 rounded-lg cursor-pointer transition-all"
			>
				<Trash2 class="h-3.5 w-3.5" />
				<span>Delete Program</span>
			</button>
		{/if}
	</div>

	{#if isLoading}
		<div class="h-96 flex flex-col justify-center items-center">
			<div class="h-10 w-10 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
			<span class="text-sm font-semibold text-slate-400 mt-4">Retrieving your workout status...</span>
		</div>
	{:else if isGenerating}
		<!-- Generating loading screen -->
		<div class="h-96 flex flex-col justify-center items-center text-center max-w-sm mx-auto">
			<div class="relative flex items-center justify-center mb-6">
				<div class="h-16 w-16 border-4 border-emerald-500/20 border-t-emerald-400 rounded-full animate-spin"></div>
				<Dumbbell class="absolute h-6 w-6 text-emerald-400 animate-bounce" />
			</div>
			<h3 class="text-lg font-bold text-slate-200 animate-pulse">{loadingMessages[loadingMsgIndex]}</h3>
			<p class="text-xs text-slate-500 mt-2 leading-relaxed">Calculating rep targets, active recovery breaks, and form adjustments based on your metrics.</p>
		</div>
	{:else}
		<!-- Main view -->
		{#if !workoutPlan}
			<div class="space-y-8 bg-slate-900/10 border border-slate-800/60 p-6 sm:p-8 rounded-2xl backdrop-blur-sm">
				<!-- Step 1: Fitness Level -->
				<div class="space-y-3">
					<span class="block text-sm font-bold text-slate-355">
						Step 1: Select Your Experience Level
					</span>
					<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
						<!-- Card 1: Beginner -->
						<button
							onclick={() => fitnessLevel = 'beginner'}
							class="p-5 rounded-xl border text-left cursor-pointer transition-all duration-300 flex flex-col justify-between h-[130px] {fitnessLevel === 'beginner' ? 'bg-emerald-950/20 border-emerald-500 shadow-md shadow-emerald-950/30' : 'bg-slate-950/40 border-slate-850 hover:border-slate-700'}"
						>
							<div class="flex justify-between items-center w-full">
								<span class="p-2 rounded-lg bg-emerald-500/10 text-emerald-400">
									<Heart class="h-5 w-5" />
								</span>
								{#if fitnessLevel === 'beginner'}
									<Check class="h-4 w-4 text-emerald-400" />
								{/if}
							</div>
							<div>
								<span class="block text-sm font-bold text-slate-200">Beginner</span>
								<span class="block text-[11px] text-slate-500 mt-0.5">Just starting out or resuming fitness</span>
							</div>
						</button>

						<!-- Card 2: Intermediate -->
						<button
							onclick={() => fitnessLevel = 'intermediate'}
							class="p-5 rounded-xl border text-left cursor-pointer transition-all duration-300 flex flex-col justify-between h-[130px] {fitnessLevel === 'intermediate' ? 'bg-emerald-950/20 border-emerald-500 shadow-md shadow-emerald-950/30' : 'bg-slate-950/40 border-slate-850 hover:border-slate-700'}"
						>
							<div class="flex justify-between items-center w-full">
								<span class="p-2 rounded-lg bg-emerald-500/10 text-emerald-400">
									<Flame class="h-5 w-5" />
								</span>
								{#if fitnessLevel === 'intermediate'}
									<Check class="h-4 w-4 text-emerald-400" />
								{/if}
							</div>
							<div>
								<span class="block text-sm font-bold text-slate-200">Intermediate</span>
								<span class="block text-[11px] text-slate-500 mt-0.5">6+ months of regular training experience</span>
							</div>
						</button>

						<!-- Card 3: Advanced -->
						<button
							onclick={() => fitnessLevel = 'advanced'}
							class="p-5 rounded-xl border text-left cursor-pointer transition-all duration-300 flex flex-col justify-between h-[130px] {fitnessLevel === 'advanced' ? 'bg-emerald-950/20 border-emerald-500 shadow-md shadow-emerald-950/30' : 'bg-slate-950/40 border-slate-850 hover:border-slate-700'}"
						>
							<div class="flex justify-between items-center w-full">
								<span class="p-2 rounded-lg bg-emerald-500/10 text-emerald-400">
									<Award class="h-5 w-5" />
								</span>
								{#if fitnessLevel === 'advanced'}
									<Check class="h-4 w-4 text-emerald-400" />
								{/if}
							</div>
							<div>
								<span class="block text-sm font-bold text-slate-200">Advanced</span>
								<span class="block text-[11px] text-slate-500 mt-0.5">2+ years of consistent heavy training</span>
							</div>
						</button>
					</div>
				</div>

				<!-- Step 2: Equipment select -->
				<div class="space-y-3">
					<span class="block text-sm font-bold text-slate-355">
						Step 2: Available Training Equipment
					</span>
					<div class="flex gap-3 flex-wrap">
						{#each [
							{ id: 'none', label: 'No Equipment (Bodyweight)' },
							{ id: 'home', label: 'Home Setup (Dumbbells/Bands)' },
							{ id: 'full_gym', label: 'Full Gym Access' }
						] as eq}
							<button
								onclick={() => equipment = eq.id as any}
								class="px-4 py-2.5 rounded-xl border text-xs font-bold transition-all cursor-pointer flex items-center gap-2 {equipment === eq.id ? 'bg-emerald-600 border-emerald-500 text-white shadow-md' : 'bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200 hover:border-slate-700'}"
							>
								{#if equipment === eq.id}
									<Check class="h-3.5 w-3.5" />
								{/if}
								<span>{eq.label}</span>
							</button>
						{/each}
					</div>
				</div>

				<!-- Step 3: Days Slider -->
				<div class="space-y-3">
					<div class="flex justify-between items-center">
						<span class="block text-sm font-bold text-slate-355">
							Step 3: Training Frequency
						</span>
						<span class="text-xs font-bold text-emerald-400 bg-emerald-950/40 border border-emerald-500/20 px-2 py-0.5 rounded-md">
							{daysPerWeek} days / week
						</span>
					</div>
					<p class="text-xs text-slate-500">How many days are you available to work out each week?</p>
					
					<div class="flex items-center gap-4 pt-2">
						<span class="text-xs text-slate-500">3 Days</span>
						<input
							type="range"
							min="3"
							max="6"
							bind:value={daysPerWeek}
							class="w-full h-2 bg-slate-950 rounded-lg appearance-none cursor-pointer accent-emerald-500 focus:outline-none"
						/>
						<span class="text-xs text-slate-500">6 Days</span>
					</div>
				</div>

				<!-- Submit -->
				<div class="pt-4 border-t border-slate-900 flex justify-end">
					<button
						onclick={handleGenerate}
						class="w-full sm:w-auto inline-flex items-center justify-center gap-2 py-3 px-6 text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-500 rounded-lg shadow-md transition-all cursor-pointer"
					>
						<span>Generate My Workout Plan</span>
						<ArrowRight class="h-4 w-4" />
					</button>
				</div>
			</div>
		{:else}
			<!-- Program view -->
			<div class="space-y-8">
				<!-- Hero section -->
				<div class="relative rounded-2xl border border-emerald-500/10 bg-emerald-950/5 p-6 flex items-center justify-between flex-wrap gap-4 overflow-hidden">
					<div class="absolute -right-16 -top-16 h-36 w-36 rounded-full bg-emerald-500/5 blur-2xl"></div>
					<div class="relative space-y-1">
						<span class="text-[10px] font-bold text-emerald-400 uppercase tracking-wider">Active Training Program</span>
						<h2 class="text-xl sm:text-2xl font-extrabold text-slate-100 capitalize">{workoutPlan.program_name}</h2>
						<p class="text-xs text-slate-400">Duration: <span class="text-slate-200 font-bold">{workoutPlan.duration_weeks} weeks</span> • Level: <span class="text-slate-200 font-bold capitalize">{workoutPlan.fitness_level}</span></p>
					</div>

					<button
						onclick={() => workoutPlan = null}
						class="relative inline-flex items-center justify-center gap-1.5 py-2 px-4 text-xs font-bold text-emerald-400 hover:text-emerald-300 bg-emerald-950/30 hover:bg-emerald-950/50 border border-emerald-500/20 rounded-lg cursor-pointer transition-all"
					>
						<RefreshCw class="h-3.5 w-3.5" />
						<span>New Configuration</span>
					</button>
				</div>

				<!-- Reasoning box -->
				<div class="p-5 rounded-xl border border-slate-800 bg-slate-900/35 leading-relaxed">
					<h3 class="text-xs font-bold text-slate-500 uppercase tracking-wider">Coach's Recommendation Notes</h3>
					<p class="text-sm text-slate-250 mt-2 italic leading-relaxed">
						"{workoutPlan.reasoning}"
					</p>
				</div>

				<!-- Schedule -->
				<div class="space-y-4">
					<div class="flex items-center justify-between">
						<h3 class="text-base font-bold text-slate-200">Weekly Training Schedule</h3>
						<span class="text-xs text-slate-500">Tap cards to inspect exercise details & form tips</span>
					</div>

					<!-- Scroll container on mobile, Grid on desktop -->
					<div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
						{#each workoutPlan.weekly_schedule as dayWorkout}
							<WorkoutDayCard {dayWorkout} />
						{/each}
					</div>
				</div>

				<!-- Notes checklist -->
				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					<!-- Progression -->
					<div class="p-5 rounded-xl border border-slate-850 bg-slate-950/40 space-y-2">
						<h4 class="text-xs font-bold uppercase tracking-wider text-emerald-450 flex items-center gap-1.5">
							<span>📈 Progression Notes</span>
						</h4>
						<p class="text-xs text-slate-350 leading-relaxed font-medium">
							{workoutPlan.progression_notes}
						</p>
					</div>

					<!-- Safety -->
					<div class="p-5 rounded-xl border border-slate-850 bg-slate-950/40 space-y-2">
						<h4 class="text-xs font-bold uppercase tracking-wider text-rose-400 flex items-center gap-1.5">
							<ShieldAlert class="h-4 w-4" />
							<span>⚠️ Safety Guidelines</span>
						</h4>
						<p class="text-xs text-slate-350 leading-relaxed font-medium">
							{workoutPlan.safety_notes}
						</p>
					</div>
				</div>

				<!-- Disclaimer -->
				<div class="text-[10px] text-slate-550 italic text-center max-w-lg mx-auto">
					🏋️ Medical Disclaimer: Please consult a physician before starting any exercise program, especially if you have chronic health conditions. Stop working out immediately if you experience dizziness or shortness of breath.
				</div>
			</div>
		{/if}
	{/if}
</div>
