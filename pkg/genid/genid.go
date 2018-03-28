package genid

import (
	"bytes"
	"errors"
	"math"
)

var seedChars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var seedCharsLen = len(seedChars)

const aChar = byte(97)

// GenerateNextID generates next Base62 ID.
func GenerateNextID(code string) (string, error) {

	if code == "" {
		return string(aChar), nil
	}
	codeBytes := []byte(code)
	codeByteLen := len(codeBytes)

	codeCharIndex := -1
	for i := (codeByteLen - 1); i >= 0; i-- {
		codeCharIndex = bytes.IndexByte(seedChars, codeBytes[i])

		// Char not in seedChars
		if codeCharIndex == -1 || codeCharIndex >= seedCharsLen {
			return "", errors.New("")
		} else if codeCharIndex == (seedCharsLen - 1) {
			codeBytes[i] = aChar
		} else {
			codeBytes[i] = seedChars[(codeCharIndex + 1)]
			break
		}
	}

	for _, byteVal := range codeBytes {
		if byteVal != aChar {
			return string(codeBytes), nil
		}
	}
	return "a" + string(codeBytes), nil
}

// IntToBase62 converts integer to
// Base62 encoding.
func IntToBase62(n int) string {
	if n == 0 {
		return "0"
	}

	// Byte array capacity
	cap := (n / seedCharsLen) + 1
	b := make([]byte, cap, cap)
	for i := (n / seedCharsLen); i >= 0; i-- {
		r := math.Mod(float64(n), float64(seedCharsLen))
		n /= seedCharsLen
		b[i] = seedChars[int(r)]

		if i == 0 && cap > 1 {
			b[i] = seedChars[int(r)] - 1
		}
	}

	return string(b)
}
