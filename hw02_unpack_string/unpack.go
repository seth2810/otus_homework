package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const escapeChar = '\\'

func Unpack(input string) (string, error) {
	if input == "" {
		return input, nil
	}

	var prevChar, char rune
	var builder strings.Builder
	var isEscaped, isPrevDigit, isDigit bool

	runes := []rune(input)
	char = runes[0]
	isDigit = unicode.IsDigit(char)

	if isDigit {
		return "", ErrInvalidString
	}

	for i := 1; i < len(runes); i++ {
		prevChar = char
		isPrevDigit = isDigit

		char = runes[i]
		isDigit = unicode.IsDigit(char)

		// handle escapes only for digits and slashes
		if prevChar == escapeChar && !isEscaped {
			if !isDigit && char != escapeChar {
				return "", ErrInvalidString
			}

			isEscaped = true
			continue
		}

		// handle digit before current char
		if isPrevDigit && !isEscaped {
			if isDigit {
				return "", ErrInvalidString
			}

			continue
		}

		if isDigit {
			count, _ := strconv.Atoi(string(char))
			builder.WriteString(strings.Repeat(string(prevChar), count))
		} else {
			builder.WriteRune(prevChar)
		}

		isEscaped = false
	}

	if isEscaped || !unicode.IsDigit(char) {
		builder.WriteRune(char)
	}

	// Place your code here
	return builder.String(), nil
}
