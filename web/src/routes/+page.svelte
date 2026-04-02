<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { loadVocabulary } from '$lib/vocabulary';
	import { conjugate } from '$lib/conjugation';
	import { ALL_FORMS } from '$lib/types';
	import type { Word, ConjugationForm, QuizQuestion } from '$lib/types';
	import FormSelector from '../components/FormSelector.svelte';
	import QuizComponent from '../components/Quiz.svelte';

	let adjectives = $state<Word[]>([]);
	let loading = $state(true);
	let error = $state('');
	let alwaysShowType = $state(false);

	// Quiz state
	let quizState = $state<'setup' | 'prompting' | 'evaluating'>('setup');
	let enabledForms = $state<Set<ConjugationForm>>(new Set(ALL_FORMS.map((f) => f.id)));
	let currentQuestion = $state<QuizQuestion | null>(null);
	let userAnswer = $state('');
	let isCorrect = $state(false);
	let typeRevealed = $state(false);

	onMount(async () => {
		try {
			adjectives = await loadVocabulary();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load vocabulary';
		} finally {
			loading = false;
		}
	});

	function start() {
		if (enabledForms.size === 0 || adjectives.length === 0) return;
		nextQuestion();
	}

	function nextQuestion() {
		const word = adjectives[Math.floor(Math.random() * adjectives.length)];
		const forms = Array.from(enabledForms);
		const formId = forms[Math.floor(Math.random() * forms.length)];
		const formInfo = ALL_FORMS.find((f) => f.id === formId)!;
		const answers = conjugate(word, formId);

		currentQuestion = {
			word,
			form: formId,
			formLabel: formInfo.label,
			correctAnswer: answers[0],
			acceptedAnswers: answers
		};
		userAnswer = '';
		isCorrect = false;
		typeRevealed = false;
		quizState = 'prompting';
	}

	function submit() {
		if (!currentQuestion || quizState !== 'prompting') return;
		isCorrect = currentQuestion.acceptedAnswers.includes(userAnswer.trim());
		quizState = 'evaluating';
	}

	function backToSetup() {
		quizState = 'setup';
		currentQuestion = null;
	}

	function toggleForm(form: ConjugationForm) {
		const next = new Set(enabledForms);
		if (next.has(form)) {
			next.delete(form);
		} else {
			next.add(form);
		}
		enabledForms = next;
	}
</script>

<svelte:head>
	<title>Conjugreater</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 py-8 px-4">
	<div class="max-w-lg mx-auto">
		<h1 class="text-3xl font-bold text-center text-gray-900 mb-8">Conjugreater</h1>

		{#if loading}
			<p class="text-center text-gray-500">Loading vocabulary...</p>
		{:else if error}
			<div class="bg-red-50 text-red-700 p-4 rounded-lg">
				<p class="font-medium">Error</p>
				<p class="text-sm">{error}</p>
			</div>
		{:else if adjectives.length === 0}
			<div class="bg-yellow-50 text-yellow-700 p-4 rounded-lg">
				<p class="font-medium">No adjectives found</p>
				<p class="text-sm">
					Run the fetcher first to populate vocabulary data.
				</p>
			</div>
		{:else if quizState === 'setup'}
			<FormSelector
				{enabledForms}
				{alwaysShowType}
				onToggle={toggleForm}
				onToggleShowType={() => (alwaysShowType = !alwaysShowType)}
				onStart={start}
			/>
		{:else if currentQuestion}
			<QuizComponent
				question={currentQuestion}
				quizState={quizState}
				{userAnswer}
				{isCorrect}
				{alwaysShowType}
				{typeRevealed}
				adjType={currentQuestion.word.pos.includes('i_adjective') ? 'い adjective' : 'な adjective'}
				onInput={(v) => (userAnswer = v)}
				onSubmit={submit}
				onNext={nextQuestion}
				onBack={backToSetup}
				onRevealType={() => (typeRevealed = true)}
			/>
		{/if}
	</div>
</div>
