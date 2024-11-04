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

type Zone byte

// DOC:
const (
	H1 Zone = iota
	H2
	H3
	H4
	Bold
	H5
	Italic
	Link
	Math
	Default
)

// DOC:
type Lexeme struct {
	Lemma   string
	Value   string
	Row     int
	Col     int
	LexType LexType
	Zone    Zone
}
