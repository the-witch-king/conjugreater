<script lang="ts">
	import { ALL_FORMS } from '$lib/types';
	import type { ConjugationForm } from '$lib/types';

	interface Props {
		enabledForms: Set<ConjugationForm>;
		alwaysShowType: boolean;
		onToggle: (form: ConjugationForm) => void;
		onToggleShowType: () => void;
		onStart: () => void;
		onOpenSettings: () => void;
	}

	let { enabledForms, alwaysShowType, onToggle, onToggleShowType, onStart, onOpenSettings }: Props = $props();
</script>

<div class="max-w-md mx-auto">
	<h2 class="text-xl font-semibold mb-4 text-gray-800">Select conjugation forms</h2>

	<div class="space-y-2 mb-6">
		{#each ALL_FORMS as form}
			<label class="flex items-center gap-3 p-2 rounded hover:bg-gray-50 cursor-pointer">
				<input
					type="checkbox"
					checked={enabledForms.has(form.id)}
					onchange={() => onToggle(form.id)}
					class="w-4 h-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
				/>
				<span class="text-gray-700">{form.label}</span>
			</label>
		{/each}
	</div>

	<h2 class="text-xl font-semibold mb-4 text-gray-800">Options</h2>

	<label class="flex items-center gap-3 p-2 rounded hover:bg-gray-50 cursor-pointer mb-6">
		<input
			type="checkbox"
			checked={alwaysShowType}
			onchange={onToggleShowType}
			class="w-4 h-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
		/>
		<span class="text-gray-700">Always show adjective type (い / な)</span>
	</label>

	<button
		onclick={onStart}
		disabled={enabledForms.size === 0}
		class="w-full py-3 px-4 bg-indigo-600 text-white font-medium rounded-lg
			hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2
			disabled:opacity-50 disabled:cursor-not-allowed transition-colors mb-3"
	>
		Start Practice
	</button>

	<button
		onclick={onOpenSettings}
		class="w-full py-2 px-4 text-gray-600 text-sm font-medium rounded-lg
			hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-2
			transition-colors"
	>
		Settings
	</button>
</div>
