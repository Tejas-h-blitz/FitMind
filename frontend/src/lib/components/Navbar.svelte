<script lang="ts">
	import { supabase } from '$lib/supabase';
	import { page } from '$app/state';
	import { Flame, Brain, LogOut, LayoutDashboard, Upload, Utensils, Dumbbell } from 'lucide-svelte';

	let { streakCount = 0 } = $props<{ streakCount?: number }>();

	async function handleLogout() {
		await supabase.auth.signOut();
		// Auth Guard will automatically redirect
		window.location.href = '/auth/login';
	}
</script>

<nav class="fixed top-0 left-0 right-0 z-50 border-b border-slate-850 bg-slate-950/75 backdrop-blur-lg transition-all duration-300">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="flex h-16 justify-between items-center">
			<div class="flex items-center gap-2">
				<a href="/dashboard" class="group flex items-center gap-2.5 font-black text-xl text-emerald-400 hover:text-emerald-300 transition-all duration-300">
					<div class="p-1.5 rounded-lg bg-emerald-500/10 text-emerald-400 group-hover:bg-emerald-500/20 group-hover:scale-110 transition-all duration-300">
						<Brain class="h-5.5 w-5.5" />
					</div>
					<span class="tracking-tight bg-gradient-to-r from-emerald-400 to-teal-300 bg-clip-text text-transparent">FitMind</span>
				</a>
			</div>

			<div class="flex items-center gap-6">
				<!-- Navigation Links -->
				<div class="hidden sm:flex items-center gap-1.5">
					<a
						href="/dashboard"
						class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-semibold transition-all duration-350 {page.url.pathname === '/dashboard' ? 'text-emerald-400 bg-emerald-950/25 border border-emerald-500/20' : 'text-slate-400 hover:text-slate-100 hover:bg-slate-900/60 border border-transparent'}"
					>
						<LayoutDashboard class="h-4 w-4" />
						<span>Dashboard</span>
					</a>
					<a
						href="/upload"
						class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-semibold transition-all duration-350 {page.url.pathname === '/upload' ? 'text-emerald-400 bg-emerald-950/25 border border-emerald-500/20' : 'text-slate-400 hover:text-slate-100 hover:bg-slate-900/60 border border-transparent'}"
					>
						<Upload class="h-4 w-4" />
						<span>Upload Files</span>
					</a>
					<a
						href="/meal-plan"
						class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-semibold transition-all duration-350 {page.url.pathname === '/meal-plan' ? 'text-emerald-400 bg-emerald-950/25 border border-emerald-500/20' : 'text-slate-400 hover:text-slate-100 hover:bg-slate-900/60 border border-transparent'}"
					>
						<Utensils class="h-4 w-4" />
						<span>Meal Plan</span>
					</a>
					<a
						href="/workout"
						class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-semibold transition-all duration-350 {page.url.pathname === '/workout' ? 'text-emerald-400 bg-emerald-950/25 border border-emerald-500/20' : 'text-slate-400 hover:text-slate-100 hover:bg-slate-900/60 border border-transparent'}"
					>
						<Dumbbell class="h-4 w-4" />
						<span>Workout</span>
					</a>
				</div>

				<!-- Streak Counter -->
				{#if streakCount > 0}
					<div class="group flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-gradient-to-r from-orange-950/30 to-amber-950/30 border border-orange-500/20 text-orange-400 text-xs font-bold shadow-sm shadow-orange-950/10 select-none transition-all duration-300 hover:border-orange-500/40">
						<Flame class="h-4 w-4 fill-orange-500 text-orange-500 animate-pulse group-hover:scale-125 transition-transform duration-300" />
						<span>{streakCount} Day Streak</span>
					</div>
				{/if}
			</div>
		</div>
	</div>
</nav>
