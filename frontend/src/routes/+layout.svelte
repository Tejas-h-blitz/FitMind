<script lang="ts">
	import './layout.css';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { supabase } from '$lib/supabase';
	import { api, type UserProfile } from '$lib/api';
	import Navbar from '$lib/components/Navbar.svelte';
	import { Toaster } from 'svelte-sonner';

	let { children } = $props();

	let session = $state<any>(null);
	let userProfile = $state<UserProfile | null>(null);
	let isInitialized = $state(false);

	// Page status helpers
	let isAuthPage = $derived(page.url.pathname.startsWith('/auth/'));
	let isLandingPage = $derived(page.url.pathname === '/');

	onMount(() => {
		// 1. Get initial session
		supabase.auth.getSession().then(({ data }) => {
			session = data.session;
			syncSessionCookie(data.session).then(() => {
				if (data.session) {
					loadUserProfile();
				}
				checkAuthGuard();
				isInitialized = true;
			});
		});

		// 2. Subscribe to auth changes
		const { data: { subscription } } = supabase.auth.onAuthStateChange(async (event, newSession) => {
			session = newSession;
			await syncSessionCookie(newSession);

			if (event === 'SIGNED_IN' || event === 'USER_UPDATED') {
				await loadUserProfile();
			} else if (event === 'SIGNED_OUT') {
				userProfile = null;
			}
			checkAuthGuard();
		});

		return () => {
			subscription.unsubscribe();
		};
	});

	async function syncSessionCookie(currSession: any) {
		await fetch('/auth/session', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ access_token: currSession?.access_token || null })
		});
	}

	async function loadUserProfile() {
		const res = await api.verifyAuth();
		if (res.success && res.data) {
			userProfile = res.data;
		}
	}

	function checkAuthGuard() {
		const loggedIn = !!session;
		if (!loggedIn && !isAuthPage && !isLandingPage) {
			goto('/auth/login');
		} else if (loggedIn && isAuthPage) {
			goto('/dashboard');
		}
	}

	// Recheck guards when route path changes
	$effect(() => {
		if (isInitialized) {
			// Trigger recheck on route path dependencies
			const _ = page.url.pathname;
			checkAuthGuard();
		}
	});
</script>

<Toaster theme="dark" position="bottom-right" richColors />

<div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans">
	{#if isInitialized}
		{#if !isAuthPage && !isLandingPage && session}
			<Navbar streakCount={userProfile?.streak_count || 0} />
		{/if}

		<main class="flex-1">
			{@render children()}
		</main>
	{:else}
		<!-- Loading state -->
		<div class="h-screen w-screen flex flex-col justify-center items-center bg-slate-950">
			<div class="h-10 w-10 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
			<span class="text-sm font-semibold text-slate-400 mt-4 animate-pulse">Initializing FitMind...</span>
		</div>
	{/if}
</div>
