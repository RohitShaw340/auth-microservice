package handlers_test

import (
	"auth-service/db"
	"auth-service/handlers"
	pb "auth-service/proto"
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// Setup Test Databases
func init() {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		panic("failed to load environment variables: " + err.Error())
	}
	// Replace with your test MongoDB and MySQL DSNs
	err = db.ConnectMongoDB(os.Getenv("MONGODB_TEST_DSN"))
	if err != nil {
		panic("failed to connect to MongoDB: " + err.Error())
	}

	err = db.ConnectMySQL(os.Getenv("MYSQL_TEST_DSN"))
	if err != nil {
		panic("failed to connect to MySQL: " + err.Error())
	}
}

// Test GenerateConsumerID
func TestGenerateConsumerID(t *testing.T) {
	server := &handlers.AuthServiceServer{}

	req := &pb.GenerateConsumerRequest{
		Name:            "Test Consumer",
		Phone:           "1234567890",
		Email:           "test@example.com",
		Schema:          map[string]string{"username": "VARCHAR(50)", "password": "VARCHAR(50)"},
		PrimaryKeyField: "username",
	}

	resp, err := server.GenerateConsumerID(context.Background(), req)
	if err != nil {
		t.Fatalf("GenerateConsumerID failed: %v", err)
	}

	if resp.ConsumerId == "" {
		t.Errorf("expected non-empty ConsumerId, got empty")
	}

	if resp.Message != "Consumer ID generated successfully" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

// Test GetConsumerID
func TestGetConsumerID(t *testing.T) {
	server := &handlers.AuthServiceServer{}

	// // First, create a new consumer to retrieve
	// genReq := &pb.GenerateConsumerRequest{
	// 	Name:            "Test Consumer",
	// 	Phone:           "1234567890",
	// 	Email:           "test2@example.com",
	// 	Schema:          map[string]string{"username": "VARCHAR(50)", "password": "VARCHAR(50)"},
	// 	PrimaryKeyField: "username",
	// }

	// _, err := server.GenerateConsumerID(context.Background(), genReq)
	// if err != nil {
	// 	t.Fatalf("failed to generate consumer ID for testing GetConsumerID: %v", err)
	// }

	req := &pb.GetConsumerRequest{Email: "test2@example.com"}
	resp, err := server.GetConsumerID(context.Background(), req)
	if err != nil {
		t.Fatalf("GetConsumerID failed: %v", err)
	}

	if resp.ConsumerId == "" {
		t.Errorf("expected non-empty ConsumerId, got empty")
	}

	if resp.Message != "Consumer ID retrieved successfully" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

// Test Signup
func TestSignup(t *testing.T) {
	server := &handlers.AuthServiceServer{}

	// Assume that "testConsumer" was created during GenerateConsumerID test
	signupReq := &pb.SignupRequest{
		ConsumerId:      "672e52cb0dee8ecc7fa2fabf",
		UserData:        map[string]string{"username": "newUser", "password": "newPassword"},
		PrimaryKeyField: "username",
	}

	resp, err := server.Signup(context.Background(), signupReq)
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}

	if resp.Message != "Signup successful" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

// Test Login
func TestLogin(t *testing.T) {
	server := &handlers.AuthServiceServer{}

	// Assume that the user was added during the Signup test
	loginReq := &pb.LoginRequest{
		ConsumerId:      "672e52cb0dee8ecc7fa2fabf",
		PrimaryKeyField: "username",
		PrimaryKeyValue: "newUser",
		Password:        "newPassword",
	}

	resp, err := server.Login(context.Background(), loginReq)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if resp.Message != "Login successful" {
		t.Errorf("unexpected message: %s", resp.Message)
	}

	if resp.UserDetails["username"] != "newUser" || resp.UserDetails["password"] != "newPassword" {
		t.Errorf("unexpected user details: %v", resp.UserDetails)
	}
}
