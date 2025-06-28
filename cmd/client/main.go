package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Erik142/veil-auth/pkg/grpc/auth"
)

var (
	serverAddress string
	logLevel      string
	authClient    pb.AuthServiceClient
	grpcConn      *grpc.ClientConn
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&serverAddress, "server-address", "localhost:50051", "gRPC server address")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error, fatal, panic)")

	viper.BindPFlag("server_address", rootCmd.PersistentFlags().Lookup("server-address"))
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	authenticateCmd.Flags().StringP("username", "u", "", "Username for authentication")
	authenticateCmd.Flags().StringP("password", "p", "", "Password for authentication")
	authenticateCmd.MarkFlagRequired("username")
	authenticateCmd.MarkFlagRequired("password")

	validateCmd.Flags().StringP("token", "t", "", "Token to validate")
	validateCmd.MarkFlagRequired("token")

	rootCmd.AddCommand(authenticateCmd)
	rootCmd.AddCommand(validateCmd)
}

func initConfig() {
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

	parsedLevel, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		logrus.Fatalf("Invalid log level: %s", viper.GetString("log_level"))
	}
	logrus.SetLevel(parsedLevel)
}

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "veil-auth gRPC client",
	Long:  `A gRPC client for the veil-auth service, providing authentication and token validation.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		grpcConn, err = grpc.NewClient(viper.GetString("server_address"), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		authClient = pb.NewAuthServiceClient(grpcConn)
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if grpcConn != nil {
			grpcConn.Close()
		}
	},
}

var authenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "Authenticate a user",
	Long:  `Authenticates a user with the provided username and password, returning an authentication token.`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		r, err := authClient.Authenticate(context.Background(), &pb.AuthenticateRequest{Username: username, Password: password})
		if err != nil {
			logrus.Fatalf("could not authenticate: %v", err)
		}
		logrus.Printf("Token: %s", r.Token)
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an authentication token",
	Long:  `Validates an authentication token, returning its validity status and associated user ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")

		r, err := authClient.Validate(context.Background(), &pb.ValidateRequest{Token: token})
		if err != nil {
			logrus.Fatalf("could not validate: %v", err)
		}
		logrus.Printf("Valid: %t, UserID: %s", r.Valid, r.UserId)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
