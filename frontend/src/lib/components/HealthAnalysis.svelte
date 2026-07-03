<script lang="ts">
	import type { DocumentAnalysis } from '$lib/api';
	import { Heart, Activity, AlertTriangle, CheckCircle, Info, Brain } from 'lucide-svelte';

	let { analysis } = $props<{ analysis: DocumentAnalysis }>();

	// Expanded state for each metric card to show plain english explanation on mobile/click
	let expandedMetrics = $state<Record<string, boolean>>({});

	function toggleExpand(name: string) {
		expandedMetrics[name] = !expandedMetrics[name];
	}

	function getStatusBadgeClass(status: string) {
		switch (status.toLowerCase()) {
			case 'normal':
				return 'bg-emerald-950/45 border-emerald-500/20 text-emerald-400';
			case 'low':
				return 'bg-sky-950/45 border-sky-500/20 text-sky-400';
			case 'borderline':
				return 'bg-amber-950/45 border-amber-500/20 text-amber-400';
			case 'high':
				return 'bg-rose-950/45 border-rose-500/20 text-rose-450';
			default:
				return 'bg-slate-900/40 border-slate-800 text-slate-400';
		}
	}
</script>

<div class="space-y-6">
	<!-- Top: Overall Status Badge -->
	<div class="flex items-center justify-between flex-wrap gap-4 p-5 rounded-2xl border border-slate-850 bg-slate-950/40 shadow-inner">
		<div class="flex items-center gap-3">
			<div class="p-2.5 rounded-xl bg-emerald-500/10 text-emerald-400">
				<Activity class="h-5 w-5" />
			</div>
			<div>
				<h3 class="text-sm font-extrabold text-slate-200 tracking-tight">Overall Health Assessment</h3>
				<p class="text-[10px] text-slate-500 font-bold mt-0.5 uppercase tracking-wider">Automated screening based on document findings.</p>
			</div>
		</div>

		<div>
			{#if analysis.overall_status === 'good'}
				<span class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-bold bg-emerald-950/60 border border-emerald-500/30 text-emerald-400">
					<CheckCircle class="h-4 w-4 fill-emerald-950 text-emerald-500" />
					All Clear
				</span>
			{:else if analysis.overall_status === 'needs_attention'}
				<span class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-bold bg-amber-950/60 border border-amber-500/30 text-amber-400 animate-pulse">
					<AlertTriangle class="h-4 w-4 fill-amber-950 text-amber-500" />
					Needs Attention
				</span>
			{:else if analysis.overall_status === 'concerning'}
				<span class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-bold bg-rose-950/60 border border-rose-550/30 text-rose-450">
					<Heart class="h-4 w-4 fill-rose-950 text-rose-500" />
					See a Doctor
				</span>
			{/if}
		</div>
	</div>

	<!-- Empty State -->
	{#if !analysis.metrics || analysis.metrics.length === 0}
		<div class="py-12 border border-dashed border-slate-850 rounded-2xl text-center shadow-inner bg-slate-950/10">
			<Info class="h-8 w-8 text-slate-500 mx-auto mb-3" />
			<p class="text-sm text-slate-350 font-semibold">No medical health metrics could be extracted from this document.</p>
			<p class="text-xs text-slate-500 mt-1 max-w-sm mx-auto leading-relaxed">This report might not be a health metric sheet, or it doesn't contain recognized test values.</p>
		</div>
	{:else}
		<!-- Metrics Grid: 3 columns on desktop, 1 on mobile -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each analysis.metrics as metric}
				<div 
					onclick={() => toggleExpand(metric.name)}
					class="group relative p-5 rounded-2xl border border-slate-850 bg-slate-950/20 hover:border-slate-700/60 hover:bg-slate-900/10 transition-all duration-200 cursor-pointer flex flex-col justify-between hover:-translate-y-0.5 shadow-sm"
				>
					<div>
						<!-- Header -->
						<div class="flex items-start justify-between gap-2">
							<span class="text-xs font-bold text-slate-300 capitalize group-hover:text-emerald-450 transition-colors line-clamp-1">
								{metric.name.replace(/_/g, ' ')}
							</span>
							<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-[9px] font-extrabold border capitalize tracking-wider {getStatusBadgeClass(metric.status)}">
								{metric.status}
							</span>
						</div>

						<!-- Value -->
						<div class="mt-4 flex items-baseline gap-1.5">
							<span class="text-2xl font-black text-slate-100 tracking-tight">{metric.value}</span>
							<span class="text-xs text-slate-500 font-bold">{metric.unit}</span>
						</div>
					</div>

					<!-- Reference range -->
					<div class="mt-4 pt-3 border-t border-slate-900 flex justify-between items-center text-[10px] font-bold">
						<span class="text-slate-500">Ref: <span class="text-slate-400 font-semibold">{metric.reference_range}</span></span>
						<span class="text-emerald-450 group-hover:underline flex items-center gap-0.5">
							{expandedMetrics[metric.name] ? 'Hide details' : 'View details'}
						</span>
					</div>

					<!-- Expandable Plain English -->
					{#if expandedMetrics[metric.name]}
						<div class="mt-3 p-3 rounded-xl bg-slate-950/60 border border-slate-900 text-[11px] text-slate-350 leading-relaxed font-medium animate-fadeIn">
							{metric.plain_english}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}

	<!-- Summary Box -->
	<div class="p-5 rounded-2xl border border-emerald-500/10 bg-emerald-950/5">
		<h4 class="text-xs font-extrabold uppercase tracking-wider text-emerald-400 flex items-center gap-2">
			<Brain class="h-4.5 w-4.5" />
			<span>Advisor Summary</span>
		</h4>
		<p class="text-xs text-slate-300 mt-2.5 leading-relaxed font-semibold">
			{analysis.summary}
		</p>
	</div>

	<!-- Disclaimer -->
	<div class="text-[10px] text-slate-500 italic text-center max-w-lg mx-auto">
		⚠️ Medical Disclaimer: This analysis is AI-generated and not a substitute for professional medical advice, diagnosis, or treatment. Always consult a healthcare provider.
	</div>
</div>

<style>
	@keyframes fadeIn {
		from { opacity: 0; transform: translateY(-3px); }
		to { opacity: 1; transform: translateY(0); }
	}
	.animate-fadeIn {
		animation: fadeIn 0.2s ease-out forwards;
	}
</style>
