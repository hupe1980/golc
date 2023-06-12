package chain

import "errors"

var (
	ErrNoInputValues        = errors.New("no input values")
	ErrInvalidInputValues   = errors.New("invalid input values")
	ErrInputValuesWrongType = errors.New("input key is of wrong type")
	ErrMultipleInputsInRun  = errors.New("run not supported in chain with more then one expected input")
	ErrMultipleOutputsInRun = errors.New("run not supported in chain with more then one expected output")
	ErrWrongOutputTypeInRun = errors.New("run not supported in chain that returns value that is not string")
)
