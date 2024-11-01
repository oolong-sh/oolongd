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
	Lemma   string
	Value   string
	Row     int
	Col     int
	LexType LexType
}
