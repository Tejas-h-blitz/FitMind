<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { api, type Document, type Message, type DocumentAnalysis } from '$lib/api';
	import ChatWindow from '$lib/components/ChatWindow.svelte';
	import HealthAnalysis from '$lib/components/HealthAnalysis.svelte';
	import { FileText, ArrowLeft, Brain, MessageSquare, ChevronDown, ChevronUp, Activity, RefreshCw } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	const docId = page.params.docId;

	let document = $state<Document | null>(null);
	let history = $state<Message[]>([]);
	let analysis = $state<DocumentAnalysis | null>(null);
	let isLoading = $state(true);
	let isAnalysisLoading = $state(false);
	let isAnalysisExpanded = $state(false);

	onMount(async () => {
		try {
			const [docsRes, historyRes] = await Promise.all([
				api.listDocuments(),
				api.getChatHistory(docId)
			]);

			if (docsRes.success && docsRes.data) {
				const found = docsRes.data.find((d) => d.id === docId);
				if (found) {
					document = found;
					if (found.status === 'analyzed') {
						isAnalysisLoading = true;
						const analysisRes = await api.getDocumentAnalysis(docId);
						if (analysisRes.success && analysisRes.data) {
							analysis = analysisRes.data;
						} else {
							console.error('Failed to load analysis:', analysisRes.error);
						}
						isAnalysisLoading = false;
					}
				}
			}

			if (historyRes.success && historyRes.data) {
				history = historyRes.data;
			}
		} catch (err) {
			console.error('Error loading chat workspace:', err);
			toast.error('Failed to load chat workspace');
		} finally {
			isLoading = false;
		}
	});
</script>

<svelte:head>
	<title>{document ? `Chat: ${document.name}` : 'Document Advisor'} - FitMind</title>
</svelte:head>

<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-6 relative z-10">
	<!-- Back link & Header details -->
	<div class="flex items-center gap-4 mb-6 pb-4 border-b border-slate-900">
		<a
			href="/dashboard"
			class="inline-flex items-center justify-center p-2.5 rounded-xl bg-slate-900/60 border border-slate-850 text-slate-400 hover:text-slate-200 hover:border-slate-700 hover:scale-105 transition-all duration-200"
			aria-label="Back to dashboard"
		>
			<ArrowLeft class="h-4 w-4" />
		</a>
		
		{#if document}
			<div class="min-w-0">
				<div class="flex items-center gap-2">
					<FileText class="h-4.5 w-4.5 text-emerald-400 shrink-0" />
					<h1 class="text-base sm:text-lg font-black text-slate-100 truncate max-w-md tracking-tight" title={document.name}>
						{document.name}
					</h1>
				</div>
				<p class="text-xs text-slate-405 font-semibold mt-1">
					Discuss and query the details of this document with your FitMind advisor.
				</p>
			</div>
		{/if}
	</div>

	{#if isLoading}
		<div class="h-[60vh] flex flex-col justify-center items-center">
			<div class="h-10 w-10 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
			<span class="text-sm font-semibold text-slate-400 mt-4 animate-pulse">Opening advisor session...</span>
		</div>
	{:else}
		{#if !document}
			<div class="rounded-2xl border border-slate-850 bg-slate-950/45 p-12 text-center max-w-md mx-auto shadow-xl">
				<h2 class="text-lg font-bold text-slate-200">Document Not Found</h2>
				<p class="text-sm text-slate-450 mt-2 font-semibold leading-relaxed">
					The document you are trying to access does not exist or has been deleted.
				</p>
				<a
					href="/dashboard"
					class="mt-6 inline-flex items-center py-2.5 px-4 text-xs font-bold text-white bg-emerald-600 hover:bg-emerald-500 rounded-lg shadow-md hover:-translate-y-0.5 transition-all"
				>
					Return to Dashboard
				</a>
			</div>
		{:else if document.status !== 'ready' && document.status !== 'analyzed'}
			<div class="rounded-2xl border border-slate-850 bg-slate-950/45 p-12 text-center max-w-md mx-auto shadow-xl">
				<div class="h-10 w-10 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
				<h2 class="text-lg font-bold text-slate-200">Analysis in progress...</h2>
				<p class="text-sm text-slate-450 mt-2 font-semibold leading-relaxed">
					This document is currently in state <strong>"{document.status}"</strong>. We are running medical metrics analysis.
				</p>
				<p class="text-xs text-slate-500 mt-4 font-semibold">
					Please wait here or return to the dashboard.
				</p>
				<a
					href="/dashboard"
					class="mt-6 inline-flex items-center py-2.5 px-4 text-xs font-bold text-white bg-emerald-600 hover:bg-emerald-500 rounded-lg shadow-md hover:-translate-y-0.5 transition-all"
				>
					Return to Dashboard
				</a>
			</div>
		{:else}
			<!-- Collapsible Health Analysis Panel -->
			{#if document.status === 'analyzed'}
				<div class="mb-6 rounded-2xl border border-slate-850 bg-slate-950/40 overflow-hidden transition-all shadow-xl backdrop-blur-md">
					<button
						onclick={() => isAnalysisExpanded = !isAnalysisExpanded}
						class="w-full flex items-center justify-between px-5 py-4.5 text-slate-200 hover:bg-slate-900/30 transition-colors font-bold text-sm cursor-pointer"
					>
						<div class="flex items-center gap-2">
							<Activity class="h-4.5 w-4.5 text-emerald-400" />
							<span>View Health Analysis Report</span>
						</div>
						{#if isAnalysisExpanded}
							<ChevronUp class="h-4.5 w-4.5 text-slate-450" />
						{:else}
							<ChevronDown class="h-4.5 w-4.5 text-slate-450" />
						{/if}
					</button>
					{#if isAnalysisExpanded}
						<div class="border-t border-slate-900 p-5 bg-slate-950/20">
							{#if isAnalysisLoading}
								<div class="py-12 flex flex-col items-center justify-center">
									<RefreshCw class="h-8 w-8 text-emerald-500 animate-spin" />
									<span class="text-xs text-slate-400 mt-3 font-semibold animate-pulse">Fetching analysis results...</span>
								</div>
							{:else if analysis}
								<HealthAnalysis {analysis} />
							{:else}
								<div class="py-4 text-center text-sm text-slate-500 italic">
									Analysis is in progress or unavailable for this file.
								</div>
							{/if}
						</div>
					{/if}
				</div>
			{/if}

			<!-- Mount ChatWindow -->
			<ChatWindow {docId} initialMessages={history} />
		{/if}
	{/if}
</div>
