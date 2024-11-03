package ngrams

// TODO: finish this function
func CalcWeights(ngmap map[string]*NGram, N int) {
	idf(ngmap, N)
	// tfidf(ngmap)
	// TODO: decide on k and b values (and allow them to be tweaked from config)
	bm25(ngmap, 1.2, 0.8)
	// CHANGE: probably take n and word length into account
	for _, ng := range ngmap {
		ng.updateWeight()
	}
}

// DOC:
func (ng *NGram) updateWeight() {
	w := 0.0
	df := 0.0
	for _, nginfo := range ng.documents {
		// TODO: decide on weighting metric
		// w += nginfo.DocumentTfIdf
		w += nginfo.DocumentBM25
		df++
	}
	// TODO: this probably shouldn't use an average, but rather a more advanced metric
	// that takes counts bm25 values into account
	// - also needs to take n-size into account (and possibly number of characters)

	ng.globalWeight = w / df
}

// TODO: maybe get rid of this function?
func FilterMeaningfulNGrams(ngmap map[string]*NGram, minDF int, maxDF int, minAvgTFIDF float64) []string {
	var out []string
	for k, ng := range ngmap {
		N := len(ng.documents)
		if N >= minDF && N <= maxDF && ng.globalWeight >= minAvgTFIDF {
			out = append(out, k)
		}
	}
	return out
}
