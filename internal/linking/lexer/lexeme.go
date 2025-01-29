package lexer

var BreakToken = "__BREAK__"

type LexType byte

// Lexeme type enum for use in lexer
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

// Zone enum used in weights calculations
// (Ordered from most significant to least significant)
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

type Lexeme struct {
	Lemma   string  // lexical root of unit (i.e. continues -> continue)
	Value   string  // lexical unit
	LexType LexType // type of lexical unit
	Zone    Zone    // document zone
}
