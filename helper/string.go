package helper

import (
	"unicode"
	"strings"
)

func CamelCase(s string) (string) {
	var result string

	words := strings.Split(s, "_")

	for _, word := range words {
		w := []rune(word)
		w[0] = unicode.ToUpper(w[0])
		result += string(w)
	}

	return result
}
