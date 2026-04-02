import type { Word, ConjugationForm } from './types';

/**
 * Conjugate a word into the specified polite form.
 * Returns an array of accepted answers (first element is the primary/displayed answer).
 */
export function conjugate(word: Word, form: ConjugationForm): string[] {
	if (word.is_exception && word.exception_id) {
		return conjugateException(word, form);
	}

	if (word.pos.includes('i_adjective')) {
		return conjugateIAdj(word.characters, form);
	}

	if (word.pos.includes('na_adjective')) {
		return conjugateNaAdj(word.characters, form);
	}

	return [];
}

function conjugateIAdj(characters: string, form: ConjugationForm): string[] {
	const stem = characters.slice(0, -1); // remove trailing い

	switch (form) {
		case 'affirmative_present':
			return [characters + 'です'];
		case 'affirmative_past':
			return [stem + 'かったです'];
		case 'negative_present_dewa':
			return [stem + 'くありません'];
		case 'negative_present_ja':
			return [stem + 'くないです'];
		case 'negative_past_dewa':
			return [stem + 'くありませんでした'];
		case 'negative_past_ja':
			return [stem + 'くなかったです'];
	}
}

function conjugateNaAdj(characters: string, form: ConjugationForm): string[] {
	switch (form) {
		case 'affirmative_present':
			return [characters + 'です'];
		case 'affirmative_past':
			return [characters + 'でした'];
		case 'negative_present_dewa':
			return [characters + 'ではありません'];
		case 'negative_present_ja':
			return [characters + 'じゃありません', characters + 'じゃないです'];
		case 'negative_past_dewa':
			return [characters + 'ではありませんでした'];
		case 'negative_past_ja':
			return [characters + 'じゃありませんでした', characters + 'じゃなかったです'];
	}
}

function conjugateException(word: Word, form: ConjugationForm): string[] {
	switch (word.exception_id) {
		case 'ii':
			return conjugateIi(form);
		case 'kakkoii':
			return conjugateKakkoii(form);
		default:
			// Unknown exception: fall back to regular i-adjective
			return conjugateIAdj(word.characters, form);
	}
}

function conjugateIi(form: ConjugationForm): string[] {
	switch (form) {
		case 'affirmative_present':
			return ['いいです'];
		case 'affirmative_past':
			return ['よかったです'];
		case 'negative_present_dewa':
			return ['よくありません'];
		case 'negative_present_ja':
			return ['よくないです'];
		case 'negative_past_dewa':
			return ['よくありませんでした'];
		case 'negative_past_ja':
			return ['よくなかったです'];
	}
}

function conjugateKakkoii(form: ConjugationForm): string[] {
	switch (form) {
		case 'affirmative_present':
			return ['かっこいいです'];
		case 'affirmative_past':
			return ['かっこよかったです'];
		case 'negative_present_dewa':
			return ['かっこよくありません'];
		case 'negative_present_ja':
			return ['かっこよくないです'];
		case 'negative_past_dewa':
			return ['かっこよくありませんでした'];
		case 'negative_past_ja':
			return ['かっこよくなかったです'];
	}
}
