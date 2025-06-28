package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/erikwahlberger/veil-auth/pkg/grpc/auth"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: client <command> [args]")
	}

	switch os.Args[1] {
	case "authenticate":
		if len(os.Args) != 4 {
			log.Fatalf("Usage: client authenticate <username> <password>")
		}
		authenticate(c, os.Args[2], os.Args[3])
	case "validate":
		if len(os.Args) != 3 {
			log.Fatalf("Usage: client validate <token>")
		}
		validate(c, os.Args[2])
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func authenticate(c pb.AuthServiceClient, username, password string) {
	r, err := c.Authenticate(context.Background(), &pb.AuthenticateRequest{Username: username, Password: password})
	if err != nil {
		log.Fatalf("could not authenticate: %v", err)
	}
	log.Printf("Token: %s", r.Token)
}

func validate(c pb.AuthServiceClient, token string) {
	r, err := c.Validate(context.Background(), &pb.ValidateRequest{Token: token})
	if err != nil {
		log.Fatalf("could not validate: %v", err)
	}
	log.Printf("Valid: %t, UserID: %s", r.Valid, r.UserId)
}
