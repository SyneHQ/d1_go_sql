/*
Package d1sql provides a database/sql driver for Cloudflare D1.

D1 is Cloudflare's native serverless SQL database. This package uses the
official Cloudflare Go SDK (github.com/cloudflare/cloudflare-go/v6) to provide
a standard database/sql interface for D1, enabling seamless integration with Go
applications.

# Installation

	go get github.com/synehq/d1_go_sql

# Basic Usage

	import (
		"database/sql"
		_ "github.com/synehq/d1_go_sql"
	)

	func main() {
		db, err := sql.Open("d1", "d1://accountID:apiToken@databaseID")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Use db with standard database/sql operations
	}

# Connection String Format

The driver accepts connection strings in the following format:

	d1://accountID:apiToken@databaseID?timeout=30s

Parameters:
  - accountID: Your Cloudflare account ID (required)
  - apiToken: Your Cloudflare API token (required)
  - databaseID: Your D1 database ID (required)
  - timeout: Request timeout duration (optional, default: 30s)

# Features

  - Full database/sql interface implementation
  - Official Cloudflare Go SDK (v6) integration
  - Connection pooling support
  - Prepared statements with parameter binding
  - Transaction support
  - Context-aware operations
  - Metadata commands and functions (LIST DATABASES, current_database, current_user, version, connection_id)
  - Comprehensive type conversion
  - Production-ready error handling

# Type Mapping

The driver supports the following type conversions:

	Go Type       D1 Type
	-------       -------
	nil           NULL
	int64         INTEGER
	float64       REAL
	bool          INTEGER (0/1)
	string        TEXT
	[]byte        BLOB
	time.Time     TEXT (RFC3339)

# Examples

Query with parameters:

	rows, err := db.Query("SELECT * FROM users WHERE age > ?", 18)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		if err := rows.Scan(&id, &name, &age); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User: %d, %s, %d\n", id, name, age)
	}

Prepared statements:

	stmt, err := db.Prepare("INSERT INTO users (name, age) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec("Alice", 30)
	if err != nil {
		log.Fatal(err)
	}

Transactions:

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO users (name) VALUES (?)", "Bob")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

# Best Practices

1. Always close resources using defer:

	db.Close()
	rows.Close()
	stmt.Close()

2. Use prepared statements for repeated queries
3. Set appropriate timeouts using context
4. Configure connection pool settings:

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

5. Handle errors appropriately

# Error Handling

The driver provides detailed error messages with context. Use standard
error handling patterns:

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Handle no rows case
		}
		// Other error handling
	}

# Thread Safety

The driver is safe for concurrent use by multiple goroutines, as required
by the database/sql package.

# Performance Considerations

- Use prepared statements for repeated queries
- Configure connection pool appropriately for your workload
- Set reasonable timeouts to avoid hanging requests
- Use transactions for bulk operations

For more information and examples, visit:
https://github.com/synehq/d1_go_sql
*/
package d1sql
