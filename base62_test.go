package base62

import (
	"math/rand"
	"testing"
)

type pair struct {
	decoded int
	encoded string
}

var pairs = []pair{
	{0, ""},
	{1, "1"},
	{62, "01"},
	{1504, "go"},
}

func TestEncode(t *testing.T) {
	for _, pair := range pairs {
		actual := Encode(pair.decoded)
		if actual != pair.encoded {
			t.Error(pair.decoded, " encoded to ", actual, " not ", pair.encoded)
		}
	}
}

func TestDecode(t *testing.T) {
	for _, pair := range pairs {
		actual := Decode(pair.encoded)
		if actual != pair.decoded {
			t.Error(pair.encoded, " decoded to ", actual, " not ", pair.decoded)
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
