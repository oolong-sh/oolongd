package lexer

func (l *Lexer) push(v lexType) {
	switch v {
	case Break:
		l.Output = append(l.Output, Lexeme{
			Value:    BreakToken,
			Location: [2]int{l.row, l.col},
		})
	case Word:
		word := l.sb.String()
		lemma := l.lemmatizer.Lemma(word)
		// if !slices.Contains(stopWords, lemma) {
		l.Output = append(l.Output, Lexeme{
			Lemma:    lemma,
			Value:    word,
			Location: [2]int{l.row, l.col},
		})
		// }
	}
	// TODO: finish this
}
