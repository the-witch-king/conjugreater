package vocabulary

// Word represents a vocabulary item in the output JSON.
type Word struct {
	Characters  string   `json:"characters"`
	Reading     string   `json:"reading"`
	Meanings    []string `json:"meanings"`
	POS         []string `json:"pos"`
	WaniKaniID  int      `json:"wanikani_id"`
	IsException bool     `json:"is_exception"`
	ExceptionID string   `json:"exception_id,omitempty"`
}

// Output is the top-level structure written to vocabulary.json.
type Output struct {
	GeneratedAt          string `json:"generated_at"`
	SubjectsUpdatedAt    string `json:"subjects_updated_at"`
	AssignmentsUpdatedAt string `json:"assignments_updated_at"`
	Words                []Word `json:"words"`
}

// posMapping maps WaniKani part_of_speech strings to our normalized tags.
var posMapping = map[string]string{
	"noun":              "noun",
	"numeral":           "numeral",
	"prefix":            "prefix",
	"suffix":            "suffix",
	"proper noun":       "proper_noun",
	"pronoun":           "pronoun",
	"adverb":            "adverb",
	"conjunction":       "conjunction",
	"expression":        "expression",
	"interjection":      "interjection",
	"counter":           "counter",
	"い adjective":       "i_adjective",
	"な adjective":       "na_adjective",
	"の adjective":       "no_adjective",
	"godan verb":        "godan_verb",
	"ichidan verb":      "ichidan_verb",
	"する verb":          "suru_verb",
	"transitive verb":   "transitive_verb",
	"intransitive verb": "intransitive_verb",
}

// NormalizePOS maps WaniKani part_of_speech values to our normalized tags.
func NormalizePOS(wanikaniPOS []string) []string {
	var result []string
	for _, pos := range wanikaniPOS {
		if mapped, ok := posMapping[pos]; ok {
			result = append(result, mapped)
		} else {
			// Unknown POS: normalize by lowercasing and replacing spaces with underscores
			normalized := normalizeUnknown(pos)
			if normalized != "" {
				result = append(result, normalized)
			}
		}
	}
	return result
}

func normalizeUnknown(s string) string {
	result := make([]byte, 0, len(s))
	for i := range len(s) {
		c := s[i]
		if c == ' ' {
			result = append(result, '_')
		} else if c >= 'A' && c <= 'Z' {
			result = append(result, c+32)
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}
