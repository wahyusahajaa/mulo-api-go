package utils

import "fmt"

type ErrMap map[string]string

// BadRequestError is used when the client sends an invalid or malformed request.
type BadReqError struct {
	Errors map[string]string
}

func (e BadReqError) Error() string {
	return "validation error"
}

// NotFoundError is used when a requested resource is not found.
type NotFoundError struct {
	Resource string
	Id       int
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %v not found.", e.Resource, e.Id)
}

// ConflictError is used when a unique constraint is violated.
type ConflictError struct {
	Resource string
	Field    string
	Value    any
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("%s with %s '%v' already exists.", e.Resource, e.Field, e.Value)
}
