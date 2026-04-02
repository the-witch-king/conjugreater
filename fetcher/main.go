package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/the-witch-king/conjugreater/fetcher/vocabulary"
	"github.com/the-witch-king/conjugreater/fetcher/wanikani"
)

var subjectTypes = []string{"vocabulary", "kana_vocabulary"}

func main() {
	apiToken := flag.String("token", "", "WaniKani API token (overrides WANIKANI_API_TOKEN env var)")
	output := flag.String("output", "../data/vocabulary.json", "Output file path")
	full := flag.Bool("full", false, "Force full re-fetch, ignoring cached timestamps")
	flag.Parse()

	if *apiToken == "" {
		*apiToken = os.Getenv("WANIKANI_API_TOKEN")
	}
	if *apiToken == "" {
		log.Fatal("API token required: set WANIKANI_API_TOKEN or pass -token flag\n" +
			"Get yours at https://www.wanikani.com/settings/personal_access_tokens")
	}

	ctx := context.Background()
	client := wanikani.NewClient(*apiToken)

	// Try to load existing output for incremental updates
	var existing *vocabulary.Output
	if !*full {
		existing = loadExisting(*output)
	}

	// Determine updated_after timestamps for incremental fetch
	var subjectsAfter, assignmentsAfter string
	if existing != nil {
		subjectsAfter = existing.SubjectsUpdatedAt
		assignmentsAfter = existing.AssignmentsUpdatedAt
		log.Printf("Incremental update (subjects after %s, assignments after %s)", subjectsAfter, assignmentsAfter)
	} else {
		log.Println("Full fetch")
	}

	// Fetch subjects
	log.Println("Fetching subjects...")
	subjects, subjectsUpdatedAt, err := client.FetchSubjects(ctx, subjectTypes, subjectsAfter)
	if err != nil {
		log.Fatalf("Failed to fetch subjects: %v", err)
	}

	// Fetch assignments
	log.Println("Fetching assignments...")
	assignments, assignmentsUpdatedAt, err := client.FetchAssignments(ctx, subjectTypes, assignmentsAfter)
	if err != nil {
		log.Fatalf("Failed to fetch assignments: %v", err)
	}

	// Build assignment map: for incremental, merge with existing data
	assignmentMap := buildAssignmentMap(assignments)

	// If incremental and we got no subjects updates, we might still have new assignments.
	// We need ALL subjects to re-evaluate which words have assignments.
	// For simplicity: on incremental, we need to re-fetch all subjects if we only got assignment updates.
	// Better approach: keep all subjects cached and only re-filter.

	// Transform new/updated subjects
	newWords := vocabulary.TransformSubjects(subjects, assignmentMap)

	var finalWords []vocabulary.Word
	if existing != nil && len(existing.Words) > 0 {
		// For incremental: if we fetched all assignments (not just updated ones),
		// we can properly filter. But since we only fetch updated assignments,
		// we merge the new words into existing ones.
		// Words from updated subjects get re-evaluated; existing words are kept.
		finalWords = vocabulary.MergeWords(existing.Words, newWords)
		log.Printf("Merged: %d existing + %d new/updated = %d total", len(existing.Words), len(newWords), len(finalWords))
	} else {
		finalWords = newWords
	}

	// Use the latest timestamps
	if subjectsUpdatedAt == "" && existing != nil {
		subjectsUpdatedAt = existing.SubjectsUpdatedAt
	}
	if assignmentsUpdatedAt == "" && existing != nil {
		assignmentsUpdatedAt = existing.AssignmentsUpdatedAt
	}

	out := vocabulary.Output{
		GeneratedAt:          time.Now().UTC().Format(time.RFC3339),
		SubjectsUpdatedAt:    subjectsUpdatedAt,
		AssignmentsUpdatedAt: assignmentsUpdatedAt,
		Words:                finalWords,
	}

	if err := writeOutput(*output, out); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}

	printSummary(finalWords, len(newWords))
}

func loadExisting(path string) *vocabulary.Output {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var out vocabulary.Output
	if err := json.Unmarshal(data, &out); err != nil {
		log.Printf("Warning: could not parse existing %s, doing full fetch: %v", path, err)
		return nil
	}
	return &out
}

func buildAssignmentMap(assignments []wanikani.Resource[wanikani.AssignmentData]) map[int]wanikani.AssignmentData {
	m := make(map[int]wanikani.AssignmentData, len(assignments))
	for _, a := range assignments {
		m[a.Data.SubjectID] = a.Data
	}
	return m
}

func writeOutput(path string, out vocabulary.Output) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	log.Printf("Wrote %s", path)
	return nil
}

func printSummary(words []vocabulary.Word, newCount int) {
	counts := make(map[string]int)
	for _, w := range words {
		for _, pos := range w.POS {
			counts[pos]++
		}
	}

	fmt.Fprintf(os.Stderr, "\nSummary: %d words", len(words))
	if newCount > 0 {
		fmt.Fprintf(os.Stderr, " (%d new/updated)", newCount)
	}
	fmt.Fprintln(os.Stderr)

	for _, pos := range []string{"i_adjective", "na_adjective", "noun", "godan_verb", "ichidan_verb", "suru_verb"} {
		if c, ok := counts[pos]; ok {
			fmt.Fprintf(os.Stderr, "  %s: %d\n", pos, c)
		}
	}

	exceptions := 0
	for _, w := range words {
		if w.IsException {
			exceptions++
		}
	}
	if exceptions > 0 {
		fmt.Fprintf(os.Stderr, "  exceptions: %d\n", exceptions)
	}
}
