// Lexer is inspired by chewxy's lexer from the lingo project: https://github.com/chewxy/lingo/blob/master/lexer/stateFn.go
package lexer

import (
	"bufio"
	"fmt"
	"io"
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

	zone        Zone
	lemmatizer  *golem.Lemmatizer
	sb          strings.Builder
	currLexType LexType

	Output []Lexeme
}

// Initialize a new lexer
func New() *Lexer {
	lemmatizer, err := golem.New(en.New())
	if err != nil {
		panic(fmt.Sprintf("failed to initialize lemmatizer: %v", err))
	}

	return &Lexer{
		pos:         1,
		start:       1,
		row:         1,
		col:         1,
		zone:        Default,
		lemmatizer:  lemmatizer,
		currLexType: Space,
		Output:      []Lexeme{},
	}
}

// Run lexer on a body of text (passed in as an io.reader for generalized handling)
// NOTE: could rewrite with regex instead of hardcoded special cases
func (l *Lexer) Lex(r io.Reader) {
	l.br = bufio.NewReader(r)
	isNewline := true
	for {
		r := l.next()
		if r == eof {
			break
		}

		if isNewline {
			l.detectZone()
			isNewline = false
		}

		switch {
		case unicode.IsSpace(r):
			if l.sb.Len() > 0 {
				l.push(l.currLexType)
				l.ignore()
			}
			if r == '\n' {
				isNewline = true
				l.push(Break)
				l.row++
				l.col = 1
			}
			l.currLexType = Space
		case unicode.IsDigit(r):
			if l.currLexType == Space {
				l.currLexType = Number
			}
			// TODO: additional number handling? (i.e. Dates?)
			l.accept()
		case r == ':':
			if l.peek() == '/' {
				l.accept()
				l.next()
				if l.peek() == '/' {
					l.currLexType = URI
					l.accept()
				} else {
					l.sb.Reset()
				}
			}
		case unicode.IsPunct(r):
			if l.currLexType == Space {
				l.currLexType = Symbol
			}

			switch r {
			// case '_':
			// TODO: should underscores be allowed as first character?
			case '-':
				if l.sb.Len() == 0 {
					l.ignore()
				} else {
					n := l.peek()
					switch {
					case n == eof:
						l.width = 1
						l.backup()
						l.width = 0
					case unicode.IsLetter(n):
						l.accept()
					}
				}
			}
		case unicode.IsSymbol(r):
			// TODO: non-punct processing
			l.ignore()
		default:
			// If no other cases were met, character should be part of a word
			if l.currLexType != URI {
				l.currLexType = Word
			}
			l.accept()
		}
	}

	// Handle any remaining content in the buffer
	if l.sb.Len() > 0 {
		l.push(l.currLexType)
	}
}
