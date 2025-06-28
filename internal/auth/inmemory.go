package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

// InMemoryAuthenticator is an in-memory implementation of the Authenticator interface.
// It stores username/password combinations in a map.
type InMemoryAuthenticator struct {
	users map[string]string
	jwtSecret string
}

// NewInMemoryAuthenticator creates a new InMemoryAuthenticator.
func NewInMemoryAuthenticator(jwtSecret string) *InMemoryAuthenticator {
	return &InMemoryAuthenticator{
		users: make(map[string]string),
		jwtSecret: jwtSecret,
	}
}

// AddUser adds a new user to the authenticator.
func (a *InMemoryAuthenticator) AddUser(username, password string) {
	a.users[username] = password
	logrus.Infof("Added user: %s", username)
}

// Authenticate authenticates a user and returns a JWT token if successful.
func (a *InMemoryAuthenticator) Authenticate(username, password string) (string, error) {
	if storedPassword, ok := a.users[username]; !ok || storedPassword != password {
		logrus.Warnf("Authentication failed for user %s: invalid username or password", username)
		return "", errors.New("invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		logrus.Errorf("Failed to sign token for user %s: %v", username, err)
		return "", err
	}

	logrus.Infof("User %s authenticated successfully", username)
	return tokenString, nil
}

// Validate validates a JWT token and returns the user ID (username) if valid.
func (a *InMemoryAuthenticator) Validate(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warnf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		logrus.Warnf("Token parsing failed: %v", err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["sub"].(string)
		logrus.Infof("Token validated successfully for user: %s", userID)
		return userID, nil
	}

	logrus.Warn("Invalid token")
	return "", errors.New("invalid token")
}
