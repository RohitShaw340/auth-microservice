package handlers

import (
	"auth-service/db"
	pb "auth-service/proto"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

// GenerateClientID generates a unique client ID for a new client.
func (s *AuthServiceServer) GenerateClientID(ctx context.Context, req *pb.GenerateClientRequest) (*pb.GenerateClientResponse, error) {
	collection := db.GetClientsCollection()
	client := map[string]interface{}{
		"name":              req.Name,
		"phone":             req.Phone,
		"email":             req.Email,
		"user_schema":       req.Schema,
		"primary_key_field": req.PrimaryKeyField,
	}
	res, err := collection.InsertOne(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to insert client: %w", err)
	}

	// Convert the MongoDB ObjectID to a string
	clientID := res.InsertedID.(primitive.ObjectID).Hex()

	err = db.CreateUserTable(clientID, req.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create user table: %w", err)
	}
	return &pb.GenerateClientResponse{ClientId: clientID, Message: "Client ID generated successfully"}, nil
}

// GetClientID retrieves the client ID using the provided email.
func (s *AuthServiceServer) GetClientID(ctx context.Context, req *pb.GetClientRequest) (*pb.GetClientResponse, error) {
	collection := db.GetClientsCollection()

	// Search for the document by email
	filter := bson.M{"email": req.Email}
	var client bson.M
	err := collection.FindOne(ctx, filter).Decode(&client)
	if err != nil {
		return nil, fmt.Errorf("failed to find client: %w", err)
	}

	// Extract the client ID
	clientID, ok := client["_id"].(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("invalid client ID format")
	}

	return &pb.GetClientResponse{
		ClientId: clientID.Hex(),
		Message:  "Client ID retrieved successfully",
	}, nil
}
