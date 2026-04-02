<script lang="ts">
	import ResultFeedback from './ResultFeedback.svelte';
	import type { QuizQuestion } from '$lib/types';

	interface Props {
		question: QuizQuestion;
		quizState: 'prompting' | 'evaluating';
		userAnswer: string;
		isCorrect: boolean;
		alwaysShowType: boolean;
		typeRevealed: boolean;
		adjType: string;
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
		onInput,
		onSubmit,
		onNext,
		onBack,
		onRevealType
	}: Props = $props();

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			if (quizState === 'prompting') {
				onSubmit();
			} else {
				onNext();
			}
		} else if (e.key === 'Escape' && quizState === 'prompting') {
			onInput('');
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="max-w-md mx-auto">
	<button onclick={onBack} class="text-sm text-gray-500 hover:text-gray-700 mb-6 block">
		&larr; Back to settings
	</button>

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
</div>
