<script lang="ts">
	let { steps = [], currentStep = 0 } = $props<{
		steps: string[];
		currentStep: number;
	}>();

	let percentage = $derived(steps.length > 1 ? (currentStep / (steps.length - 1)) * 100 : 0);
</script>

<div class="w-full py-4 max-w-xl mx-auto select-none relative z-10">
	<div class="relative flex items-center justify-between">
		<!-- Progress Line Background -->
		<div class="absolute left-4 right-4 top-1/2 -translate-y-1/2 h-0.5 bg-slate-900/80 -z-10">
			<!-- Fill Line -->
			<div
				class="h-full bg-gradient-to-r from-emerald-500 to-teal-500 transition-all duration-500 ease-out"
				style="width: {percentage}%"
			></div>
		</div>

		<!-- Step circles -->
		{#each steps as step, idx}
			<div class="flex flex-col items-center">
				<div
					class="h-7.5 w-7.5 rounded-lg border flex items-center justify-center text-[10px] font-bold transition-all duration-500
						{idx <= currentStep
							? 'bg-emerald-950/85 border-emerald-500 text-emerald-400 shadow-md shadow-emerald-550/10 scale-105'
							: 'bg-slate-950 border-slate-850 text-slate-500'}"
				>
					{idx + 1}
				</div>
				<span
					class="text-[8px] font-extrabold uppercase tracking-widest mt-2 transition-colors duration-500
						{idx <= currentStep ? 'text-slate-200' : 'text-slate-500'}"
				>
					{step}
				</span>
			</div>
		{/each}
	</div>
</div>
