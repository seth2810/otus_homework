package hw03frequencyanalysis

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

func Top10(input string) []string {
	wordsSplitRegExp := regexp.MustCompile(`[\s\n,.!?]+`)
	words := wordsSplitRegExp.Split(strings.ToLower(input), -1)
	counter := map[string]int{}
	var topWords []string

	for _, word := range words {
		if len(word) == 0 || word == "-" {
			continue
		}

		counter[word]++

		if counter[word] == 1 {
			topWords = append(topWords, word)
		}
	}

	sort.Slice(topWords, func(i, j int) bool {
		a := topWords[i]
		b := topWords[j]

		if counter[a] == counter[b] {
			return strings.Compare(a, b) == -1
		}

		return counter[a] > counter[b]
	})

	resLen := int(math.Min(float64(len(topWords)), 10))

	// Place your code here
	return topWords[:resLen]
}
