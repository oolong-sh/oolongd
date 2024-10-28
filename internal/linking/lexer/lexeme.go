package lexer

var BreakToken = "__BREAK__"

type lexType byte

const (
	EOF lexType = iota
	Word
	URI
	Time
	Date
	Number
	Symbol
	Puntuation
	Space
	Break
)

// DOC:
type Lexeme struct {
	Lemma    string
	Value    string
	Location [2]int // location (useful for potential LSP implementation)
}
