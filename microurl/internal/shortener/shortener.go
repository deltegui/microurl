package shortener

import (
	"math"
)

var alphabet = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i',
	'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A',
	'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
	'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1',
	'2', '3', '4', '5', '6', '7', '8', '9',
}

var mod = len(alphabet)

type Base62 struct{}

func (hasher Base62) Shorten(id int) string {
	digits := []int{}
	if id == 0 {
		return "a"
	}
	for id > 0 {
		rem := id % mod
		digits = append(digits, rem)
		id = id / mod
	}
	reverse(digits)
	return mapDigits(digits)
}

func reverse[T any](arr []T) {
	for i, n := range arr {
		j := len(arr) - 1 - i
		if i == j {
			break
		}
		other := arr[j]
		arr[j] = n
		arr[i] = other
	}
}

func mapDigits(arr []int) string {
	out := []rune{}
	for _, n := range arr {
		out = append(out, alphabet[n])
	}
	return string(out)
}

func (hasher Base62) Unwrap(shorten string) (int, error) {
	result, err := mapString(shorten)
	if err != nil {
		return 0, err
	}
	out := 0
	for _, n := range result {
		out += n
	}
	return out, nil
}

func mapString(str string) ([]int, error) {
	in := []rune(str)
	reverse(in)
	out := []int{}
	for pos, urlChar := range in {
		for i, char := range alphabet {
			if char == urlChar {
				pow := intPow(mod, pos)
				out = append(out, i*pow)
				break
			}
		}
	}
	return out, nil
}

func intPow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
