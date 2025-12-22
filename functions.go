package d1sql

import (
	"database/sql/driver"
	"regexp"
	"strings"
)

// MetadataFunction represents a supported metadata function
type MetadataFunction struct {
	Pattern *regexp.Regexp
	Handler func(*Conn) (driver.Rows, error)
}

var (
	// Supported metadata functions
	metadataFunctions = []MetadataFunction{
		{
			Pattern: regexp.MustCompile(`(?i)^\s*SELECT\s+current_database\s*\(\s*\)\s*;?\s*$`),
			Handler: handleCurrentDatabase,
		},
		{
			Pattern: regexp.MustCompile(`(?i)^\s*SELECT\s+database\s*\(\s*\)\s*;?\s*$`),
			Handler: handleCurrentDatabase,
		},
		{
			Pattern: regexp.MustCompile(`(?i)^\s*SELECT\s+current_user\s*\(\s*\)\s*;?\s*$`),
			Handler: handleCurrentUser,
		},
		{
			Pattern: regexp.MustCompile(`(?i)^\s*SELECT\s+user\s*\(\s*\)\s*;?\s*$`),
			Handler: handleCurrentUser,
		},
		{
			Pattern: regexp.MustCompile(`(?i)^\s*SELECT\s+version\s*\(\s*\)\s*;?\s*$`),
			Handler: handleVersion,
		},
		{
			Pattern: regexp.MustCompile(`(?i)^\s*SELECT\s+connection_id\s*\(\s*\)\s*;?\s*$`),
			Handler: handleConnectionID,
		},
	}
)

// tryMetadataFunction attempts to execute a metadata function
// Returns (rows, true) if it's a metadata function, (nil, false) otherwise
func tryMetadataFunction(conn *Conn, query string) (driver.Rows, bool) {
	query = strings.TrimSpace(query)

	for _, fn := range metadataFunctions {
		if fn.Pattern.MatchString(query) {
			rows, err := fn.Handler(conn)
			if err != nil {
				// If handler fails, return empty result but still indicate it was handled
				return &Rows{
					columns: []string{"error"},
					rows:    [][]interface{}{},
					index:   -1,
				}, true
			}
			return rows, true
		}
	}

	return nil, false
}

// handleCurrentDatabase returns the current database ID
func handleCurrentDatabase(conn *Conn) (driver.Rows, error) {
	return &Rows{
		columns: []string{"current_database()"},
		rows: [][]interface{}{
			{conn.config.DatabaseID},
		},
		index: -1,
	}, nil
}

// handleCurrentUser returns the current account ID
func handleCurrentUser(conn *Conn) (driver.Rows, error) {
	return &Rows{
		columns: []string{"current_user()"},
		rows: [][]interface{}{
			{conn.config.AccountID},
		},
		index: -1,
	}, nil
}

// handleVersion returns the driver version and D1 info
func handleVersion(conn *Conn) (driver.Rows, error) {
	version := "D1 Go SQL Driver v0.1.0 (Cloudflare D1 - SQLite compatible)"
	return &Rows{
		columns: []string{"version()"},
		rows: [][]interface{}{
			{version},
		},
		index: -1,
	}, nil
}

// handleConnectionID returns a pseudo connection ID based on database and account
func handleConnectionID(conn *Conn) (driver.Rows, error) {
	// Generate a deterministic connection ID from account and database
	connID := conn.config.AccountID + ":" + conn.config.DatabaseID
	return &Rows{
		columns: []string{"connection_id()"},
		rows: [][]interface{}{
			{connID},
		},
		index: -1,
	}, nil
}
