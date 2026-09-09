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

// bindParameters substitutes anonymous SQLite placeholders in one pass over the
// original SQL. Quoted text and comments are never interpreted as parameters,
// and substituted values are never rescanned.
// The installed D1 REST SDK only exposes string parameters; retain SQL value
// types here until native binding supports the full database/sql value set.
func bindParameters(query string, args []driver.Value) (string, error) {
	var result strings.Builder
	result.Grow(len(query))
	argument := 0
	for i := 0; i < len(query); {
		start := i
		switch query[i] {
		case '\'', '"', '`', '[':
			delimiter := query[i]
			if delimiter == '[' {
				delimiter = ']'
			}
			i++
			for i < len(query) {
				if query[i] == delimiter {
					i++
					if delimiter != ']' && i < len(query) && query[i] == delimiter {
						i++
						continue
					}
					break
				}
				i++
			}
			result.WriteString(query[start:i])
		case '-':
			if i+1 < len(query) && query[i+1] == '-' {
				i += 2
				for i < len(query) && query[i] != '\n' {
					i++
				}
				result.WriteString(query[start:i])
			} else {
				result.WriteByte(query[i])
				i++
			}
		case '/':
			if i+1 < len(query) && query[i+1] == '*' {
				i += 2
				for i < len(query) {
					if i+1 < len(query) && query[i] == '*' && query[i+1] == '/' {
						i += 2
						break
					}
					i++
				}
				result.WriteString(query[start:i])
			} else {
				result.WriteByte(query[i])
				i++
			}
		case '?':
			if i+1 < len(query) && query[i+1] >= '0' && query[i+1] <= '9' {
				return "", fmt.Errorf("numbered parameters are not supported, use anonymous ? placeholders")
			}
			if argument < len(args) {
				literal, err := valueToSQLLiteral(args[argument])
				if err != nil {
					return "", err
				}
				result.WriteString(literal)
			}
			argument++
			i++
		default:
			result.WriteByte(query[i])
			i++
		}
	}
	if argument != len(args) {
		return "", fmt.Errorf("parameter count mismatch: query has %d placeholders, but %d arguments provided", argument, len(args))
	}
	return result.String(), nil
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
		return "", fmt.Errorf("unsupported SQL parameter type %T", value)
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
