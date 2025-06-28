package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Erik142/veil-auth/pkg/grpc/auth"
)

func init() {
	viper.SetDefault("server_address", "localhost:50051")
	viper.SetDefault("log_level", "info")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Warn("No config file found, using defaults and environment variables")
		} else {
			logrus.Fatalf("Error reading config file: %v", err)
		}
	}

	logLevel := viper.GetString("log_level")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("Invalid log level: %s", logLevel)
	}
	logrus.SetLevel(level)
}

func main() {
	serverAddress := viper.GetString("server_address")

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)

	if len(os.Args) < 2 {
		logrus.Fatalf("Usage: client <command> [args]")
	}

	switch os.Args[1] {
	case "authenticate":
		if len(os.Args) != 4 {
			logrus.Fatalf("Usage: client authenticate <username> <password>")
		}
		authenticate(c, os.Args[2], os.Args[3])
	case "validate":
		if len(os.Args) != 3 {
			logrus.Fatalf("Usage: client validate <token>")
		}
		validate(c, os.Args[2])
	default:
		logrus.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func authenticate(c pb.AuthServiceClient, username, password string) {
	r, err := c.Authenticate(context.Background(), &pb.AuthenticateRequest{Username: username, Password: password})
	if err != nil {
		logrus.Fatalf("could not authenticate: %v", err)
	}
	logrus.Printf("Token: %s", r.Token)
}

func validate(c pb.AuthServiceClient, token string) {
	r, err := c.Validate(context.Background(), &pb.ValidateRequest{Token: token})
	if err != nil {
		logrus.Fatalf("could not validate: %v", err)
	}
	logrus.Printf("Valid: %t, UserID: %s", r.Valid, r.UserId)
}
