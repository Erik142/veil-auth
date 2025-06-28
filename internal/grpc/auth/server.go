package auth

import (
	"context"

	"github.com/erikwahlberger/veil-auth/internal/auth"
	pb "github.com/erikwahlberger/veil-auth/pkg/grpc/auth"
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
	token, err := s.auth.Authenticate(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.AuthenticateResponse{Token: token}, nil
}

// Validate handles a validation request.
func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	userID, err := s.auth.Validate(req.Token)
	if err != nil {
		return &pb.ValidateResponse{Valid: false}, nil
	}
	return &pb.ValidateResponse{Valid: true, UserId: userID}, nil
}
