<script lang="ts">
	import { Brain, ShieldAlert, Heart, Activity, ArrowRight, MessageSquareCode } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { supabase } from '$lib/supabase';

	let isLoggedIn = $state(false);

	onMount(() => {
		supabase.auth.getSession().then(({ data }) => {
			isLoggedIn = !!data.session;
		});
	});
</script>

<div class="relative bg-slate-950 min-h-screen text-slate-100 overflow-hidden flex flex-col justify-between">
	<!-- Background grid overlay -->
	<div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-emerald-950/20 via-slate-950 to-slate-950 -z-10"></div>
	<div class="absolute inset-0 bg-[linear-gradient(to_right,#0f172a_1px,transparent_1px),linear-gradient(to_bottom,#0f172a_1px,transparent_1px)] bg-[size:4rem_4rem] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_100%)] opacity-35 -z-10"></div>

	<!-- Header Brand -->
	<header class="max-w-7xl mx-auto px-6 w-full h-20 flex justify-between items-center">
		<div class="flex items-center gap-2 font-bold text-xl text-emerald-400">
			<Brain class="h-6 w-6" />
			<span>FitMind</span>
		</div>
		<div>
			<a
				href={isLoggedIn ? "/dashboard" : "/auth/login"}
				class="inline-flex items-center justify-center py-2 px-4 text-sm font-semibold rounded-lg border border-slate-800 hover:border-slate-700 bg-slate-900/50 hover:bg-slate-900 transition-all cursor-pointer"
			>
				{isLoggedIn ? "Dashboard" : "Log In"}
			</a>
		</div>
	</header>

	<!-- Hero Section -->
	<main class="max-w-5xl mx-auto px-6 py-20 flex-1 flex flex-col justify-center items-center text-center">
		<div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/25 text-emerald-400 text-xs font-semibold mb-6 select-none">
			<Heart class="h-3.5 w-3.5 fill-emerald-500 text-emerald-500" />
			<span>AI-Powered Personal Health Intelligence</span>
		</div>

		<h1 class="text-4xl sm:text-6xl font-extrabold tracking-tight bg-gradient-to-r from-white via-slate-200 to-slate-400 bg-clip-text text-transparent leading-tight sm:leading-none">
			Understand Your Health <br class="hidden sm:inline">
			Like Never Before.
		</h1>

		<p class="mt-6 text-base sm:text-lg text-slate-400 max-w-2xl leading-relaxed">
			Upload blood reports, workout logs, medical summaries, and diet plans. 
			Get instant, personalized, and structured insights backed by advanced AI document retrieval.
		</p>

		<!-- Action CTA -->
		<div class="mt-10 flex flex-col sm:flex-row gap-4 justify-center items-center">
			<a
				href={isLoggedIn ? "/dashboard" : "/auth/signup"}
				class="w-full sm:w-auto inline-flex items-center justify-center gap-2 py-3 px-6 text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-500 border border-emerald-500/25 hover:border-emerald-400/25 rounded-lg shadow-lg shadow-emerald-950/50 transition-all cursor-pointer"
			>
				<span>Get Started Free</span>
				<ArrowRight class="h-4 w-4" />
			</a>
			<a
				href="#features"
				class="w-full sm:w-auto inline-flex items-center justify-center py-3 px-6 text-sm font-semibold text-slate-300 hover:text-white hover:bg-slate-900/60 rounded-lg transition-all"
			>
				Learn More
			</a>
		</div>

		<!-- Features Grid -->
		<section id="features" class="grid grid-cols-1 md:grid-cols-3 gap-6 w-full mt-24 text-left">
			<!-- Feature 1 -->
			<div class="rounded-xl border border-slate-900 bg-slate-900/20 p-6 backdrop-blur-sm">
				<div class="p-3 bg-emerald-500/10 rounded-lg text-emerald-400 w-fit mb-4">
					<MessageSquareCode class="h-5 w-5" />
				</div>
				<h3 class="font-bold text-slate-100 text-base">RAG-Based Health Chat</h3>
				<p class="text-xs text-slate-400 mt-2 leading-relaxed">
					Chat directly with your documents. Ask questions and get precise answers containing line references and source citations.
				</p>
			</div>

			<!-- Feature 2 -->
			<div class="rounded-xl border border-slate-900 bg-slate-900/20 p-6 backdrop-blur-sm">
				<div class="p-3 bg-emerald-500/10 rounded-lg text-emerald-400 w-fit mb-4">
					<Activity class="h-5 w-5" />
				</div>
				<h3 class="font-bold text-slate-100 text-base">BMI Tracker & History</h3>
				<p class="text-xs text-slate-400 mt-2 leading-relaxed">
					Log height and weight, calculate BMI instantly, and track progress trends using interactive historical graphics.
				</p>
			</div>

			<!-- Feature 3 -->
			<div class="rounded-xl border border-slate-900 bg-slate-900/20 p-6 backdrop-blur-sm">
				<div class="p-3 bg-emerald-500/10 rounded-lg text-emerald-400 w-fit mb-4">
					<ShieldAlert class="h-5 w-5" />
				</div>
				<h3 class="font-bold text-slate-100 text-base">Privacy & RLS Security</h3>
				<p class="text-xs text-slate-400 mt-2 leading-relaxed">
					Your data belongs to you. Fully secure storage under Supabase Auth and isolated collections with Qdrant.
				</p>
			</div>
		</section>
	</main>

	<!-- Footer -->
	<footer class="border-t border-slate-900 py-6 text-center text-xs text-slate-500">
		<p>&copy; {new Date().getFullYear()} FitMind AI. All rights reserved. For informational purposes only.</p>
	</footer>
</div>
