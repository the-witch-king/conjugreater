package vocabulary

import (
	"github.com/the-witch-king/conjugreater/fetcher/wanikani"
)

// TransformSubjects converts WaniKani subjects into our Word format.
// Only subjects that have an assignment (exist in assignmentsBySubjectID) are included.
func TransformSubjects(
	subjects []wanikani.Resource[wanikani.SubjectData],
	assignmentsBySubjectID map[int]wanikani.AssignmentData,
) []Word {
	var words []Word

	for _, subj := range subjects {
		// Only include words the user has unlocked (has an assignment)
		if _, hasAssignment := assignmentsBySubjectID[subj.ID]; !hasAssignment {
			continue
		}

		word := Word{
			Characters: subj.Data.Characters,
			Reading:    primaryReading(subj.Data.Readings),
			Meanings:   acceptedMeanings(subj.Data.Meanings),
			POS:        NormalizePOS(subj.Data.PartOfSpeech),
			WaniKaniID: subj.ID,
		}

		if exc, ok := Exceptions[subj.Data.Characters]; ok {
			word.IsException = true
			word.ExceptionID = exc.ID
		}

		words = append(words, word)
	}

	return words
}

// MergeWords merges updated words into the existing word list.
// Updated words replace existing ones by WaniKaniID; new words are appended.
func MergeWords(existing, updated []Word) []Word {
	byID := make(map[int]Word, len(existing))
	order := make([]int, 0, len(existing))

	for _, w := range existing {
		byID[w.WaniKaniID] = w
		order = append(order, w.WaniKaniID)
	}

	for _, w := range updated {
		if _, exists := byID[w.WaniKaniID]; !exists {
			order = append(order, w.WaniKaniID)
		}
		byID[w.WaniKaniID] = w
	}

	result := make([]Word, 0, len(order))
	for _, id := range order {
		result = append(result, byID[id])
	}
	return result
}

func primaryReading(readings []wanikani.Reading) string {
	for _, r := range readings {
		if r.Primary {
			return r.Reading
		}
	}
	if len(readings) > 0 {
		return readings[0].Reading
	}
	return ""
}

func acceptedMeanings(meanings []wanikani.Meaning) []string {
	var result []string
	for _, m := range meanings {
		if m.AcceptedAnswer {
			result = append(result, m.Meaning)
		}
	}
	return result
}
