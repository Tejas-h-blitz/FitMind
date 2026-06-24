<script lang="ts">
	import type { Message } from '$lib/api';
	import SourceCitation from './SourceCitation.svelte';
	import { Brain, User } from 'lucide-svelte';

	let { message } = $props<{ message: Message }>();
	const isUser = message.role === 'user';
</script>

<div class="flex gap-3 w-full py-4 {isUser ? 'justify-end' : 'justify-start'}">
	<!-- Avatar for Assistant -->
	{#if !isUser}
		<div class="flex h-9 w-9 shrink-0 select-none items-center justify-center rounded-lg border border-emerald-500/20 bg-slate-900 text-emerald-400">
			<Brain class="h-5 w-5" />
		</div>
	{/if}

	<div class="flex flex-col gap-2 max-w-[85%] sm:max-w-[70%]">
		<!-- Message bubble -->
		<div class="rounded-xl px-4 py-3 text-sm leading-relaxed shadow-md
			{isUser
				? 'bg-emerald-600 text-white rounded-tr-none border border-emerald-500/25'
				: 'bg-slate-900 text-slate-100 rounded-tl-none border border-slate-800'
			}"
		>
			<div class="whitespace-pre-wrap">{message.content}</div>
		</div>

		<!-- Citations Grid (For Assistant Messages) -->
		{#if !isUser && message.sources && message.sources.length > 0}
			<div class="mt-2 flex flex-col gap-2">
				<p class="text-[10px] uppercase tracking-wider font-semibold text-slate-500 ml-1">
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
		<div class="flex h-9 w-9 shrink-0 select-none items-center justify-center rounded-lg border border-slate-700 bg-slate-850 text-slate-300">
			<User class="h-5 w-5" />
		</div>
	{/if}
</div>
