package ngrams

import (
	"math"
)

// Calculate term frequency
func tf(ngmap map[string]*NGram, path string) {
	// totalCount := 0
	// for _, ng := range ngmap {
	// 	totalCount += ng.documents[path].DocumentCount
	// }

	for _, ng := range ngmap {
		nginfo := ng.documents[path]
		// normalize by document token count
		// nginfo.DocumentTF = float64(nginfo.DocumentCount) / float64(totalCount)
		nginfo.DocumentTF = float64(nginfo.DocumentCount)
	}
}

// Calculate inverse document frequency of NGrams
// N is the total number of documents in the text corpus
//
// This calculation must happen after all document NGram maps are merged
func idf(ngmap map[string]*NGram, N int) {
	for _, ng := range ngmap {
		// normal idf
		// ng.inverseDocFreq = math.Log(float64(N) / float64(1+len(ng.documents)))

		// smoothed idf
		// ng.idf = math.Log(1 + float64(N)/float64(1+len(ng.documents)))

		// okapi idf
		n := float64(len(ng.documents))
		ng.idf = math.Log((0.5+float64(N)-n)/(0.5+n) + 1)
	}
}

// Calculate term frequency-inverse document frequency of NGrams
// Requires both TF and IDF to be computed beforehand
func tfidf(ngmap map[string]*NGram) {
	for _, ng := range ngmap {
		for _, nginfo := range ng.documents {
			nginfo.DocumentTfIdf = nginfo.DocumentTF * ng.idf
		}
	}
}

// Best Matching 25 -- Alternative matching function that doesn't downweight common terms as much
// k1: controls saturation of TF (normally between 1.2 and 2)
// b: controls document length normalization (0 is no normaliztion)
// TODO: add bm25f modifications to account for zones -- add zone tracking to lexer (zones affect b, k1, idf)
func bm25(ngmap map[string]*NGram) {
	d := make(map[string]float64)
	totalLength := 0.0

	// calculate document lengths and avg
	for _, ng := range ngmap {
		for path, nginfo := range ng.documents {
			if _, ok := d[path]; !ok {
				d[path] = 0
			}
			d[path] += float64(nginfo.DocumentCount)
			totalLength += float64(nginfo.DocumentCount)
		}
	}
	davg := totalLength / float64(len(d))

	// calculate bm25 per ngram per document
	var b, k1 float64
	for _, ng := range ngmap {
		b = zoneB[ng.zone]
		k1 = zoneK1[ng.zone]
		for path, nginfo := range ng.documents {
			nginfo.DocumentBM25 = ng.idf * ((nginfo.DocumentTF * (k1 + 1)) / (nginfo.DocumentTF + k1*(1-b+b*(d[path]/davg))))
		}
	}
}
