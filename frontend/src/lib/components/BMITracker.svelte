<script lang="ts">
	import type { HealthMetric } from '$lib/api';
	import { Activity, Calculator, Plus, TrendingUp } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import AnimatedGauge from './AnimatedGauge.svelte';

	let { metrics = [], onLog } = $props<{
		metrics: HealthMetric[];
		onLog: (height: number, weight: number) => Promise<void>;
	}>();

	let height = $state<number | ''>('');
	let weight = $state<number | ''>('');
	let isSubmitting = $state(false);

	// Get latest metric for current display
	let latestMetric = $derived(metrics.length > 0 ? metrics[metrics.length - 1] : null);

	// Determine BMI Category and Color
	let bmiCategory = $derived(() => {
		if (!latestMetric) return { label: 'No Data', color: 'text-slate-400', bg: 'bg-slate-950/40', border: 'border-slate-850' };
		const val = latestMetric.bmi;
		if (val < 18.5) return { label: 'Underweight', color: 'text-yellow-400', bg: 'bg-yellow-500/5', border: 'border-yellow-500/20' };
		if (val < 25) return { label: 'Healthy Weight', color: 'text-emerald-400', bg: 'bg-emerald-500/5', border: 'border-emerald-500/20' };
		if (val < 30) return { label: 'Overweight', color: 'text-orange-400', bg: 'bg-orange-500/5', border: 'border-orange-500/20' };
		return { label: 'Obese', color: 'text-red-400', bg: 'bg-red-500/5', border: 'border-red-500/20' };
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		if (!height || !weight || isSubmitting) return;
		isSubmitting = true;

		try {
			await onLog(Number(height), Number(weight));
			toast.success('BMI entry recorded successfully');
			weight = '';
		} catch (err: any) {
			toast.error(err.message || 'Failed to log metrics');
		} finally {
			isSubmitting = false;
		}
	}

	// SVG Graph helper logic
	let svgPoints = $derived(() => {
		if (metrics.length < 2) return '';
		const width = 460;
		const height = 110;
		const padding = 15;

		const bmis = metrics.map((m) => m.bmi);
		const maxVal = Math.max(...bmis) + 2;
		const minVal = Math.max(0, Math.min(...bmis) - 2);
		const valRange = maxVal - minVal || 1;

		return metrics
			.map((m, idx) => {
				const x = padding + (idx / (metrics.length - 1)) * (width - 2 * padding);
				const y = height - padding - ((m.bmi - minVal) / valRange) * (height - 2 * padding);
				return `${x},${y}`;
			})
			.join(' ');
	});
</script>

<div class="rounded-2xl border border-slate-850 bg-slate-950/30 p-6 shadow-xl backdrop-blur-md">
	<div class="flex items-center gap-2.5 mb-6 pb-3 border-b border-slate-900/60">
		<div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400">
			<Activity class="h-5 w-5" />
		</div>
		<h2 class="text-base font-extrabold text-slate-100 tracking-tight">BMI Tracker & History</h2>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 gap-8">
		<!-- Calculator Form / Current Stats -->
		<div class="space-y-6">
			{#if latestMetric}
				<div class="flex flex-col sm:flex-row items-center gap-6 p-4 rounded-xl border border-slate-850 bg-slate-950/45 shadow-inner">
					<AnimatedGauge bmi={latestMetric.bmi} />
					
					<div class="flex-1 text-center sm:text-left space-y-2 border-t sm:border-t-0 sm:border-l border-slate-900 pt-4 sm:pt-0 sm:pl-4">
						<span class="text-[9px] font-extrabold uppercase tracking-widest text-slate-500">Physical Metrics</span>
						<div class="space-y-1 font-semibold text-xs">
							<p class="text-slate-400">Height: <span class="text-slate-205 font-black">{latestMetric.height} cm</span></p>
							<p class="text-slate-400">Weight: <span class="text-slate-205 font-black">{latestMetric.weight} kg</span></p>
						</div>
					</div>
				</div>
			{/if}

			<!-- Entry Form -->
			<form onsubmit={handleSubmit} class="space-y-4">
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="height" class="block text-xs font-bold text-slate-450 mb-1.5">Height (cm)</label>
						<input
							id="height"
							type="number"
							step="0.1"
							bind:value={height}
							placeholder="e.g. 175"
							required
							class="w-full px-3.5 py-2 text-xs font-semibold rounded-lg border border-slate-850 bg-slate-950/60 text-slate-100 focus:outline-none focus:border-emerald-500/80 transition-all placeholder-slate-650"
						/>
					</div>
					<div>
						<label for="weight" class="block text-xs font-bold text-slate-450 mb-1.5">Weight (kg)</label>
						<input
							id="weight"
							type="number"
							step="0.1"
							bind:value={weight}
							placeholder="e.g. 70"
							required
							class="w-full px-3.5 py-2 text-xs font-semibold rounded-lg border border-slate-850 bg-slate-950/60 text-slate-100 focus:outline-none focus:border-emerald-500/80 transition-all placeholder-slate-650"
						/>
					</div>
				</div>

				<button
					type="submit"
					disabled={isSubmitting || !height || !weight}
					class="w-full flex items-center justify-center gap-2 py-2.5 px-4 bg-emerald-600 hover:bg-emerald-500 text-white font-bold rounded-lg text-xs transition-all cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed hover:-translate-y-0.5 active:translate-y-0"
				>
					<Plus class="h-4 w-4" />
					<span>Log Health Metric</span>
				</button>
			</form>
		</div>

		<!-- History Chart -->
		<div class="flex flex-col justify-between rounded-xl border border-slate-850 bg-slate-950/20 p-4">
			<div class="flex justify-between items-center mb-4">
				<span class="text-xs font-bold text-slate-350 flex items-center gap-2">
					<TrendingUp class="h-4 w-4 text-emerald-400" />
					<span>BMI Progress Trend</span>
				</span>
			</div>

			<div class="flex-1 flex items-center justify-center">
				{#if metrics.length === 0}
					<p class="text-xs text-slate-500 italic p-6 text-center">No health logs found. Add weight and height metrics to plot.</p>
				{:else if metrics.length === 1}
					<div class="flex flex-col items-center justify-center p-6 text-center text-xs text-slate-400">
						<p class="font-semibold text-slate-300">Initial log registered: {metrics[0].bmi.toFixed(1)} BMI</p>
						<p class="text-slate-500 mt-1">Add one more measurement to draw dynamic trend lines.</p>
					</div>
				{:else}
					<!-- SVG Polyline Graph -->
					<div class="w-full h-full min-h-[110px]">
						<svg viewBox="0 0 460 110" class="w-full h-full overflow-visible">
							<!-- Grid Lines -->
							<line x1="15" y1="15" x2="445" y2="15" stroke="#111827" stroke-dasharray="3" />
							<line x1="15" y1="55" x2="445" y2="55" stroke="#111827" stroke-dasharray="3" />
							<line x1="15" y1="95" x2="445" y2="95" stroke="#111827" stroke-dasharray="3" />

							<!-- Trend line -->
							<polyline
								fill="none"
								stroke="#10b981"
								stroke-width="2.5"
								stroke-linecap="round"
								stroke-linejoin="round"
								points={svgPoints()}
							/>

							<!-- Data points -->
							{#each metrics as m, idx}
								{@const width = 460}
								{@const height = 110}
								{@const padding = 15}
								{@const bmis = metrics.map(x => x.bmi)}
								{@const maxVal = Math.max(...bmis) + 2}
								{@const minVal = Math.max(0, Math.min(...bmis) - 2)}
								{@const valRange = maxVal - minVal || 1}
								{@const cx = padding + (idx / (metrics.length - 1)) * (width - 2 * padding)}
								{@const cy = height - padding - ((m.bmi - minVal) / valRange) * (height - 2 * padding)}
								
								<g class="group/point cursor-pointer">
									<circle
										{cx}
										{cy}
										r="3.5"
										class="fill-slate-950 stroke-emerald-400 stroke-2 hover:r-5 hover:fill-emerald-450 transition-all duration-200"
									/>
									<title>
										BMI: {m.bmi.toFixed(1)} ({new Date(m.recorded_at).toLocaleDateString()})
									</title>
								</g>
							{/each}
						</svg>
					</div>
				{/if}
			</div>

			<div class="flex justify-between items-center mt-3 pt-2 border-t border-slate-900 text-[10px] text-slate-500 font-bold tracking-tight">
				<span>Underweight (&lt;18.5)</span>
				<span>Normal (18.5-24.9)</span>
				<span>Overweight (&gt;25)</span>
			</div>
		</div>
	</div>
</div>
