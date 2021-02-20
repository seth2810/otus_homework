package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"math"
	"regexp"
	"sort"
)

func Top10(input string) []string {
	wordsSplitRegExp := regexp.MustCompile(`[\s\n]+`)
	words := wordsSplitRegExp.Split(input, -1)
	counter := map[string]int{}

	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		counter[word]++
	}

	topWords := make([]string, 0, len(counter))

	for word := range counter {
		topWords = append(topWords, word)
	}

	sort.Slice(topWords, func(i, j int) bool {
		return counter[topWords[i]] > counter[topWords[j]]
	})

	resLen := int(math.Min(float64(len(topWords)), 10))

	// Place your code here
	return topWords[:resLen]
}
