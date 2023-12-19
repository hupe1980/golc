package schema

import "errors"

var (
	ErrInvalidChainValues  = errors.New("invalid chain values")
	ErrChainValueWrongType = errors.New("chain value is of wrong type")
)
