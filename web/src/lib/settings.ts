import type { ConjugationForm } from './types';
import { ALL_FORMS } from './types';

const STORAGE_KEY = 'conjugreater-settings';

export interface Settings {
	enabledForms: ConjugationForm[];
	alwaysShowType: boolean;
	targetCount: number;
	includeIAdjectives: boolean;
	includeNaAdjectives: boolean;
}

export const defaults: Settings = {
	enabledForms: ALL_FORMS.map((f) => f.id),
	alwaysShowType: false,
	targetCount: 20,
	includeIAdjectives: true,
	includeNaAdjectives: true
};

export function loadSettings(): Settings {
	try {
		const raw = localStorage.getItem(STORAGE_KEY);
		if (!raw) return { ...defaults };
		const parsed = JSON.parse(raw);
		return {
			enabledForms: Array.isArray(parsed.enabledForms) ? parsed.enabledForms : defaults.enabledForms,
			alwaysShowType: typeof parsed.alwaysShowType === 'boolean' ? parsed.alwaysShowType : defaults.alwaysShowType,
			targetCount: typeof parsed.targetCount === 'number' ? parsed.targetCount : defaults.targetCount,
			includeIAdjectives: typeof parsed.includeIAdjectives === 'boolean' ? parsed.includeIAdjectives : defaults.includeIAdjectives,
			includeNaAdjectives: typeof parsed.includeNaAdjectives === 'boolean' ? parsed.includeNaAdjectives : defaults.includeNaAdjectives
		};
	} catch {
		return { ...defaults };
	}
}

export function saveSettings(settings: Settings): void {
	try {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
	} catch {
		// Storage full or unavailable — silently ignore
	}
}

export function resetSettings(): void {
	try {
		localStorage.removeItem(STORAGE_KEY);
	} catch {
		// Silently ignore
	}
}
