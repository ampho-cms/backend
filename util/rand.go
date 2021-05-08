// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package util contains various useful utilities
package util

import (
	"math/rand"
	"strings"
)

const (
	AsciiNum      = "0123456789"
	AsciiSpecials = "~=+%^*/()[]{}/!@#$?|"
	AsciiUpper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AsciiLower    = "abcdefghijklmnopqrstuvwxyz"
	AsciiAlpha    = AsciiUpper + AsciiLower
	AsciiAlphaNum = AsciiAlpha + AsciiNum
	AsciiFull     = AsciiAlphaNum + AsciiSpecials
)

// RandStr generates a random string,
func RandStr(alphabet []rune, length int) string {
	var b strings.Builder

	// Since
	for i := 0; i < length; i++ {
		b.WriteRune(alphabet[rand.Intn(len(alphabet))])
	}

	return b.String()
}

// RandAscii generates a random string using full ASCII alphabet.
func RandAscii(length int) string {
	return RandStr([]rune(AsciiFull), length)
}

// RandAsciiAlpha generates a random string using only ASCII letters.
func RandAsciiAlpha(length int) string {
	return RandStr([]rune(AsciiAlpha), length)
}

// RandAsciiAlphaNum generates a random string using ASCII letters and numbers.
func RandAsciiAlphaNum(length int) string {
	return RandStr([]rune(AsciiAlphaNum), length)
}
