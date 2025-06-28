package auth

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/Erik142/veil-auth/internal/auth"
	pb "github.com/Erik142/veil-auth/pkg/grpc/auth"
)

// Server is the gRPC server for the AuthService.
type Server struct {
	pb.UnimplementedAuthServiceServer
	auth *auth.InMemoryAuthenticator
}

// NewServer creates a new gRPC server.
func NewServer(auth *auth.InMemoryAuthenticator) *Server {
	return &Server{auth: auth}
}

// Authenticate handles an authentication request.
func (s *Server) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	logrus.Infof("Authenticating user: %s", req.Username)
	token, err := s.auth.Authenticate(req.Username, req.Password)
	if err != nil {
		logrus.Errorf("Authentication failed for user %s: %v", req.Username, err)
		return nil, err
	}
	logrus.Infof("User %s authenticated successfully", req.Username)
	return &pb.AuthenticateResponse{Token: token}, nil
}

// Validate handles a validation request.
func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	logrus.Infof("Validating token")
	userID, err := s.auth.Validate(req.Token)
	if err != nil {
		logrus.Errorf("Token validation failed: %v", err)
		return &pb.ValidateResponse{Valid: false}, nil
	}
	logrus.Infof("Token validated successfully for user: %s", userID)
	return &pb.ValidateResponse{Valid: true, UserId: userID}, nil
}
