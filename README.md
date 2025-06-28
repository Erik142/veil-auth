# Veil Auth

Veil Auth is a robust and efficient authentication service built with Go, leveraging gRPC for high-performance communication. It provides a secure and scalable solution for managing user authentication within your applications, featuring an in-memory authentication mechanism for rapid prototyping and development.

## Features

-   **gRPC API:** Exposes a well-defined gRPC interface for seamless integration with various client applications.
-   **In-Memory Authentication:** Offers a simple yet effective in-memory store for user credentials, ideal for development, testing, and scenarios where persistent storage is not required.
-   **Client-Server Architecture:** Clearly separated client and server components for modular development and deployment.
-   **Go-Native:** Developed entirely in Go, ensuring excellent performance, concurrency, and ease of deployment.

## Project Structure

The repository is organized into the following key directories:

-   `cmd/`: Contains the main entry points for the client and server applications.
    -   `client/`: The client application that interacts with the authentication service.
    -   `server/`: The gRPC authentication server.
-   `internal/`: Houses internal, non-exportable packages and implementations.
    -   `auth/`: Core authentication logic, including the in-memory authentication provider.
    -   `grpc/`: gRPC-specific implementations, such as the authentication server.
-   `pkg/`: Contains public, exportable packages intended for use by other projects or modules.
    -   `auth/`: Common authentication interfaces and types.
    -   `grpc/`: Generated gRPC protobuf files and service definitions.

## Getting Started

To get started with Veil Auth, follow these steps:

### Prerequisites

-   Go (version 1.16 or higher recommended)
-   Protocol Buffers compiler (`protoc`)
-   `protoc-gen-go` and `protoc-gen-go-grpc` plugins

### Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/Erik142/veil-auth.git
    cd veil-auth
    ```

2.  Install dependencies:

    ```bash
    go mod tidy
    ```

### Running the Server

To start the authentication server:

```bash
go run cmd/server/main.go
```

The server will typically listen on `localhost:50051` (or a configured port).

### Running the Client

To build the client:

```bash
go build -o cmd/client/client cmd/client/main.go
```

To run the client and interact with the server:

```bash
./cmd/client/client
```

This will display the available commands and flags.

## Usage

The Veil Auth client now uses Cobra for command-line interface management.

### Global Flags

-   `--server-address <address>`: gRPC server address (default: `localhost:50051`)
-   `--log-level <level>`: Log level (debug, info, warn, error, fatal, panic) (default: `info`)

### Authenticate Command

Authenticates a user with the provided username and password, returning an authentication token.

```bash
./cmd/client/client authenticate -u <username> -p <password>
```

**Flags:**

-   `-u, --username <username>`: Username for authentication (required)
-   `-p, --password <password>`: Password for authentication (required)

### Validate Command

Validates an authentication token, returning its validity status and associated user ID.

```bash
./cmd/client/client validate -t <token>
```

**Flags:**

-   `-t, --token <token>`: Token to validate (required)

Example:

```bash
# Authenticate a user
./cmd/client/client authenticate -u testuser -p testpassword

# Validate a token
./cmd/client/client validate -t <received_token>
```

Refer to the `pkg/grpc/auth/auth.proto` file for the complete service definition and message structures.

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details. (Note: A `LICENSE` file is assumed to exist or will be added.)
