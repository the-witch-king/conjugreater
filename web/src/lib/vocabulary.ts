import type { Word } from './types';

interface VocabularyData {
	words: Word[];
}

export async function loadVocabulary(): Promise<Word[]> {
	const response = await fetch('/vocabulary.json');
	if (!response.ok) {
		throw new Error(`Failed to load vocabulary: ${response.status}`);
	}
	const data: VocabularyData = await response.json();
	return data.words.filter(isAdjective);
}

function isAdjective(word: Word): boolean {
	if (!word.pos) return false;
	return word.pos.includes('i_adjective') || word.pos.includes('na_adjective');
}
