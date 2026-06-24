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
				? 'border-emerald-400 bg-emerald-950/20 shadow-lg shadow-emerald-900/10'
				: selectedFile
					? 'border-slate-800 bg-slate-900/30'
					: 'border-slate-800 bg-slate-950 hover:border-emerald-500/40 hover:bg-slate-900/20'
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
			<div class="p-4 rounded-full bg-slate-900 border border-slate-850 text-slate-400 mb-4 transition-transform group-hover:scale-110">
				<UploadCloud class="h-8 w-8 text-emerald-400" />
			</div>
			<h3 class="font-semibold text-slate-200 text-base">Drag & drop your health document</h3>
			<p class="text-sm text-slate-500 mt-1 max-w-sm">
				Or <span class="text-emerald-400 font-medium">browse your files</span> to upload blood reports, diet plans, or workout logs.
			</p>
			<p class="text-xs text-slate-600 mt-6">
				Supports PDF files up to 15MB
			</p>
		{:else}
			<div class="flex flex-col items-center w-full">
				<div class="p-3.5 rounded-xl bg-emerald-500/10 text-emerald-400 mb-4 border border-emerald-500/20">
					<FileText class="h-8 w-8" />
				</div>
				<h4 class="font-bold text-slate-100 break-all px-6">
					{selectedFile.name}
				</h4>
				<p class="text-xs text-slate-400 mt-1">
					{(selectedFile.size / (1024 * 1024)).toFixed(2)} MB
				</p>

				<!-- Buttons -->
				<div class="flex items-center gap-4 mt-6 w-full max-w-xs">
					<button
						onclick={(e) => { e.stopPropagation(); clearFile(); }}
						disabled={isUploading}
						class="flex-1 py-2 px-4 border border-slate-700 hover:border-slate-600 text-slate-300 hover:text-white rounded-lg text-sm font-semibold transition-all cursor-pointer disabled:opacity-50"
					>
						Cancel
					</button>
					<button
						onclick={(e) => { e.stopPropagation(); handleSubmit(); }}
						disabled={isUploading}
						class="flex-1 py-2 px-4 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg text-sm font-semibold transition-all cursor-pointer disabled:opacity-50 flex justify-center items-center gap-2"
					>
						{#if isUploading}
							<div class="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
							<span>Uploading</span>
						{:else}
							<span>Analyze File</span>
						{/if}
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>
