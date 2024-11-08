package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var MySQLClient *sql.DB

func ConnectMySQL(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	MySQLClient = db
	return nil
}

func CreateUserTable(consumerID string, schema map[string]string) error {
	// Create a new MySQL database for this consumer
	dbName := fmt.Sprintf("consumer_%s", consumerID)
	_, err := MySQLClient.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	fmt.Println("Database created successfully")

	// Switch to the new database
	_, err = MySQLClient.Exec(fmt.Sprintf("USE %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to switch to database %s: %w", dbName, err)
	}

	// Construct the CREATE TABLE statement
	tableSchema := "CREATE TABLE IF NOT EXISTS users ("
	for field, dataType := range schema {
		tableSchema += fmt.Sprintf("%s %s,", field, dataType)
	}
	tableSchema = tableSchema[:len(tableSchema)-1] + ")" // Remove trailing comma and add closing parenthesis

	// Execute the CREATE TABLE statement
	_, err = MySQLClient.Exec(tableSchema)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("User table created successfully")
	return nil
}
