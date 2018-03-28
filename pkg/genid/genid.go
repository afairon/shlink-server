package genid

import (
	"bytes"
	"errors"
)

var seedChars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var seedCharsLen = len(seedChars)
var aChar = byte(97)

// GenerateNextID generates next base62 ID
func GenerateNextID(code string) (string, error) {
	if code == "" {
		return string(aChar), nil
	}
	codeBytes := []byte(code)
	codeByteLen := len(codeBytes)

	codeCharIndex := -1
	for i := (codeByteLen - 1); i >= 0; i-- {
		codeCharIndex = bytes.IndexByte(seedChars, codeBytes[i])
		if codeCharIndex == -1 || codeCharIndex >= seedCharsLen {
			return "", errors.New("")
		} else if codeCharIndex == (seedCharsLen - 1) {
			codeBytes[i] = aChar
		} else {
			codeBytes[i] = seedChars[(codeCharIndex + 1)]
		}
	}
	for _, byteVal := range codeBytes {
		if byteVal != aChar {
			return string(codeBytes), nil
		}
	}
	return "a" + string(codeBytes), nil
}
