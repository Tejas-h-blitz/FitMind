<script lang="ts">
	import type { DocumentAnalysis } from '$lib/api';
	import { Heart, Activity, AlertTriangle, CheckCircle, Info } from 'lucide-svelte';

	let { analysis } = $props<{ analysis: DocumentAnalysis }>();

	// Expanded state for each metric card to show plain english explanation on mobile/click
	let expandedMetrics = $state<Record<string, boolean>>({});

	function toggleExpand(name: string) {
		expandedMetrics[name] = !expandedMetrics[name];
	}

	function getStatusBadgeClass(status: string) {
		switch (status.toLowerCase()) {
			case 'normal':
				return 'bg-emerald-950/40 border-emerald-500/30 text-emerald-400';
			case 'low':
				return 'bg-sky-950/40 border-sky-500/30 text-sky-400';
			case 'borderline':
				return 'bg-amber-950/40 border-amber-500/30 text-amber-400';
			case 'high':
				return 'bg-rose-950/40 border-rose-500/30 text-rose-400';
			default:
				return 'bg-slate-900/40 border-slate-700/30 text-slate-400';
		}
	}
</script>

<div class="space-y-6">
	<!-- Top: Overall Status Badge -->
	<div class="flex items-center justify-between flex-wrap gap-4 p-4 rounded-xl border border-slate-800 bg-slate-950/30">
		<div class="flex items-center gap-3">
			<div class="p-2.5 rounded-lg bg-emerald-500/10 text-emerald-400">
				<Activity class="h-5 w-5" />
			</div>
			<div>
				<h3 class="text-sm font-bold text-slate-200">Overall Health Assessment</h3>
				<p class="text-[11px] text-slate-400 mt-0.5">Automated screening based on document findings.</p>
			</div>
		</div>

		<div>
			{#if analysis.overall_status === 'good'}
				<span class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold bg-emerald-950/60 border border-emerald-500 text-emerald-400">
					<CheckCircle class="h-4 w-4 fill-emerald-950" />
					All Clear
				</span>
			{:else if analysis.overall_status === 'needs_attention'}
				<span class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold bg-amber-950/60 border border-amber-500 text-amber-400 animate-pulse">
					<AlertTriangle class="h-4 w-4 fill-amber-950" />
					Needs Attention
				</span>
			{:else if analysis.overall_status === 'concerning'}
				<span class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold bg-rose-950/60 border border-rose-550 text-rose-400">
					<Heart class="h-4 w-4 fill-rose-950" />
					See a Doctor
				</span>
			{/if}
		</div>
	</div>

	<!-- Empty State -->
	{#if !analysis.metrics || analysis.metrics.length === 0}
		<div class="py-12 border border-dashed border-slate-800 rounded-xl text-center">
			<Info class="h-8 w-8 text-slate-500 mx-auto mb-3" />
			<p class="text-sm text-slate-400">No medical health metrics could be extracted from this document.</p>
			<p class="text-xs text-slate-500 mt-1">This report might not be a health metric sheet, or it doesn't contain recognized test values.</p>
		</div>
	{:else}
		<!-- Metrics Grid: 3 columns on desktop, 1 on mobile -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each analysis.metrics as metric}
				<div 
					onclick={() => toggleExpand(metric.name)}
					class="group relative p-4 rounded-xl border border-slate-800 bg-slate-900/30 hover:border-slate-700/60 hover:bg-slate-900/50 transition-all duration-200 cursor-pointer flex flex-col justify-between"
				>
					<div>
						<!-- Header -->
						<div class="flex items-start justify-between gap-2">
							<span class="text-xs font-bold text-slate-300 capitalize group-hover:text-emerald-400 transition-colors line-clamp-1">
								{metric.name.replace(/_/g, ' ')}
							</span>
							<span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-bold border capitalize {getStatusBadgeClass(metric.status)}">
								{metric.status}
							</span>
						</div>

						<!-- Value -->
						<div class="mt-3 flex items-baseline gap-1">
							<span class="text-2xl font-extrabold text-slate-100">{metric.value}</span>
							<span class="text-xs text-slate-400 font-semibold">{metric.unit}</span>
						</div>
					</div>

					<!-- Reference range -->
					<div class="mt-4 pt-3 border-t border-slate-900 flex justify-between items-center text-[10px]">
						<span class="text-slate-500">Ref: <span class="text-slate-400 font-medium">{metric.reference_range}</span></span>
						<span class="text-emerald-400 font-bold group-hover:underline flex items-center gap-0.5">
							{expandedMetrics[metric.name] ? 'Hide explanation' : 'Show explanation'}
						</span>
					</div>

					<!-- Expandable Plain English -->
					{#if expandedMetrics[metric.name]}
						<div class="mt-3 p-2.5 rounded-lg bg-slate-950/60 border border-slate-800 text-[11px] text-slate-300 leading-relaxed animate-fadeIn">
							{metric.plain_english}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}

	<!-- Summary Box -->
	<div class="p-4.5 rounded-xl border border-emerald-500/20 bg-emerald-950/10">
		<h4 class="text-xs font-bold uppercase tracking-wider text-emerald-400 flex items-center gap-1.5">
			<Brain class="h-4 w-4" />
			<span>Advisor Summary</span>
		</h4>
		<p class="text-sm text-slate-200 mt-2 leading-relaxed">
			{analysis.summary}
		</p>
	</div>

	<!-- Disclaimer -->
	<div class="text-[10px] text-slate-500 italic text-center max-w-lg mx-auto">
		⚠️ Medical Disclaimer: This analysis is AI-generated and not a substitute for professional medical advice, diagnosis, or treatment. Always consult a healthcare provider.
	</div>
</div>

<style>
	:global(@keyframes fadeIn) {
		from { opacity: 0; transform: translateY(-3px); }
		to { opacity: 1; transform: translateY(0); }
	}
	.animate-fadeIn {
		animation: fadeIn 0.2s ease-out forwards;
	}
</style>
