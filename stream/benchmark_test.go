package stream

import (
	"strings"
	"testing"
)

func mapper(word string) string {
	return isWord.FindString(strings.ToLower(word))
}

func BenchmarkWordCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EagerFromFile("../assets/words.txt").GroupByString(mapper).CountValues()
	}
}
