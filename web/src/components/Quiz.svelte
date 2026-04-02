<script lang="ts">
	import ResultFeedback from './ResultFeedback.svelte';
	import type { QuizQuestion, AnswerRecord } from '$lib/types';

	interface Props {
		question: QuizQuestion;
		quizState: 'prompting' | 'evaluating';
		userAnswer: string;
		isCorrect: boolean;
		alwaysShowType: boolean;
		typeRevealed: boolean;
		adjType: string;
		history: AnswerRecord[];
		targetCount: number;
		onInput: (value: string) => void;
		onSubmit: () => void;
		onNext: () => void;
		onBack: () => void;
		onRevealType: () => void;
	}

	let {
		question,
		quizState,
		userAnswer,
		isCorrect,
		alwaysShowType,
		typeRevealed,
		adjType,
		history,
		targetCount,
		onInput,
		onSubmit,
		onNext,
		onBack,
		onRevealType
	}: Props = $props();

	let inputEl: HTMLInputElement | undefined;

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			if (quizState === 'prompting') {
				onSubmit();
			} else {
				onNext();
				// Focus after Svelte re-enables the input on next tick
				requestAnimationFrame(() => inputEl?.focus());
			}
		} else if (e.key === 'Escape' && quizState === 'prompting') {
			onInput('');
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="max-w-md mx-auto">
	<div class="flex items-center justify-between mb-4">
		<button onclick={onBack} class="text-sm text-gray-500 hover:text-gray-700">
			&larr; Back to settings
		</button>

		<div class="flex items-center gap-3 text-sm">
			<span class="text-gray-500 font-medium">{history.length}/{targetCount}</span>
			{#if history.length > 0}
				{@const correctCount = history.filter((r) => r.correct).length}
				<span class="text-green-600 font-medium">{correctCount} correct</span>
				<span class="text-gray-300">|</span>
				<span class="text-red-500 font-medium">{history.length - correctCount} wrong</span>
			{/if}
		</div>
	</div>

	<!-- Progress bar -->
	<div class="w-full bg-gray-200 rounded-full h-1.5 mb-6">
		<div
			class="bg-indigo-600 h-1.5 rounded-full transition-all duration-300"
			style="width: {Math.round((history.length / targetCount) * 100)}%"
		></div>
	</div>

	<div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
		<!-- Word display -->
		<div class="text-center mb-6">
			<p class="text-4xl font-bold text-gray-900 mb-1">{question.word.characters}</p>
			<p class="text-lg text-gray-500 mb-1">{question.word.reading}</p>
			<p class="text-sm text-gray-400">{question.word.meanings.join(', ')}</p>

			<!-- Adjective type hint -->
			{#if alwaysShowType || typeRevealed}
				<p class="mt-2 text-xs text-indigo-500 font-medium">{adjType}</p>
			{:else}
				<button
					onclick={onRevealType}
					class="mt-2 text-xs text-gray-300 hover:text-gray-500 transition-colors"
				>
					show type
				</button>
			{/if}
		</div>

		<!-- Target form -->
		<div class="text-center mb-6">
			<span
				class="inline-block px-3 py-1 bg-indigo-50 text-indigo-700 rounded-full text-sm font-medium"
			>
				{question.formLabel}
			</span>
		</div>

		<!-- Input -->
		<div>
			<input
				bind:this={inputEl}
				type="text"
				lang="ja"
				value={userAnswer}
				oninput={(e) => onInput(e.currentTarget.value)}
				disabled={quizState === 'evaluating'}
				placeholder="Type your answer..."
				class="w-full px-4 py-3 text-lg border border-gray-300 rounded-lg
					focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500
					disabled:bg-gray-50 disabled:text-gray-500"
			/>

			{#if quizState === 'prompting'}
				<button
					onclick={onSubmit}
					disabled={userAnswer.trim() === ''}
					class="w-full mt-3 py-2 px-4 bg-indigo-600 text-white font-medium rounded-lg
						hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2
						disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
				>
					Submit
				</button>
			{/if}

			{#if quizState === 'evaluating'}
				<ResultFeedback
					{isCorrect}
					correctAnswer={question.correctAnswer}
					{userAnswer}
					{onNext}
				/>
			{/if}
		</div>
	</div>

	<!-- Answer history -->
	{#if history.length > 0}
		<details class="mt-6">
			<summary class="text-sm text-gray-500 cursor-pointer hover:text-gray-700">
				History ({history.length} answers)
			</summary>
			<div class="mt-2 space-y-1">
				{#each [...history].reverse() as record}
					<div
						class="flex items-center justify-between text-sm p-2 rounded {record.correct
							? 'bg-green-50'
							: 'bg-red-50'}"
					>
						<div>
							<span class="font-medium">{record.word.characters}</span>
							<span class="text-gray-400 mx-1">&rarr;</span>
							<span class="text-gray-600">{record.formLabel}</span>
						</div>
						<div>
							{#if record.correct}
								<span class="text-green-600">{record.userAnswer}</span>
							{:else}
								<span class="text-red-500 line-through">{record.userAnswer}</span>
								<span class="text-green-600 ml-1">{record.correctAnswer}</span>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		</details>
	{/if}
</div>
