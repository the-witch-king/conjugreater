<script lang="ts">
	import { ALL_FORMS } from '$lib/types';
	import type { ConjugationForm } from '$lib/types';

	interface Props {
		enabledForms: Set<ConjugationForm>;
		alwaysShowType: boolean;
		targetCount: number;
		includeIAdjectives: boolean;
		includeNaAdjectives: boolean;
		onToggle: (form: ConjugationForm) => void;
		onToggleShowType: () => void;
		onTargetChange: (value: number) => void;
		onToggleIAdjectives: () => void;
		onToggleNaAdjectives: () => void;
		onStart: () => void;
		onReset: () => void;
		onOpenSettings: () => void;
	}

	let { enabledForms, alwaysShowType, targetCount, includeIAdjectives, includeNaAdjectives, onToggle, onToggleShowType, onTargetChange, onToggleIAdjectives, onToggleNaAdjectives, onStart, onReset, onOpenSettings }: Props = $props();
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

	<h2 class="text-xl font-semibold mb-4 text-gray-800">Adjective types</h2>

	<div class="space-y-2 mb-6">
		<label class="flex items-center gap-3 p-2 rounded hover:bg-gray-50 cursor-pointer">
			<input
				type="checkbox"
				checked={includeIAdjectives}
				onchange={onToggleIAdjectives}
				class="w-4 h-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
			/>
			<span class="text-gray-700">い adjectives</span>
		</label>
		<label class="flex items-center gap-3 p-2 rounded hover:bg-gray-50 cursor-pointer">
			<input
				type="checkbox"
				checked={includeNaAdjectives}
				onchange={onToggleNaAdjectives}
				class="w-4 h-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
			/>
			<span class="text-gray-700">な adjectives</span>
		</label>
	</div>

	<h2 class="text-xl font-semibold mb-4 text-gray-800">Options</h2>

	<div class="flex items-center gap-3 p-2 mb-2">
		<label for="target-count" class="text-gray-700">Questions</label>
		<input
			id="target-count"
			type="number"
			min="1"
			max="999"
			value={targetCount}
			oninput={(e) => onTargetChange(Math.max(1, parseInt(e.currentTarget.value) || 1))}
			class="w-20 px-2 py-1 border border-gray-300 rounded-lg text-center
				focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
		/>
	</div>

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
		disabled={enabledForms.size === 0 || (!includeIAdjectives && !includeNaAdjectives)}
		class="w-full py-3 px-4 bg-indigo-600 text-white font-medium rounded-lg
			hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2
			disabled:opacity-50 disabled:cursor-not-allowed transition-colors mb-3"
	>
		Start Practice
	</button>

	<div class="flex gap-3">
		<button
			onclick={onOpenSettings}
			class="flex-1 py-2 px-4 text-gray-600 text-sm font-medium rounded-lg
				hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-2
				transition-colors"
		>
			Settings
		</button>
		<button
			onclick={onReset}
			class="flex-1 py-2 px-4 text-gray-600 text-sm font-medium rounded-lg
				hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-2
				transition-colors"
		>
			Reset Defaults
		</button>
	</div>
</div>
