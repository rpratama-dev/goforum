package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type BaseClaim struct {
	Name  string `json:"name"`
	UserID uuid.UUID   `json:"uid"`
	SessionID uuid.UUID `json:"sid"`
}

type ClaimPayload struct {
	UserName  string `json:"user_name"`
	ExpiresAt int64  `json:"exp,omitempty"`
	BaseClaim
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time.
type Claims struct {
	BaseClaim
	jwt.StandardClaims
}

// generateJWT generates a new JWT using the RS256 algorithm
func GenerateJWT(claim ClaimPayload) (string, *Claims, error) {
	// Load the private key
	privateKeyBytes := ReadFile(privateKeyPath)
	if (privateKeyBytes == nil) {
		return "", nil, fmt.Errorf("failed to read private key file")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create the JWT claims
	claims := Claims{}
	claims.Name = claim.Name;
	claims.UserID = claim.UserID;
	claims.SessionID = claim.SessionID;
	claims.Audience = claim.UserName;
	claims.IssuedAt = time.Now().Unix();
	claims.Issuer = "My Movie Company, PT";
	claims.Subject = "My Movie App";
	claims.ExpiresAt = claim.ExpiresAt;

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the JWT token with the private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign JWT: %w", err);
	}

	return tokenString, &claims, nil
}

// verifyJWT verifies and parses the JWT using the RS256 algorithm
func VerifyJWT(tokenString string) (*Claims, error) {
	// Load the public key
	publicKeyBytes := ReadFile(publicKeyPath)
	if publicKeyBytes == nil {
		return nil, fmt.Errorf("failed to read public key file")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Parse and verify the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse and verify JWT token: %w", err)
	}

	// Map the claims into the struct
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("JWT token is not valid")
	}

	return claims, nil
}
