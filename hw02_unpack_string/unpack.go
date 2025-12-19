package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(iString string) (string, error) {
	if iString == "" {
		return "", nil
	}

	var prevRune rune // предыдущая руна
	var builder strings.Builder

	for i, Irune := range iString {

		if unicode.IsDigit(Irune) {
			// нельзя число в начале
			if i == 0 && unicode.IsDigit(Irune) {
				return "", ErrInvalidString
			}
			// нельзя цифру!
			if i > 0 && unicode.IsDigit(Irune) && unicode.IsDigit(rune(iString[i-1])) {
				return "", ErrInvalidString
			}

			// конвертируем в int
			num, _ := strconv.Atoi(string(Irune))

			if num == 0 {
				// если ноль, то убираем символ
				runes := []rune(builder.String())
				builder.Reset()
				builder.WriteString(string(runes[:len(runes)-1]))
			} else {
				// повторяем символ
				builder.WriteString(strings.Repeat(string(prevRune), num-1))
			}
		} else {
			// записываем символ
			builder.WriteRune(Irune)
			prevRune = Irune
		}

	}

	return builder.String(), nil
}
