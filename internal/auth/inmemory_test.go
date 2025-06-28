package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryAuthenticator(t *testing.T) {
	secret := "test_secret"
	auth := NewInMemoryAuthenticator(secret)

	assert.NotNil(t, auth)
	assert.NotNil(t, auth.users)
	assert.Equal(t, secret, auth.jwtSecret)
}

func TestAddUser(t *testing.T) {
	auth := NewInMemoryAuthenticator("test_secret")
	username := "testuser"
	password := "testpassword"

	auth.AddUser(username, password)

	assert.Equal(t, password, auth.users[username])
}

func TestAuthenticate_Success(t *testing.T) {
	secret := "test_secret"
	auth := NewInMemoryAuthenticator(secret)
	username := "testuser"
	password := "testpassword"
	auth.AddUser(username, password)

	tokenString, err := auth.Authenticate(username, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Validate the token to ensure it's correct
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, username, claims["sub"])
}

func TestAuthenticate_InvalidCredentials(t *testing.T) {
	auth := NewInMemoryAuthenticator("test_secret")
	auth.AddUser("testuser", "testpassword")

	tokenString, err := auth.Authenticate("wronguser", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, tokenString)
	assert.EqualError(t, err, "invalid username or password")
}

func TestValidate_Success(t *testing.T) {
	secret := "test_secret"
	auth := NewInMemoryAuthenticator(secret)
	username := "testuser"
	password := "testpassword"
	auth.AddUser(username, password)

	tokenString, err := auth.Authenticate(username, password)
	assert.NoError(t, err)

	userID, err := auth.Validate(tokenString)

	assert.NoError(t, err)
	assert.Equal(t, username, userID)
}

func TestValidate_InvalidToken(t *testing.T) {
	auth := NewInMemoryAuthenticator("test_secret")

	userID, err := auth.Validate("invalid.token.string")

	assert.Error(t, err)
	assert.Empty(t, userID)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestValidate_ExpiredToken(t *testing.T) {
	secret := "test_secret"
	auth := NewInMemoryAuthenticator(secret)
	username := "testuser"

	// Create an expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(-time.Hour).Unix(), // Expired an hour ago
	})
	tokenString, _ := token.SignedString([]byte(secret))

	userID, err := auth.Validate(tokenString)

	assert.Error(t, err)
	assert.Empty(t, userID)
	assert.Contains(t, err.Error(), "Token is expired")
}

func TestValidate_WrongSecret(t *testing.T) {
	secret := "test_secret"
	auth := NewInMemoryAuthenticator(secret)
	username := "testuser"
	password := "testpassword"
	auth.AddUser(username, password)

	// Authenticate with the correct secret to get a valid token
	tokenString, err := auth.Authenticate(username, password)
	assert.NoError(t, err)

	// Create a new authenticator with a different secret
	wrongAuth := NewInMemoryAuthenticator("wrong_secret")

	// Try to validate the token with the wrong secret
	userID, err := wrongAuth.Validate(tokenString)

	assert.Error(t, err)
	assert.Empty(t, userID)
	assert.Contains(t, err.Error(), "signature is invalid")
}