// Lexer is heavily inspired by chewxy's lexer from the lingo project: https://github.com/chewxy/lingo/blob/master/lexer/stateFn.go
package lexer

import (
	"io"
)

var eof rune = -1

func (l *Lexer) push(v LexType) {
	switch v {
	case Break:
		l.Output = append(l.Output, Lexeme{
			Value:   BreakToken,
			Row:     l.row,
			Col:     l.col,
			LexType: Break,
		})
	case Word:
		word := l.sb.String()
		lemma := l.lemmatizer.Lemma(word)
		// if !slices.Contains(stopWords, lemma) {
		l.Output = append(l.Output, Lexeme{
			Lemma:   lemma,
			Value:   word,
			Row:     l.row,
			Col:     l.col,
			LexType: Word,
		})
		// }
	}
	// TODO: finish this
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