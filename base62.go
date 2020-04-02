package main

import (
	"errors"
	"strings"
)

// Characters defines the character set for base 62 encoding
const Characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const base = len(Characters)

var digits = make(map[rune]int)

func init() {
	for i, char := range Characters {
		digits[char] = i
	}
}

// Encode a non-negative integer into a base 62 symbol string
func Encode(id int) (string, error) {
	if id < 0 {
		return "", errors.New("Argument to Encode must be non-negative")
	}
	var sb strings.Builder
	for id > 0 {
		rem := id % base
		sb.WriteByte(Characters[rem])
		id = id / base
	}
	return sb.String(), nil
}

// Decode a base 62 encoded string fragment
func Decode(fragment string) (int, error) {
	id := 0
	coeff := 1
	for _, char := range fragment {
		digit := digits[char]
		if char != '0' && digit == 0 {
			return 0, errors.New("Fragment contains illegal character(s)")
		}
		id += coeff * digit
		coeff *= base
	}
	return id, nil
}
