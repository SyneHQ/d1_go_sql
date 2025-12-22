# Basic Example

This example demonstrates the basic usage of the D1 Go SQL driver.

## Prerequisites

1. A Cloudflare account with D1 database access
2. A D1 database created
3. Cloudflare API token with D1 permissions

## Setup

Set the required environment variables:

```bash
export CLOUDFLARE_ACCOUNT_ID="your-account-id"
export CLOUDFLARE_API_TOKEN="your-api-token"
export D1_DATABASE_ID="your-database-id"
```

## Running the Example

```bash
cd examples/basic
go run main.go
```

## What This Example Demonstrates

1. **Connection Management**: Opening and configuring database connections
2. **Table Creation**: Creating tables with DDL statements
3. **Data Insertion**: Inserting data with parameterized queries
4. **Data Querying**: Querying single and multiple rows
5. **Prepared Statements**: Using prepared statements for efficient repeated queries
6. **Transactions**: Managing transactions with commit/rollback

## Expected Output

```
✓ Connected to D1 database successfully

--- Create Table Example ---
✓ Table 'users' created successfully

--- Insert Data Example ---
✓ Inserted user with ID: 1 (rows affected: 1)

--- Query Data Example ---
✓ User #1: Alice Johnson (alice@example.com)

✓ All users:
  - [1] Alice Johnson (alice@example.com) - Age: 28
✓ Retrieved 1 users

--- Prepared Statement Example ---
✓ Inserted Bob Smith with ID: 2
✓ Inserted Carol White with ID: 3
✓ Inserted David Brown with ID: 4

--- Transaction Example ---
✓ Transaction committed successfully (2 users inserted)

✓ All examples completed successfully!
```

