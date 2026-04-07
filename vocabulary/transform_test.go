package vocabulary

import (
	"testing"

	"github.com/the-witch-king/conjugreater/wanikani"
)

func TestNormalizePOS(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "i adjective",
			input: []string{"い adjective"},
			want:  []string{"i_adjective"},
		},
		{
			name:  "na adjective",
			input: []string{"な adjective"},
			want:  []string{"na_adjective"},
		},
		{
			name:  "noun and suru verb",
			input: []string{"noun", "する verb"},
			want:  []string{"noun", "suru_verb"},
		},
		{
			name:  "godan verb transitive",
			input: []string{"godan verb", "transitive verb"},
			want:  []string{"godan_verb", "transitive_verb"},
		},
		{
			name:  "unknown pos normalized",
			input: []string{"Some New Type"},
			want:  []string{"some_new_type"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizePOS(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("NormalizePOS(%v) = %v, want %v", tt.input, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("NormalizePOS(%v)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestTransformSubjects(t *testing.T) {
	subjects := []wanikani.Resource[wanikani.SubjectData]{
		{
			ID: 100,
			Data: wanikani.SubjectData{
				Characters:   "大きい",
				PartOfSpeech: []string{"い adjective"},
				Readings:     []wanikani.Reading{{Reading: "おおきい", Primary: true}},
				Meanings: []wanikani.Meaning{
					{Meaning: "big", Primary: true, AcceptedAnswer: true},
					{Meaning: "large", Primary: false, AcceptedAnswer: true},
					{Meaning: "huge", Primary: false, AcceptedAnswer: false},
				},
			},
		},
		{
			ID: 200,
			Data: wanikani.SubjectData{
				Characters:   "いい",
				PartOfSpeech: []string{"い adjective"},
				Readings:     []wanikani.Reading{{Reading: "いい", Primary: true}},
				Meanings:     []wanikani.Meaning{{Meaning: "good", Primary: true, AcceptedAnswer: true}},
			},
		},
		{
			ID: 300,
			Data: wanikani.SubjectData{
				Characters:   "元気",
				PartOfSpeech: []string{"な adjective", "noun"},
				Readings:     []wanikani.Reading{{Reading: "げんき", Primary: true}},
				Meanings:     []wanikani.Meaning{{Meaning: "energetic", Primary: true, AcceptedAnswer: true}},
			},
		},
		{
			ID: 400, // No assignment — should be excluded
			Data: wanikani.SubjectData{
				Characters:   "静か",
				PartOfSpeech: []string{"な adjective"},
				Readings:     []wanikani.Reading{{Reading: "しずか", Primary: true}},
				Meanings:     []wanikani.Meaning{{Meaning: "quiet", Primary: true, AcceptedAnswer: true}},
			},
		},
	}

	assignments := map[int]wanikani.AssignmentData{
		100: {SubjectID: 100, SRSStage: 5},
		200: {SubjectID: 200, SRSStage: 9},
		300: {SubjectID: 300, SRSStage: 3},
		// 400 has no assignment
	}

	words := TransformSubjects(subjects, assignments)

	if len(words) != 3 {
		t.Fatalf("got %d words, want 3", len(words))
	}

	// Check 大きい — regular i-adjective
	if words[0].Characters != "大きい" {
		t.Errorf("word[0].Characters = %q, want %q", words[0].Characters, "大きい")
	}
	if words[0].Reading != "おおきい" {
		t.Errorf("word[0].Reading = %q, want %q", words[0].Reading, "おおきい")
	}
	if words[0].IsException {
		t.Error("word[0] should not be an exception")
	}
	// Should include "big" and "large" (accepted) but not "huge" (not accepted)
	if len(words[0].Meanings) != 2 {
		t.Errorf("word[0] has %d meanings, want 2", len(words[0].Meanings))
	}

	// Check いい — exception
	if !words[1].IsException {
		t.Error("word[1] (いい) should be an exception")
	}
	if words[1].ExceptionID != "ii" {
		t.Errorf("word[1].ExceptionID = %q, want %q", words[1].ExceptionID, "ii")
	}

	// Check 元気 — na-adjective + noun
	if len(words[2].POS) != 2 || words[2].POS[0] != "na_adjective" || words[2].POS[1] != "noun" {
		t.Errorf("word[2].POS = %v, want [na_adjective noun]", words[2].POS)
	}
}

func TestTransformSubjects_NoAssignments(t *testing.T) {
	subjects := []wanikani.Resource[wanikani.SubjectData]{
		{ID: 100, Data: wanikani.SubjectData{Characters: "test"}},
	}
	words := TransformSubjects(subjects, map[int]wanikani.AssignmentData{})
	if len(words) != 0 {
		t.Errorf("got %d words, want 0 (no assignments)", len(words))
	}
}

func TestMergeWords(t *testing.T) {
	existing := []Word{
		{WaniKaniID: 1, Characters: "大きい", Reading: "おおきい"},
		{WaniKaniID: 2, Characters: "小さい", Reading: "ちいさい"},
	}
	updated := []Word{
		{WaniKaniID: 2, Characters: "小さい", Reading: "ちいさい", Meanings: []string{"small"}}, // updated
		{WaniKaniID: 3, Characters: "新しい", Reading: "あたらしい"},                               // new
	}

	result := MergeWords(existing, updated)

	if len(result) != 3 {
		t.Fatalf("got %d words, want 3", len(result))
	}
	// Order: 1, 2 (updated), 3 (new)
	if result[0].WaniKaniID != 1 {
		t.Errorf("result[0].WaniKaniID = %d, want 1", result[0].WaniKaniID)
	}
	if result[1].WaniKaniID != 2 || len(result[1].Meanings) != 1 {
		t.Error("result[1] should be the updated version of word 2")
	}
	if result[2].WaniKaniID != 3 {
		t.Errorf("result[2].WaniKaniID = %d, want 3", result[2].WaniKaniID)
	}
}
