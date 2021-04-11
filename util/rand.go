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

// RandStrAscii generates a random string using full ASCII alphabet.
func RandStrAscii(length int) string {
	return RandStr([]rune(AsciiFull), length)
}

// RandStrAsciiAlpha generates a random string using ASCII letters.
func RandStrAsciiAlpha(length int) string {
	return RandStr([]rune(AsciiAlpha), length)
}

// RandStrAsciiAlphaNum generates a random string using ASCII letters and numbers.
func RandStrAsciiAlphaNum(length int) string {
	return RandStr([]rune(AsciiAlphaNum), length)
}
