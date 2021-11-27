package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	var (
		previousRune rune
		builder      strings.Builder
	)
	runes := []rune(str)
	length := len(runes)
	for index, currentRune := range runes {
		if isFirstDigitalRune(index, currentRune) {
			return "", ErrInvalidString
		}

		if isDoubleDigitalRune(currentRune, previousRune) {
			return "", ErrInvalidString
		}

		if isFirstRune(index) {
			previousRune = currentRune
			continue
		}

		if isDigitalRune(currentRune) && isNotDigitalRune(previousRune) {
			addRepeatedRune(&builder, previousRune, currentRune)
		}

		if isNotDigitalRune(currentRune) && isNotDigitalRune(previousRune) {
			addRune(&builder, previousRune)
		}

		if isLastRune(index, length) {
			addRune(&builder, currentRune)
		}

		previousRune = currentRune
	}

	return builder.String(), nil
}

func isFirstDigitalRune(index int, r rune) bool {
	return isFirstRune(index) && isDigitalRune(r)
}

func isDoubleDigitalRune(r1 rune, r2 rune) bool {
	return isDigitalRune(r1) && isDigitalRune(r2)
}

func isFirstRune(index int) bool {
	return index == 0
}

func isLastRune(index int, length int) bool {
	return index == length-1
}

func isDigitalRune(r rune) bool {
	return unicode.IsDigit(r)
}

func isNotDigitalRune(r rune) bool {
	return !unicode.IsDigit(r)
}

func addRune(builder *strings.Builder, r rune) {
	builder.WriteString(string(r))
}

func addRepeatedRune(builder *strings.Builder, r rune, numberRune rune) {
	number, _ := strconv.Atoi(string(numberRune))
	builder.WriteString(strings.Repeat(string(r), number))
}
