package linking

// TODO:
// - tf, idf, tf-idf functions
// - also custom weighting based on special tags (i.e. markdown headers, bold text, latex sections)

func tf(d *Document) map[string]float32 {
	// TODO: f_td / (sum f_t'd)
	// f_td is occurence count of ngram in a document
	// - denom is sum of all counts -> could be substituted for total count of tokens in the document? -> easier implementation
	// -
	// - iterate through list of ngrams

	wgts := make(map[string]float32)

	total := float32(len(d.ngrams))

	for _, ng := range d.ngrams {
		// TODO: decide how to calculate tf by ngram, depends on ngram type
		if _, ok := wgts[ng.ngram]; ok {
			continue
		}

		wgts[ng.ngram] = float32(ng.count) / total
	}

	return wgts
}
