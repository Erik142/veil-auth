package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/Erik142/veil-auth/internal/auth"
	grpcAuth "github.com/Erik142/veil-auth/internal/grpc/auth"
	pb "github.com/Erik142/veil-auth/pkg/grpc/auth"
)

func init() {
	viper.SetDefault("port", 50051)
	viper.SetDefault("secret_key", "your-secret-key")
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
	port := viper.GetInt("port")
	secretKey := viper.GetString("secret_key")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	authenticator := auth.NewInMemoryAuthenticator(secretKey)
	authenticator.AddUser("testuser", "testpassword")

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, grpcAuth.NewServer(authenticator))

	logrus.Infof("gRPC server listening on :%d", port)
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
