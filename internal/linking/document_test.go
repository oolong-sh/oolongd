package linking

import (
	"fmt"
	"strings"
	"testing"
)

func TestReadDocumentSimple(t *testing.T) {
	// Basic test
	s := "Hello world!"
	reader := strings.NewReader(s)
	doc, err := readDocument(reader)
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}
	fmt.Println("Input:", s, "Output:", doc.contents)
	if len(doc.contents) != 2 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d", 2, len(doc.contents))
	}

	// basic test with newlines
	s = "Hello, \nworld!"
	reader = strings.NewReader(s)
	doc, err = readDocument(reader)
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}
	fmt.Println("Input:", s, " Output:", doc.contents)
	if len(doc.contents) != 2 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d", 2, len(doc.contents))
	}

	// test with many newlines and multiple single-line lexemes
	s = "\nHello, \nworld! Foo-bar baz   \n\n foo"
	reader = strings.NewReader(s)
	doc, err = readDocument(reader)
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}
	fmt.Println("Input:", s, " Output:", doc.contents)
	if len(doc.contents) != 5 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d", 5, len(doc.contents))
	}
}
