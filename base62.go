package base62

import (
	"strings"
)

const characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base = len(characters)

var digits = make(map[rune]int)

func init() {
	for i, char := range characters {
		digits[char] = i
	}
}

func Encode(id int) string {
	var sb strings.Builder
	for id > 0 {
		rem := id % base
		sb.WriteByte(characters[rem])
		id = id / base
	}
	return sb.String()
}

func Decode(fragment string) int {
	id := 0
	coeff := 1
	for _, char := range fragment {
		id += coeff * digits[char]
		coeff *= base
	}
	return id
}
