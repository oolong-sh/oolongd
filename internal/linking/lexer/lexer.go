// Lexer is heavily inspired by chewxy's lexer from the lingo project: https://github.com/chewxy/lingo/blob/master/lexer/stateFn.go
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

	zone       Zone
	lemmatizer *golem.Lemmatizer
	sb         strings.Builder

	Output []Lexeme
}

// Initialize a new lexer
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
		zone:       Default,
		lemmatizer: lemmatizer,
		Output:     []Lexeme{},
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
			// TODO:
			if l.sb.Len() > 0 {
				l.push(Word)
				l.ignore()
			}
			if r == '\n' {
				isNewline = true
				l.push(Break)
				l.row++
				l.col = 1
			}
		case unicode.IsDigit(r):
			// TODO: number handling?
			l.accept()
		case r == ':':
			// get this case from lingo
			if l.peek() == '/' {
				// l.accept() // accept ':'
				// l.next()
				// if l.peek() == '/' {
				// 	l.accept()
				// 	return lexURI
				// }
				// // otherwise...
				// l.backup()
				// // "unaccept". since '/' has a width of 1 we can do the following
				// l.buf.Truncate(l.buf.Len() - 1)
			}
			// fn = lexPunctuation
			// TODO: possible url handling?
		case unicode.IsPunct(r):
			// TODO:
			switch r {
			case '_':
				l.accept()
			case '-':
				n := l.peek()
				switch {
				case n == eof:
					l.width = 1
					l.backup()
					l.width = 0
					// t = Symbol
					// TODO: something?
				case unicode.IsLetter(n):
					l.accept()
				}
			default:
				// l.ignore()
			}
		case unicode.IsSymbol(r):
			// TODO: non-punct processing
			l.ignore()
		default:
			l.accept()
		}
	}

	// Handle any remaining content in the buffer
	if l.sb.Len() > 0 {
		l.push(Word) // CHANGE: needs to be able to handle the other types as well
	}
}
