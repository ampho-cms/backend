// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package security

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"ampho.xyz/core/config"
)

var (
	signingMethod   jwt.SigningMethod
	hmacKey         []byte
	rsaPrivateKey   *rsa.PrivateKey
	rsaPublicKey    *rsa.PublicKey
	ecdsaPrivateKey *ecdsa.PrivateKey
	ecdsaPublicKey  *ecdsa.PublicKey
)

// SetSigningMethod sets a default signing method.
func SetSigningMethod(alg string) error {
	method := jwt.GetSigningMethod(alg)
	if method == nil {
		return errors.New("unknown signing method algorithm: " + alg)
	}

	signingMethod = method

	return nil
}

// GetSigningMethod returns a default signing method.
func GetSigningMethod() jwt.SigningMethod {
	if signingMethod == nil {
		_ = SetSigningMethod("HS256")
	}

	return signingMethod
}

// SetHMACKey sets an HMAC key.
func SetHMACKey(k []byte) {
	hmacKey = k
}

// SetRSAPrivateKey sets an RSA private key.
func SetRSAPrivateKey(k []byte) error {
	var err error

	if rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(k); err != nil {
		return err
	}

	return nil
}

// SetRSAPublicKey sets an RSA public key.
func SetRSAPublicKey(k []byte) error {
	var err error

	if rsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(k); err != nil {
		return err
	}

	return nil
}

// SetECDSAPrivateKey sets an ECDSA private key.
func SetECDSAPrivateKey(k []byte) error {
	var err error

	if ecdsaPrivateKey, err = jwt.ParseECPrivateKeyFromPEM(k); err != nil {
		return err
	}

	return nil
}

// SetECDSAPublicKey sets an ECDSA public key.
func SetECDSAPublicKey(k []byte) error {
	var err error

	if ecdsaPublicKey, err = jwt.ParseECPublicKeyFromPEM(k); err != nil {
		return err
	}

	return nil
}

// GetPrivateKey returns a private key configured for a signing method.
func GetPrivateKey(signingMethod string) (interface{}, error) {
	if strings.HasPrefix(signingMethod, "HS") {
		if hmacKey == nil {
			return nil, errors.New("HMAC key is not set, check your configuration")
		}
		return hmacKey, nil
	} else if strings.HasPrefix(signingMethod, "RS") {
		if rsaPrivateKey == nil {
			return nil, errors.New("RSA private key is not set, check your configuration")
		}
		return rsaPrivateKey, nil
	} else if strings.HasPrefix(signingMethod, "ES") {
		if ecdsaPrivateKey == nil {
			return nil, errors.New("ECDSA private key is not set, check your configuration")
		}
		return ecdsaPrivateKey, nil
	}

	return nil, errors.New("unknown signing method: " + signingMethod)
}

// GetPublicKey returns a public key configured for a signing method.
func GetPublicKey(signingMethod string) (interface{}, error) {
	if strings.HasPrefix(signingMethod, "HS") {
		if hmacKey == nil {
			return nil, errors.New("HMAC key is not set, check your configuration")
		}
		return hmacKey, nil
	} else if strings.HasPrefix(signingMethod, "RS") {
		if rsaPublicKey == nil {
			return nil, errors.New("RSA private key is not set, check your configuration")
		}
		return rsaPublicKey, nil
	} else if strings.HasPrefix(signingMethod, "ES") {
		if ecdsaPublicKey == nil {
			return nil, errors.New("ECDSA private key is not set, check your configuration")
		}
		return ecdsaPublicKey, nil
	}

	return nil, errors.New("unknown signing method: " + signingMethod)
}

// InitFromConfig initializes security from config.
func InitFromConfig(cfg config.Config) error {
	var err error

	cfg.SetDefault("security.signingMethod", "HS256")

	// Signing method
	if err = SetSigningMethod(cfg.GetString("security.signingMethod")); err != nil {
		return fmt.Errorf("failed to set signing method: %v", err)
	}

	// HMAC key
	hmacKey := cfg.GetString("security.hmac.key")
	if hmacKey != "" {
		SetHMACKey([]byte(hmacKey))
	}

	// RSA keys
	rsaPrvKey := cfg.GetString("security.rsa.privateKey")
	if rsaPrvKey != "" {
		if err = SetRSAPrivateKey([]byte(rsaPrvKey)); err != nil {
			return fmt.Errorf("failed to load private RSA key: %v", err)
		}
	}
	rsaPubKey := cfg.GetString("security.rsa.publicKey")
	if rsaPubKey != "" {
		if err = SetRSAPublicKey([]byte(rsaPubKey)); err != nil {
			return fmt.Errorf("failed to load public RSA key: %v", err)
		}
	}

	// ECDSA keys
	ecdsaPrvKey := cfg.GetString("security.ecdsa.privateKey")
	if ecdsaPrvKey != "" {
		if err = SetECDSAPrivateKey([]byte(ecdsaPrvKey)); err != nil {
			return fmt.Errorf("failed to load private ECDSA key: %v", err)
		}
	}
	ecdsaPubKey := cfg.GetString("security.ecdsa.publicKey")
	if ecdsaPubKey != "" {
		if err = SetECDSAPublicKey([]byte(ecdsaPubKey)); err != nil {
			return fmt.Errorf("failed to load public ECDSA key: %v", err)
		}
	}

	return nil
}
