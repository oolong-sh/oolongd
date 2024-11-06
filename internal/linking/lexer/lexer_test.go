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

var cfg config.OolongConfig

func init() {
	var err error
	cfg, err = config.Setup("../../../examples/oolong.json")
	if err != nil {
		panic(err)
	}
}

func TestReadDocumentSimple(t *testing.T) {
	// Basic test
	s := "Hello world!"
	var rd io.Reader = strings.NewReader(s)
	l := lexer.New()
	l.Lex(rd)

	fmt.Println("Input:", s, "Output:", l.Output)
	if len(l.Output) != 2 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d\n", 2, len(l.Output))
	}
	expectedTokens := []lexer.Lexeme{
		{
			Lemma: "hello",
			Value: "Hello",
			// Row:     1,
			// Col:     1,
			LexType: lexer.Word,
			Zone:    lexer.Default,
		}, {
			Lemma: "world",
			Value: "world",
			// Row:     1,
			// Col:     7,
			LexType: lexer.Word,
			Zone:    lexer.Default,
		},
	}
	if !slices.Equal(l.Output, expectedTokens) {
		t.Fatalf("Unexpected lexer output. Expected %+v, got %+v\n", expectedTokens, l.Output)
	}

	// basic test with newlines (should contain `breakToken`)
	s = "Hello, \nworld!"
	rd = strings.NewReader(s)
	l = lexer.New()
	l.Lex(rd)
	fmt.Println("Input:", s, " Output:", l.Output)
	if len(l.Output) != 3 {
		t.Fatalf("Incorrect output length. Expected %d, got %d", 3, len(l.Output))
	}
	expectedTokens = []lexer.Lexeme{
		{
			Lemma: "hello",
			Value: "Hello",
			// Row:     1,
			// Col:     1,
			LexType: lexer.Word,
			Zone:    lexer.Default,
		},
		{
			Value: lexer.BreakToken,
			// Row:     1,
			// Col:     8,
			LexType: lexer.Break,
			Zone:    lexer.Default,
		},
		{
			Lemma: "world",
			Value: "world",
			// Row:     2,
			// Col:     1,
			LexType: lexer.Word,
			Zone:    lexer.Default,
		},
	}
	if !slices.Equal(l.Output, expectedTokens) {
		t.Fatalf("Unexpected lexer output. Expected %+v, got %+v\n", expectedTokens, l.Output)
	}

	// test with many newlines and multiple single-line lexemes
	s = "\nHello, \nworld! Foo-bar baz   \n\n foo"
	rd = strings.NewReader(s)
	l = lexer.New()
	l.Lex(rd)
	fmt.Println("Input:", s, " Output:", l.Output)
	if len(l.Output) != 9 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d", 5, len(l.Output))
	}
}

// TODO: tests with hyphens and other specially handled cases
