package d1sql

import (
	"context"
	"testing"
	"time"
)

func TestMetadataFunctions(t *testing.T) {
	// Create a test connection
	config := &Config{
		AccountID:  "test-account-123",
		APIToken:   "test-token",
		DatabaseID: "test-db-456",
		Timeout:    30 * time.Second,
	}

	conn := &Conn{
		config: config,
		ctx:    context.Background(),
	}

	tests := []struct {
		name           string
		query          string
		shouldHandle   bool
		expectedColumn string
		expectedValue  string
	}{
		{
			name:           "LIST DATABASES",
			query:          "LIST DATABASES",
			shouldHandle:   true,
			expectedColumn: "name",
			expectedValue:  "", // Will check columns instead
		},
		{
			name:           "LIST DATABASES with semicolon",
			query:          "LIST DATABASES;",
			shouldHandle:   true,
			expectedColumn: "name",
			expectedValue:  "",
		},
		{
			name:           "current_database() lowercase",
			query:          "SELECT current_database()",
			shouldHandle:   true,
			expectedColumn: "current_database()",
			expectedValue:  "test-db-456",
		},
		{
			name:           "current_database() uppercase",
			query:          "SELECT CURRENT_DATABASE()",
			shouldHandle:   true,
			expectedColumn: "current_database()",
			expectedValue:  "test-db-456",
		},
		{
			name:           "current_database() with semicolon",
			query:          "SELECT current_database();",
			shouldHandle:   true,
			expectedColumn: "current_database()",
			expectedValue:  "test-db-456",
		},
		{
			name:           "database() alias",
			query:          "SELECT database()",
			shouldHandle:   true,
			expectedColumn: "current_database()",
			expectedValue:  "test-db-456",
		},
		{
			name:           "current_user()",
			query:          "SELECT current_user()",
			shouldHandle:   true,
			expectedColumn: "current_user()",
			expectedValue:  "test-account-123",
		},
		{
			name:           "user() alias",
			query:          "SELECT user()",
			shouldHandle:   true,
			expectedColumn: "current_user()",
			expectedValue:  "test-account-123",
		},
		{
			name:           "version()",
			query:          "SELECT version()",
			shouldHandle:   true,
			expectedColumn: "version()",
			expectedValue:  "D1 Go SQL Driver v0.1.0 (Cloudflare D1 - SQLite compatible)",
		},
		{
			name:           "connection_id()",
			query:          "SELECT connection_id()",
			shouldHandle:   true,
			expectedColumn: "connection_id()",
			expectedValue:  "test-account-123:test-db-456",
		},
		{
			name:         "regular query not handled",
			query:        "SELECT * FROM users",
			shouldHandle: false,
		},
		{
			name:         "function in WHERE clause not handled",
			query:        "SELECT * FROM users WHERE current_database() = 'test'",
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, handled := tryMetadataFunction(conn, tt.query)

			if handled != tt.shouldHandle {
				t.Errorf("tryMetadataFunction() handled = %v, want %v", handled, tt.shouldHandle)
				return
			}

			if !tt.shouldHandle {
				return // Test passed for non-handled queries
			}

			// Check column name
			cols := rows.(*Rows).columns
			
			// Special handling for LIST DATABASES
			if tt.query == "LIST DATABASES" || tt.query == "LIST DATABASES;" {
				if len(cols) != 3 {
					t.Errorf("Expected 3 columns for LIST DATABASES, got %d", len(cols))
					return
				}
				if cols[0] != "name" || cols[1] != "uuid" || cols[2] != "version" {
					t.Errorf("Expected columns [name, uuid, version], got %v", cols)
				}
				// For tests with nil client, rows should be empty
				rowData := rows.(*Rows).rows
				if len(rowData) != 0 {
					t.Errorf("Expected empty rows for LIST DATABASES in test, got %d rows", len(rowData))
				}
				return
			}
			
			if len(cols) != 1 {
				t.Errorf("Expected 1 column, got %d", len(cols))
				return
			}
			if cols[0] != tt.expectedColumn {
				t.Errorf("Expected column %q, got %q", tt.expectedColumn, cols[0])
			}

			// Check value
			rowData := rows.(*Rows).rows
			if len(rowData) != 1 {
				t.Errorf("Expected 1 row, got %d", len(rowData))
				return
			}
			if len(rowData[0]) != 1 {
				t.Errorf("Expected 1 value in row, got %d", len(rowData[0]))
				return
			}
			if rowData[0][0] != tt.expectedValue {
				t.Errorf("Expected value %q, got %q", tt.expectedValue, rowData[0][0])
			}
		})
	}
}

func TestMetadataFunctionsWithSpaces(t *testing.T) {
	config := &Config{
		AccountID:  "account-123",
		DatabaseID: "db-456",
	}

	conn := &Conn{
		config: config,
		ctx:    context.Background(),
	}

	queries := []string{
		"  SELECT current_database()  ",
		"SELECT   current_database  (  )  ",
		"\nSELECT current_database()\n",
		"\t\tSELECT current_database();\t",
	}

	for _, query := range queries {
		t.Run("query_with_whitespace", func(t *testing.T) {
			rows, handled := tryMetadataFunction(conn, query)
			if !handled {
				t.Errorf("Query %q should be handled", query)
				return
			}

			rowData := rows.(*Rows).rows
			if len(rowData) != 1 || rowData[0][0] != "db-456" {
				t.Errorf("Query %q returned unexpected result", query)
			}
		})
	}
}

func TestMetadataFunctionsCaseInsensitive(t *testing.T) {
	config := &Config{
		AccountID:  "account-123",
		DatabaseID: "db-456",
	}

	conn := &Conn{
		config: config,
		ctx:    context.Background(),
	}

	queries := []string{
		"select current_database()",
		"SELECT CURRENT_DATABASE()",
		"SeLeCt CuRrEnT_DaTaBaSe()",
	}

	for _, query := range queries {
		t.Run("case_insensitive", func(t *testing.T) {
			rows, handled := tryMetadataFunction(conn, query)
			if !handled {
				t.Errorf("Query %q should be handled", query)
				return
			}

			rowData := rows.(*Rows).rows
			if len(rowData) != 1 || rowData[0][0] != "db-456" {
				t.Errorf("Query %q returned unexpected result", query)
			}
		})
	}
}
