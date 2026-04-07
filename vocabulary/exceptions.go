package vocabulary

// ExceptionInfo describes an irregular conjugation pattern.
type ExceptionInfo struct {
	ID string
}

// Exceptions maps characters to their exception info.
// These are words with irregular conjugation patterns.
var Exceptions = map[string]ExceptionInfo{
	// いい uses よ- stem for all conjugations except affirmative present
	"いい": {ID: "ii"},
	// かっこいい is a compound with いい; uses かっこよ- stem
	"かっこいい": {ID: "kakkoii"},
}
