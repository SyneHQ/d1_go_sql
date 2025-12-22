# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of D1 Go SQL Driver
- Full `database/sql` interface implementation using official Cloudflare Go SDK (v6)
- Support for Cloudflare D1 database
- Connection pooling support
- Prepared statements with parameter binding
- Transaction support with statement batching (Begin, Commit, Rollback)
- Context-aware operations
- Comprehensive type conversion (string, int64, float64, bool, []byte, time.Time, nil)
- DSN parsing and validation
- Error handling with detailed error messages
- Complete test coverage (36.6% coverage)
- Example applications
- Comprehensive documentation (README, ARCHITECTURE, CONTRIBUTING)
- Production-grade code following DRY principles

### Features
- Official Cloudflare Go SDK (v6.5.0+) integration
- Auto-pagination support via SDK iterators for large result sets
- Thread-safe concurrent operations
- `driver.Driver` implementation
- `driver.DriverContext` implementation
- `driver.Connector` implementation
- `driver.Conn` implementation
- `driver.ConnPrepareContext` implementation
- `driver.ExecerContext` implementation
- `driver.QueryerContext` implementation
- `driver.Stmt` implementation
- `driver.StmtExecContext` implementation
- `driver.StmtQueryContext` implementation
- `driver.Rows` implementation
- `driver.Result` implementation
- `driver.Tx` implementation

### Technical Requirements
- Go 1.22 or later
- `github.com/cloudflare/cloudflare-go/v6` v6.5.0+
- Cloudflare D1 database
- Cloudflare API token with D1 access

### Documentation
- README with usage examples
- API documentation with godoc
- Contributing guidelines
- Example code
- Test suite

## [0.1.0] - 2025-12-22

### Added
- Initial development version
- Core driver functionality
- Basic CRUD operations
- Transaction support
- Type conversion utilities
- Test suite
- Examples

[Unreleased]: https://github.com/synehq/d1_go_sql/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/synehq/d1_go_sql/releases/tag/v0.1.0

