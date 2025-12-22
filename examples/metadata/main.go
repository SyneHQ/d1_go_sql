package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/synehq/d1_go_sql"
)

func main() {
	// Example DSN - replace with your actual credentials
	dsn := "d1://your-account-id:your-api-token@your-database-id?timeout=30s"

	db, err := sql.Open("d1", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("=== D1 Metadata Functions Demo ===")
	fmt.Println()

	// Get current database ID
	var dbID string
	err = db.QueryRow("SELECT current_database()").Scan(&dbID)
	if err != nil {
		log.Fatal("Failed to get current_database():", err)
	}
	fmt.Printf("Current Database ID: %s\n", dbID)

	// Get current database ID using alias
	var dbID2 string
	err = db.QueryRow("SELECT database()").Scan(&dbID2)
	if err != nil {
		log.Fatal("Failed to get database():", err)
	}
	fmt.Printf("Database ID (alias):  %s\n", dbID2)

	// Get current account ID (user)
	var accountID string
	err = db.QueryRow("SELECT current_user()").Scan(&accountID)
	if err != nil {
		log.Fatal("Failed to get current_user():", err)
	}
	fmt.Printf("\nCurrent Account ID:   %s\n", accountID)

	// Get current user using alias
	var accountID2 string
	err = db.QueryRow("SELECT user()").Scan(&accountID2)
	if err != nil {
		log.Fatal("Failed to get user():", err)
	}
	fmt.Printf("Account ID (alias):   %s\n", accountID2)

	// Get driver version
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatal("Failed to get version():", err)
	}
	fmt.Printf("\nDriver Version:       %s\n", version)

	// Get connection ID
	var connID string
	err = db.QueryRow("SELECT connection_id()").Scan(&connID)
	if err != nil {
		log.Fatal("Failed to get connection_id():", err)
	}
	fmt.Printf("\nConnection ID:        %s\n", connID)

	// Demo: Case insensitivity
	fmt.Println()
	fmt.Println("=== Case Insensitivity Demo ===")
	fmt.Println()

	var dbIDUpper string
	err = db.QueryRow("SELECT CURRENT_DATABASE()").Scan(&dbIDUpper)
	if err != nil {
		log.Fatal("Failed to get CURRENT_DATABASE():", err)
	}
	fmt.Printf("CURRENT_DATABASE():   %s\n", dbIDUpper)

	var dbIDMixed string
	err = db.QueryRow("SeLeCt CuRrEnT_DaTaBaSe()").Scan(&dbIDMixed)
	if err != nil {
		log.Fatal("Failed to get mixed case:", err)
	}
	fmt.Printf("Mixed case:           %s\n", dbIDMixed)

	fmt.Println("\n=== All metadata functions executed successfully! ===")
}
