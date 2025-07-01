package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ContainsPEMHeaders checks if the string contains PEM headers
func ContainsPEMHeaders(s string) bool {
	return strings.Contains(s, "-----BEGIN") && strings.Contains(s, "-----END")
}

// JWTConfig contains the configuration needed to generate JWT tokens
type JWTConfig struct {
	KeyID          string
	IssuerID       string
	PrivateKeyPEM  string
	Expiration     time.Duration
}

// GenerateToken creates a new JWT token for App Store Connect API authentication
func GenerateToken(config JWTConfig) (string, error) {
	
	// The private key from Apple might already be in PEM format
	privateKeyData := config.PrivateKeyPEM
	
	// Check if the key already has PEM headers
	if !ContainsPEMHeaders(privateKeyData) {
		privateKeyData = fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----", privateKeyData)
	}
	
	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyData))
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
	// According to Apple documentation, the JWT must include specific claims
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": config.IssuerID,                       // Issuer
		"iat": now.Unix(),                            // Issued at
		"exp": exp.Unix(),                            // Expiration
		"aud": "appstoreconnect-v1",                  // Audience
	}

	// Create the token with ES256 (ECDSA using P-256 and SHA-256)
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = config.KeyID
	token.Header["typ"] = "JWT"                       // Explicitly set token type


	// Sign the token
	tokenString, err := token.SignedString(ecdsaKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}


	return tokenString, nil
}
