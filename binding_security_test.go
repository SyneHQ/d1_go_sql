package d1sql

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

func TestBindingDoesNotRescanUntrustedValues(t *testing.T) {
	query := "SELECT ? AS first_value, ? AS second_value"
	got, err := bindParameters(query, []driver.Value{"?", " OR 1=1 --"})
	want := "SELECT '?' AS first_value, ' OR 1=1 --' AS second_value"
	if err != nil || got != want {
		t.Fatalf("parameter escaped its literal: got %q err=%v want %q", got, err, want)
	}
}

func TestBindingRecognizesSQLiteQuotesAndComments(t *testing.T) {
	query := "SELECT '?', 'it''s ?', \"col?\", `col?`, [col?], ? -- ignored ?\n/* ignored ? */ WHERE value = ?"
	want := "SELECT '?', 'it''s ?', \"col?\", `col?`, [col?], 'O''Brien?' -- ignored ?\n/* ignored ? */ WHERE value = X'003fff'"
	got, err := bindParameters(query, []driver.Value{"O'Brien?", []byte{0, 63, 255}})
	if err != nil || got != want {
		t.Fatalf("quoted text or types changed: got %q err=%v want %q", got, err, want)
	}
	for _, query := range []string{"SELECT \"escaped\"\"?\", ?", "SELECT `escaped``?`, ?"} {
		got, err := bindParameters(query, []driver.Value{int64(7)})
		if err != nil || !strings.HasSuffix(got, ", 7") {
			t.Fatalf("escaped identifier delimiter mishandled: %q %v", got, err)
		}
	}
}

func TestBindingRejectsArgumentMismatchAndUnsupportedValues(t *testing.T) {
	tests := []struct {
		query string
		args  []driver.Value
	}{
		{"SELECT ?", nil},
		{"SELECT '?' -- ?", []driver.Value{int64(1)}},
		{"SELECT ?, ?", []driver.Value{int64(1)}},
		{"SELECT ?", []driver.Value{int64(1), int64(2)}},
		{"SELECT ?1", []driver.Value{int64(1)}},
		{"SELECT ?", []driver.Value{struct{ Value string }{"' OR 1=1 --"}}},
	}
	for _, test := range tests {
		if got, err := bindParameters(test.query, test.args); err == nil || got != "" {
			t.Fatalf("invalid parameters accepted for %q: %q %v", test.query, got, err)
		}
	}
}

func TestHostileValuesStayQuotedInExecQueryAndTransactionRequests(t *testing.T) {
	var submitted []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			SQL   string `json:"sql"`
			Batch []struct {
				SQL string `json:"sql"`
			} `json:"batch"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body.SQL != "" {
			submitted = append(submitted, body.SQL)
		}
		for _, item := range body.Batch {
			submitted = append(submitted, item.SQL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"errors":[],"messages":[],"result":[{"success":true,"results":[],"meta":{"rows_written":1}}]}`))
	}))
	defer server.Close()
	client := cloudflare.NewClient(option.WithBaseURL(server.URL), option.WithAPIToken("test-token"))
	conn := &Conn{client: client, config: &Config{AccountID: "test-account", DatabaseID: "test-database"}, ctx: context.Background()}
	query := "SELECT ? AS first_value, ? AS second_value"
	args := []driver.NamedValue{{Ordinal: 1, Value: "?"}, {Ordinal: 2, Value: " OR 1=1 --"}}
	if _, err := conn.ExecContext(conn.ctx, query, args); err != nil {
		t.Fatal(err)
	}
	rows, err := conn.QueryContext(conn.ctx, query, args)
	if err != nil {
		t.Fatal(err)
	}
	_ = rows.Close()
	tx, err := conn.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.ExecContext(conn.ctx, query, args); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if len(submitted) != 3 {
		t.Fatalf("got %d SQL requests, want exec, query, and transaction", len(submitted))
	}
	for _, sql := range submitted {
		if sql != "SELECT '?' AS first_value, ' OR 1=1 --' AS second_value" {
			t.Fatalf("unsafe SQL reached the D1 API: %q", sql)
		}
	}
}
