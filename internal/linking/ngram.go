package linking

// DOC:
type NGram struct {
	ngram  string
	weight float32 // weight of NGram across all documents
	count  int     // count across all documents  NOTE: possibly replace this with a map of ngram->int
	n      int

	// store all documents ngram is present in and counts within the document
	// - maps document path to count of ngram in the document
	documents map[string]int

	// TODO: store per-weight documents here?
}

// NGram implements Keyword interface
func (d *NGram) Weight() float32 { return d.weight }
func (d *NGram) Keyword() string { return d.ngram }

// TODO: ngram calculation functions
