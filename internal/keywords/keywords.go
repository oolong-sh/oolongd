package keywords

import (
	"encoding/json"
	"os"

	"github.com/oolong-sh/oolong/internal/linking/ngrams"
)

var keywordsFile = "./oolong-keywords.json"

type Keyword struct {
	Keyword string  `json:"keyword"`
	Weight  float64 `json:"weight"`
}

// DOC:
func SerializeNGrams(ngmap map[string]*ngrams.NGram) error {
	keywords := NGramsToKeywords(ngmap)

	err := serializeKeywords(keywords)
	if err != nil {
		return err
	}

	return nil
}

func serializeKeywords(keywords []Keyword) error {
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
func NGramsToKeywords(ngmap map[string]*ngrams.NGram) []Keyword {
	// keywords := make([]keyword, len(ngmap))
	keywords := []Keyword{}
	threshold := 8.0

	for k, v := range ngmap {
		w := v.Weight()

		if w > threshold {
			keywords = append(keywords, Keyword{
				Keyword: k,
				Weight:  w,
			})
		}
	}

	return keywords
}

func NGramsToKeywordsMap(ngmap map[string]*ngrams.NGram) map[string]Keyword {
	keywords := map[string]Keyword{}
	threshold := 8.0

	for k, v := range ngmap {
		w := v.Weight()

		if w > threshold {
			keywords[k] = Keyword{
				Keyword: k,
				Weight:  w,
			}
		}
	}

	return keywords
}
