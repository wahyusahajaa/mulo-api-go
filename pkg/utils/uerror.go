package utils

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// 400: BadRequestError is used when the client sends an invalid or malformed request.
type BadReqError struct {
	Errors map[string]string
}

func (e BadReqError) Error() string {
	return "Invalid body request."
}

// 404: NotFoundError is used when a requested resource is not found with return message with specific Id.
type NotFoundError struct {
	Resource string
	Id       int
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %v not found.", e.Resource, e.Id)
}

// 404: NotFoundError is used when a requested resource is not found with return message with field.
type NotFoundErrorWithCustomField struct {
	Resource string
	Field    string
	Value    any
}

func (e NotFoundErrorWithCustomField) Error() string {
	return fmt.Sprintf("%s with %s '%v' not found.", e.Resource, e.Field, e.Value)
}

// 404: NotFoundError is used when a requested resource is not found with return the custome meessage
type NotFoundErrorWithCustomMsg struct {
	Message string
}

func (e NotFoundErrorWithCustomMsg) Error() string {
	return e.Message
}

// 409: ConflictError is used when a unique constraint is violated.
type ConflictError struct {
	Resource string
	Field    string
	Value    any
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("%s with %s '%v' already exists.", e.Resource, e.Field, e.Value)
}

// 410: GoneError is used when the requested resource is no longer available on the origin server and this condition is expected to be permanent
type GoneError struct {
	Resource string
	Field    string
	Value    any
}

func (e GoneError) Error() string {
	return fmt.Sprintf("%s with %s '%v' has expired.", e.Resource, e.Field, e.Value)
}

// Handle http error
func HandleHTTPError(c *fiber.Ctx, log *logrus.Logger, layer, operation string, err error) error {
	var badReq BadReqError
	switch {
	case errors.As(err, &badReq):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
			"errors":  badReq.Errors,
		})
	case errors.As(err, &NotFoundError{}):
		LogWarn(log, c.Context(), layer, operation, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	case errors.As(err, &NotFoundErrorWithCustomField{}):
		LogWarn(log, c.Context(), layer, operation, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	case errors.As(err, &NotFoundErrorWithCustomMsg{}):
		LogWarn(log, c.Context(), layer, operation, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	case errors.As(err, &ConflictError{}):
		LogWarn(log, c.Context(), layer, operation, err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	case errors.As(err, &GoneError{}):
		LogWarn(log, c.Context(), layer, operation, err)
		return c.Status(fiber.StatusGone).JSON(fiber.Map{
			"message": err.Error(),
		})
	default:
		LogError(log, c.Context(), layer, operation, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": GetRequestId(c.Context()),
		})
	}
}
