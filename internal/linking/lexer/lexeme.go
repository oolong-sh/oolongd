package lexer

var BreakToken = "__BREAK__"

type LexType byte

const (
	EOF LexType = iota
	Word
	URI
	Time
	Date
	Number
	Symbol
	Punctuation
	Space
	Break
)

// DOC:
type Lexeme struct {
	Lemma    string
	Value    string
	Location [2]int // location (useful for potential LSP implementation)
	LexType  LexType
}
