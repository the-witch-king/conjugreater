import { describe, it, expect } from 'vitest';
import { conjugate } from './conjugation';
import type { Word, ConjugationForm } from './types';

function makeWord(overrides: Partial<Word> & Pick<Word, 'characters' | 'reading' | 'pos'>): Word {
	return {
		meanings: ['test'],
		wanikani_id: 1,
		is_exception: false,
		...overrides
	};
}

describe('i-adjective conjugation', () => {
	const word = makeWord({
		characters: '大きい',
		reading: 'おおきい',
		pos: ['i_adjective']
	});

	const cases: [ConjugationForm, string][] = [
		['affirmative_present', '大きいです'],
		['affirmative_past', '大きかったです'],
		['negative_present_dewa', '大きくありません'],
		['negative_present_ja', '大きくないです'],
		['negative_past_dewa', '大きくありませんでした'],
		['negative_past_ja', '大きくなかったです']
	];

	for (const [form, expected] of cases) {
		it(`${form} → ${expected}`, () => {
			const result = conjugate(word, form);
			expect(result[0]).toBe(expected);
		});
	}
});

describe('i-adjective conjugation (another word: 小さい)', () => {
	const word = makeWord({
		characters: '小さい',
		reading: 'ちいさい',
		pos: ['i_adjective']
	});

	it('affirmative present', () => {
		expect(conjugate(word, 'affirmative_present')[0]).toBe('小さいです');
	});

	it('affirmative past', () => {
		expect(conjugate(word, 'affirmative_past')[0]).toBe('小さかったです');
	});

	it('negative present (では)', () => {
		expect(conjugate(word, 'negative_present_dewa')[0]).toBe('小さくありません');
	});
});

describe('na-adjective conjugation', () => {
	const word = makeWord({
		characters: '元気',
		reading: 'げんき',
		pos: ['na_adjective', 'noun']
	});

	const cases: [ConjugationForm, string][] = [
		['affirmative_present', '元気です'],
		['affirmative_past', '元気でした'],
		['negative_present_dewa', '元気ではありません'],
		['negative_present_ja', '元気じゃありません'],
		['negative_past_dewa', '元気ではありませんでした'],
		['negative_past_ja', '元気じゃありませんでした']
	];

	for (const [form, expected] of cases) {
		it(`${form} → ${expected}`, () => {
			const result = conjugate(word, form);
			expect(result[0]).toBe(expected);
		});
	}

	it('negative present (じゃ) also accepts じゃないです', () => {
		const result = conjugate(word, 'negative_present_ja');
		expect(result).toContain('元気じゃないです');
	});

	it('negative past (じゃ) also accepts じゃなかったです', () => {
		const result = conjugate(word, 'negative_past_ja');
		expect(result).toContain('元気じゃなかったです');
	});
});

describe('na-adjective conjugation (しんせつ)', () => {
	const word = makeWord({
		characters: '親切',
		reading: 'しんせつ',
		pos: ['na_adjective']
	});

	it('negative present (では)', () => {
		expect(conjugate(word, 'negative_present_dewa')[0]).toBe('親切ではありません');
	});

	it('negative present (じゃ)', () => {
		expect(conjugate(word, 'negative_present_ja')[0]).toBe('親切じゃありません');
	});
});

describe('exception: いい', () => {
	const word = makeWord({
		characters: 'いい',
		reading: 'いい',
		pos: ['i_adjective'],
		is_exception: true,
		exception_id: 'ii'
	});

	const cases: [ConjugationForm, string][] = [
		['affirmative_present', 'いいです'],
		['affirmative_past', 'よかったです'],
		['negative_present_dewa', 'よくありません'],
		['negative_present_ja', 'よくないです'],
		['negative_past_dewa', 'よくありませんでした'],
		['negative_past_ja', 'よくなかったです']
	];

	for (const [form, expected] of cases) {
		it(`${form} → ${expected}`, () => {
			const result = conjugate(word, form);
			expect(result[0]).toBe(expected);
		});
	}
});

describe('exception: かっこいい', () => {
	const word = makeWord({
		characters: 'かっこいい',
		reading: 'かっこいい',
		pos: ['i_adjective'],
		is_exception: true,
		exception_id: 'kakkoii'
	});

	const cases: [ConjugationForm, string][] = [
		['affirmative_present', 'かっこいいです'],
		['affirmative_past', 'かっこよかったです'],
		['negative_present_dewa', 'かっこよくありません'],
		['negative_present_ja', 'かっこよくないです'],
		['negative_past_dewa', 'かっこよくありませんでした'],
		['negative_past_ja', 'かっこよくなかったです']
	];

	for (const [form, expected] of cases) {
		it(`${form} → ${expected}`, () => {
			const result = conjugate(word, form);
			expect(result[0]).toBe(expected);
		});
	}
});

describe('いい-compound: 運がいい', () => {
	const word = makeWord({
		characters: '運がいい',
		reading: 'うんがいい',
		pos: ['expression', 'i_adjective'],
		is_exception: true,
		exception_id: 'ii_compound'
	});

	const cases: [ConjugationForm, string][] = [
		['affirmative_present', '運がいいです'],
		['affirmative_past', '運がよかったです'],
		['negative_present_dewa', '運がよくありません'],
		['negative_present_ja', '運がよくないです'],
		['negative_past_dewa', '運がよくありませんでした'],
		['negative_past_ja', '運がよくなかったです']
	];

	for (const [form, expected] of cases) {
		it(`${form} → ${expected}`, () => {
			const result = conjugate(word, form);
			expect(result[0]).toBe(expected);
		});
	}
});

describe('i-adjective ending in いい but not 良い-related: かわいい', () => {
	const word = makeWord({
		characters: 'かわいい',
		reading: 'かわいい',
		pos: ['i_adjective']
	});

	const cases: [ConjugationForm, string][] = [
		['affirmative_present', 'かわいいです'],
		['affirmative_past', 'かわいかったです'],
		['negative_present_dewa', 'かわいくありません'],
		['negative_present_ja', 'かわいくないです'],
		['negative_past_dewa', 'かわいくありませんでした'],
		['negative_past_ja', 'かわいくなかったです']
	];

	for (const [form, expected] of cases) {
		it(`${form} → ${expected}`, () => {
			const result = conjugate(word, form);
			expect(result[0]).toBe(expected);
		});
	}
});

describe('edge cases', () => {
	it('returns empty array for unknown POS', () => {
		const word = makeWord({
			characters: '食べる',
			reading: 'たべる',
			pos: ['ichidan_verb']
		});
		expect(conjugate(word, 'affirmative_present')).toEqual([]);
	});

	it('unknown exception falls back to regular i-adj', () => {
		const word = makeWord({
			characters: '新しい',
			reading: 'あたらしい',
			pos: ['i_adjective'],
			is_exception: true,
			exception_id: 'unknown'
		});
		expect(conjugate(word, 'affirmative_past')[0]).toBe('新しかったです');
	});
});
