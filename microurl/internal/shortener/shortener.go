package shortener

import (
	"fmt"
	"strconv"
)

var alphabet = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i',
	'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '1',
	'2', '3', '4', '5', '6', '7', '8', '9', '0',
}

var mod = len(alphabet)

type Custom struct{}

func (hasher Custom) Shorten(id int) string {
	digits := []int{}
	for id > 0 {
		rem := id % mod
		digits = append(digits, rem)
		id = id / mod
	}
	reverse(digits)
	return mapDigits(digits)
}

func reverse(arr []int) {
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

func (hasher Custom) Unwrap(shorten string) (int, error) {
	result, err := mapString(shorten)
	if err != nil {
		return 0, err
	}
	str := ""
	for _, n := range result {
		str = fmt.Sprintf("%s%d", str, n)
	}
	return strconv.Atoi(str)
}

func mapString(str string) ([]int, error) {
	in := []rune(str)
	out := []int{}
	for pos, urlChar := range in {
		for i, char := range alphabet {
			if char == urlChar {
				out = append(out, i*mod^pos)
				break
			}
			return nil, fmt.Errorf("malformed url string")
		}
	}
	return out, nil
}
