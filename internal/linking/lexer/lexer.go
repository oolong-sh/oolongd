package lexer

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
)

var allowedSpecialChars = []rune{'-', '_', '\''}

type Lexer struct {
	br *bufio.Reader

	r     rune
	width int
	pos   int
	start int
	row   int
	col   int

	lemmatizer *golem.Lemmatizer
	sb         strings.Builder

	Output []Lexeme
}

// DOC:
func New() *Lexer {
	lemmatizer, err := golem.New(en.New())
	if err != nil {
		panic(fmt.Sprintf("failed to initialize lemmatizer: %v", err))
	}

	return &Lexer{
		pos:        1,
		start:      1,
		row:        1,
		col:        1,
		lemmatizer: lemmatizer,
		Output:     []Lexeme{},
	}
}

// DOC:
// NOTE: could rewrite with regex instead of hardcoded special cases
func (l *Lexer) Lex(r io.Reader, stage int) error {
	l.br = bufio.NewReader(r)

	for {
		r := l.next()
		if r == eof {
			break
		}

		// TODO: figure out where lowercasing should be happening
		c := processChar(r, stage)

		// switch {
		// case unicode.IsSpace(r):
		// 	// TODO:
		// case unicode.IsDigit(r):
		// 	// TODO:
		// case r == ':':
		// 	// TODO:
		// case unicode.IsPunct(r):
		// 	// TODO:
		// case unicode.IsSymbol(r):
		// 	// TODO:
		// }

		// REFACTOR: change to switch case, possibly use a state machine
		if c == 0 {
			if l.sb.Len() > 0 {
				l.push(Word)
				// REFACTOR:
				// l.sb.Reset()
				l.ignore()
			}

			if r == '\n' {
				// TODO: call push correctly
				l.push(Break)
				l.row++
			}
		} else {
			// TODO: replace with peek
			if c == '-' {
				n := l.peek()
				switch {
				case n == eof:
					l.width = 1
					l.backup()
					l.width = 0
					// t = Symbol
				case unicode.IsLetter(n):
					l.accept()
				}
			} else {
				l.accept()
			}
		}
	}

	// Handle any remaining content in the buffer
	if l.sb.Len() > 0 {
		l.push(Word) // CHANGE: needs to be able to handle the other types as well
	}

	return nil
}

// DOC:
func processChar(c rune, stage int) rune {
	if unicode.IsSpace(c) {
		return 0
	}

	// stage 0
	if stage == 0 {
		return c
	}

	// stage 1+
	if stage > 0 {
		c = unicode.ToLower(c)
	}

	// stage 2+
	if stage > 1 {
		if unicode.IsLetter(c) || unicode.IsNumber(c) || slices.Contains(allowedSpecialChars, c) {
			c = unicode.ToLower(c)
		} else {
			return 0
		}
	}

	return c
}