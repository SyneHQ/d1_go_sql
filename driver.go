package d1sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

const (
	// DriverName is the name used to register this driver with database/sql
	DriverName = "d1"

	// DefaultTimeout is the default timeout for database operations
	DefaultTimeout = 30 * time.Second
)

// init registers the driver with database/sql
func init() {
	sql.Register(DriverName, &Driver{})
}

// Driver implements database/sql/driver.Driver and driver.DriverContext
type Driver struct{}

// Open opens a new connection to the D1 database
// DSN format: d1://accountID:apiToken@databaseID?timeout=30s
func (d *Driver) Open(dsn string) (driver.Conn, error) {
	connector, err := d.OpenConnector(dsn)
	if err != nil {
		return nil, err
	}
	return connector.Connect(context.Background())
}

// OpenConnector creates a new connector for the D1 database
func (d *Driver) OpenConnector(dsn string) (driver.Connector, error) {
	config, err := ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	return &Connector{
		config: config,
		driver: d,
	}, nil
}

// Config holds the configuration for connecting to D1
type Config struct {
	AccountID  string
	APIToken   string
	DatabaseID string
	Timeout    time.Duration
}

// ParseDSN parses a DSN string into a Config
// Format: d1://accountID:apiToken@databaseID?timeout=30s
func ParseDSN(dsn string) (*Config, error) {
	if dsn == "" {
		return nil, errors.New("empty DSN")
	}

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid DSN format: %w", err)
	}

	if u.Scheme != DriverName {
		return nil, fmt.Errorf("invalid scheme: expected '%s', got '%s'", DriverName, u.Scheme)
	}

	config := &Config{
		Timeout: DefaultTimeout,
	}

	// Parse account ID and API token
	if u.User != nil {
		config.AccountID = u.User.Username()
		config.APIToken, _ = u.User.Password()
	}

	if config.AccountID == "" {
		return nil, errors.New("missing account ID in DSN")
	}

	if config.APIToken == "" {
		return nil, errors.New("missing API token in DSN")
	}

	// Parse database ID from host
	config.DatabaseID = u.Host
	if config.DatabaseID == "" {
		return nil, errors.New("missing database ID in DSN")
	}

	// Parse query parameters
	query := u.Query()
	if timeoutStr := query.Get("timeout"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout parameter: %w", err)
		}
		config.Timeout = timeout
	}

	return config, nil
}

// FormatDSN formats a Config into a DSN string
func (c *Config) FormatDSN() string {
	u := &url.URL{
		Scheme: DriverName,
		User:   url.UserPassword(c.AccountID, c.APIToken),
		Host:   c.DatabaseID,
	}

	query := url.Values{}
	if c.Timeout != DefaultTimeout {
		query.Set("timeout", c.Timeout.String())
	}

	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	return u.String()
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.AccountID == "" {
		return errors.New("account ID is required")
	}
	if c.APIToken == "" {
		return errors.New("API token is required")
	}
	if c.DatabaseID == "" {
		return errors.New("database ID is required")
	}
	if c.Timeout <= 0 {
		return errors.New("timeout must be positive")
	}
	return nil
}

// Connector implements driver.Connector
type Connector struct {
	config *Config
	driver *Driver
}

// Connect returns a new connection to the database
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	if err := c.config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create Cloudflare client with API token
	client := cloudflare.NewClient(
		option.WithAPIToken(c.config.APIToken),
		option.WithRequestTimeout(c.config.Timeout),
	)

	return &Conn{
		client: client,
		config: c.config,
		ctx:    ctx,
	}, nil
}

// Driver returns the driver
func (c *Connector) Driver() driver.Driver {
	return c.driver
}
