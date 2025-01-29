// Lexer is inspired by chewxy's lexer from the lingo project: https://github.com/chewxy/lingo/blob/master/lexer/stateFn.go
package lexer

import (
	"io"
	"regexp"
	"strings"
)

var eof rune = -1

var (
	// Heading (e.g., # Heading, ## Heading) - only matches standalone heading lines
	// (?m) is required to allow matching from start/end of line rather than start/end of string
	// FIX: these capture groups are sometimes wrapping around lines (probably abandon regex and use more advanced lexer logic)
	h1Pattern = regexp.MustCompile(`(?m)^(#)\s+(.+?)\s*$`)
	h2Pattern = regexp.MustCompile(`(?m)^(#{2})\s+(.+?)\s*$`)
	h3Pattern = regexp.MustCompile(`(?m)^(#{3})\s+(.+?)\s*$`)
	h4Pattern = regexp.MustCompile(`(?m)^(#{4})\s+(.+?)\s*$`)
	h5Pattern = regexp.MustCompile(`(?m)^(#{5})\s+(.+?)\s*$`)
	// Bold text (e.g., **bold** or __bold__) - matches inline without lookaheads/behinds
	boldPattern = regexp.MustCompile(`\*\*(.+?)\*\*|__(.+?)__`)
	// Italic text (e.g., *italic* or _italic_) - matches inline without lookaheads/behinds
	italicPattern = regexp.MustCompile(`(?:^|[^\w])(\*(\w+?)\*|_(\w+?)_)(?:[^\w]|$)`)
	// Link (e.g., [text](url))
	linkPattern = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	// Image (e.g., ![alt text](url))
	imagePattern = regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
)

func (l *Lexer) push(v LexType) {
	switch v {
	case Break:
		l.Output = append(l.Output, Lexeme{
			Value:   BreakToken,
			LexType: Break,
			Zone:    l.zone,
		})
	case Word:
		word := l.sb.String()
		lemma := l.lemmatizer.Lemma(word)

		// handle lemmatizer bug where first dict entry contains <feff> character
		lemma = strings.TrimPrefix(lemma, "\ufeff")

		l.Output = append(l.Output, Lexeme{
			Lemma:   lemma,
			Value:   word,
			LexType: Word,
			Zone:    l.zone,
		})
	}
	// TODO: handle other types as necessary (mainly urls)
}

func (l *Lexer) next() rune {
	var err error
	l.r, l.width, err = l.br.ReadRune()
	if err == io.EOF {
		l.width = 1
		return eof
	}
	l.col += l.width
	l.pos += l.width
	// l.col++
	// l.pos++

	return l.r
}

func (l *Lexer) backup() {
	l.br.UnreadRune()
	l.pos -= l.width
	l.col -= l.width
}

func (l *Lexer) peek() rune {
	backup := l.r
	pos := l.pos
	col := l.col

	r := l.next()
	l.backup()

	l.r = backup
	l.pos = pos
	l.col = col
	return r
}

func (l *Lexer) accept() {
	l.sb.WriteRune(l.r)
}

func (l *Lexer) ignore() {
	l.start = l.pos
	l.sb.Reset()
}

func (l *Lexer) detectZone() {
	peekBuffer, _ := l.br.Peek(32)

	switch {
	// TODO: handle remaining cases
	// - add capture group for code blocks (might just need a boolean flag for them)
	case h1Pattern.Match(peekBuffer):
		l.zone = H1
	case h2Pattern.Match(peekBuffer):
		l.zone = H2
	case h3Pattern.Match(peekBuffer):
		l.zone = H3
	case h4Pattern.Match(peekBuffer):
		l.zone = H4
	case h5Pattern.Match(peekBuffer):
		l.zone = H5
	// case sectionPattern.Match(peekBuffer):
	// 	l.zone = Default
	// case boldPattern.Match(peekBuffer):
	// 	l.zone = Bold
	// case italicPattern.Match(peekBuffer):
	// 	l.zone = Italic
	// case inlineMathPattern.Match(peekBuffer):
	// 	l.zone = Math
	// case displayMathPattern.Match(peekBuffer):
	// 	l.zone = "math-display"
	default:
		l.zone = Default
	}
}
