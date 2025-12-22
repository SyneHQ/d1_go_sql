package d1sql

import (
	"errors"
)

var (
	// ErrInvalidDSN is returned when the DSN is invalid
	ErrInvalidDSN = errors.New("invalid DSN")

	// ErrConnectionClosed is returned when attempting to use a closed connection
	ErrConnectionClosed = errors.New("connection is closed")

	// ErrInTransaction is returned when an operation is not allowed during a transaction
	ErrInTransaction = errors.New("operation not allowed in transaction")

	// ErrNotInTransaction is returned when commit/rollback is called outside a transaction
	ErrNotInTransaction = errors.New("not in transaction")

	// ErrParameterCountMismatch is returned when the number of parameters doesn't match
	ErrParameterCountMismatch = errors.New("parameter count mismatch")

	// ErrNamedParametersNotSupported is returned when named parameters are used
	ErrNamedParametersNotSupported = errors.New("named parameters are not supported")

	// ErrInvalidConfiguration is returned when the configuration is invalid
	ErrInvalidConfiguration = errors.New("invalid configuration")

	// ErrQueryFailed is returned when a query fails
	ErrQueryFailed = errors.New("query failed")
)
