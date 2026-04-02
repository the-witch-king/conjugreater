import { conjugate } from './conjugation';
import { ALL_FORMS } from './types';
import type { Word, ConjugationForm, QuizQuestion } from './types';

export type QuizState = 'setup' | 'prompting' | 'evaluating';

export class Quiz {
	state = $state<QuizState>('setup');
	enabledForms = $state<Set<ConjugationForm>>(new Set(ALL_FORMS.map((f) => f.id)));
	currentQuestion = $state<QuizQuestion | null>(null);
	userAnswer = $state('');
	isCorrect = $state(false);

	private words: Word[];

	constructor(words: Word[]) {
		this.words = words;
	}

	start() {
		if (this.enabledForms.size === 0 || this.words.length === 0) return;
		this.nextQuestion();
	}

	nextQuestion() {
		const word = this.words[Math.floor(Math.random() * this.words.length)];
		const forms = Array.from(this.enabledForms);
		const formId = forms[Math.floor(Math.random() * forms.length)];
		const formInfo = ALL_FORMS.find((f) => f.id === formId)!;
		const answers = conjugate(word, formId);

		this.currentQuestion = {
			word,
			form: formId,
			formLabel: formInfo.label,
			correctAnswer: answers[0],
			acceptedAnswers: answers
		};
		this.userAnswer = '';
		this.isCorrect = false;
		this.state = 'prompting';
	}

	submit() {
		if (!this.currentQuestion || this.state !== 'prompting') return;
		const normalized = this.userAnswer.trim();
		this.isCorrect = this.currentQuestion.acceptedAnswers.includes(normalized);
		this.state = 'evaluating';
	}

	backToSetup() {
		this.state = 'setup';
		this.currentQuestion = null;
	}

	toggleForm(form: ConjugationForm) {
		const next = new Set(this.enabledForms);
		if (next.has(form)) {
			next.delete(form);
		} else {
			next.add(form);
		}
		this.enabledForms = next;
	}
}
