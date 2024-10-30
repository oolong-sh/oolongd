package documents

type Document interface {
	Path() string
	Keywords() map[string]float32
}
