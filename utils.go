package d1sql

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// namedValuesToValues converts []driver.NamedValue to []driver.Value
func namedValuesToValues(named []driver.NamedValue) ([]driver.Value, error) {
	values := make([]driver.Value, len(named))
	for i, nv := range named {
		if nv.Name != "" {
			return nil, fmt.Errorf("named parameters are not supported, use positional parameters instead")
		}
		values[i] = nv.Value
	}
	return values, nil
}

// namedValuesFromValues converts []driver.Value to []driver.NamedValue
func namedValuesFromValues(values []driver.Value) []driver.NamedValue {
	named := make([]driver.NamedValue, len(values))
	for i, v := range values {
		named[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   v,
		}
	}
	return named
}

// bindParameters replaces ? placeholders with actual values in SQL query
func bindParameters(query string, args []driver.Value) (string, error) {
	if len(args) == 0 {
		return query, nil
	}

	// Count the number of placeholders
	placeholderCount := strings.Count(query, "?")
	if placeholderCount != len(args) {
		return "", fmt.Errorf("parameter count mismatch: query has %d placeholders, but %d arguments provided", placeholderCount, len(args))
	}

	// Replace each ? with the corresponding value
	result := query
	for _, arg := range args {
		// Find the first occurrence of ?
		index := strings.Index(result, "?")
		if index == -1 {
			break
		}

		// Convert the argument to SQL literal
		literal, err := valueToSQLLiteral(arg)
		if err != nil {
			return "", err
		}

		// Replace the ? with the literal
		result = result[:index] + literal + result[index+1:]
	}

	return result, nil
}

// valueToSQLLiteral converts a driver.Value to SQL literal string
func valueToSQLLiteral(value driver.Value) (string, error) {
	if value == nil {
		return "NULL", nil
	}

	switch v := value.(type) {
	case int64:
		return fmt.Sprintf("%d", v), nil
	case float64:
		return fmt.Sprintf("%g", v), nil
	case bool:
		if v {
			return "1", nil
		}
		return "0", nil
	case []byte:
		return fmt.Sprintf("X'%x'", v), nil
	case string:
		// Escape single quotes by doubling them
		escaped := strings.ReplaceAll(v, "'", "''")
		return fmt.Sprintf("'%s'", escaped), nil
	case time.Time:
		// Format as RFC3339 which is compatible with SQLite
		return fmt.Sprintf("'%s'", v.Format(time.RFC3339)), nil
	default:
		// Try to convert to string
		return fmt.Sprintf("'%v'", v), nil
	}
}

// convertValue converts a value from D1 response to driver.Value
func convertValue(value interface{}) driver.Value {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case float32:
		return float64(v)
	case float64:
		return v
	case bool:
		return v
	case string:
		return v
	case []byte:
		return v
	default:
		// Convert to string as fallback
		return fmt.Sprintf("%v", v)
	}
}

