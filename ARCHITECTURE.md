# D1 Go SQL Driver - Architecture

This document provides an overview of the architectural design and implementation details of the D1 Go SQL Driver.

## Overview

The D1 Go SQL Driver is a production-grade wrapper around Cloudflare's D1 database, leveraging the official Cloudflare Go SDK (`github.com/cloudflare/cloudflare-go/v6`) to provide seamless integration with Go's standard `database/sql` package. The implementation follows Go best practices, the DRY principle, and is optimized for performance and maintainability.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│                  (database/sql package)                     │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   D1 SQL Driver (d1sql)                     │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Driver    │  │  Connector  │  │     Conn    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │    Stmt     │  │    Rows     │  │   Result    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐                         │
│  │     Tx      │  │   Utils     │                         │
│  └─────────────┘  └─────────────┘                         │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│         Cloudflare Go SDK (cloudflare-go/v6)               │
│              D1 Database API Client                         │
│      (https://api.cloudflare.com/client/v4)                │
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. Driver (`driver.go`)

**Responsibilities:**
- Implements `database/sql/driver.Driver` interface
- Implements `driver.DriverContext` interface
- Handles driver registration with `database/sql`
- Parses and validates DSN (Data Source Name)
- Creates connectors for database connections

**Key Features:**
- DSN parsing with comprehensive validation
- Support for connection parameters (timeout, etc.)
- Configuration validation
- Error handling with detailed messages

**DSN Format:**
```
d1://accountID:apiToken@databaseID?timeout=30s
```

### 2. Connector (`driver.go`)

**Responsibilities:**
- Implements `driver.Connector` interface
- Factory for creating connections
- Manages configuration

**Key Features:**
- Connection pooling support
- Configuration validation
- Cloudflare SDK client creation with proper timeouts

### 3. Connection (`conn.go`)

**Responsibilities:**
- Implements `driver.Conn` interface
- Implements `driver.ConnPrepareContext` interface
- Implements `driver.ExecerContext` interface
- Implements `driver.QueryerContext` interface
- Manages database connection lifecycle
- Handles transactions

**Key Features:**
- Context-aware operations
- Transaction support with Begin/Commit/Rollback
- Prepared statement creation
- Query execution (both Exec and Query)
- Connection pooling compatibility
- Cloudflare SDK integration for D1 operations

**Implementation Details:**
- Uses official Cloudflare Go SDK (`github.com/cloudflare/cloudflare-go/v6`)
- Leverages SDK's D1 Database Query API
- Handles authentication via SDK's API token option
- Automatic pagination support via SDK's AutoPaging iterators
- Manages transaction state with statement queuing

### 4. Statement (`stmt.go`)

**Responsibilities:**
- Implements `driver.Stmt` interface
- Implements `driver.StmtExecContext` interface
- Implements `driver.StmtQueryContext` interface
- Manages prepared statement lifecycle

**Key Features:**
- Parameter binding support
- Context-aware execution
- Reusable for multiple executions
- Proper resource cleanup

### 5. Rows (`rows.go`)

**Responsibilities:**
- Implements `driver.Rows` interface
- Manages result set iteration
- Handles column metadata

**Key Features:**
- Efficient row iteration
- Column name access
- Type conversion for result values
- Proper resource cleanup

### 6. Result (`result.go`)

**Responsibilities:**
- Implements `driver.Result` interface
- Provides execution metadata

**Key Features:**
- Last insert ID tracking
- Rows affected counting

### 7. Transaction (`conn.go`)

**Responsibilities:**
- Implements `driver.Tx` interface
- Manages transaction lifecycle

**Key Features:**
- Statement queuing during transaction
- Batch execution on commit using D1's batch API
- Rollback support (discards buffered statements)
- Nested transaction prevention

**Implementation:**
- D1 doesn't support SQL `BEGIN TRANSACTION`/`COMMIT` statements
- Instead, uses `DatabaseQueryParamsBodyMultipleQueries` with an array of SQL statements
- All statements within a transaction are buffered in memory
- On `Commit()`, statements are sent as a batch array to D1's batch API
- On `Rollback()`, buffered statements are discarded without API calls
- This provides atomic execution of all statements in the transaction

### 8. Utilities (`utils.go`)

**Responsibilities:**
- Type conversion between Go and SQL types
- Parameter binding
- Value formatting for SQL

**Key Functions:**
- `bindParameters()`: Replaces placeholders with actual values
- `valueToSQLLiteral()`: Converts Go values to SQL literals
- `convertValue()`: Converts D1 response values to Go types
- `namedValuesToValues()`: Converts named parameters to positional
- `namedValuesFromValues()`: Converts positional to named parameters

**Type Mapping:**
```
Go Type       →  SQL Literal
---------        ------------
nil           →  NULL
int64         →  42
float64       →  3.14
bool          →  1 or 0
string        →  'value'
[]byte        →  X'hex'
time.Time     →  '2023-01-01T12:00:00Z'
```

## Data Flow

### Query Execution Flow

```
1. Application calls db.Query(sql, args...)
2. Driver converts args to driver.Value
3. Statement binds parameters to SQL
4. Connection executes query via Cloudflare SDK
5. SDK calls D1 API and returns structured response
6. Driver parses SDK response to Rows
7. Application iterates through Rows
8. Driver converts values to Go types
9. Application receives final results
```

### Transaction Flow

```
1. Application calls db.Begin()
2. Driver creates Tx and sets inTx flag
3. Application executes statements via tx.Exec()
4. Driver queues statements (doesn't execute)
5. Application calls tx.Commit()
6. Driver batches all statements
7. Driver executes batch via Cloudflare SDK
8. SDK sends batch to D1 API for processing
9. Driver returns result to application
```

## Cloudflare SDK Integration

The driver uses the official Cloudflare Go SDK (`github.com/cloudflare/cloudflare-go/v6`) for all D1 operations.

### SDK Usage

```go
// Create client with API token
client := cloudflare.NewClient(
    option.WithAPIToken(apiToken),
    option.WithRequestTimeout(timeout),
)

// Execute query
params := d1.DatabaseQueryParams{
    AccountID: cloudflare.F(accountID),
    Body: d1.DatabaseQueryParamsBodyD1SingleQuery{
        Sql: cloudflare.F(query),
    },
}

// Use autopaging iterator for results
iter := client.D1.Database.QueryAutoPaging(ctx, databaseID, params)
for iter.Next() {
    queryResult := iter.Current()
    // Process results
}
```

### API Endpoint (Internal)

The SDK communicates with:
```
POST https://api.cloudflare.com/client/v4/accounts/{accountID}/d1/database/{databaseID}/query
```

### Response Structure

The SDK returns structured `d1.QueryResult` objects containing:
- `Results`: Slice of result rows (each row is a map[string]interface{})
- `Meta.RowsRead`: Number of rows read
- `Meta.RowsWritten`: Number of rows written
- `Meta.LastRowID`: Last inserted row ID
- `Meta.Duration`: Query execution duration

## Error Handling

The driver implements comprehensive error handling at multiple levels:

1. **Configuration Errors**: DSN parsing, validation
2. **Connection Errors**: Network issues, authentication failures
3. **Query Errors**: SQL syntax errors, constraint violations
4. **Transaction Errors**: Commit failures, isolation level issues

All errors are wrapped with context using `fmt.Errorf` with `%w` for proper error unwrapping.

## Performance Considerations

### Connection Pooling

The driver is fully compatible with `database/sql` connection pooling:

```go
db.SetMaxOpenConns(25)      // Max connections
db.SetMaxIdleConns(5)       // Max idle connections
db.SetConnMaxLifetime(5*time.Minute)  // Connection lifetime
db.SetConnMaxIdleTime(1*time.Minute)  // Idle connection lifetime
```

### Prepared Statements

Prepared statements improve performance for repeated queries:
- Statement is parsed once
- Parameters are bound efficiently
- Reduces parsing overhead

### Query Optimization

- Use parameterized queries to prevent SQL injection
- Use appropriate indexes in D1 database
- Batch operations using transactions
- Set appropriate timeouts

## Security

### Authentication

- Uses Cloudflare API token authentication via SDK
- Token managed by SDK's `option.WithAPIToken()`
- SDK handles secure Authorization header transmission
- No credentials stored in connection strings after parsing

### SQL Injection Prevention

- All queries use parameterized binding
- User input is properly escaped
- No string concatenation for query building

### Connection Security

- SDK ensures HTTPS for all API communications
- TLS/SSL encryption in transit
- Proper timeout configuration via SDK options
- SDK handles certificate validation

## Testing Strategy

### Unit Tests

- DSN parsing and validation
- Type conversion functions
- Parameter binding
- Configuration validation

### Integration Tests

- Would require actual D1 database
- End-to-end query execution
- Transaction handling
- Connection pooling

### Test Coverage

Current test coverage includes:
- All utility functions
- DSN parsing and validation
- Configuration validation
- Type conversions
- Parameter binding

## Best Practices

### Usage

1. **Always close resources**:
   ```go
   defer db.Close()
   defer rows.Close()
   defer stmt.Close()
   ```

2. **Use contexts with timeouts**:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   ```

3. **Configure connection pool**:
   ```go
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(5)
   ```

4. **Use prepared statements for repeated queries**:
   ```go
   stmt, _ := db.Prepare("INSERT INTO users (name) VALUES (?)")
   defer stmt.Close()
   ```

5. **Handle errors properly**:
   ```go
   if err != nil {
       if errors.Is(err, sql.ErrNoRows) {
           // Handle no rows case
       }
       return fmt.Errorf("query failed: %w", err)
   }
   ```

## Future Enhancements

### Potential Improvements

1. **Batch Operations**: Native support for bulk inserts
2. **Query Builder**: Type-safe query builder
3. **Migrations**: Schema migration support
4. **Metrics**: Built-in performance metrics
5. **Retry Logic**: Automatic retry for transient failures
6. **Connection Multiplexing**: HTTP/2 connection reuse
7. **Query Caching**: Cache frequently used queries
8. **Mock Interface**: Testing utilities for unit tests

## Dependencies

The driver has minimal dependencies:
- Go 1.22+ standard library
- `github.com/cloudflare/cloudflare-go/v6` (official Cloudflare Go SDK)

## Thread Safety

The driver is fully thread-safe and can be used concurrently by multiple goroutines, as required by the `database/sql` package specification.

## Compatibility

- **Go Version**: 1.22+
- **Database**: Cloudflare D1
- **Cloudflare SDK**: v6.5.0+
- **API Version**: v4
- **SQL Dialect**: SQLite-compatible

## Maintenance

### Code Organization

- `driver.go`: Driver registration and configuration
- `conn.go`: Connection management and API communication
- `stmt.go`: Statement handling
- `rows.go`: Result set iteration
- `result.go`: Execution results
- `utils.go`: Utility functions
- `errors.go`: Error definitions
- `doc.go`: Package documentation

### Code Quality

- Follows Go Code Review Comments guidelines
- Uses `gofmt` and `goimports` for formatting
- Passes `go vet` without errors
- Includes comprehensive unit tests
- Documented public APIs

## Conclusion

The D1 Go SQL Driver provides a robust, production-ready solution for integrating Cloudflare D1 with Go applications. Its clean architecture, comprehensive error handling, and adherence to Go best practices make it suitable for production use while maintaining ease of use for developers.

