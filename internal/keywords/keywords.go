package keywords

import (
	"encoding/json"

	"github.com/oolong-sh/oolong/internal/linking/ngrams"
)

type Keyword struct {
	Keyword string  `json:"keyword"`
	Weight  float64 `json:"weight"`
}

// DOC:
func SerializeNGrams(ngmap map[string]*ngrams.NGram) ([]byte, error) {
	keywords := NGramsToKeywords(ngmap)

	b, err := json.Marshal(keywords)
	if err != nil {
		return nil, err
	}

	return b, nil
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
