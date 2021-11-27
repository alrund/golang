package hw02unpackstring

import (
	"errors"
	"strings"
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
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
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
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestIsFirstDigitalRune(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		r        rune
		expected bool
	}{
		{name: "0:1", index: 0, r: '1', expected: true},
		{name: "0:x", index: 0, r: 'x', expected: false},
		{name: "1:1", index: 1, r: '1', expected: false},
		{name: "1:x", index: 1, r: 'x', expected: false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, isFirstDigitalRune(tc.index, tc.r), tc.expected)
		})
	}
}

func TestIsDoubleDigitalRune(t *testing.T) {
	tests := []struct {
		name     string
		r1       rune
		r2       rune
		expected bool
	}{
		{name: "1:1", r1: '1', r2: '1', expected: true},
		{name: "1:x", r1: '1', r2: 'x', expected: false},
		{name: "x:1", r1: 'x', r2: '1', expected: false},
		{name: "x:x", r1: 'x', r2: 'x', expected: false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, isDoubleDigitalRune(tc.r1, tc.r2), tc.expected)
		})
	}
}

func TestAddRepeatedRune(t *testing.T) {
	tests := []struct {
		name     string
		r1       rune
		r2       rune
		expected string
	}{
		{name: "z:5", r1: 'z', r2: '5', expected: "zzzzz"},
		{name: "y:0", r1: 'y', r2: '0', expected: ""},
		{name: "x:x", r1: 'x', r2: 'x', expected: ""},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var builder strings.Builder
			addRepeatedRune(&builder, tc.r1, tc.r2)
			require.Equal(t, tc.expected, builder.String())
		})
	}
}
