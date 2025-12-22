package d1sql

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestValueToSQLLiteral(t *testing.T) {
	tests := []struct {
		name    string
		value   driver.Value
		want    string
		wantErr bool
	}{
		{
			name:  "nil value",
			value: nil,
			want:  "NULL",
		},
		{
			name:  "int64 value",
			value: int64(42),
			want:  "42",
		},
		{
			name:  "float64 value",
			value: float64(3.14),
			want:  "3.14",
		},
		{
			name:  "bool true",
			value: true,
			want:  "1",
		},
		{
			name:  "bool false",
			value: false,
			want:  "0",
		},
		{
			name:  "string value",
			value: "hello",
			want:  "'hello'",
		},
		{
			name:  "string with quotes",
			value: "it's a test",
			want:  "'it''s a test'",
		},
		{
			name:  "byte slice",
			value: []byte{0x01, 0x02, 0x03},
			want:  "X'010203'",
		},
		{
			name:  "time value",
			value: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			want:  "'2023-01-01T12:00:00Z'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valueToSQLLiteral(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("valueToSQLLiteral() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("valueToSQLLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBindParameters(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		args    []driver.Value
		want    string
		wantErr bool
	}{
		{
			name:  "no parameters",
			query: "SELECT * FROM users",
			args:  []driver.Value{},
			want:  "SELECT * FROM users",
		},
		{
			name:  "single parameter",
			query: "SELECT * FROM users WHERE id = ?",
			args:  []driver.Value{int64(1)},
			want:  "SELECT * FROM users WHERE id = 1",
		},
		{
			name:  "multiple parameters",
			query: "INSERT INTO users (name, age) VALUES (?, ?)",
			args:  []driver.Value{"Alice", int64(30)},
			want:  "INSERT INTO users (name, age) VALUES ('Alice', 30)",
		},
		{
			name:  "parameter with quotes",
			query: "INSERT INTO users (name) VALUES (?)",
			args:  []driver.Value{"O'Brien"},
			want:  "INSERT INTO users (name) VALUES ('O''Brien')",
		},
		{
			name:    "parameter count mismatch - too many args",
			query:   "SELECT * FROM users WHERE id = ?",
			args:    []driver.Value{int64(1), int64(2)},
			wantErr: true,
		},
		{
			name:    "parameter count mismatch - too few args",
			query:   "SELECT * FROM users WHERE id = ? AND age = ?",
			args:    []driver.Value{int64(1)},
			wantErr: true,
		},
		{
			name:  "null parameter",
			query: "INSERT INTO users (name, email) VALUES (?, ?)",
			args:  []driver.Value{"Alice", nil},
			want:  "INSERT INTO users (name, email) VALUES ('Alice', NULL)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bindParameters(tt.query, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("bindParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("bindParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertValue(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  driver.Value
	}{
		{
			name:  "nil value",
			value: nil,
			want:  nil,
		},
		{
			name:  "int to int64",
			value: 42,
			want:  int64(42),
		},
		{
			name:  "int32 to int64",
			value: int32(42),
			want:  int64(42),
		},
		{
			name:  "int64",
			value: int64(42),
			want:  int64(42),
		},
		{
			name:  "float32 to float64",
			value: float32(3.14),
			want:  float64(float32(3.14)),
		},
		{
			name:  "float64",
			value: float64(3.14),
			want:  float64(3.14),
		},
		{
			name:  "bool",
			value: true,
			want:  true,
		},
		{
			name:  "string",
			value: "hello",
			want:  "hello",
		},
		{
			name:  "byte slice",
			value: []byte{0x01, 0x02},
			want:  []byte{0x01, 0x02},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertValue(tt.value)
			
			// Special handling for byte slices
			if gotBytes, ok := got.([]byte); ok {
				wantBytes, ok := tt.want.([]byte)
				if !ok {
					t.Errorf("convertValue() = %v ([]byte), want %v", got, tt.want)
					return
				}
				if len(gotBytes) != len(wantBytes) {
					t.Errorf("convertValue() = %v, want %v", got, tt.want)
					return
				}
				for i := range gotBytes {
					if gotBytes[i] != wantBytes[i] {
						t.Errorf("convertValue() = %v, want %v", got, tt.want)
						return
					}
				}
				return
			}

			if got != tt.want {
				t.Errorf("convertValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamedValuesToValues(t *testing.T) {
	tests := []struct {
		name    string
		named   []driver.NamedValue
		want    []driver.Value
		wantErr bool
	}{
		{
			name:  "empty slice",
			named: []driver.NamedValue{},
			want:  []driver.Value{},
		},
		{
			name: "positional parameters",
			named: []driver.NamedValue{
				{Ordinal: 1, Value: "Alice"},
				{Ordinal: 2, Value: int64(30)},
			},
			want: []driver.Value{"Alice", int64(30)},
		},
		{
			name: "named parameters not supported",
			named: []driver.NamedValue{
				{Name: "name", Value: "Alice"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := namedValuesToValues(tt.named)
			if (err != nil) != tt.wantErr {
				t.Errorf("namedValuesToValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("namedValuesToValues() length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("namedValuesToValues()[%d] = %v, want %v", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestNamedValuesFromValues(t *testing.T) {
	values := []driver.Value{"Alice", int64(30), true}
	named := namedValuesFromValues(values)

	if len(named) != len(values) {
		t.Errorf("namedValuesFromValues() length = %v, want %v", len(named), len(values))
		return
	}

	for i, nv := range named {
		if nv.Ordinal != i+1 {
			t.Errorf("namedValuesFromValues()[%d].Ordinal = %v, want %v", i, nv.Ordinal, i+1)
		}
		if nv.Value != values[i] {
			t.Errorf("namedValuesFromValues()[%d].Value = %v, want %v", i, nv.Value, values[i])
		}
	}
}

