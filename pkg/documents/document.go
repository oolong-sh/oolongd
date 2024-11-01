package documents

type Document interface {
	Path() string
	KeywordWeights() map[string]float32
}
