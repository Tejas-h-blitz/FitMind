<script lang="ts">
	import { Brain, ShieldAlert, Heart, Activity, ArrowRight, MessageSquareCode, LogOut } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { supabase } from '$lib/supabase';
	import { fade, fly } from 'svelte/transition';

	let isLoggedIn = $state(false);
	let isMounted = $state(false);

	onMount(() => {
		isMounted = true;
		supabase.auth.getSession().then(({ data }) => {
			isLoggedIn = !!data.session;
		});
	});

	async function handleLogout() {
		await supabase.auth.signOut();
		isLoggedIn = false;
		window.location.reload();
	}
</script>

<svelte:head>
	<title>FitMind AI - Personal Health Advisor</title>
</svelte:head>

<div class="relative bg-slate-950 min-h-screen text-slate-100 overflow-hidden flex flex-col justify-between">
	<!-- Aurora gradient mesh background -->
	<div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-emerald-950/20 via-slate-955 to-slate-955 -z-10"></div>
	<div class="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] rounded-full bg-emerald-500/5 blur-[120px] pointer-events-none animate-float-slow"></div>
	<div class="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] rounded-full bg-teal-500/5 blur-[120px] pointer-events-none animate-float-slower"></div>
	<div class="absolute inset-0 bg-grid-pattern [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_100%)] opacity-35 -z-10"></div>

	<!-- Header Brand -->
	<header class="fixed top-0 left-0 right-0 z-50 border-b border-slate-850 bg-slate-955/75 backdrop-blur-lg h-20 flex items-center px-6">
		<div class="max-w-7xl mx-auto w-full flex justify-between items-center">
			<a href="/" class="group flex items-center gap-2.5 font-black text-xl text-emerald-400 select-none">
				<div class="p-1.5 rounded-lg bg-emerald-500/10 text-emerald-400 transition-all duration-300">
					<Brain class="h-5.5 w-5.5" />
				</div>
				<span class="inline-block py-0.5 tracking-tight bg-gradient-to-r from-emerald-400 to-teal-350 bg-clip-text text-transparent font-black">FitMind</span>
			</a>
			<div class="flex items-center gap-3">
				{#if isLoggedIn}
					<a
						href="/dashboard"
						class="inline-flex items-center justify-center py-2 px-4 text-xs font-bold rounded-xl border border-transparent bg-slate-950/40 hover:bg-slate-900/30 hover:border-slate-700 transition-all duration-300 active:scale-[0.97] cursor-pointer"
					>
						Dashboard
					</a>
					<button
						onclick={handleLogout}
						class="inline-flex items-center justify-center gap-2 py-2 px-4 text-xs font-bold rounded-xl border border-slate-800 text-slate-350 hover:text-rose-400 hover:border-rose-500/30 hover:bg-rose-955/15 transition-all active:scale-[0.97] cursor-pointer"
					>
						<LogOut class="h-3.5 w-3.5" />
						<span>Logout</span>
					</button>
				{:else}
					<a
						href="/auth/login"
						class="inline-flex items-center justify-center py-2 px-4 text-xs font-bold rounded-xl border border-transparent bg-slate-950/40 hover:bg-slate-900/30 hover:border-slate-700 transition-all duration-300 active:scale-[0.97] cursor-pointer"
					>
						Log In
					</a>
				{/if}
			</div>
		</div>
	</header>

	<!-- Hero Section -->
	<main class="max-w-5xl mx-auto px-6 pt-32 pb-20 flex-1 flex flex-col justify-center items-center text-center">
		{#if isMounted}

			<!-- Headline -->
			<h1 
				in:fly={{ y: 15, duration: 800, delay: 200 }}
				class="text-4xl sm:text-6xl font-black tracking-tight bg-gradient-to-r from-white via-slate-100 to-slate-400 bg-clip-text text-transparent leading-tight sm:leading-none"
			>
				Understand Your Health <br class="hidden sm:inline">
				Like Never Before.
			</h1>

			<!-- Subtitle -->
			<p 
				in:fly={{ y: 15, duration: 800, delay: 300 }}
				class="mt-6 text-sm sm:text-base text-slate-400 max-w-2xl leading-relaxed font-semibold"
			>
				Upload blood reports, workout logs, medical summaries, and diet plans. 
				Get instant, personalized, and structured insights backed by advanced AI document retrieval.
			</p>

			<!-- Action CTA -->
			<div 
				in:fly={{ y: 15, duration: 800, delay: 400 }}
				class="mt-10 flex flex-col sm:flex-row gap-4 justify-center items-center"
			>
				<a
					href={isLoggedIn ? "/dashboard" : "/auth/signup"}
					class="w-full sm:w-auto inline-flex items-center justify-center gap-2 py-3 px-6 text-xs font-extrabold text-white bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-550 hover:to-teal-550 rounded-xl shadow-lg shadow-emerald-950/40 hover:shadow-emerald-950/60 hover:-translate-y-0.5 active:scale-[0.97] transition-all cursor-pointer border border-emerald-500/10"
				>
					<span>Get Started Free</span>
					<ArrowRight class="h-4 w-4" />
				</a>
				<a
					href="#features"
					class="w-full sm:w-auto inline-flex items-center justify-center py-3 px-6 text-xs font-bold text-slate-350 hover:text-white hover:bg-slate-900/40 rounded-xl border border-transparent hover:border-slate-700 transition-all duration-300"
				>
					Learn More
				</a>
			</div>

			<!-- Features Grid -->
			<section 
				in:fly={{ y: 20, duration: 1000, delay: 500 }}
				id="features" 
				class="grid grid-cols-1 md:grid-cols-3 gap-6 w-full mt-24 text-left"
			>
				<!-- Feature 1 -->
				<div class="group rounded-2xl border border-slate-850 bg-slate-950/20 p-6 backdrop-blur-md transition-all duration-300 hover:border-emerald-500/25 hover:bg-slate-900/15 hover:-translate-y-1 shadow-md hover:shadow-emerald-950/5">
					<div class="p-3 bg-emerald-500/10 rounded-xl text-emerald-400 w-fit mb-5 border border-emerald-500/10 group-hover:bg-emerald-500/20 transition-colors">
						<MessageSquareCode class="h-5 w-5" />
					</div>
					<h3 class="font-extrabold text-slate-100 text-sm sm:text-base tracking-tight">RAG-Based Health Chat</h3>
					<p class="text-xs text-slate-450 mt-2 leading-relaxed font-semibold">
						Chat directly with your documents. Ask questions and get precise answers containing line references and source citations.
					</p>
				</div>

				<!-- Feature 2 -->
				<div class="group rounded-2xl border border-slate-850 bg-slate-950/20 p-6 backdrop-blur-md transition-all duration-300 hover:border-emerald-500/25 hover:bg-slate-900/15 hover:-translate-y-1 shadow-md hover:shadow-emerald-950/5">
					<div class="p-3 bg-emerald-500/10 rounded-xl text-emerald-400 w-fit mb-5 border border-emerald-500/10 group-hover:bg-emerald-500/20 transition-colors">
						<Activity class="h-5 w-5" />
					</div>
					<h3 class="font-extrabold text-slate-100 text-sm sm:text-base tracking-tight">BMI Tracker & History</h3>
					<p class="text-xs text-slate-455 mt-2 leading-relaxed font-semibold">
						Log height and weight, calculate BMI instantly, and track progress trends using interactive historical graphics.
					</p>
				</div>

				<!-- Feature 3 -->
				<div class="group rounded-2xl border border-slate-850 bg-slate-950/20 p-6 backdrop-blur-md transition-all duration-300 hover:border-emerald-500/25 hover:bg-slate-900/15 hover:-translate-y-1 shadow-md hover:shadow-emerald-950/5">
					<div class="p-3 bg-emerald-500/10 rounded-xl text-emerald-400 w-fit mb-5 border border-emerald-500/10 group-hover:bg-emerald-500/20 transition-colors">
						<ShieldAlert class="h-5 w-5" />
					</div>
					<h3 class="font-extrabold text-slate-100 text-sm sm:text-base tracking-tight">Privacy & RLS Security</h3>
					<p class="text-xs text-slate-455 mt-2 leading-relaxed font-semibold">
						Your data belongs to you. Fully secure storage under Supabase Auth and isolated collections with Qdrant.
					</p>
				</div>
			</section>
		{/if}
	</main>

	<!-- Footer -->
	<footer class="border-t border-slate-900 py-6 text-center text-[10px] font-bold tracking-tight text-slate-500">
		<p>&copy; {new Date().getFullYear()} FitMind AI. All rights reserved. For informational purposes only.</p>
	</footer>
</div>

<style>
	@keyframes floatSlow {
		0%, 100% { transform: translateY(0) scale(1); }
		50% { transform: translateY(-10px) scale(1.05); }
	}
	:global(.animate-float-slow) {
		animation: floatSlow 6s ease-in-out infinite;
	}
	:global(.animate-float-slower) {
		animation: floatSlow 8s ease-in-out infinite;
	}
</style>
