package auth

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/adfer-dev/analock-api/models"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

// GenerateToken generates and returns a new token of the selected kind.
// The possible token kinds are:
//   - Access token
//   - Refresh token
//
// Returns error if there was a problem signing the token
func GenerateToken(user models.User, kind models.TokenKind) (string, error) {
	secretKey, envErr := getSecretKey()

	if envErr != nil {
		return "", nil
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	var expiration int64

	if kind == models.Access {
		expiration = time.Now().Add(1 * time.Hour).Unix()
	} else {
		expiration = time.Now().Add(8766 * time.Hour).Unix()
	}

	claims["exp"] = expiration
	claims["email"] = user.Email

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken validates the token string passed and returns an error if it's not valid
func ValidateToken(tokenString string) error {

	secretKey, envErr := getSecretKey()

	if envErr != nil {
		return envErr
	}

	//First check if token is in the database, since if it's not it would be revoked
	token, parseErr := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("signing method not valid")
		}

		return secretKey, nil
	})

	if parseErr != nil || !token.Valid {
		return parseErr
	}

	return nil
}

// GetClaims parses the given JWT and returns a map containing its claims.
// It returns error if could not get the SECRET_KEY env variable or there was a problem parsing the token.
func GetClaims(tokenString string) (jwt.MapClaims, error) {
	secretKey, envErr := getSecretKey()

	if envErr != nil {
		return nil, envErr
	}

	jwtToken, parseErr := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if parseErr != nil {
		return nil, parseErr
	}

	claims := jwtToken.Claims.(jwt.MapClaims)

	return claims, nil
}

// getSecretKey returns the SECRET_KEY env variable.
// It returns error if the variable could not be loaded.
func getSecretKey() ([]byte, error) {
	if secretKey == nil {
		secretKey = make([]byte, 32)
		_, err := rand.Read(secretKey)

		if err != nil {
			return nil, err
		}
	}

	return secretKey, nil
}
