package hw02_unpack_string //nolint:golint,stylecheck

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

	var isEscaped bool
	var prevChar, char rune
	var builder strings.Builder

	runes := []rune(input)

	for i := 1; i < len(runes); i++ {
		prevChar = runes[i-1]
		isPrevDigit := unicode.IsDigit(prevChar)

		if i == 1 && isPrevDigit {
			return "", ErrInvalidString
		}

		char = runes[i]
		isDigit := unicode.IsDigit(char)

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
