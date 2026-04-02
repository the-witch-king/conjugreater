export interface Word {
	characters: string;
	reading: string;
	meanings: string[];
	pos: string[];
	wanikani_id: number;
	is_exception: boolean;
	exception_id?: string;
}

export type ConjugationForm =
	| 'affirmative_present'
	| 'affirmative_past'
	| 'negative_present_dewa'
	| 'negative_present_ja'
	| 'negative_past_dewa'
	| 'negative_past_ja';

export interface FormInfo {
	id: ConjugationForm;
	label: string;
}

export const ALL_FORMS: FormInfo[] = [
	{ id: 'affirmative_present', label: 'Affirmative present' },
	{ id: 'affirmative_past', label: 'Affirmative past' },
	{ id: 'negative_present_dewa', label: 'Negative present (では)' },
	{ id: 'negative_present_ja', label: 'Negative present (じゃ)' },
	{ id: 'negative_past_dewa', label: 'Negative past (では)' },
	{ id: 'negative_past_ja', label: 'Negative past (じゃ)' }
];

export interface QuizQuestion {
	word: Word;
	form: ConjugationForm;
	formLabel: string;
	correctAnswer: string;
	acceptedAnswers: string[];
}
