<script lang="ts">
	import { supabase } from '$lib/supabase';
	import { goto } from '$app/navigation';
	import { Brain, UserPlus, Lock, Mail, User, AlertTriangle } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let fullName = $state('');
	let email = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let errorMsg = $state('');

	async function handleSignup(e: Event) {
		e.preventDefault();
		if (!fullName || !email || !password || isLoading) return;

		isLoading = true;
		errorMsg = '';

		try {
			const { error } = await supabase.auth.signUp({
				email,
				password,
				options: {
					data: {
						full_name: fullName
					}
				}
			});

			if (error) {
				errorMsg = error.message;
				toast.error(error.message);
			} else {
				toast.success('Registration successful! Please check your email for confirmation or log in.');
				goto('/auth/login');
			}
		} catch (err: any) {
			errorMsg = err.message || 'An unexpected error occurred';
			toast.error(errorMsg);
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="relative bg-slate-950 min-h-screen text-slate-100 flex flex-col justify-center items-center px-4">
	<!-- Background grid overlay -->
	<div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-emerald-950/10 via-slate-950 to-slate-950 -z-10"></div>

	<div class="w-full max-w-md bg-slate-900/40 border border-slate-800/80 rounded-2xl p-8 backdrop-blur-md shadow-2xl">
		<!-- Brand Header -->
		<div class="flex flex-col items-center mb-8">
			<div class="p-3 bg-emerald-500/10 rounded-xl text-emerald-400 mb-3">
				<Brain class="h-8 w-8" />
			</div>
			<h2 class="text-2xl font-extrabold text-slate-100">Create Account</h2>
			<p class="text-slate-400 text-sm mt-1">Get started with your health companion</p>
		</div>

		<!-- Error Message Alert -->
		{#if errorMsg}
			<div class="flex items-center gap-2 p-3.5 mb-5 rounded-lg bg-red-500/10 border border-red-500/25 text-red-400 text-xs">
				<AlertTriangle class="h-4 w-4 shrink-0" />
				<span>{errorMsg}</span>
			</div>
		{/if}

		<!-- Sign Up Form -->
		<form onsubmit={handleSignup} class="space-y-4">
			<div>
				<label for="fullName" class="block text-xs font-semibold text-slate-400 mb-1.5">Full Name</label>
				<div class="relative">
					<span class="absolute inset-y-0 left-0 pl-3 flex items-center text-slate-500 pointer-events-none">
						<User class="h-4 w-4" />
					</span>
					<input
						id="fullName"
						type="text"
						bind:value={fullName}
						placeholder="John Doe"
						required
						disabled={isLoading}
						class="w-full pl-10 pr-4 py-2.5 text-sm rounded-lg border border-slate-800 bg-slate-950 text-slate-100 focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all disabled:opacity-50"
					/>
				</div>
			</div>

			<div>
				<label for="email" class="block text-xs font-semibold text-slate-400 mb-1.5">Email Address</label>
				<div class="relative">
					<span class="absolute inset-y-0 left-0 pl-3 flex items-center text-slate-500 pointer-events-none">
						<Mail class="h-4 w-4" />
					</span>
					<input
						id="email"
						type="email"
						bind:value={email}
						placeholder="you@example.com"
						required
						disabled={isLoading}
						class="w-full pl-10 pr-4 py-2.5 text-sm rounded-lg border border-slate-800 bg-slate-950 text-slate-100 focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all disabled:opacity-50"
					/>
				</div>
			</div>

			<div>
				<label for="password" class="block text-xs font-semibold text-slate-400 mb-1.5">Password</label>
				<div class="relative">
					<span class="absolute inset-y-0 left-0 pl-3 flex items-center text-slate-500 pointer-events-none">
						<Lock class="h-4 w-4" />
					</span>
					<input
						id="password"
						type="password"
						bind:value={password}
						placeholder="Min. 8 characters"
						minlength="8"
						required
						disabled={isLoading}
						class="w-full pl-10 pr-4 py-2.5 text-sm rounded-lg border border-slate-800 bg-slate-950 text-slate-100 focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all disabled:opacity-50"
					/>
				</div>
			</div>

			<button
				type="submit"
				disabled={isLoading}
				class="w-full flex items-center justify-center gap-2 py-3 px-4 bg-emerald-600 hover:bg-emerald-500 text-white font-bold rounded-lg text-sm shadow-md transition-all cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed mt-6"
			>
				{#if isLoading}
					<div class="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
					<span>Creating Account...</span>
				{:else}
					<UserPlus class="h-4 w-4" />
					<span>Register</span>
				{/if}
			</button>
		</form>

		<div class="mt-8 pt-6 border-t border-slate-800/60 text-center">
			<p class="text-xs text-slate-400">
				Already have an account? 
				<a href="/auth/login" class="text-emerald-400 font-semibold hover:underline ml-1">Log in here</a>
			</p>
		</div>
	</div>
</div>
