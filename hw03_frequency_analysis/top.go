package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	if text == "" {
		return nil
	}

	// разделение по пробелам
	words := strings.Fields(text)

	// подсчет частоты
	freqCount := make(map[string]int)
	for _, word := range words {
		freqCount[word]++
	}

	// формирование слайса
	type wordSlice struct {
		word  string
		count int
	}
	slice := make([]wordSlice, 0, len(freqCount))
	for w, c := range freqCount {
		slice = append(slice, wordSlice{w, c})
	}

	// сортировка слов (по убыванию частоты, лексикографически)
	sort.Slice(slice, func(i, j int) bool {
		if slice[i].count != slice[j].count {
			return slice[i].count > slice[j].count
		}
		return slice[i].word < slice[j].word
	})

	topWords := make([]string, 0, 10)
	for i, w := range slice {
		if i >= 10 {
			break
		}
		topWords = append(topWords, w.word)
	}

	return topWords
}
