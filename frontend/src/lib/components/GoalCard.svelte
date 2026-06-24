<script lang="ts">
	import type { Goal } from '$lib/api';
	import { CheckSquare, Square, Calendar, CheckCircle2 } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let { goal, onToggle } = $props<{
		goal: Goal;
		onToggle: (id: string, nextStatus: 'active' | 'completed') => Promise<void>;
	}>();

	let isToggling = $state(false);
	const isCompleted = $derived(goal.status === 'completed');

	async function handleToggleClick() {
		if (isToggling) return;
		isToggling = true;
		const nextStatus = isCompleted ? 'active' : 'completed';

		try {
			await onToggle(goal.id, nextStatus);
			toast.success(nextStatus === 'completed' ? 'Goal completed! Keep it up!' : 'Goal marked active');
		} catch (err: any) {
			toast.error(err.message || 'Failed to update goal');
		} finally {
			isToggling = false;
		}
	}

	function formatDate(dateStr: string) {
		if (!dateStr) return 'No target date';
		const date = new Date(dateStr);
		return date.toLocaleDateString(undefined, {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function isOverdue(dateStr: string) {
		if (!dateStr || isCompleted) return false;
		const target = new Date(dateStr);
		const today = new Date();
		today.setHours(0,0,0,0);
		return target < today;
	}
</script>

<div
	class="flex items-center justify-between rounded-xl border p-4 transition-all duration-350 select-none
		{isCompleted
			? 'border-slate-800 bg-slate-900/10 opacity-70'
			: isOverdue(goal.target_date)
				? 'border-red-900/30 bg-red-950/5 hover:border-red-500/20'
				: 'border-slate-850 bg-slate-950/60 hover:border-slate-700/50'
		}"
>
	<div class="flex items-center gap-3.5 flex-1 min-w-0">
		<!-- Checkbox Toggle -->
		<button
			onclick={handleToggleClick}
			disabled={isToggling}
			class="flex items-center justify-center p-1 rounded-lg hover:bg-slate-900 text-slate-400 hover:text-emerald-400 transition-colors disabled:opacity-50 cursor-pointer"
			aria-label="Toggle goal status"
		>
			{#if isCompleted}
				<CheckCircle2 class="h-5.5 w-5.5 text-emerald-500 fill-emerald-500/10" />
			{:else}
				<Square class="h-5.5 w-5.5 text-slate-600" />
			{/if}
		</button>

		<div class="min-w-0 flex-1">
			<span
				class="block text-sm font-semibold text-slate-200 truncate leading-snug
					{isCompleted ? 'line-through text-slate-500' : ''}"
			>
				{goal.title}
			</span>
			
			{#if goal.target_date}
				<span class="inline-flex items-center gap-1 mt-1 text-[11px] font-medium
					{isOverdue(goal.target_date) ? 'text-red-400' : 'text-slate-400'}"
				>
					<Calendar class="h-3 w-3" />
					Target: {formatDate(goal.target_date)}
					{#if isOverdue(goal.target_date)}
						• Overdue
					{/if}
				</span>
			{/if}
		</div>
	</div>

	<!-- Status Badge -->
	<div class="ml-4 shrink-0">
		{#if isCompleted}
			<span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-950/40 border border-emerald-500/25 text-emerald-400 uppercase tracking-wider">
				Done
			</span>
		{:else if isOverdue(goal.target_date)}
			<span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-red-950/40 border border-red-500/25 text-red-400 uppercase tracking-wider">
				Overdue
			</span>
		{:else}
			<span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-slate-900 border border-slate-800 text-slate-400 uppercase tracking-wider">
				Active
			</span>
		{/if}
	</div>
</div>
