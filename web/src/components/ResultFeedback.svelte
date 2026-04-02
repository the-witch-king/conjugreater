<script lang="ts">
	interface Props {
		isCorrect: boolean;
		correctAnswer: string;
		userAnswer: string;
		onNext: () => void;
	}

	let { isCorrect, correctAnswer, userAnswer, onNext }: Props = $props();
</script>

<div class="mt-4 space-y-3">
	{#if isCorrect}
		<div class="flex items-center gap-2 text-green-700 bg-green-50 p-3 rounded-lg">
			<span class="text-2xl">&#10004;</span>
			<span class="font-medium">Correct!</span>
		</div>
	{:else}
		<div class="space-y-2">
			<div class="flex items-center gap-2 text-red-700 bg-red-50 p-3 rounded-lg">
				<span class="text-2xl">&#10008;</span>
				<span class="font-medium">Incorrect</span>
			</div>
			<div class="text-sm text-gray-600">
				<p>Your answer: <span class="font-medium text-red-600">{userAnswer}</span></p>
				<p>Correct answer: <span class="font-medium text-green-700">{correctAnswer}</span></p>
			</div>
		</div>
	{/if}

	<button
		onclick={onNext}
		class="w-full py-2 px-4 bg-gray-100 text-gray-700 font-medium rounded-lg
			hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-2
			transition-colors"
	>
		Next
	</button>
</div>
