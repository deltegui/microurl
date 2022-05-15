package shortener_test

import (
	"microurl/internal/shortener"
	"testing"
)

func TestShouldShortenIDs(t *testing.T) {
	type data struct {
		name     string
		inputID  int
		expected string
	}
	cases := []data{
		{
			name:     "zero",
			inputID:  0,
			expected: "a",
		},
		{
			name:     "one",
			inputID:  1,
			expected: "b",
		},
		{
			name:     "one digit",
			inputID:  4,
			expected: "e",
		},
		{
			name:     "two digits",
			inputID:  76,
			expected: "ob",
		},
		{
			name:     "three digits",
			inputID:  323,
			expected: "nf",
		},
		{
			name:     "Sequence",
			inputID:  12345,
			expected: "dnh",
		},
		{
			name:     "Large",
			inputID:  8973490812734,
			expected: "cO7678Hc",
		},
	}
	sh := shortener.Base62{}
	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := sh.Shorten(c.inputID)
			if result != c.expected {
				t.Errorf(
					"[%d] With ID: %d, expected: %s but get %s",
					i,
					c.inputID,
					c.expected,
					result,
				)
			}
		})
	}
}

func TestShouldUnwrapStrings(t *testing.T) {
	type data struct {
		name     string
		short    string
		expected int
	}
	cases := []data{
		{
			name:     "tested",
			expected: 19158,
			short:    "e9a",
		},
		{
			name:     "zero",
			expected: 0,
			short:    "a",
		},
		{
			name:     "one",
			expected: 1,
			short:    "b",
		},
		{
			name:     "one digit",
			expected: 4,
			short:    "e",
		},
		{
			name:     "two digits",
			expected: 76,
			short:    "ob",
		},
		{
			name:     "three digits",
			expected: 323,
			short:    "nf",
		},
		{
			name:     "Sequence",
			expected: 12345,
			short:    "dnh",
		},
		{
			name:     "Large",
			expected: 8973490812734,
			short:    "cO7678Hc",
		},
	}
	sh := shortener.Base62{}
	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := sh.Unwrap(c.short)
			if err != nil {
				t.Errorf("[%d] Error while unwrapping: %s", i, err)
			}
			if result != c.expected {
				t.Errorf(
					"[%d] expected: %d but get %d",
					i,
					c.expected,
					result,
				)
			}
		})
	}
}
