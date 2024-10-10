# MPC Wallet Backend Service

This project provides a backend service for a Multi-Party Computation (MPC) wallet, offering secure key management and transaction signing capabilities.

## Features

- Secure MPC-based key generation and management
- Transaction signing using MPC
- RESTful API for wallet operations
- Docker support for easy deployment

## Prerequisites

- Go 1.21 or higher
- Docker (optional)

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/vietddude/mpcoin.git
   cd mpcoin
   ```

2. Install dependencies:
   ```
   go mod download
   ```

## Usage

### Running locally

1. Build and run the application:

   ```
   go build -o main .
   ./main
   ```

2. The service will be available at `http://localhost:8080` (adjust if using a different port).

### Using Docker

1. Build the Docker image:

   ```
   docker build -t mpc-wallet-backend .
   ```

2. Run the container:
   ```
   docker run -p 8080:8080 mpc-wallet-backend
   ```

For detailed API documentation, please refer to our [API Documentation](API.md).

## Security

This service implements MPC protocols to ensure that no single party has access to the full private key. However, proper key management and secure deployment practices are crucial for maintaining the security of the system.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
