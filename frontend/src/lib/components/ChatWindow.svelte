<script lang="ts">
	import { tick } from 'svelte';
	import type { Message, SourceChunk } from '$lib/api';
	import { api } from '$lib/api';
	import { supabase } from '$lib/supabase';
	import { PUBLIC_API_BASE_URL } from '$env/static/public';
	import MessageBubble from './MessageBubble.svelte';
	import { Send, Trash2, ArrowDown, Brain, AlertCircle } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let { docId, initialMessages = [] } = $props<{
		docId: string;
		initialMessages?: Message[];
	}>();

	let messages = $state<Message[]>([...initialMessages]);
	let queryText = $state('');
	let isGenerating = $state(false);
	let chatContainer = $state<HTMLDivElement | null>(null);

	// Temp message for showing real-time streaming tokens
	let streamingMessage = $state<Message | null>(null);

	async function scrollToBottom() {
		await tick();
		if (chatContainer) {
			chatContainer.scrollTop = chatContainer.scrollHeight;
		}
	}

	$effect(() => {
		// Scroll to bottom on initial load and when message list changes
		if (messages.length > 0 || streamingMessage) {
			scrollToBottom();
		}
	});

	async function handleSend(e: Event) {
		e.preventDefault();
		if (!queryText.trim() || isGenerating) return;

		const currentQuery = queryText;
		queryText = '';
		isGenerating = true;

		// 1. Add user message locally for immediate UI update
		const userMsg: Message = {
			id: 'temp-user-id-' + Date.now(),
			doc_id: docId,
			user_id: '',
			role: 'user',
			content: currentQuery,
			sources: [],
			created_at: new Date().toISOString()
		};
		messages = [...messages, userMsg];
		await scrollToBottom();

		// 2. Initialize the empty assistant streaming message
		streamingMessage = {
			id: 'temp-assistant-id-' + Date.now(),
			doc_id: docId,
			user_id: '',
			role: 'assistant',
			content: '',
			sources: [],
			created_at: new Date().toISOString()
		};

		try {
			const { data: sessionData } = await supabase.auth.getSession();
			const token = sessionData.session?.access_token || '';

			// 3. Make the query request to Go backend
			const url = `${PUBLIC_API_BASE_URL}/api/chat/${docId}/query`;
			const response = await fetch(url, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${token}`
				},
				body: JSON.stringify({ query: currentQuery })
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ error: 'Query execution failed' }));
				throw new Error(errorData.error || 'Server error occurred');
			}

			const reader = response.body?.getReader();
			const decoder = new TextDecoder('utf-8');
			if (!reader) {
				throw new Error('Response stream not readable');
			}

			// 4. Stream parser
			let buffer = '';
			while (true) {
				const { done, value } = await reader.read();
				if (done) break;

				buffer += decoder.decode(value, { stream: true });
				const events = buffer.split('\n\n');
				buffer = events.pop() || ''; // Keep partial events in buffer

				for (const event of events) {
					if (!event.trim()) continue;

					// Extract event details
					const lines = event.split('\n');
					let eventType = '';
					let dataStr = '';

					for (const line of lines) {
						if (line.startsWith('event: ')) {
							eventType = line.replace('event: ', '').trim();
						} else if (line.startsWith('data: ')) {
							dataStr = line.replace('data: ', '').trim();
						}
					}

					if (eventType === 'sources' && dataStr) {
						try {
							const sources = JSON.parse(dataStr) as SourceChunk[];
							if (streamingMessage) {
								streamingMessage.sources = sources;
							}
						} catch (e) {
							console.error('Failed to parse sources', e);
						}
					} else if (eventType === 'text' && dataStr) {
						try {
							const token = JSON.parse(dataStr) as string;
							if (streamingMessage) {
								streamingMessage.content += token;
							}
						} catch (e) {
							console.error('Failed to parse token', e);
						}
					} else if (eventType === 'error' && dataStr) {
						try {
							const errPayload = JSON.parse(dataStr);
							throw new Error(errPayload.message || 'Error from RAG service');
						} catch (e: any) {
							throw new Error(e.message || 'Stream generation error');
						}
					}
				}
			}

			// Add finalized assistant message to lists
			if (streamingMessage) {
				// Fetch actual messages to get backend-generated IDs
				const historyRes = await api.getChatHistory(docId);
				if (historyRes.success && historyRes.data) {
					messages = historyRes.data;
				} else {
					messages = [...messages, streamingMessage];
				}
			}

		} catch (error: any) {
			console.error('Streaming error:', error);
			toast.error(error.message || 'Failed to get answer from AI advisor');
			// Remove the user message since the turn failed
			messages = messages.filter((m) => m.id !== userMsg.id);
		} finally {
			streamingMessage = null;
			isGenerating = false;
			await scrollToBottom();
		}
	}

	async function handleClear() {
		if (confirm('Are you sure you want to clear the entire chat history for this document?')) {
			const res = await api.clearChatHistory(docId);
			if (res.success) {
				messages = [];
				toast.success('Chat history cleared successfully');
			} else {
				toast.error(res.error || 'Failed to clear chat history');
			}
		}
	}
</script>

<div class="flex flex-col h-[calc(100vh-12rem)] min-h-[450px] rounded-xl border border-slate-800 bg-slate-950 overflow-hidden shadow-2xl">
	<!-- Chat Header -->
	<div class="flex justify-between items-center px-6 py-4 border-b border-slate-800 bg-slate-900/40">
		<div class="flex items-center gap-3">
			<div class="h-2 w-2 rounded-full bg-emerald-500 animate-ping"></div>
			<div>
				<h2 class="text-sm font-semibold text-slate-100">AI Medical Document Analyst</h2>
				<p class="text-xs text-slate-400">Context-aware RAG advisor</p>
			</div>
		</div>
		
		{#if messages.length > 0}
			<button
				onclick={handleClear}
				class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-slate-400 hover:text-red-400 rounded-lg hover:bg-red-950/20 border border-transparent hover:border-red-500/25 transition-all cursor-pointer"
			>
				<Trash2 class="h-3.5 w-3.5" />
				<span>Clear Chat</span>
			</button>
		{/if}
	</div>

	<!-- Messages Scroll View -->
	<div
		bind:this={chatContainer}
		class="flex-1 overflow-y-auto px-6 py-4 space-y-4 scrollbar-thin scrollbar-thumb-slate-800"
	>
		{#if messages.length === 0 && !streamingMessage}
			<div class="h-full flex flex-col justify-center items-center text-center p-8 max-w-md mx-auto">
				<div class="p-4 rounded-full bg-emerald-500/10 text-emerald-400 mb-4 animate-bounce">
					<Brain class="h-10 w-10" />
				</div>
				<h3 class="font-bold text-slate-100 text-lg">Start Health Analysis</h3>
				<p class="text-slate-400 text-sm mt-2 leading-relaxed">
					Ask any questions about this document. For example: <br>
					<span class="text-emerald-400 italic">"Are my cholesterol levels within range?"</span> or <br>
					<span class="text-emerald-400 italic">"What is the recommended protein intake in this diet?"</span>
				</p>
			</div>
		{:else}
			{#each messages as message}
				<MessageBubble {message} />
			{/each}

			{#if streamingMessage}
				<MessageBubble message={streamingMessage} />
			{/if}

			{#if isGenerating && (!streamingMessage || !streamingMessage.content)}
				<div class="flex gap-3 py-4 justify-start">
					<div class="flex h-9 w-9 items-center justify-center rounded-lg border border-emerald-500/20 bg-slate-900 text-emerald-400">
						<Brain class="h-5 w-5 animate-pulse" />
					</div>
					<div class="flex flex-col gap-2 max-w-[70%]">
						<div class="rounded-xl px-4 py-3 bg-slate-900 border border-slate-800 rounded-tl-none">
							<div class="flex items-center gap-1.5 py-1">
								<span class="h-2 w-2 bg-slate-500 rounded-full animate-bounce" style="animation-delay: 0ms"></span>
								<span class="h-2 w-2 bg-slate-500 rounded-full animate-bounce" style="animation-delay: 150ms"></span>
								<span class="h-2 w-2 bg-slate-500 rounded-full animate-bounce" style="animation-delay: 300ms"></span>
							</div>
						</div>
					</div>
				</div>
			{/if}
		{/if}
	</div>

	<!-- Input Area -->
	<div class="p-4 border-t border-slate-800 bg-slate-900/20">
		<form onsubmit={handleSend} class="flex gap-3">
			<input
				type="text"
				bind:value={queryText}
				placeholder={isGenerating ? "AI is typing..." : "Ask FitMind advisor..."}
				disabled={isGenerating}
				class="flex-1 px-4 py-3 rounded-lg border border-slate-800 bg-slate-950 text-slate-100 placeholder-slate-500 text-sm focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all disabled:opacity-50"
			/>
			<button
				type="submit"
				disabled={isGenerating || !queryText.trim()}
				class="flex items-center justify-center h-11 w-11 shrink-0 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
				aria-label="Send query"
			>
				<Send class="h-4 w-4" />
			</button>
		</form>
	</div>
</div>
