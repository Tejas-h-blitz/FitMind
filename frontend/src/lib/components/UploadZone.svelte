<script lang="ts">
	import { UploadCloud, FileText, X, AlertCircle } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let { onUpload } = $props<{
		onUpload: (file: File) => Promise<void>;
	}>();

	let isDragging = $state(false);
	let fileInput = $state<HTMLInputElement | null>(null);
	let selectedFile = $state<File | null>(null);
	let isUploading = $state(false);

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		isDragging = true;
	}

	function handleDragLeave() {
		isDragging = false;
	}

	function validateAndSetFile(file: File) {
		if (!file.name.toLowerCase().endsWith('.pdf')) {
			toast.error('Only PDF documents are supported');
			return;
		}

		if (file.size > 15 * 1024 * 1024) {
			toast.error('File size cannot exceed 15MB');
			return;
		}

		selectedFile = file;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
		if (e.dataTransfer?.files && e.dataTransfer.files.length > 0) {
			validateAndSetFile(e.dataTransfer.files[0]);
		}
	}

	function handleFileChange(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files && target.files.length > 0) {
			validateAndSetFile(target.files[0]);
		}
	}

	function triggerFileInput() {
		fileInput?.click();
	}

	function clearFile() {
		selectedFile = null;
		if (fileInput) fileInput.value = '';
	}

	async function handleSubmit() {
		if (!selectedFile || isUploading) return;
		isUploading = true;

		try {
			await onUpload(selectedFile);
			selectedFile = null;
			if (fileInput) fileInput.value = '';
		} catch (err: any) {
			console.error(err);
		} finally {
			isUploading = false;
		}
	}
</script>

<div class="w-full max-w-2xl mx-auto">
	<div
		role="button"
		tabindex="0"
		onkeydown={(e) => e.key === 'Enter' && triggerFileInput()}
		onclick={selectedFile ? undefined : triggerFileInput}
		ondragover={handleDragOver}
		ondragleave={handleDragLeave}
		ondrop={handleDrop}
		class="relative rounded-2xl border-2 border-dashed p-10 flex flex-col items-center justify-center text-center cursor-pointer transition-all duration-300
			{isDragging
				? 'border-emerald-400 bg-emerald-950/20 shadow-lg shadow-emerald-950/5'
				: selectedFile
					? 'border-slate-850 bg-slate-900/10'
					: 'border-slate-850 bg-slate-950/40 hover:border-emerald-500/30 hover:bg-slate-900/15'
			}"
	>
		<input
			type="file"
			bind:this={fileInput}
			onchange={handleFileChange}
			accept=".pdf"
			class="hidden"
		/>

		{#if !selectedFile}
			<div class="p-4.5 rounded-full bg-slate-950 border border-slate-900 text-slate-400 mb-4 transition-transform duration-300 hover:scale-105">
				<UploadCloud class="h-8 w-8 text-emerald-400" />
			</div>
			<h3 class="font-extrabold text-slate-200 text-base tracking-tight">Drag & drop your health document</h3>
			<p class="text-sm text-slate-500 mt-2 max-w-sm font-semibold">
				Or <span class="text-emerald-400 font-bold">browse files</span> to upload reports, diets, or workouts.
			</p>
			<p class="text-[10px] text-slate-600 font-extrabold uppercase tracking-wider mt-6">
				PDF documents up to 15MB
			</p>
		{:else}
			<div class="flex flex-col items-center w-full">
				<div class="p-3.5 rounded-xl bg-emerald-500/10 text-emerald-400 mb-4 border border-emerald-500/20">
					<FileText class="h-8 w-8" />
				</div>
				<h4 class="font-black text-sm text-slate-100 break-all px-6 leading-relaxed">
					{selectedFile.name}
				</h4>
				<p class="text-[11px] text-slate-500 font-bold mt-1.5">
					{(selectedFile.size / (1024 * 1024)).toFixed(2)} MB
				</p>

				<!-- Buttons -->
				<div class="flex items-center gap-4 mt-6 w-full max-w-xs">
					<button
						onclick={(e) => { e.stopPropagation(); clearFile(); }}
						disabled={isUploading}
						class="flex-1 py-2.5 px-4 border border-slate-800 hover:border-slate-700 text-slate-400 hover:text-slate-200 rounded-xl text-xs font-bold transition-all cursor-pointer disabled:opacity-50"
					>
						Cancel
					</button>
					<button
						onclick={(e) => { e.stopPropagation(); handleSubmit(); }}
						disabled={isUploading}
						class="flex-1 py-2.5 px-4 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl text-xs font-bold transition-all cursor-pointer disabled:opacity-50 flex justify-center items-center gap-2"
					>
						{#if isUploading}
							<div class="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
							<span>Uploading...</span>
						{:else}
							<span>Analyze File</span>
						{/if}
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>
