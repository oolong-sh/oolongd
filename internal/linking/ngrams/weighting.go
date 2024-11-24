package ngrams

import (
	"math"

	"github.com/oolong-sh/oolong/internal/linking/lexer"
)

// Weight adjustment multipliers for NGram sizes
var nadj map[int]float64 = map[int]float64{
	1: 0.6,
	2: 1.5,
	3: 1.4,
	4: 0.7,
	5: 0.6,
}

// Zone adjustments for k1 in bm25f calculations
var zoneK1 map[lexer.Zone]float64 = map[lexer.Zone]float64{
	lexer.Default: 1.1,
	lexer.H1:      2.2,
	lexer.H2:      2.0,
	lexer.H3:      1.9,
	lexer.H4:      1.8,
	lexer.H5:      1.7,
	lexer.Bold:    1.5,
	lexer.Italic:  1.3,
	lexer.Link:    1.2,
}

// Zone adjustments for b in bm25f calculations
var zoneB map[lexer.Zone]float64 = map[lexer.Zone]float64{
	lexer.Default: 0.9,
	lexer.H1:      0.2,
	lexer.H2:      0.3,
	lexer.H3:      0.35,
	lexer.H4:      0.5,
	lexer.H5:      0.5,
	lexer.Bold:    0.75,
	lexer.Italic:  0.75,
	lexer.Link:    0.8,
}

// Calculate weights for ngrams found in a map of NGrams across N documents
// To be used after all document NGram maps are merged
func CalcWeights(ngmap map[string]*NGram, N int) {
	idf(ngmap, N)
	// tfidf(ngmap)
	bm25(ngmap)

	for _, ng := range ngmap {
		ng.updateWeight()
	}
}

// Update weights for an individual NGram instance
func (ng *NGram) updateWeight() {
	w := 0.0
	df := 0.01

	// TODO: these numbers are subject to change
	// - document and count adjustments are too high for n=1
	ladj := math.Min(0.12*float64(len(ng.keyword)), 1.2)                 // length adjustment
	cadj := math.Min(0.08*float64(ng.n)*float64(ng.globalCount), 1.5)    // count adjustment
	dadj := math.Min(0.08*float64(ng.n)*float64(len(ng.documents)), 1.8) // document occurence adjustment
	// TODO: heavily prefer count / len(dg.documents) > 1
	// cdadj := math.Min(0.5*float64(ng.globalCount)/float64(len(ng.documents)), 2)

	adjustment := ladj * cadj * nadj[ng.n] * dadj // * cdadj

	for _, nginfo := range ng.documents {
		// documentWeight will be bm25 or tf-idf before this point
		// apply adjustment to document weight
		nginfo.DocumentWeight = nginfo.DocumentWeight * adjustment
		w += nginfo.DocumentWeight
		df++
	}

	ng.globalWeight = w * adjustment / df
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
