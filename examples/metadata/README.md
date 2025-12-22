# Metadata Functions Example

This example demonstrates the use of special metadata functions in the D1 Go SQL driver.

## Metadata Functions

The driver supports several metadata functions that return connection information without making API calls to D1:

### Available Functions

1. **`LIST DATABASES`**
   - Lists all D1 databases in the account
   - Returns columns: name, uuid, version
   - Makes an API call
   - Example: `LIST DATABASES`

2. **`current_database()` / `database()`**
   - Returns the current D1 database name
   - Makes an API call to fetch the actual name
   - Example: `SELECT current_database()`

3. **`current_user()` / `user()`**
   - Returns the Cloudflare account ID
   - Example: `SELECT current_user()`

4. **`version()`**
   - Returns driver and database version information
   - No API call
   - Example: `SELECT version()`

5. **`connection_id()`**
   - Returns a unique connection identifier (format: `account-id:database-id`)
   - No API call
   - Example: `SELECT connection_id()`

## Features

- **Case Insensitive**: All functions work with any case (lowercase, uppercase, mixed)
- **API Calls**: `LIST DATABASES` and `current_database()` make API calls; others execute instantly
- **Standard SQL Interface**: Use them just like regular SQL queries

## Usage

```go
// List all databases
rows, err := db.Query("LIST DATABASES")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var name, uuid, version string
    rows.Scan(&name, &uuid, &version)
    fmt.Printf("Database: %s (UUID: %s)\n", name, uuid)
}

// Get database name
var dbName string
db.QueryRow("SELECT current_database()").Scan(&dbName)

// Get account ID
var accountID string
db.QueryRow("SELECT current_user()").Scan(&accountID)

// Get version
var version string
db.QueryRow("SELECT version()").Scan(&version)

// Get connection ID
var connID string
db.QueryRow("SELECT connection_id()").Scan(&connID)
```

## Running the Example

1. Update the DSN in `main.go` with your credentials:

   ```go
   dsn := "d1://your-account-id:your-api-token@your-database-id"
   ```

2. Run the example:

   ```bash
   go run main.go
   ```

## Expected Output

```text
=== D1 Metadata Functions Demo ===

Current Database ID: your-database-id
Database ID (alias):  your-database-id

Current Account ID:   your-account-id
Account ID (alias):   your-account-id

Driver Version:       D1 Go SQL Driver v0.1.0 (Cloudflare D1 - SQLite compatible)

Connection ID:        your-account-id:your-database-id

=== Case Insensitivity Demo ===

CURRENT_DATABASE():   your-database-id
Mixed case:           your-database-id

=== All metadata functions executed successfully! ===
```

## Use Cases

1. **Connection Validation**: Verify you're connected to the correct database
2. **Logging**: Include database/account info in application logs
3. **Multi-tenancy**: Identify which database is being used
4. **Debugging**: Quick access to connection information
5. **Monitoring**: Track which databases are in use
