package main

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
		actual, err := Encode(pair.decoded)
		if err != nil {
			t.Fatal(err)
		}
		if actual != pair.encoded {
			t.Errorf("%d encoded to %s, not %s", pair.decoded, actual, pair.encoded)
		}
	}
}

func TestEncodeNegative(t *testing.T) {
	arg := -1
	encoded, err := Encode(arg)
	if err == nil {
		t.Errorf("Negative argument %d encoded to %s", arg, encoded)
	}
}

func TestDecode(t *testing.T) {
	for _, pair := range pairs {
		actual, err := Decode(pair.encoded)
		if err != nil {
			t.Fatal(err)
		}
		if actual != pair.decoded {
			t.Errorf("%s decoded to %d, not %d", pair.encoded, actual, pair.decoded)
		}
	}
}

func TestDecodeIllegalCharacter(t *testing.T) {
	arg := "!llegal"
	decoded, err := Decode(arg)
	if err == nil {
		t.Errorf("Malformed fragment %s decoded to %d", arg, decoded)
	}
}

func TestEncodeDecode(t *testing.T) {
	testCases := []int{0, 1, 62}
	for i := 0; i < 20; i++ {
		testCases = append(testCases, rand.Int())
	}
	for _, value := range testCases {
		encoded, err := Encode(value)
		if err != nil {
			t.Fatal(err)
		}
		actual, err := Decode(encoded)
		if err != nil {
			t.Fatal(err)
		}
		if actual != value {
			t.Errorf("%d does not roundtrip", value)
		}
	}
}
