package util

import "golang.org/x/crypto/bcrypt"

// HashPassword takes plaintext password as input and returns
// its hashed value.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

// MatchesWithHash returns true if password matches with hash.
func MatchesWithHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
