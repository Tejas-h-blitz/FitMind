<script lang="ts">
	import UploadZone from '$lib/components/UploadZone.svelte';
	import PipelineVisualizer from '$lib/components/PipelineVisualizer.svelte';
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { Brain, FilePlus2, Sparkles, ShieldCheck } from 'lucide-svelte';

	let uploadStatus = $state<'idle' | 'uploading' | 'processing' | 'ready'>('idle');
	let uploadedDocName = $state('');

	async function handleUpload(file: File) {
		uploadStatus = 'uploading';
		uploadedDocName = file.name;

		try {
			const res = await api.uploadDocument(file);
			if (res.success && res.data) {
				const docId = res.data.id;
				uploadStatus = 'processing';
				toast.info('Document uploaded. Commencing chunking & vector indexing...');

				// Start polling until ready
				await pollDocumentStatus(docId);
			} else {
				uploadStatus = 'idle';
				toast.error(res.error || 'Failed to upload document');
				throw new Error(res.error);
			}
		} catch (err: any) {
			uploadStatus = 'idle';
			console.error(err);
		}
	}

	async function pollDocumentStatus(docId: string) {
		const pollInterval = setInterval(async () => {
			try {
				const res = await api.listDocuments();
				if (res.success && res.data) {
					const doc = res.data.find(d => d.id === docId);
					if (doc) {
						if (doc.status === 'ready' || doc.status === 'analyzed') {
							clearInterval(pollInterval);
							uploadStatus = 'ready';
							toast.success(`"${uploadedDocName}" successfully indexed in Qdrant!`);
							
							// Brief delay so they can appreciate the pipeline complete visual
							setTimeout(() => {
								goto('/dashboard');
							}, 1500);
						} else if (doc.status === 'failed') {
							clearInterval(pollInterval);
							uploadStatus = 'idle';
							toast.error('RAG vector indexing failed for this document.');
						}
					}
				}
			} catch (err) {
				console.error(err);
			}
		}, 2000);
	}
</script>

<svelte:head>
	<title>Upload Health Documents - FitMind</title>
</svelte:head>

<div class="max-w-4xl mx-auto px-6 py-12 relative z-10">
	<div class="text-center mb-12">
		<h1 class="text-3xl font-black text-slate-100 tracking-tight flex items-center justify-center gap-3">
			<span class="p-2 rounded-xl bg-emerald-500/10 text-emerald-400">
				<FilePlus2 class="h-6 w-6" />
			</span>
			<span>Add Health Intelligence Data</span>
		</h1>
		<p class="text-slate-400 text-sm mt-3 max-w-md mx-auto">
			Upload workout schedules, diet menus, or blood reports to index them in the AI engine.
		</p>
	</div>

	<div class="space-y-8">
		<!-- Upload Zone Card -->
		<div class="bg-slate-955/45 border border-slate-850 p-8 rounded-2xl shadow-xl backdrop-blur-md">
			{#if uploadStatus === 'idle'}
				<UploadZone onUpload={handleUpload} />
			{:else}
				<div class="flex flex-col items-center justify-center py-10 text-center">
					<div class="h-10 w-10 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin mb-4"></div>
					<h3 class="font-extrabold text-slate-200 text-base capitalize">
						{#if uploadStatus === 'uploading'}
							Uploading "{uploadedDocName}"...
						{:else if uploadStatus === 'processing'}
							Indexing RAG Knowledge Base...
						{:else}
							Ingestion Complete!
						{/if}
					</h3>
					<p class="text-xs text-slate-500 mt-2 max-w-xs font-semibold leading-relaxed">
						{#if uploadStatus === 'uploading'}
							Sending file bytes to FitMind secure data server.
						{:else if uploadStatus === 'processing'}
							Running OpenAI embeddings & saving vector indices to Qdrant.
						{:else}
							Redirecting to health workspace...
						{/if}
					</p>
				</div>
			{/if}
		</div>

		<!-- Live Pipeline Visualizer -->
		{#if uploadStatus !== 'idle'}
			<PipelineVisualizer status={uploadStatus} />
		{/if}
	</div>

	<!-- Extra info / Help cards -->
	{#if uploadStatus === 'idle'}
		<div class="grid grid-cols-1 sm:grid-cols-2 gap-6 mt-10">
			<div class="flex gap-3.5 p-5 rounded-2xl border border-slate-850 bg-slate-900/15 backdrop-blur-sm">
				<div class="p-2 bg-emerald-500/10 rounded-xl text-emerald-400 shrink-0 h-fit">
					<Sparkles class="h-4.5 w-4.5" />
				</div>
				<div>
					<h3 class="text-xs font-extrabold text-slate-200 uppercase tracking-wider">How indexing works</h3>
					<p class="text-[11px] text-slate-450 mt-1.5 leading-relaxed font-semibold">
						Once uploaded, your PDF is downloaded and split into small chunks. OpenAI embeds these chunks using text-embedding models, and registers them in a dedicated Qdrant database vector collection.
					</p>
				</div>
			</div>

			<div class="flex gap-3.5 p-5 rounded-2xl border border-slate-850 bg-slate-900/15 backdrop-blur-sm">
				<div class="p-2 bg-emerald-500/10 rounded-xl text-emerald-400 shrink-0 h-fit">
					<ShieldCheck class="h-4.5 w-4.5" />
				</div>
				<div>
					<h3 class="text-xs font-extrabold text-slate-200 uppercase tracking-wider">Your privacy is isolated</h3>
					<p class="text-[11px] text-slate-455 mt-1.5 leading-relaxed font-semibold">
						Row-level security policies ensure only you have access to your database records. Vector stores are separated by user ID so that search queries never overlap across profiles.
					</p>
				</div>
			</div>
		</div>
	{/if}
</div>
