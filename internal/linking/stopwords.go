package linking

const breakToken = "__BREAK__"

// TODO: add more stop words
var stopWords = []string{
	"a", "an", "and", "are", "as", "at", "be", "by", "for", "from", "has",
	"have", "had", "he", "her", "here", "hers", "herself", "him", "himself",
	"his", "how", "i", "if", "in", "into", "it", "its", "itself", "just",
	"like", "me", "might", "more", "most", "must", "my", "myself", "no",
	"not", "now", "of", "off", "on", "once", "only", "or", "other", "our",
	"ours", "ourselves", "out", "over", "own", "re", "s", "same", "she",
	"should", "so", "some", "such", "t", "than", "that", "the", "their",
	"theirs", "them", "themselves", "then", "there", "these", "they",
	"this", "those", "through", "to", "too", "under", "until", "up", "ve",
	"very", "was", "wasn", "we", "were", "what", "when", "where", "which",
	"while", "who", "whom", "why", "will", "with", "won", "would", "y",
	"you", "your", "yours", "yourself", "yourselves",
}
