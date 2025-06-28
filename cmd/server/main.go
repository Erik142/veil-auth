package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/erikwahlberger/veil-auth/internal/auth"
	grpcAuth "github.com/erikwahlberger/veil-auth/internal/grpc/auth"
	pb "github.com/erikwahlberger/veil-auth/pkg/grpc/auth"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	authenticator := auth.NewInMemoryAuthenticator("your-secret-key")
	authenticator.AddUser("testuser", "testpassword")

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, grpcAuth.NewServer(authenticator))

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
