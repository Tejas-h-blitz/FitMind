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

<nav class="sticky top-0 z-50 border-b border-slate-800 bg-slate-950/80 backdrop-blur-md">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="flex h-16 justify-between items-center">
			<div class="flex items-center gap-2">
				<a href="/dashboard" class="flex items-center gap-2 font-bold text-xl text-emerald-400 hover:text-emerald-300 transition-colors">
					<Brain class="h-6 w-6" />
					<span>FitMind</span>
				</a>
			</div>

			<div class="flex items-center gap-6">
				<!-- Navigation Links -->
				<div class="hidden sm:flex gap-4">
					<a
						href="/dashboard"
						class="flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors {page.url.pathname === '/dashboard' ? 'text-emerald-400 bg-emerald-950/30' : 'text-slate-300 hover:text-white hover:bg-slate-900'}"
					>
						<LayoutDashboard class="h-4 w-4" />
						<span>Dashboard</span>
					</a>
					<a
						href="/upload"
						class="flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors {page.url.pathname === '/upload' ? 'text-emerald-400 bg-emerald-950/30' : 'text-slate-300 hover:text-white hover:bg-slate-900'}"
					>
						<Upload class="h-4 w-4" />
						<span>Upload Files</span>
					</a>
					<a
						href="/meal-plan"
						class="flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors {page.url.pathname === '/meal-plan' ? 'text-emerald-400 bg-emerald-950/30' : 'text-slate-300 hover:text-white hover:bg-slate-900'}"
					>
						<Utensils class="h-4 w-4" />
						<span>Meal Plan</span>
					</a>
					<a
						href="/workout"
						class="flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors {page.url.pathname === '/workout' ? 'text-emerald-400 bg-emerald-950/30' : 'text-slate-300 hover:text-white hover:bg-slate-900'}"
					>
						<Dumbbell class="h-4 w-4" />
						<span>Workout</span>
					</a>
				</div>

				<!-- Streak Counter -->
				{#if streakCount > 0}
					<div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-orange-950/35 border border-orange-500/25 text-orange-400 text-xs font-semibold select-none animate-pulse">
						<Flame class="h-4 w-4 fill-orange-500 text-orange-500" />
						<span>{streakCount} Day Streak</span>
					</div>
				{/if}

				<!-- Logout Button -->
				<button
					onclick={handleLogout}
					class="flex items-center gap-2 px-3.5 py-1.5 rounded-lg border border-slate-700 text-slate-300 hover:text-red-400 hover:border-red-500/30 hover:bg-red-950/20 text-sm font-medium transition-all cursor-pointer"
				>
					<LogOut class="h-4 w-4" />
					<span class="hidden sm:inline">Logout</span>
				</button>
			</div>
		</div>
	</div>
</nav>
