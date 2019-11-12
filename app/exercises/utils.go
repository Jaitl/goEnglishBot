package exercises

import (
	"github.com/jaitl/goEnglishBot/app/phrase"
	"regexp"
	"strings"
)

func computeVariants(text []string, curVariants []string) []string {
	m := make(map[string]bool)
	variants := make([]string, 0, len(text))

	for _, val := range text {
		if _, ok := m[val]; !ok {
			m[val] = true
		}
	}

	for _, val := range curVariants {
		if _, ok := m[val]; ok {
			variants = append(variants, val)
		}
	}

	return variants
}

func ClearText(text string) string {
	text = phrase.Clear(text)
	reg := regexp.MustCompile(`[^a-zA-Z1-9\s\\'\-]+`)

	return reg.ReplaceAllString(strings.ToLower(text), "")
}

func unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
