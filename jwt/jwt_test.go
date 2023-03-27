package jwt

import "testing"

func TestValidate(t *testing.T) {
	// Set the static JWT token to validate
	tokenString := ""

	// Call the Validate method with the static JWT token
	valid := Validate(tokenString)

	// Assert that the token is valid
	if !valid {
		t.Errorf("Token validation failed, expected true but got false")
	}
}

func TestFindValue(t *testing.T) {
	// Set the static JWT token with claim value
	tokenString := ""

	// Set the key that is in claim
	keyInClaim := ""

	// Set the value that expected to get
	valueExpected := ""

	tokenValue := FindValue(tokenString, keyInClaim)
	if tokenValue == nil {
		t.Errorf("Token find value failed, expected got a value but got nil")
	}
	if valueExpected != tokenValue {
		t.Errorf("Token find value failed, value returnd is wrong")
	}
}

func TestOpen(t *testing.T) {
	// Set the static JWT token with claim value
	tokenString := ""

	claim := Open(tokenString)
	if claim == nil {
		t.Errorf("open token failed, expected to got a claim value but got nil")
	}
}
