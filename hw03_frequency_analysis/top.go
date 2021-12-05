package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type Word struct {
	Name string
	Num  int
}

var (
	reSpecSymbols = regexp.MustCompile(`[^a-zA-Zа-яА-Я- ]`)
	reSpace       = regexp.MustCompile(`\s+`)
	reDash        = regexp.MustCompile(`\s-\s`)
)

func Top10(source string) []string {
	if len(source) == 0 {
		return []string{}
	}

	var words []Word
	words = getWords(source)
	words = sortWords(words)

	return getSlice(words, 10)
}

func getSlice(words []Word, sliceLength int) []string {
	if length := len(words); length < sliceLength {
		sliceLength = length
	}

	result := make([]string, 0, sliceLength)
	for _, word := range words[:sliceLength] {
		result = append(result, word.Name)
	}

	return result
}

func sortWords(words []Word) []Word {
	sort.Slice(words, func(i, j int) bool {
		if words[i].Num == words[j].Num {
			return words[i].Name < words[j].Name
		}
		return words[i].Num > words[j].Num
	})

	return words
}

func getWords(source string) []Word {
	numbers := getNumbers(source)
	words := make([]Word, 0, len(numbers))
	for str, num := range numbers {
		words = append(words, Word{str, num})
	}

	return words
}

func getNumbers(source string) map[string]int {
	numbers := make(map[string]int)
	for _, str := range strings.Fields(cleanSource(source)) {
		numbers[str]++
	}

	return numbers
}

func cleanSource(source string) string {
	source = reSpecSymbols.ReplaceAllString(source, " ")
	source = reDash.ReplaceAllString(source, " ")
	source = reSpace.ReplaceAllString(source, " ")

	return strings.ToLower(strings.Trim(source, " "))
}
