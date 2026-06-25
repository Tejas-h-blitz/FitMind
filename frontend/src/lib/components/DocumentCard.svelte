<script lang="ts">
	import type { Document } from '$lib/api';
	import { FileText, MessageSquare, Trash2, RefreshCw, AlertCircle, CheckCircle } from 'lucide-svelte';
	import { createEventDispatcher } from 'svelte';

	let { document, onDelete } = $props<{
		document: Document;
		onDelete: (id: string) => Promise<void>;
	}>();

	let isDeleting = $state(false);

	function formatBytes(bytes: number, decimals = 1) {
		if (bytes === 0) return '0 Bytes';
		const k = 1024;
		const dm = decimals < 0 ? 0 : decimals;
		const sizes = ['Bytes', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
	}

	function formatDate(dateStr: string) {
		const date = new Date(dateStr);
		return date.toLocaleDateString(undefined, {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	async function handleDeleteClick() {
		if (confirm(`Are you sure you want to delete "${document.name}"? This will permanently delete all indexed chunks and chat history.`)) {
			isDeleting = true;
			try {
				await onDelete(document.id);
			} catch (err) {
				isDeleting = false;
			}
		}
	}
</script>

<div class="group relative rounded-xl border border-slate-800 bg-slate-900/50 p-5 hover:border-emerald-500/40 hover:bg-slate-900/80 transition-all duration-300">
	<div class="flex items-start justify-between">
		<div class="flex items-center gap-3">
			<div class="p-2.5 rounded-lg bg-emerald-500/10 text-emerald-400 group-hover:bg-emerald-500/20 transition-all duration-300">
				<FileText class="h-6 w-6" />
			</div>
			<div>
				<h3 class="font-semibold text-slate-100 line-clamp-1 break-all" title={document.name}>
					{document.name}
				</h3>
				<p class="text-xs text-slate-400 mt-0.5">
					{formatBytes(document.size)} • {formatDate(document.created_at)}
				</p>
			</div>
		</div>

		<!-- Status Badge -->
		<div>
			{#if document.status === 'analyzed'}
				<span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-950/40 border border-emerald-500/30 text-emerald-400">
					<CheckCircle class="h-3 w-3" />
					Analyzed
				</span>
			{:else if document.status === 'ready'}
				<span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-950/40 border border-emerald-500/30 text-emerald-400">
					<CheckCircle class="h-3 w-3" />
					Ready
				</span>
			{:else if document.status === 'processing'}
				<span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-blue-950/40 border border-blue-500/30 text-blue-400">
					<RefreshCw class="h-3 w-3 animate-spin" />
					Processing
				</span>
			{:else if document.status === 'pending'}
				<span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-amber-950/40 border border-amber-500/30 text-amber-400 animate-pulse">
					<RefreshCw class="h-3 w-3" />
					Pending
				</span>
			{:else}
				<span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-red-950/40 border border-red-500/30 text-red-400">
					<AlertCircle class="h-3 w-3" />
					Failed
				</span>
			{/if}
		</div>
	</div>

	<!-- Actions Bar -->
	<div class="mt-6 flex justify-end gap-3 border-t border-slate-800/60 pt-4">
		<!-- Delete Button -->
		<button
			disabled={isDeleting}
			onclick={handleDeleteClick}
			class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-slate-400 hover:text-red-400 rounded-md hover:bg-red-950/10 border border-transparent hover:border-red-500/15 transition-all cursor-pointer disabled:opacity-50"
			aria-label="Delete document"
		>
			<Trash2 class="h-3.5 w-3.5" />
			<span>Delete</span>
		</button>

		<!-- Chat Action -->
		{#if document.status === 'ready' || document.status === 'analyzed'}
			<a
				href="/chat/{document.id}"
				class="flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-semibold text-white bg-emerald-600 hover:bg-emerald-500 border border-emerald-500/20 hover:border-emerald-400/20 rounded-md transition-all cursor-pointer shadow-sm shadow-emerald-950/50"
			>
				<MessageSquare class="h-3.5 w-3.5" />
				<span>Chat Advisor</span>
			</a>
		{:else}
			<button
				disabled
				class="flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-semibold text-slate-500 bg-slate-800/40 border border-slate-800/20 rounded-md cursor-not-allowed"
			>
				<MessageSquare class="h-3.5 w-3.5" />
				<span>Processing...</span>
			</button>
		{/if}
	</div>
</div>
