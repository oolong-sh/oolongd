package linking

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestReadDocumentSimple(t *testing.T) {
	// Basic test
	s := "Hello world!"
	reader := strings.NewReader(s)
	doc, err := readDocument(reader, "")
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}
	fmt.Println("Input:", s, "Output:", doc.tokens)
	if len(doc.tokens) != 2 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d\n", 2, len(doc.tokens))
	}
	expectedTokens := []token{
		{
			token:    "hello",
			location: 0,
		}, {
			token:    "world",
			location: 0,
		},
	}
	if !slices.Equal(doc.tokens, expectedTokens) {
		t.Fatalf("Incorrect order in output tokens slice. Expected %+v, got %+v\n", doc.tokens, expectedTokens)
	}

	// basic test with newlines (should contain `breakToken`)
	s = "Hello, \nworld!"
	reader = strings.NewReader(s)
	doc, err = readDocument(reader, "")
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}
	fmt.Println("Input:", s, " Output:", doc.tokens)
	if len(doc.tokens) != 3 {
		t.Fatalf("Incorrect Document.tokens length. Expected %d, got %d", 2, len(doc.tokens))
	}
	expectedTokens = []token{
		{
			token:    "hello",
			location: 0,
		},
		{
			token:    breakToken,
			location: 0,
		},
		{
			token:    "world",
			location: 1,
		},
	}
	if !slices.Equal(doc.tokens, expectedTokens) {
		t.Fatalf("Incorrect order in output tokens slice. Expected %+v, got %+v\n", doc.tokens, expectedTokens)
	}

	// test with many newlines and multiple single-line lexemes
	s = "\nHello, \nworld! Foo-bar baz   \n\n foo"
	reader = strings.NewReader(s)
	doc, err = readDocument(reader, "")
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}
	fmt.Println("Input:", s, " Output:", doc.tokens)
	if len(doc.tokens) != 9 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d", 5, len(doc.tokens))
	}
}

// TODO: tests with hyphens and other specially handled cases
