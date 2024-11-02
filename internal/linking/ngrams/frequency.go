package ngrams

import "math"

// Calculate term frequency (f_td / (sum f_t'd))
func tf(ngmap map[string]*NGram, path string) {
	totalCount := 0
	for _, ng := range ngmap {
		totalCount += ng.documents[path].DocumentCount
	}

	for _, ng := range ngmap {
		nginfo := ng.documents[path]
		nginfo.DocumentTermFreq = float64(nginfo.DocumentCount) / float64(totalCount)
	}
}

// Calculate inverse document frequency of NGrams
// N is the total number of documents in the text corpus
//
// This calculation must happen after all document NGram maps are merged
func idf(ngmap map[string]*NGram, N int) {
	for _, ng := range ngmap {
		ng.inverseDocFreq = math.Log(float64(N) / float64(1+len(ng.documents)))
	}
}

// Calculate term frequency-inverse document frequency of NGrams
// Requires both TF and IDF to be computed beforehand
func tfidf(ngmap map[string]*NGram) {
	for _, ng := range ngmap {
		for _, nginfo := range ng.documents {
			nginfo.DocumentTfIdf = nginfo.DocumentTermFreq * ng.inverseDocFreq
		}
	}
}
