<script lang="ts">
	let { bmi = 0 } = $props<{ bmi: number }>();

	const size = 150;
	const strokeWidth = 12;
	const radius = (size - strokeWidth) / 2;
	const circumference = 2 * Math.PI * radius;

	// Map BMI range (15 - 35) to (0 - 100%)
	let percentage = $derived(Math.min(100, Math.max(0, ((bmi - 15) / 20) * 100)));
	let strokeDashoffset = $derived(bmi > 0 ? circumference - (percentage / 100) * circumference : circumference);

	function getCategory(val: number) {
		if (val === 0) return { label: 'No Data', color: 'text-slate-500', stroke: 'url(#emeraldGrad)' };
		if (val < 18.5) return { label: 'Underweight', color: 'text-yellow-400', stroke: 'url(#yellowGrad)', shadowColor: '#fbbf24' };
		if (val < 25) return { label: 'Healthy Weight', color: 'text-emerald-450', stroke: 'url(#emeraldGrad)', shadowColor: '#10b981' };
		if (val < 30) return { label: 'Overweight', color: 'text-orange-400', stroke: 'url(#orangeGrad)', shadowColor: '#f97316' };
		return { label: 'Obese', color: 'text-rose-455', stroke: 'url(#roseGrad)', shadowColor: '#f43f5e' };
	}

	let cat = $derived(getCategory(bmi));
</script>

<div class="flex flex-col items-center justify-center p-2 relative select-none">
	<svg width={size} height={size} viewBox="0 0 {size} {size}" class="transform -rotate-90 overflow-visible">
		<defs>
			<linearGradient id="emeraldGrad" x1="0%" y1="0%" x2="100%" y2="100%">
				<stop offset="0%" stop-color="#10b981" />
				<stop offset="100%" stop-color="#059669" />
			</linearGradient>
			<linearGradient id="yellowGrad" x1="0%" y1="0%" x2="100%" y2="100%">
				<stop offset="0%" stop-color="#fbbf24" />
				<stop offset="100%" stop-color="#d97706" />
			</linearGradient>
			<linearGradient id="orangeGrad" x1="0%" y1="0%" x2="100%" y2="100%">
				<stop offset="0%" stop-color="#f97316" />
				<stop offset="100%" stop-color="#ea580c" />
			</linearGradient>
			<linearGradient id="roseGrad" x1="0%" y1="0%" x2="100%" y2="100%">
				<stop offset="0%" stop-color="#f43f5e" />
				<stop offset="100%" stop-color="#e11d48" />
			</linearGradient>
			{#if bmi > 0}
				<filter id="glow" x="-20%" y="-20%" width="140%" height="140%">
					<feDropShadow dx="0" dy="0" stdDeviation="6" flood-color={cat.shadowColor} flood-opacity="0.25" />
				</filter>
			{/if}
		</defs>

		<!-- Background Circle -->
		<circle
			cx={size / 2}
			cy={size / 2}
			r={radius}
			class="stroke-slate-900/80"
			stroke-width={strokeWidth}
			fill="none"
		/>

		<!-- Foreground Progress Circle -->
		{#if bmi > 0}
			<circle
				cx={size / 2}
				cy={size / 2}
				r={radius}
				stroke={cat.stroke}
				stroke-width={strokeWidth}
				stroke-dasharray={circumference}
				stroke-dashoffset={strokeDashoffset}
				stroke-linecap="round"
				fill="none"
				class="transition-all duration-1000 ease-out"
				filter="url(#glow)"
			/>
		{/if}
	</svg>

	<!-- Centered details -->
	<div class="absolute inset-0 flex flex-col items-center justify-center text-center">
		{#if bmi > 0}
			<span class="text-3xl font-black text-slate-100 tracking-tighter mt-1">{bmi.toFixed(1)}</span>
			<span class="text-[9px] font-extrabold uppercase tracking-widest mt-0.5 {cat.color}">{cat.label}</span>
		{:else}
			<span class="text-lg font-black text-slate-650 tracking-tight">--.-</span>
			<span class="text-[9px] font-extrabold text-slate-600 uppercase tracking-wider mt-0.5">No Records</span>
		{/if}
	</div>
</div>
