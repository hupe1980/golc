package rag

import "errors"

var (
	ErrNoInputValues  = errors.New("no input values")
	ErrNoOutputParser = errors.New("no output parser")
)
