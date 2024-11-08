// handlers/signup.go
package handlers

import (
	"auth-service/db"
	pb "auth-service/proto"
	"context"
	"fmt"
)

func (s *AuthServiceServer) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	dbName := fmt.Sprintf("client_%s", req.ClientId)
	tableName := "users"

	fields := ""
	placeholders := ""
	values := []interface{}{}
	for field, value := range req.UserData {
		fields += fmt.Sprintf("%s, ", field)
		placeholders += "?, "
		values = append(values, value)
	}
	fields = fields[:len(fields)-2]
	placeholders = placeholders[:len(placeholders)-2]

	query := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s)", dbName, tableName, fields, placeholders)
	_, err := db.MySQLClient.Exec(query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &pb.SignupResponse{
		Message: "Signup successful",
	}, nil
}
