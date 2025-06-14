package auth

import (
	"crypto/rand"
)

var secretKey []byte

// GetSecretKey returns the SECRET_KEY env variable.
// It returns error if the variable could not be loaded.
func GetSecretKey() ([]byte, error) {
	if secretKey == nil {
		secretKey = make([]byte, 32)
		_, err := rand.Read(secretKey)

		if err != nil {
			return nil, err
		}
	}

	return secretKey, nil
}
