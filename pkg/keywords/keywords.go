package keywords

import (
	"encoding/json"
	"os"

	"github.com/oolong-sh/oolong/internal/linking/ngrams"
)

var keywordsFile = "./oolong-keywords.json"

type keyword struct {
	Keyword string  `json:"keyword"`
	Weight  float64 `json:"weight"`
}

// DOC:
func SerializeNGrams(ngmap map[string]*ngrams.NGram) error {
	keywords := ngramsToKeywords(ngmap)

	err := serializeKeywords(keywords)
	if err != nil {
		return err
	}

	return nil
}

func serializeKeywords(keywords []keyword) error {
	b, err := json.Marshal(keywords)
	if err != nil {
		return err
	}

	err = os.WriteFile(keywordsFile, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

// TODO: parameterize filtering threshold (maybe a percentage?)
func ngramsToKeywords(ngmap map[string]*ngrams.NGram) []keyword {
	// keywords := make([]keyword, len(ngmap))
	keywords := []keyword{}
	threshold := 8.0

	for k, v := range ngmap {
		w := v.Weight()

		if w > threshold {
			keywords = append(keywords, keyword{
				Keyword: k,
				Weight:  w,
			})
		}
	}

	return keywords
}
