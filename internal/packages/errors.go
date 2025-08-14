package packages

import "errors"

var (
	ErrInvalid         = errors.New("invalid items")
	NoPackSizes        = errors.New("no pack sizes")
	NoCombinationError = errors.New("no feasible combination found")
)
