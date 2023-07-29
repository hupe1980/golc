package schema

import "errors"

var (
	ErrInvalidInputValues   = errors.New("invalid input values")
	ErrInputValuesWrongType = errors.New("input value is of wrong type")
)
