package wanikani

// Collection represents a paginated WaniKani API response.
type Collection[T any] struct {
	Object        string `json:"object"`
	URL           string `json:"url"`
	Pages         Pages  `json:"pages"`
	TotalCount    int    `json:"total_count"`
	DataUpdatedAt string `json:"data_updated_at"`
	Data          []T    `json:"data"`
}

// Pages contains pagination URLs for a collection.
type Pages struct {
	NextURL     *string `json:"next_url"`
	PreviousURL *string `json:"previous_url"`
	PerPage     int     `json:"per_page"`
}

// Resource wraps an individual WaniKani API resource.
type Resource[T any] struct {
	ID            int    `json:"id"`
	Object        string `json:"object"`
	URL           string `json:"url"`
	DataUpdatedAt string `json:"data_updated_at"`
	Data          T      `json:"data"`
}

// SubjectData represents vocabulary or kana_vocabulary subject data.
type SubjectData struct {
	Characters   string    `json:"characters"`
	Meanings     []Meaning `json:"meanings"`
	Readings     []Reading `json:"readings"`
	PartOfSpeech []string  `json:"part_of_speech"`
	Level        int       `json:"level"`
	Slug         string    `json:"slug"`
}

// Meaning represents a single meaning entry.
type Meaning struct {
	Meaning        string `json:"meaning"`
	Primary        bool   `json:"primary"`
	AcceptedAnswer bool   `json:"accepted_answer"`
}

// Reading represents a single reading entry.
type Reading struct {
	Reading        string `json:"reading"`
	Primary        bool   `json:"primary"`
	AcceptedAnswer bool   `json:"accepted_answer"`
}

// AssignmentData represents a user's assignment (progress) for a subject.
type AssignmentData struct {
	SubjectID   int     `json:"subject_id"`
	SubjectType string  `json:"subject_type"`
	SRSStage    int     `json:"srs_stage"`
	UnlockedAt  *string `json:"unlocked_at"`
	StartedAt   *string `json:"started_at"`
	PassedAt    *string `json:"passed_at"`
	BurnedAt    *string `json:"burned_at"`
}
