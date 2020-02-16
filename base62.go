package base62

import (
	"strings"
)

// Characters defines the character set for base 62 encoding.
const Characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const base = len(Characters)

var digits = make(map[rune]int)

func init() {
	for i, char := range Characters {
		digits[char] = i
	}
}

// Encode a non-negative integer into a base 62 symbol string.
// Panics if argument is positive.
func Encode(id int) string {
	if id < 0 {
		panic("Argument must be non-negative")
	}
	var sb strings.Builder
	for id > 0 {
		rem := id % base
		sb.WriteByte(Characters[rem])
		id = id / base
	}
	return sb.String()
}

// Decode a base 62 encoded string fragment.
// Panics if the fragment contains illegal characters.
func Decode(fragment string) int {
	id := 0
	coeff := 1
	for _, char := range fragment {
		digit := digits[char]
		if char != '0' && digit == 0 {
			panic("Fragment contains illegal character(s)")
		}
		id += coeff * digit
		coeff *= base
	}
	return id
}
