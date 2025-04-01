package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig contains the configuration needed to generate JWT tokens
type JWTConfig struct {
	KeyID          string
	IssuerID       string
	PrivateKeyPEM  string
	Expiration     time.Duration
}

// GenerateToken creates a new JWT token for App Store Connect API authentication
func GenerateToken(config JWTConfig) (string, error) {
	// Parse the private key
	block, _ := pem.Decode([]byte(config.PrivateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the private key")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	ecdsaKey, ok := privKey.(*ecdsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key is not an ECDSA key")
	}

	// Set the expiration time
	exp := time.Now().Add(config.Expiration)
	if config.Expiration == 0 {
		// Default to 20 minutes if not specified
		exp = time.Now().Add(20 * time.Minute)
	}

	// Create the JWT claims
	claims := jwt.RegisteredClaims{
		Issuer:    config.IssuerID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(exp),
		Audience:  jwt.ClaimStrings{"appstoreconnect-v1"},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = config.KeyID

	// Sign the token
	tokenString, err := token.SignedString(ecdsaKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
