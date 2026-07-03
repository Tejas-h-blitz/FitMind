<script lang="ts">
	import type { Document } from '$lib/api';
	import { FileText, MessageSquare, Trash2, RefreshCw, AlertCircle, CheckCircle } from 'lucide-svelte';

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

<div class="group relative rounded-xl border border-slate-850 bg-slate-950/20 p-5 hover:border-emerald-500/35 hover:bg-slate-900/25 hover:-translate-y-0.5 hover:shadow-xl hover:shadow-emerald-950/5 transition-all duration-300 flex flex-col justify-between h-[160px]">
	<div class="flex items-start justify-between gap-3">
		<div class="flex items-center gap-3 min-w-0">
			<div class="p-2.5 rounded-lg bg-emerald-500/10 text-emerald-400 group-hover:bg-emerald-500/20 transition-all duration-300 shrink-0">
				<FileText class="h-6 w-6" />
			</div>
			<div class="min-w-0">
				<h3 class="font-bold text-sm text-slate-200 truncate" title={document.name}>
					{document.name}
				</h3>
				<p class="text-[10px] text-slate-500 font-bold mt-1 tracking-tight">
					{formatBytes(document.size)} • {formatDate(document.created_at)}
				</p>
			</div>
		</div>

		<!-- Status Badge -->
		<div class="shrink-0">
			{#if document.status === 'analyzed'}
				<span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[9px] font-extrabold bg-emerald-950/45 border border-emerald-500/20 text-emerald-400 uppercase tracking-wider">
					<CheckCircle class="h-3 w-3" />
					Analyzed
				</span>
			{:else if document.status === 'ready'}
				<span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[9px] font-extrabold bg-emerald-950/45 border border-emerald-500/20 text-emerald-400 uppercase tracking-wider">
					<CheckCircle class="h-3 w-3" />
					Ready
				</span>
			{:else if document.status === 'processing'}
				<span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[9px] font-extrabold bg-blue-950/45 border border-blue-500/20 text-blue-400 uppercase tracking-wider animate-pulse">
					<RefreshCw class="h-3 w-3 animate-spin" />
					Processing
				</span>
			{:else if document.status === 'pending'}
				<span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[9px] font-extrabold bg-amber-950/45 border border-amber-500/20 text-amber-400 uppercase tracking-wider animate-pulse">
					<RefreshCw class="h-3 w-3" />
					Pending
				</span>
			{:else}
				<span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[9px] font-extrabold bg-red-950/45 border border-red-500/20 text-red-400 uppercase tracking-wider">
					<AlertCircle class="h-3 w-3" />
					Failed
				</span>
			{/if}
		</div>
	</div>

	<!-- Actions Bar -->
	<div class="flex justify-end gap-3 border-t border-slate-900/60 pt-3 mt-4 shrink-0">
		<!-- Delete Button -->
		<button
			disabled={isDeleting}
			onclick={handleDeleteClick}
			class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold text-slate-400 hover:text-rose-450 rounded-lg hover:bg-rose-950/10 border border-transparent hover:border-rose-500/15 transition-all cursor-pointer disabled:opacity-50"
			aria-label="Delete document"
		>
			<Trash2 class="h-3.5 w-3.5" />
			<span>Delete</span>
		</button>

		<!-- Chat Action -->
		{#if document.status === 'ready' || document.status === 'analyzed'}
			<a
				href="/chat/{document.id}"
				class="flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-bold text-white bg-emerald-600 hover:bg-emerald-500 border border-emerald-500/20 hover:border-emerald-400/20 rounded-lg transition-all cursor-pointer shadow-md shadow-emerald-950/30 hover:-translate-y-0.5"
			>
				<MessageSquare class="h-3.5 w-3.5" />
				<span>Chat Advisor</span>
			</a>
		{:else}
			<button
				disabled
				class="flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-bold text-slate-500 bg-slate-900 border border-slate-800 rounded-lg cursor-not-allowed"
			>
				<MessageSquare class="h-3.5 w-3.5" />
				<span>Processing...</span>
			</button>
		{/if}
	</div>
</div>
