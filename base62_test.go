package base62

import (
	"math/rand"
	"testing"
)

func TestEncode(t *testing.T) {
	type pair struct {
		value    int
		expected string
	}
	pairs := []pair{
		{0, ""},
		{1, "1"},
		{62, "01"},
	}
	for _, pair := range pairs {
		actual := Encode(pair.value)
		if actual != pair.expected {
			t.Error(pair.value)
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	testCases := 20
	for i := 0; i < testCases; i++ {
		value := rand.Int()
		actual := Decode(Encode(value))
		if actual != value {
			t.Error(value, " does not roundtrip")
		}
	}
}
