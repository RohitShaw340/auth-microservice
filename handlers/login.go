// handlers/login.go
package handlers

import (
	"auth-service/db"
	pb "auth-service/proto"
	"context"
	"database/sql"
	"fmt"
)

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Construct the database and table name
	dbName := fmt.Sprintf("consumer_%s", req.ConsumerId)
	tableName := "users"
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE %s = ? AND password = ?", dbName, tableName, req.PrimaryKeyField)

	// Query for user data
	rows, err := db.MySQLClient.Query(query, req.PrimaryKeyValue, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to execute login query: %w", err)
	}
	defer rows.Close()

	// Check if user exists
	if !rows.Next() {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Get columns and map the result into a map[string]string
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve columns: %w", err)
	}

	userData := make(map[string]string)
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(sql.NullString)
	}

	// Scan the result into values and map it to userData
	if err := rows.Scan(values...); err != nil {
		return nil, fmt.Errorf("failed to scan user data: %w", err)
	}
	for i, col := range columns {
		if val, ok := values[i].(*sql.NullString); ok && val.Valid {
			userData[col] = val.String
		} else {
			userData[col] = ""
		}
	}

	return &pb.LoginResponse{
		UserDetails: userData,
		Message:     "Login successful",
	}, nil
}
