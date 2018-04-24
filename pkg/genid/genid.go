package genid

import (
	"math"
)

var seedChars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var seedCharsLen = len(seedChars)

// IntToBase62 converts integer to
// Base62 encoding.
func IntToBase62(n int) string {
	if n <= 0 {
		return "0"
	}

	b := make([]byte, 0, 64)
	for n > 0 {
		r := math.Mod(float64(n), float64(seedCharsLen))
		n /= seedCharsLen
		b = append(b, seedChars[int(r)])
	}
	b = reverse(b)

	return string(b)
}

// reverse reverses slice of bytes.
func reverse(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - i - 1
		b[i], b[j] = b[j], b[i]
	}

	return b
}
