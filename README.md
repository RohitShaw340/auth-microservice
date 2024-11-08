# auth-microservice

## Setup

### Prerequisites

- Go 1.16 or higher
- Docker (optional, for containerization)
- Make sure you have `GOPATH` set up correctly

### Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/auth-microservice.git
   cd auth-microservice
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Run the service:
   ```sh
   go run main.go
   ```

### Environment Variables

Create a `.env` file in the root directory and add the following variables:

```env
SQL_TEST_DSN="username:password@/"
MONGODB_TEST_DSN="MongoDB URL"
```

### Docker

To run the service in a Docker container:

1. Build the Docker image:

   ```sh
   docker build -t auth-microservice .
   ```

2. Run the Docker container:
   ```sh
   docker run -d -p 8080:8080 --env-file .env auth-microservice
   ```

### Testing

Run tests using:

```sh
go test ./...
```

### Contributing

Feel free to submit issues, fork the repository and send pull requests!

### License

This project is licensed under the MIT License.
