<script lang="ts">
	interface Props {
		onClose: () => void;
		onVocabularyUpdated: () => void;
	}

	let { onClose, onVocabularyUpdated }: Props = $props();

	let token = $state('');
	let tokenLoaded = $state(false);
	let saving = $state(false);
	let fetching = $state(false);
	let message = $state('');
	let messageType = $state<'success' | 'error'>('success');

	async function loadConfig() {
		try {
			const res = await fetch('/api/config');
			const data = await res.json();
			if (data.token) token = data.token;
		} catch {
			// ignore
		} finally {
			tokenLoaded = true;
		}
	}

	async function saveToken() {
		saving = true;
		message = '';
		try {
			const res = await fetch('/api/config', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ token: token.trim() })
			});
			if (!res.ok) {
				throw new Error(await res.text());
			}
			showMessage('success', 'API token saved');
		} catch (e) {
			showMessage('error', e instanceof Error ? e.message : 'Failed to save');
		} finally {
			saving = false;
		}
	}

	async function fetchVocabulary() {
		fetching = true;
		message = '';
		try {
			const res = await fetch('/api/fetch', { method: 'POST' });
			if (!res.ok) {
				throw new Error(await res.text());
			}
			const data = await res.json();
			showMessage('success', `Fetched ${data.total_words} words (${data.new_words} new/updated)`);
			onVocabularyUpdated();
		} catch (e) {
			showMessage('error', e instanceof Error ? e.message : 'Fetch failed');
		} finally {
			fetching = false;
		}
	}

	function showMessage(type: 'success' | 'error', text: string) {
		messageType = type;
		message = text;
	}

	loadConfig();
</script>

<div class="max-w-md mx-auto">
	<div class="flex items-center justify-between mb-6">
		<h2 class="text-xl font-semibold text-gray-800">Settings</h2>
		<button
			onclick={onClose}
			class="text-gray-400 hover:text-gray-600 text-2xl leading-none"
		>
			&times;
		</button>
	</div>

	{#if !tokenLoaded}
		<p class="text-gray-500">Loading...</p>
	{:else}
		<div class="space-y-6">
			<div>
				<label for="api-token" class="block text-sm font-medium text-gray-700 mb-1">
					WaniKani API Token
				</label>
				<input
					id="api-token"
					type="password"
					bind:value={token}
					placeholder="Enter your API token"
					class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
				/>
				<p class="mt-1 text-xs text-gray-500">
					Get yours at
					<a href="https://www.wanikani.com/settings/personal_access_tokens" target="_blank" class="text-indigo-600 hover:underline">
						wanikani.com/settings/personal_access_tokens
					</a>
				</p>
				<button
					onclick={saveToken}
					disabled={saving || !token.trim()}
					class="mt-2 px-4 py-2 bg-gray-600 text-white text-sm font-medium rounded-lg
						hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2
						disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
				>
					{saving ? 'Saving...' : 'Save Token'}
				</button>
			</div>

			<div>
				<button
					onclick={fetchVocabulary}
					disabled={fetching || !token.trim()}
					class="w-full py-3 px-4 bg-indigo-600 text-white font-medium rounded-lg
						hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2
						disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
				>
					{fetching ? 'Fetching from WaniKani...' : 'Fetch Vocabulary'}
				</button>
				{#if fetching}
					<p class="mt-2 text-sm text-gray-500">This may take a minute...</p>
				{/if}
			</div>

			{#if message}
				<div class={messageType === 'success' ? 'bg-green-50 text-green-700 p-3 rounded-lg text-sm' : 'bg-red-50 text-red-700 p-3 rounded-lg text-sm'}>
					{message}
				</div>
			{/if}
		</div>
	{/if}
</div>
