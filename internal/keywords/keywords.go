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

type SearchKeyword struct {
	Weight    float64                      `json:"weight"`
	Count     int                          `json:"count"`
	Documents map[string]*ngrams.NGramInfo `json:"documents"`
}

func SerializeNGrams(ngmap map[string]*ngrams.NGram) ([]byte, error) {
	keywords := NGramsToKeywords(ngmap)

	b, err := json.Marshal(keywords)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func SearchByKeyword(s string, ngmap map[string]*ngrams.NGram) (SearchKeyword, bool) {
	ng, exist := ngmap[s]
	if !exist {
		return SearchKeyword{}, false
	}

	return SearchKeyword{
		Weight:    ng.Weight(),
		Count:     ng.Count(),
		Documents: ng.Documents(),
	}, true
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

// Note: NGramsToKeywordsMap is being used by the API
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
