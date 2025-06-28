package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	

	mockauth "github.com/Erik142/veil-auth/internal/auth/mocks"
	pb "github.com/Erik142/veil-auth/pkg/grpc/auth"
)

func TestNewServer(t *testing.T) {
	mockAuthenticator := new(mockauth.Authenticator)
	server := NewServer(mockAuthenticator)

	assert.NotNil(t, server)
	assert.Equal(t, mockAuthenticator, server.auth)
}

func TestAuthenticate_Success(t *testing.T) {
	mockAuthenticator := new(mockauth.Authenticator)
	server := NewServer(mockAuthenticator)

	username := "testuser"
	password := "testpassword"
	token := "testtoken"

	mockAuthenticator.On("Authenticate", username, password).Return(token, nil).Once()

	req := &pb.AuthenticateRequest{
		Username: username,
		Password: password,
	}

	resp, err := server.Authenticate(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, token, resp.Token)
	mockAuthenticator.AssertExpectations(t)
}

func TestAuthenticate_Failure(t *testing.T) {
	mockAuthenticator := new(mockauth.Authenticator)
	server := NewServer(mockAuthenticator)

	username := "testuser"
	password := "testpassword"
	expectedErr := errors.New("authentication failed")

	mockAuthenticator.On("Authenticate", username, password).Return("", expectedErr).Once()

	req := &pb.AuthenticateRequest{
		Username: username,
		Password: password,
	}

	resp, err := server.Authenticate(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockAuthenticator.AssertExpectations(t)
}

func TestValidate_Success(t *testing.T) {
	mockAuthenticator := new(mockauth.Authenticator)
	server := NewServer(mockAuthenticator)

	token := "testtoken"
	userID := "testuser"

	mockAuthenticator.On("Validate", token).Return(userID, nil).Once()

	req := &pb.ValidateRequest{
		Token: token,
	}

	resp, err := server.Validate(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, userID, resp.UserId)
	mockAuthenticator.AssertExpectations(t)
}

func TestValidate_Failure(t *testing.T) {
	mockAuthenticator := new(mockauth.Authenticator)
	server := NewServer(mockAuthenticator)

	token := "invalidtoken"
	expectedErr := errors.New("invalid token")

	mockAuthenticator.On("Validate", token).Return("", expectedErr).Once()

	req := &pb.ValidateRequest{
		Token: token,
	}

	resp, err := server.Validate(context.Background(), req)

	assert.NoError(t, err) // gRPC server returns nil error even if validation fails
	assert.NotNil(t, resp)
	assert.False(t, resp.Valid)
	assert.Empty(t, resp.UserId)
	mockAuthenticator.AssertExpectations(t)
}
