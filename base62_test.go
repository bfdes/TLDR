package base62

import (
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	tests := []int{0, 1, 20, 62, 95, 200}
	for _, value := range tests {
		actual := Decode(Encode(value))
		if actual != value {
			t.Error(value, " does not roundtrip")
		}
	}
}
