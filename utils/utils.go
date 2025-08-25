package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of a plain password.
// Call this **once** when you create the user (e.g. via an admin script).
func HashPassword(plain string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost) // cost = 10
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	return string(hash) // store this string in the DB
}
