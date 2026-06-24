<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { api, type Document, type Message } from '$lib/api';
	import ChatWindow from '$lib/components/ChatWindow.svelte';
	import { FileText, ArrowLeft, Brain, MessageSquare } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	const docId = page.params.docId;

	let document = $state<Document | null>(null);
	let history = $state<Message[]>([]);
	let isLoading = $state(true);

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

<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<!-- Back link & Header details -->
	<div class="flex items-center gap-4 mb-6">
		<a
			href="/dashboard"
			class="inline-flex items-center justify-center p-2 rounded-lg bg-slate-900 border border-slate-800 text-slate-400 hover:text-white hover:border-slate-700 transition-all"
			aria-label="Back to dashboard"
		>
			<ArrowLeft class="h-4 w-4" />
		</a>
		
		{#if document}
			<div class="min-w-0">
				<div class="flex items-center gap-2">
					<FileText class="h-4.5 w-4.5 text-emerald-400 shrink-0" />
					<h1 class="text-base sm:text-lg font-bold text-slate-100 truncate max-w-md" title={document.name}>
						{document.name}
					</h1>
				</div>
				<p class="text-xs text-slate-400 mt-0.5">
					Discuss and query the details of this document with your FitMind advisor.
				</p>
			</div>
		{/if}
	</div>

	{#if isLoading}
		<div class="h-96 flex flex-col justify-center items-center">
			<div class="h-10 w-10 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
			<span class="text-sm font-semibold text-slate-400 mt-4">Opening advisor session...</span>
		</div>
	{:else}
		{#if !document}
			<div class="rounded-2xl border border-slate-900 bg-slate-900/10 p-12 text-center max-w-md mx-auto">
				<h2 class="text-lg font-bold text-slate-200">Document Not Found</h2>
				<p class="text-sm text-slate-400 mt-2">
					The document you are trying to access does not exist or has been deleted.
				</p>
				<a
					href="/dashboard"
					class="mt-6 inline-flex items-center py-2.5 px-4 text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-500 rounded-lg shadow-md"
				>
					Return to Dashboard
				</a>
			</div>
		{:else if document.status !== 'ready'}
			<div class="rounded-2xl border border-slate-900 bg-slate-900/10 p-12 text-center max-w-md mx-auto">
				<h2 class="text-lg font-bold text-slate-200">Analysis In Progress</h2>
				<p class="text-sm text-slate-400 mt-2">
					This document is currently in state <strong>"{document.status}"</strong> and cannot be queried yet.
				</p>
				<p class="text-xs text-slate-500 mt-4">
					Please return to the dashboard and wait for processing to finish.
				</p>
				<a
					href="/dashboard"
					class="mt-6 inline-flex items-center py-2.5 px-4 text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-500 rounded-lg shadow-md"
				>
					Return to Dashboard
				</a>
			</div>
		{:else}
			<!-- Mount ChatWindow -->
			<ChatWindow {docId} initialMessages={history} />
		{/if}
	{/if}
</div>
