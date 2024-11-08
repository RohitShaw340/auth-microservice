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

// GenerateConsumerID generates a unique consumer ID for a new consumer.
func (s *AuthServiceServer) GenerateConsumerID(ctx context.Context, req *pb.GenerateConsumerRequest) (*pb.GenerateConsumerResponse, error) {
	collection := db.GetConsumersCollection()
	consumer := map[string]interface{}{
		"name":              req.Name,
		"phone":             req.Phone,
		"email":             req.Email,
		"user_schema":       req.Schema,
		"primary_key_field": req.PrimaryKeyField,
	}
	res, err := collection.InsertOne(ctx, consumer)
	if err != nil {
		return nil, fmt.Errorf("failed to insert consumer: %w", err)
	}

	// Convert the MongoDB ObjectID to a string
	consumerID := res.InsertedID.(primitive.ObjectID).Hex()

	err = db.CreateUserTable(consumerID, req.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create user table: %w", err)
	}
	return &pb.GenerateConsumerResponse{ConsumerId: consumerID, Message: "Consumer ID generated successfully"}, nil
}

// GetConsumerID retrieves the consumer ID using the provided email.
func (s *AuthServiceServer) GetConsumerID(ctx context.Context, req *pb.GetConsumerRequest) (*pb.GetConsumerResponse, error) {
	collection := db.GetConsumersCollection()

	// Search for the document by email
	filter := bson.M{"email": req.Email}
	var consumer bson.M
	err := collection.FindOne(ctx, filter).Decode(&consumer)
	if err != nil {
		return nil, fmt.Errorf("failed to find consumer: %w", err)
	}

	// Extract the consumer ID
	consumerID, ok := consumer["_id"].(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("invalid consumer ID format")
	}

	return &pb.GetConsumerResponse{
		ConsumerId: consumerID.Hex(),
		Message:    "Consumer ID retrieved successfully",
	}, nil
}
