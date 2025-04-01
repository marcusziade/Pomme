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
	fmt.Println("Debug: Starting JWT token generation")
	fmt.Println("Debug: Key ID:", config.KeyID)
	fmt.Println("Debug: Issuer ID:", config.IssuerID)
	fmt.Println("Debug: PEM key length:", len(config.PrivateKeyPEM))
	
	// Print first few characters of key to debug
	keyStart := ""
	if len(config.PrivateKeyPEM) > 20 {
		keyStart = config.PrivateKeyPEM[:20]
	} else {
		keyStart = config.PrivateKeyPEM
	}
	fmt.Printf("Debug: Key starts with: %s\n", keyStart)
	
	// The private key from Apple might already be in PEM format
	privateKeyData := config.PrivateKeyPEM
	
	// Check if the key already has PEM headers
	if !ContainsPEMHeaders(privateKeyData) {
		fmt.Println("Debug: Adding PEM headers to private key")
		privateKeyData = fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----", privateKeyData)
	}
	
	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyData))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the private key")
	}

	fmt.Println("Debug: Successfully parsed PEM block")
	
	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	fmt.Println("Debug: Successfully parsed PKCS8 private key")
	
	ecdsaKey, ok := privKey.(*ecdsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key is not an ECDSA key")
	}
	
	fmt.Println("Debug: Successfully cast to ECDSA key")

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

	fmt.Println("Debug: Created token with claims and key ID")

	// Sign the token
	tokenString, err := token.SignedString(ecdsaKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	fmt.Println("Debug: Successfully signed token")
	fmt.Printf("Debug: Complete token: %s\n", tokenString)

	return tokenString, nil
}
