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
	br    *bufio.Reader
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
	return &Lexer{
		pos:    1,
		start:  1,
		row:    1,
		col:    1,
		Output: []Lexeme{},
	}
}

// DOC:
// NOTE: could rewrite with regex instead of hardcoded special cases
func (l *Lexer) Lex(r io.Reader, stage int) error {
	hyphenFlag := false // REFACTOR: remove

	l.br = bufio.NewReader(r)

	var err error
	l.lemmatizer, err = golem.New(en.New())
	if err != nil {
		return fmt.Errorf("failed to initialize lemmatizer: %v", err)
	}

	for {
		r, _, err := l.br.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		c := processChar(r, stage)
		if c == 0 {
			// REFACTOR: change hyphenflag to a peek/lookback system
			// - would help with special char handling, particularly with uris
			hyphenFlag = false
			if l.sb.Len() > 0 {
				l.push(Word)
				l.sb.Reset()
			}

			if r == '\n' {
				// TODO: call push
				l.push(Break)
				l.row++
			}
		} else {
			if c == '-' {
				hyphenFlag = true
				continue
			}

			if hyphenFlag && l.sb.Len() > 0 {
				hyphenFlag = false
				l.sb.WriteRune('-')
			}

			l.sb.WriteRune(c)
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
