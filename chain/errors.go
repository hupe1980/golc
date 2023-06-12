package chain

import "errors"

var (
	ErrNoInputValues        = errors.New("no input values")
	ErrInvalidInputValues   = errors.New("invalid input values")
	ErrInputValuesWrongType = errors.New("input key is of wrong type")
)
