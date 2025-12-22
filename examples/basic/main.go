package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/synehq/d1_go_sql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Get credentials from environment variables
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	databaseID := os.Getenv("D1_DATABASE_ID")

	if accountID == "" || apiToken == "" || databaseID == "" {
		log.Fatal("Please set CLOUDFLARE_ACCOUNT_ID, CLOUDFLARE_API_TOKEN, and D1_DATABASE_ID environment variables")
	}

	// Build DSN
	dsn := fmt.Sprintf("d1://%s:%s@%s?timeout=30s", accountID, apiToken, databaseID)

	// Open database connection
	db, err := sql.Open("d1", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✓ Connected to D1 database successfully")

	// Run examples
	if err := createTableExample(db); err != nil {
		log.Fatalf("Create table example failed: %v", err)
	}

	if err := insertDataExample(db); err != nil {
		log.Fatalf("Insert data example failed: %v", err)
	}

	if err := queryDataExample(db); err != nil {
		log.Fatalf("Query data example failed: %v", err)
	}

	if err := preparedStatementExample(db); err != nil {
		log.Fatalf("Prepared statement example failed: %v", err)
	}

	if err := transactionExample(db); err != nil {
		log.Fatalf("Transaction example failed: %v", err)
	}

	fmt.Println("\n✓ All examples completed successfully!")
}

func createTableExample(db *sql.DB) error {
	fmt.Println("\n--- Create Table Example ---")

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("✓ Table 'users' created successfully")
	return nil
}

func insertDataExample(db *sql.DB) error {
	fmt.Println("\n--- Insert Data Example ---")

	result, err := db.Exec(
		"INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		"Alice Johnson",
		fmt.Sprintf("alice_%d@example.com", time.Now().Unix()),
		28,
	)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	lastID, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()

	fmt.Printf("✓ Inserted user with ID: %d (rows affected: %d)\n", lastID, rowsAffected)
	return nil
}

func queryDataExample(db *sql.DB) error {
	fmt.Println("\n--- Query Data Example ---")

	// Query single row
	var name string
	var email string
	err := db.QueryRow("SELECT name, email FROM users WHERE id = ?", 1).Scan(&name, &email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to query single row: %w", err)
	}

	if err == sql.ErrNoRows {
		fmt.Println("No user found with ID 1")
	} else {
		fmt.Printf("✓ User #1: %s (%s)\n", name, email)
	}

	// Query multiple rows
	rows, err := db.Query("SELECT id, name, email, age FROM users ORDER BY id LIMIT 10")
	if err != nil {
		return fmt.Errorf("failed to query multiple rows: %w", err)
	}
	defer rows.Close()

	fmt.Println("\n✓ All users:")
	count := 0
	for rows.Next() {
		var id, age int
		var name, email string
		if err := rows.Scan(&id, &name, &email, &age); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		fmt.Printf("  - [%d] %s (%s) - Age: %d\n", id, name, email, age)
		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	fmt.Printf("✓ Retrieved %d users\n", count)
	return nil
}

func preparedStatementExample(db *sql.DB) error {
	fmt.Println("\n--- Prepared Statement Example ---")

	stmt, err := db.Prepare("INSERT INTO users (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	users := []struct {
		name  string
		email string
		age   int
	}{
		{"Bob Smith", fmt.Sprintf("bob_%d@example.com", time.Now().Unix()), 32},
		{"Carol White", fmt.Sprintf("carol_%d@example.com", time.Now().Unix()), 25},
		{"David Brown", fmt.Sprintf("david_%d@example.com", time.Now().Unix()), 40},
	}

	for _, user := range users {
		result, err := stmt.Exec(user.name, user.email, user.age)
		if err != nil {
			return fmt.Errorf("failed to execute prepared statement: %w", err)
		}

		lastID, _ := result.LastInsertId()
		fmt.Printf("✓ Inserted %s with ID: %d\n", user.name, lastID)
	}

	return nil
}

func transactionExample(db *sql.DB) error {
	fmt.Println("\n--- Transaction Example ---")

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Insert multiple users in a transaction
	_, err = tx.Exec(
		"INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		"Eve Wilson",
		fmt.Sprintf("eve_%d@example.com", time.Now().Unix()),
		35,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert in transaction: %w", err)
	}

	_, err = tx.Exec(
		"INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		"Frank Miller",
		fmt.Sprintf("frank_%d@example.com", time.Now().Unix()),
		29,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert in transaction: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Println("✓ Transaction committed successfully (2 users inserted)")
	return nil
}
