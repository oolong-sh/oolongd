package keywords

type Keyword interface {
	Keyword() string
	Weight() float32
}
