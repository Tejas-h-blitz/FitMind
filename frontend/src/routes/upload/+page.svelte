<script lang="ts">
	import UploadZone from '$lib/components/UploadZone.svelte';
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { Brain, FilePlus2, Sparkles, ShieldCheck } from 'lucide-svelte';
</script>

<svelte:head>
	<title>Upload Health Documents - FitMind</title>
</svelte:head>

<div class="max-w-4xl mx-auto px-6 py-12">
	<div class="text-center mb-10">
		<h1 class="text-2xl sm:text-3xl font-extrabold text-slate-100 flex items-center justify-center gap-2">
			<FilePlus2 class="text-emerald-400" />
			<span>Add Health Intelligence Data</span>
		</h1>
		<p class="text-slate-400 text-sm mt-2 max-w-lg mx-auto">
			Upload workout schedules, diet menus, or blood reports to index them in the AI engine.
		</p>
	</div>

	<!-- Upload Zone Card -->
	<div class="bg-slate-900/10 border border-slate-900 p-8 rounded-2xl shadow-xl">
		<UploadZone
			onUpload={async (file) => {
				const res = await api.uploadDocument(file);
				if (res.success && res.data) {
					toast.success(`Successfully uploaded "${file.name}"! Starting RAG background indexing...`);
					goto('/dashboard');
				} else {
					toast.error(res.error || 'Failed to upload document');
					throw new Error(res.error);
				}
			}}
		/>
	</div>

	<!-- Extra info / Help cards -->
	<div class="grid grid-cols-1 sm:grid-cols-2 gap-6 mt-10">
		<div class="flex gap-3 p-4 rounded-xl border border-slate-900 bg-slate-900/10">
			<div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400 shrink-0 h-fit">
				<Sparkles class="h-4 w-4" />
			</div>
			<div>
				<h3 class="text-xs font-bold text-slate-200">How indexing works</h3>
				<p class="text-[11px] text-slate-400 mt-1 leading-relaxed">
					Once uploaded, your PDF is downloaded and split into small chunks. OpenAI embeds these chunks using text-embedding models, and registers them in a dedicated Qdrant database vector collection.
				</p>
			</div>
		</div>

		<div class="flex gap-3 p-4 rounded-xl border border-slate-900 bg-slate-900/10">
			<div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400 shrink-0 h-fit">
				<ShieldCheck class="h-4 w-4" />
			</div>
			<div>
				<h3 class="text-xs font-bold text-slate-200">Your privacy is isolated</h3>
				<p class="text-[11px] text-slate-400 mt-1 leading-relaxed">
					Row-level security policies ensure only you have access to your database records. Vector stores are separated by user ID so that search queries never overlap across profiles.
				</p>
			</div>
		</div>
	</div>
</div>
