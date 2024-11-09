package integration_test

import (
	pb "auth-service/proto"
	"context"
	"fmt"
	"log"
	"testing"

	"google.golang.org/grpc"
)

// var c pb.AuthServiceClient
// var conn *grpc.ClientConn

// Setup Test Databases
func createConnection() (pb.AuthServiceClient, *grpc.ClientConn, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, nil, err
	}
	// defer conn.Close()

	c := pb.NewAuthServiceClient(conn)
	fmt.Println("Connected to the server")

	return c, conn, nil

}

// Test GenerateClientID
func TestGenerateClientID(t *testing.T) {
	// Contact the server and print out its response.
	c, conn, err := createConnection()
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	response, err := c.GenerateClientID(context.Background(), &pb.GenerateClientRequest{
		Name:            "Test integartion Client",
		Phone:           "1234567890",
		Email:           "integaration_test@example.com",
		Schema:          map[string]string{"username": "VARCHAR(50)", "password": "VARCHAR(50)"},
		PrimaryKeyField: "username",
	})

	fmt.Println(response)

	if err != nil {
		t.Fatalf("GenerateClientID failed: %v", err)
	}

	if response.ClientId == "" {
		t.Errorf("expected non-empty ClientId, got empty")
	}

	if response.Message != "Client ID generated successfully" {
		t.Errorf("unexpected message: %s", response.Message)
	}
	t.Logf("Client ID: %s", response.ClientId)

}

// Test GetClientID
func TestGetClientID(t *testing.T) {
	// Contact the server and print out its response.
	c, conn, err := createConnection()
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	req := &pb.GetClientRequest{Email: "integaration_test@example.com"}
	resp, err := c.GetClientID(context.Background(), req)
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
	// Contact the server and print out its response.
	c, conn, err := createConnection()
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Assume that "testClient" was created during GenerateClientID test
	signupReq := &pb.SignupRequest{
		ClientId:        "672fb2f21459bef0383323e6",
		UserData:        map[string]string{"username": "newUser", "password": "newPassword"},
		PrimaryKeyField: "username",
	}

	resp, err := c.Signup(context.Background(), signupReq)
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
	// Contact the server and print out its response.
	c, conn, err := createConnection()
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Assume that the user was added during the Signup test
	loginReq := &pb.LoginRequest{
		ClientId:        "672e6755878f1dd94d4aa61d",
		PrimaryKeyField: "username",
		PrimaryKeyValue: "newUser",
		Password:        "newPassword",
	}

	resp, err := c.Login(context.Background(), loginReq)
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
