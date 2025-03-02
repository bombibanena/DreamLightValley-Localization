package translator

import (
	"fmt"
	"strings"
)

func normalizePlaceholders(text string) string {
	res := text
	for _, template := range []string{"S_0", "S_o", "S_O", "s_0", "s_o", "s_O"} {
		for _, lSpace := range []string{"", " "} {
			for _, rSpace := range []string{"", " "} {
				res = strings.ReplaceAll(res, fmt.Sprintf("(%s%s%s)", lSpace, template, rSpace), "(S_0)")
			}
		}
	}

	return res
}

func escapeSpecialChars(text string) string {
	text = strings.ReplaceAll(text, "\"", "(Q_0)")
	text = strings.ReplaceAll(text, "*", "(Q_1)")
	text = strings.ReplaceAll(text, "[", "(Q_2)")
	text = strings.ReplaceAll(text, "]", "(Q_3)")
	return text
}

func unescapeSpecialChars(text string) string {
	res := text

	for _, template := range []string{"Q_0", "Q_o", "Q_O", "q_0", "q_o", "q_O"} {
		for _, lSpace := range []string{"", " "} {
			for _, rSpace := range []string{"", " "} {
				res = strings.ReplaceAll(res, fmt.Sprintf("(%s%s%s)", lSpace, template, rSpace), "(Q_0)")
			}
		}
	}
	res = strings.ReplaceAll(res, "(Q_0)", "\"")

	for _, template := range []string{"Q_1", "q_1"} {
		for _, lSpace := range []string{"", " "} {
			for _, rSpace := range []string{"", " "} {
				res = strings.ReplaceAll(res, fmt.Sprintf("(%s%s%s)", lSpace, template, rSpace), "(Q_1)")
			}
		}
	}
	res = strings.ReplaceAll(res, "(Q_1)", "*")
	// res = strings.ReplaceAll(res, " *", "*")
	res = strings.ReplaceAll(res, "(Q_2)", "[")
	res = strings.ReplaceAll(res, "(Q_3)", "]")
	// res = strings.ReplaceAll(res, "* ", "*")

	return res
}

func splitText(text string, maxSize int) []string {
	var result []string
	start := 0

	for start < len(text) {
		end := start + maxSize
		if end > len(text) {
			end = len(text)
		}

		lastDelimiter := -1
		for i := end - 1; i >= start; i-- {
			if isDelimiter(rune(text[i])) || text[i] == ' ' {
				lastDelimiter = i
				break
			}
		}

		if lastDelimiter != -1 {
			end = lastDelimiter + 1
		}

		result = append(result, text[start:end])
		start = end
	}

	return result
}

func isDelimiter(char rune) bool {
	delimiters := []rune{'.', ',', '!', '?', ';', ':', ' '}
	for _, delim := range delimiters {
		if char == delim {
			return true
		}
	}
	return false
}
