package linking

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// TODO: multi-document read
// TODO: multi-document ngram merge (maps.copy(dest, src))

var allowedExts = []string{".md", ".mdx", ".tex", ".typ", ".txt"}

// DOC:
func ReadNotesDir(notesDirPaths ...string) ([]*Document, error) {
	documents := []*Document{}

	for _, notesDirPath := range notesDirPaths {
		// expand home dir shorthand
		if strings.HasPrefix(notesDirPath, "~/") || notesDirPath == "~" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			if notesDirPath == "~" {
				notesDirPath = home
			}
			notesDirPath = filepath.Join(home, notesDirPath[2:])
		}

		// extract all note file paths from notes directory
		notePaths := []string{}
		if err := filepath.WalkDir(notesDirPath, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			// REFACTOR: probably change this to a blacklist
			// disallow binaries and allow users to blacklist filetypes
			if slices.Contains(allowedExts, filepath.Ext(path)) {
				notePaths = append(notePaths, path)
			}

			return nil
		}); err != nil {
			return nil, err
		}

		// TODO: this could probably be parallelized? --> mutex on documents array
		for _, notePath := range notePaths {
			doc, err := ReadDocument(notePath)
			if err != nil {
				return nil, err
			}
			documents = append(documents, doc)
		}
	}

	return documents, nil
}
