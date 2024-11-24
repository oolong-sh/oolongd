package keywords

import (
	"encoding/json"

	"github.com/oolong-sh/oolong/internal/config"
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

func NGramsToKeywordsMap(ngmap map[string]*ngrams.NGram) map[string]Keyword {
	keywords := map[string]Keyword{}
	minThresh := config.WeightThresholds().MinNodeWeight

	for k, v := range ngmap {
		w := v.Weight()

		if w > minThresh {
			keywords[k] = Keyword{
				Keyword: k,
				Weight:  w,
			}
		}
	}

	return keywords
}

// Note: NGramsToKeywords is being used
func NGramsToKeywords(ngmap map[string]*ngrams.NGram) []Keyword {
	keywords := []Keyword{}
	minThresh := config.WeightThresholds().MinNodeWeight

	for k, v := range ngmap {
		w := v.Weight()

		if w > minThresh {
			keywords = append(keywords, Keyword{
				Keyword: k,
				Weight:  w,
			})
		}
	}

	return keywords
}
