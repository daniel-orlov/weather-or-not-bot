package types

import "github.com/pkg/errors"

var (
	Err = errors.New("some new error")
)

const (
	ErrOnHandling = "cannot handle '%s'"
)
