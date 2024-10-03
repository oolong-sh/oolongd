package document

import (
	"fmt"
	"strings"
	"testing"
)

func TestReadDocumentSimple(t *testing.T) {
	s := "Hello world!"
	reader := strings.NewReader(s)
	doc, err := ReadDocument(reader, "")
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}
	fmt.Println("Input:", s, "Output:", doc.Contents)
	if len(doc.Contents) != 2 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d", 2, len(doc.Contents))
	}

	s = "Hello, \nworld!"
	reader = strings.NewReader(s)
	doc, err = ReadDocument(reader, "")
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}
	fmt.Println("Input:", s, " Output:", doc.Contents)
	if len(doc.Contents) != 2 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d", 2, len(doc.Contents))
	}

	s = "\nHello, \nworld! Foo-bar baz   \n\n foo"
	reader = strings.NewReader(s)
	doc, err = ReadDocument(reader, "")
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}
	fmt.Println("Input:", s, " Output:", doc.Contents)
	if len(doc.Contents) != 5 {
		t.Fatalf("Incorrect Document.Content length. Expected %d, got %d", 5, len(doc.Contents))
	}
}
