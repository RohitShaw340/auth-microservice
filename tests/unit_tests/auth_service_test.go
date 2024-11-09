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

// Test GenerateClientID
func TestGenerateClientID(t *testing.T) {
	server := &handlers.AuthServiceServer{}

	req := &pb.GenerateClientRequest{
		Name:            "Test Client 4",
		Phone:           "1234567890",
		Email:           "test4@example.com",
		Schema:          map[string]string{"username": "VARCHAR(50)", "password": "VARCHAR(50)"},
		PrimaryKeyField: "username",
	}

	resp, err := server.GenerateClientID(context.Background(), req)

	if err != nil {
		t.Fatalf("GenerateClientID failed: %v", err)
	}

	if resp.ClientId == "" {
		t.Errorf("expected non-empty ClientId, got empty")
	}

	if resp.Message != "Client ID generated successfully" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
	t.Logf("Client ID: %s", resp.ClientId)
}

// Test GetClientID
func TestGetClientID(t *testing.T) {
	server := &handlers.AuthServiceServer{}

	// // First, create a new client to retrieve
	// genReq := &pb.GenerateClientRequest{
	// 	Name:            "Test Client",
	// 	Phone:           "1234567890",
	// 	Email:           "test2@example.com",
	// 	Schema:          map[string]string{"username": "VARCHAR(50)", "password": "VARCHAR(50)"},
	// 	PrimaryKeyField: "username",
	// }

	// _, err := server.GenerateClientID(context.Background(), genReq)
	// if err != nil {
	// 	t.Fatalf("failed to generate client ID for testing GetClientID: %v", err)
	// }

	req := &pb.GetClientRequest{Email: "test4@example.com"}
	resp, err := server.GetClientID(context.Background(), req)
	if err != nil {
		t.Fatalf("GetClientID failed: %v", err)
	}

	if resp.ClientId == "" {
		t.Errorf("expected non-empty ClientId, got empty")
	}

	if resp.Message != "Client ID retrieved successfully" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
	t.Logf("Client ID: %s", resp.ClientId)
}

// Test Signup
func TestSignup(t *testing.T) {
	server := &handlers.AuthServiceServer{}

	// Assume that "testClient" was created during GenerateClientID test
	signupReq := &pb.SignupRequest{
		ClientId:        "672e6755878f1dd94d4aa61d",
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
	t.Logf("Signup response: %s", resp.Message)
}

// Test Login
func TestLogin(t *testing.T) {
	server := &handlers.AuthServiceServer{}

	// Assume that the user was added during the Signup test
	loginReq := &pb.LoginRequest{
		ClientId:        "672e6755878f1dd94d4aa61d",
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
	t.Logf("User details: %v", resp.UserDetails)
}
