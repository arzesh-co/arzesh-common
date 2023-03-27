package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

// Validate is a method that validates a JWT using an RSA public key
func Validate(token string) bool {
	// Read the RSA public key from file
	publicKey, _ := ioutil.ReadFile("privateRsa.rsa.pub")

	// Parse the RSA public key
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		// Return false if parsing fails
		return false
	}

	// Parse the JWT token using the RSA public key
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			// Return an error if the token's signing method is unexpected
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		// Return the RSA public key for signature verification
		return key, nil
	})
	if err != nil {
		// Return false if token parsing or signature verification fails
		return false
	}

	// Check whether the token's claims are valid and return true if they are
	_, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return false
	}
	return true
}

// FindValue find a specific value in token claim
func FindValue(token, key string) interface{} {
	// Read the RSA public key from file
	publicKey, _ := ioutil.ReadFile("privateRsa.rsa.pub")

	// Parse the RSA public key
	rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		// Return nil if parsing fails
		return nil
	}

	// Parse the JWT token using the RSA public key
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			// Return an error if the token's signing method is unexpected
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		// Return the RSA public key for signature verification
		return rsaKey, nil
	})
	if err != nil {
		// Return nil if token parsing or signature verification fails
		return nil
	}
	// Extract the claims from the parsed token
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		// If the token is invalid or the claims can't be extracted, return nil
		return nil
	}
	// Try to get the value associated with the given key from the token's claims
	value, ok := claims[key]
	if !ok {
		// If the key doesn't exist in the claims, return nil
		return nil
	}
	return value
}

// Open extract token claim and return claim values
func Open(token string) map[string]interface{} {
	// Read the RSA public key from file
	publicKey, _ := ioutil.ReadFile("privateRsa.rsa.pub")

	// Parse the RSA public key
	rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		// Return nil if parsing fails
		return nil
	}

	// Parse the JWT token using the RSA public key
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			// Return an error if the token's signing method is unexpected
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		// Return the RSA public key for signature verification
		return rsaKey, nil
	})
	if err != nil {
		// Return nil if token parsing or signature verification fails
		return nil
	}
	// Extract the claims from the parsed token
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		// If the token is invalid or the claims can't be extracted, return nil
		return nil
	}
	return claims
}
