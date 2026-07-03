<script lang="ts">
	import { FileText, Split, Cpu, Database } from 'lucide-svelte';

	let { status = 'idle' } = $props<{
		status: 'idle' | 'uploading' | 'processing' | 'ready';
	}>();

	// Step indices:
	// 0: Uploading (PDF)
	// 1: Chunking
	// 2: Embedding (OpenAI)
	// 3: Indexing (Qdrant)
	let activeStep = $derived(() => {
		if (status === 'uploading') return 0;
		if (status === 'processing') return 2; 
		if (status === 'ready') return 3;
		return -1;
	});
</script>

<div class="w-full py-8 px-4 bg-slate-950/45 border border-slate-850 rounded-2xl shadow-xl backdrop-blur-md relative overflow-hidden">
	<div class="absolute -right-24 -bottom-24 h-48 w-48 rounded-full bg-emerald-500/5 blur-3xl pointer-events-none"></div>

	<!-- Title -->
	<div class="text-center mb-8">
		<span class="text-[9px] font-extrabold uppercase tracking-widest text-emerald-450">Active RAG Pipeline Status</span>
		<h4 class="text-sm font-black text-slate-200 mt-1">AI Health Document Ingestion Pipeline</h4>
	</div>

	<!-- Flow Container -->
	<div class="flex flex-col md:flex-row items-center justify-between gap-6 md:gap-4 max-w-2xl mx-auto relative">
		<!-- Steps -->
		<!-- Step 1: Upload PDF -->
		<div class="flex flex-col items-center z-10">
			<div class="h-14 w-14 rounded-2xl flex items-center justify-center border transition-all duration-500
				{status === 'ready' || activeStep() >= 0
					? 'bg-emerald-950/20 border-emerald-500 text-emerald-400 shadow-lg shadow-emerald-500/10 animate-pulse-glow'
					: 'bg-slate-900/60 border-slate-800 text-slate-500'}"
			>
				<FileText class="h-6 w-6" />
			</div>
			<span class="text-[10px] font-bold text-slate-200 mt-2.5">PDF Ingestion</span>
			<span class="text-[9px] font-bold text-slate-500 mt-0.5">File Upload</span>
		</div>

		<!-- Link 1 -->
		<div class="flex-1 h-0.5 w-12 md:w-full bg-slate-900 relative">
			<div class="absolute inset-0 bg-gradient-to-r from-emerald-500 to-emerald-400 transition-all duration-550 origin-left
				{status === 'ready' || activeStep() >= 1 ? 'scale-x-100' : 'scale-x-0'}"></div>
			{#if status === 'uploading' || status === 'processing'}
				<div class="absolute top-1/2 -translate-y-1/2 h-1.5 w-1.5 rounded-full bg-emerald-400 animate-slide-right"></div>
			{/if}
		</div>

		<!-- Step 2: Chunking -->
		<div class="flex flex-col items-center z-10">
			<div class="h-14 w-14 rounded-2xl flex items-center justify-center border transition-all duration-500
				{status === 'ready' || activeStep() >= 1
					? 'bg-emerald-950/20 border-emerald-500 text-emerald-400 shadow-lg shadow-emerald-500/10 animate-pulse-glow'
					: 'bg-slate-900/60 border-slate-800 text-slate-500'}"
			>
				<Split class="h-6 w-6" />
			</div>
			<span class="text-[10px] font-bold text-slate-200 mt-2.5">Recursive Splitting</span>
			<span class="text-[9px] font-bold text-slate-500 mt-0.5">500 Token Chunks</span>
		</div>

		<!-- Link 2 -->
		<div class="flex-1 h-0.5 w-12 md:w-full bg-slate-900 relative">
			<div class="absolute inset-0 bg-gradient-to-r from-emerald-500 to-emerald-400 transition-all duration-550 origin-left
				{status === 'ready' || activeStep() >= 2 ? 'scale-x-100' : 'scale-x-0'}"></div>
			{#if status === 'processing'}
				<div class="absolute top-1/2 -translate-y-1/2 h-1.5 w-1.5 rounded-full bg-emerald-400 animate-slide-right"></div>
			{/if}
		</div>

		<!-- Step 3: Embeddings -->
		<div class="flex flex-col items-center z-10">
			<div class="h-14 w-14 rounded-2xl flex items-center justify-center border transition-all duration-500
				{status === 'ready' || activeStep() >= 2
					? 'bg-emerald-950/20 border-emerald-500 text-emerald-400 shadow-lg shadow-emerald-500/10 animate-pulse-glow'
					: 'bg-slate-900/60 border-slate-800 text-slate-500'}"
			>
				<Cpu class="h-6 w-6" />
			</div>
			<span class="text-[10px] font-bold text-slate-200 mt-2.5">Vector Embeddings</span>
			<span class="text-[9px] font-bold text-slate-500 mt-0.5">OpenAI text-3-large</span>
		</div>

		<!-- Link 3 -->
		<div class="flex-1 h-0.5 w-12 md:w-full bg-slate-900 relative">
			<div class="absolute inset-0 bg-gradient-to-r from-emerald-500 to-emerald-400 transition-all duration-550 origin-left
				{status === 'ready' || activeStep() >= 3 ? 'scale-x-100' : 'scale-x-0'}"></div>
			{#if status === 'processing'}
				<div class="absolute top-1/2 -translate-y-1/2 h-1.5 w-1.5 rounded-full bg-emerald-400 animate-slide-right" style="animation-delay: 0.5s"></div>
			{/if}
		</div>

		<!-- Step 4: Index Qdrant -->
		<div class="flex flex-col items-center z-10">
			<div class="h-14 w-14 rounded-2xl flex items-center justify-center border transition-all duration-500
				{status === 'ready' || activeStep() >= 3
					? 'bg-emerald-950/20 border-emerald-500 text-emerald-400 shadow-lg shadow-emerald-500/10 animate-pulse-glow'
					: 'bg-slate-900/60 border-slate-800 text-slate-500'}"
			>
				<Database class="h-6 w-6" />
			</div>
			<span class="text-[10px] font-bold text-slate-200 mt-2.5">Qdrant Vector Database</span>
			<span class="text-[9px] font-bold text-slate-500 mt-0.5">Isolated Collection</span>
		</div>
	</div>
</div>

<style>
	@keyframes slideRight {
		0% { left: 0%; opacity: 0; }
		10% { opacity: 1; }
		90% { opacity: 1; }
		100% { left: 100%; opacity: 0; }
	}
	:global(.animate-slide-right) {
		animation: slideRight 1.5s linear infinite;
	}
	@keyframes pulseGlow {
		0%, 100% { box-shadow: 0 0 10px 0px rgba(16, 185, 129, 0.15); border-color: rgba(16, 185, 129, 0.4); }
		50% { box-shadow: 0 0 20px 4px rgba(16, 185, 129, 0.35); border-color: rgba(16, 185, 129, 0.8); }
	}
	:global(.animate-pulse-glow) {
		animation: pulseGlow 2s ease-in-out infinite;
	}
</style>
