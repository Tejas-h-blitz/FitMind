<script lang="ts">
	import type { Message } from '$lib/api';
	import SourceCitation from './SourceCitation.svelte';
	import { Brain, User } from 'lucide-svelte';

	let { message } = $props<{ message: Message }>();
	const isUser = message.role === 'user';
</script>

<div class="flex gap-3 w-full py-3 {isUser ? 'justify-end' : 'justify-start'}">
	<!-- Avatar for Assistant -->
	{#if !isUser}
		<div class="flex h-9 w-9 shrink-0 select-none items-center justify-center rounded-lg border border-emerald-550/20 bg-slate-900 text-emerald-400 shadow-sm shadow-emerald-950/20">
			<Brain class="h-5 w-5" />
		</div>
	{/if}

	<div class="flex flex-col gap-2.5 max-w-[85%] sm:max-w-[70%]">
		<!-- Message bubble -->
		<div class="rounded-2xl px-4 py-3 text-xs leading-relaxed font-semibold shadow-md transition-all duration-300
			{isUser
				? 'bg-emerald-600 text-white rounded-tr-none border border-emerald-500/20 hover:bg-emerald-550'
				: 'bg-slate-900/40 text-slate-100 rounded-tl-none border border-slate-850 hover:bg-slate-900/60'
			}"
		>
			<div class="whitespace-pre-wrap">{message.content}</div>
		</div>

		<!-- Citations Grid (For Assistant Messages) -->
		{#if !isUser && message.sources && message.sources.length > 0}
			<div class="mt-1 flex flex-col gap-2">
				<p class="text-[9px] uppercase tracking-wider font-extrabold text-slate-500 ml-1">
					References
				</p>
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
					{#each message.sources as source}
						<SourceCitation {source} />
					{/each}
				</div>
			</div>
		{/if}
	</div>

	<!-- Avatar for User -->
	{#if isUser}
		<div class="flex h-9 w-9 shrink-0 select-none items-center justify-center rounded-lg border border-slate-800 bg-slate-900 text-slate-300 shadow-sm shadow-slate-950/25">
			<User class="h-5 w-5" />
		</div>
	{/if}
</div>
