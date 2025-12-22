# D1 Go SQL Driver

A production-grade Go `database/sql` driver for Cloudflare D1, enabling seamless integration with Go's standard database interfaces.

## Features

- ✅ Full `database/sql` interface implementation
- ✅ Official Cloudflare Go SDK (v6) integration
- ✅ Connection pooling support
- ✅ Prepared statements with parameter binding
- ✅ Transaction support
- ✅ Proper type conversion (string, int64, float64, bool, []byte, time.Time, nil)
- ✅ Context-aware operations
- ✅ Metadata functions (current_database, current_user, version, connection_id)
- ✅ Comprehensive error handling
- ✅ Production-ready with best practices

## Requirements

- Go 1.22 or later

## Installation

```bash
go get github.com/synehq/d1_go_sql
```

## Quick Start

```go
package main

import (
    "database/sql"
    "log"
    
    _ "github.com/synehq/d1_go_sql"
)

func main() {
    // Open database connection
    // Format: d1://accountID:apiToken@databaseID
    db, err := sql.Open("d1", "d1://your-account-id:your-api-token@your-database-id")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Execute queries
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
        log.Printf("User: %d, %s, %d\n", id, name, age)
    }
}
```

## Connection String Format

The driver accepts connection strings in the following format:

```text
d1://accountID:apiToken@databaseID
```

Or using DSN parameters:

```text
d1://accountID:apiToken@databaseID?timeout=30s
```

### Parameters

- `accountID`: Your Cloudflare account ID
- `apiToken`: Your Cloudflare API token
- `databaseID`: Your D1 database ID
- `timeout`: (Optional) Request timeout duration (default: 30s)

## Usage Examples

### Basic Queries

```go
// Query single row
var name string
err := db.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&name)

// Query multiple rows
rows, err := db.Query("SELECT id, name FROM users")
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    rows.Scan(&id, &name)
}
```

### Prepared Statements

```go
stmt, err := db.Prepare("INSERT INTO users (name, age) VALUES (?, ?)")
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()

result, err := stmt.Exec("Alice", 30)
if err != nil {
    log.Fatal(err)
}

lastID, _ := result.LastInsertId()
affected, _ := result.RowsAffected()
```

### Transactions

```go
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
```

### Context Support

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

rows, err := db.QueryContext(ctx, "SELECT * FROM users")
```

### Metadata Functions

The driver supports special metadata functions that return connection information without querying D1:

```go
// Get current database ID
var dbName string
db.QueryRow("SELECT current_database()").Scan(&dbName)
// Returns: your-database-id

// Get current account ID (user)
var accountID string
db.QueryRow("SELECT current_user()").Scan(&accountID)
// Returns: your-account-id

// Get driver version
var version string
db.QueryRow("SELECT version()").Scan(&version)
// Returns: D1 Go SQL Driver v0.1.0 (Cloudflare D1 - SQLite compatible)

// Get connection ID
var connID string
db.QueryRow("SELECT connection_id()").Scan(&connID)
// Returns: account-id:database-id
```

**Supported Functions:**
- `current_database()` / `database()` - Returns the D1 database ID
- `current_user()` / `user()` - Returns the Cloudflare account ID
- `version()` - Returns driver and database version information
- `connection_id()` - Returns a unique connection identifier

These functions are case-insensitive and execute instantly without API calls.

## Architecture

The driver implements the following `database/sql/driver` interfaces:

- `Driver` - Main driver registration
- `DriverContext` - Context-aware driver operations
- `Connector` - Connection factory
- `Conn` - Database connection
- `ConnPrepareContext` - Prepared statement creation
- `Stmt` - Prepared statement
- `StmtExecContext` - Statement execution with context
- `StmtQueryContext` - Statement querying with context
- `Rows` - Result set iteration
- `Result` - Execution results

## Type Mapping

| Go Type | D1 Type |
|---------|---------|
| `nil` | NULL |
| `int`, `int64` | INTEGER |
| `float64` | REAL |
| `bool` | INTEGER (0/1) |
| `string` | TEXT |
| `[]byte` | BLOB |
| `time.Time` | TEXT (RFC3339) |

## Error Handling

The driver provides detailed error messages with context:

```go
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        // Handle no rows case
    }
    // Other error handling
}
```

## Best Practices

1. **Always close resources**: Use `defer` to close connections, statements, and rows
2. **Use prepared statements**: For repeated queries to improve performance
3. **Use contexts**: Set appropriate timeouts for operations
4. **Connection pooling**: Configure `SetMaxOpenConns()` and `SetMaxIdleConns()`
5. **Error handling**: Always check and handle errors appropriately

## Configuration

```go
db, err := sql.Open("d1", dsn)
if err != nil {
    log.Fatal(err)
}

// Configure connection pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(1 * time.Minute)
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue on GitHub.
