package base62

import (
	"math/rand"
	"testing"
)

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
