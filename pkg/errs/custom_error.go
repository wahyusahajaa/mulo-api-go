package errs

import "fmt"

type BaseError struct {
	Message string
	Code    int
	Cause   error // the underlying error
}

func (e *BaseError) Error() string {
	return e.Message
}

func (e *BaseError) StatusCode() int {
	return e.Code
}

func (e *BaseError) Unwrap() error {
	return e.Cause
}

// 400: BadRequest Error
type BadRequestError struct {
	*BaseError
	Errors map[string]string // optional
}

func NewBadRequestError(message string, errors map[string]string, cause ...error) *BadRequestError {
	var underlying error
	if len(cause) > 0 {
		underlying = cause[0]
	}

	if message == "" {
		message = "Invalid body request."
	}

	return &BadRequestError{
		BaseError: &BaseError{
			Message: message,
			Code:    400,
			Cause:   underlying,
		},
		Errors: errors,
	}
}

// 404 NotFoundError represents a 404 error with resource and field info.
type NotFoundError struct {
	*BaseError
}

func NewNotFoundError(resource, field string, value any, cause ...error) *NotFoundError {
	var underlying error
	if len(cause) > 0 {
		underlying = cause[0]
	}

	return &NotFoundError{
		BaseError: &BaseError{
			Message: fmt.Sprintf("%s with %s '%v' not found.", resource, field, value),
			Code:    400,
			Cause:   underlying, // store the original cause
		},
	}
}

func NewNotFoundErrorWithMsg(Message string) *NotFoundError {
	return &NotFoundError{
		BaseError: &BaseError{
			Message: Message,
			Code:    404,
		},
	}
}

type ConflictError struct {
	*BaseError
}

func NewConflictError(resource, field string, value any, cause ...error) *ConflictError {
	var underlying error
	if len(cause) > 0 {
		underlying = cause[0]
	}

	return &ConflictError{
		BaseError: &BaseError{
			Message: fmt.Sprintf("%s with %s '%v' already exists.", resource, field, value),
			Code:    409,
			Cause:   underlying,
		},
	}
}

func NewConflictErrorWithMsg(Message string) *NotFoundError {
	return &NotFoundError{
		BaseError: &BaseError{
			Message: Message,
			Code:    404,
		},
	}
}

// GoneError represents a 410 error for expired resources.
type GoneError struct {
	*BaseError
}

func NewGoneError(resource, field string, value any, cause ...error) *GoneError {
	var underlying error
	if len(cause) > 0 {
		underlying = cause[0]
	}

	return &GoneError{
		BaseError: &BaseError{
			Message: fmt.Sprintf("%s with %s '%v' has expired.", resource, field, value),
			Code:    410,
			Cause:   underlying,
		},
	}
}

// Forbidden client error response status code indicates that the server understood the request but refused to process it
type Fobidden struct {
	*BaseError
}

func NewForbiddenError(message string, cause ...error) *Fobidden {
	var underlying error
	if len(cause) > 0 {
		underlying = cause[0]
	}

	return &Fobidden{
		BaseError: &BaseError{
			Message: message,
			Code:    403,
			Cause:   underlying,
		},
	}
}

type Unauthorized struct {
	*BaseError
}

func NewUnauthorizedError(message string, cause ...error) *Unauthorized {
	var underlying error
	if len(cause) > 0 {
		underlying = cause[0]
	}

	return &Unauthorized{
		BaseError: &BaseError{
			Message: message,
			Code:    401,
			Cause:   underlying,
		},
	}
}
