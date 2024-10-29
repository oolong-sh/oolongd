package lexer_test

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/linking/lexer"
)

var (
	cfg   config.OolongConfig
	stage = 3
)

func init() {
	var err error
	cfg, err = config.Setup("")
	if err != nil {
		panic(err)
	}
}

func TestReadDocumentSimple(t *testing.T) {
	// Basic test
	s := "Hello world!"
	var rd io.Reader = strings.NewReader(s)
	l := lexer.New()
	err := l.Lex(rd, stage)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	fmt.Println("Input:", s, "Output:", l.Output)
	if len(l.Output) != 2 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d\n", 2, len(l.Output))
	}
	expectedTokens := []lexer.Lexeme{
		{
			Value:   "hello",
			Row:     1,
			Col:     1,
			LexType: lexer.Word,
		}, {
			Value:   "world",
			Row:     1,
			Col:     7,
			LexType: lexer.Word,
		},
	}
	if !slices.Equal(l.Output, expectedTokens) {
		t.Fatalf("Unexepcted lexer output. Expected %+v, got %+v\n", expectedTokens, l.Output)
	}

	// basic test with newlines (should contain `breakToken`)
	s = "Hello, \nworld!"
	rd = strings.NewReader(s)
	l.Lex(rd, stage)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	fmt.Println("Input:", s, " Output:", l.Output)
	if len(l.Output) != 3 {
		t.Fatalf("Incorrect Document.tokens length. Expected %d, got %d", 2, len(l.Output))
	}
	expectedTokens = []lexer.Lexeme{
		{
			Value:   "hello",
			Row:     1,
			Col:     1,
			LexType: lexer.Word,
		},
		{
			Value:   lexer.BreakToken,
			Row:     1,
			Col:     8,
			LexType: lexer.Break,
		},
		{
			Value:   "world",
			Row:     2,
			Col:     1,
			LexType: lexer.Word,
		},
	}
	if !slices.Equal(l.Output, expectedTokens) {
		t.Fatalf("Unexepcted lexer output. Expected %+v, got %+v\n", expectedTokens, l.Output)
	}

	// test with many newlines and multiple single-line lexemes
	s = "\nHello, \nworld! Foo-bar baz   \n\n foo"
	rd = strings.NewReader(s)
	l.Lex(rd, stage)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	fmt.Println("Input:", s, " Output:", l.Output)
	if len(l.Output) != 9 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d", 5, len(l.Output))
	}
}

// TODO: tests with hyphens and other specially handled cases
