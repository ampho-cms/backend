// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package security

import (
	"github.com/dgrijalva/jwt-go"
)

// NewToken creates a new JWT token without claims section.
func NewToken() *jwt.Token {
	return jwt.New(GetSigningMethod())
}

// NewTokenWithClaims creates a new JWT token with claims section.
func NewTokenWithClaims(claims jwt.Claims) *jwt.Token {
	return jwt.NewWithClaims(GetSigningMethod(), claims)
}

// GetTokenSignedString returns a token as a signed string.
func GetTokenSignedString(t *jwt.Token) (string, error) {
	var err error

	key, err := GetPrivateKey(t.Method.Alg())
	if err != nil {
		return "", err
	}

	s, err := t.SignedString(key)
	if err != nil {
		return "", err
	}

	return s, nil
}

// ParseToken parses a signed token string without claims section.
func ParseToken(s string) (*jwt.Token, error) {
	return jwt.Parse(s, func(t *jwt.Token) (interface{}, error) {
		return GetPublicKey(t.Method.Alg())
	})
}

// ParseTokenWithClaims parses a signed token string with claims section.
func ParseTokenWithClaims(s string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(s, claims, func(t *jwt.Token) (interface{}, error) {
		return GetPublicKey(t.Method.Alg())
	})
}
