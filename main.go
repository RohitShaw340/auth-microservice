package main

import (
	"auth-service/db"
	"auth-service/handlers"
	pb "auth-service/proto"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	// MongoDB Connection
	if err := db.ConnectMongoDB(os.Getenv("MONGODB_TEST_DSN")); err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	// MySQL Connection
	if err := db.ConnectMySQL(os.Getenv("MYSQL_TEST_DSN")); err != nil {
		log.Fatalf("MySQL connection failed: %v", err)
	}

	// gRPC Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &handlers.AuthServiceServer{})
	log.Printf("Server is listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
