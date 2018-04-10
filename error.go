package goorm

import (
	"errors"
)

// Error messages.
var (
	ErrNoMoreRows = errors.New(`orm: no more rows in this result set`)
)
