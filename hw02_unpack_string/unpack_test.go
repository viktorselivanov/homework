package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "aaф0b", expected: "aab"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestUnpackWithIncorrect(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{

		{"a4bc2d5e", "aaaabccddddde", false}, //  повторения
		{"abcd", "abcd", false},              // без повторений
		{"aaa0b", "aab", false},              //ноль
		{"ra0b0c", "rc", false},              // несколько нулей
		{"a1r1c", "arc", false},              // повтор 1
		{"🙂🙃3", "🙂🙃🙃🙃", false},               // смайлы
		{"к2е3", "ккеее", false},             // UTF-8 символы
		{"аыамы0я", "аыамя", false},          // UTF-8 символы
		{"a\n3b", "a\n\n\nb", false},         // спецсимволы
		{"\t2x", "\t\tx", false},             // спецсимволы
		{"", "", false},                      // пустая строка
		{"3abc", "", true},                   // некорректная строка
		{"45", "", true},                     // некорректная строка
		{"aaa10b", "", true},                 // некорректная строка
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			if tc.expectError {
				require.Error(t, err)
				require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}
