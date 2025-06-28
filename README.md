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

To run the example client that interacts with the server:

```bash
go run cmd/client/main.go
```

*(Note: The client's functionality will depend on its implementation, likely demonstrating authentication requests.)*

## Usage

The Veil Auth service exposes a gRPC API for authentication operations. You can interact with it by generating client stubs in your preferred language using the `.proto` files located in `pkg/grpc/auth/`.

Example gRPC service definition (from `pkg/grpc/auth/auth.proto`):

```protobuf
// Simplified example
service AuthService {
  rpc Authenticate (AuthRequest) returns (AuthResponse);
}
```

Refer to the `pkg/grpc/auth/auth.proto` file for the complete service definition and message structures.

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details. (Note: A `LICENSE` file is assumed to exist or will be added.)
