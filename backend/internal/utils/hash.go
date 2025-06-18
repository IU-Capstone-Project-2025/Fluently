package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the given plain-text password using bcrypt and returns the resulting hash.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plain-text password with its bcrypt hash and returns true if they match.
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
