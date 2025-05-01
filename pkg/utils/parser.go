package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ParseUserId(c *fiber.Ctx) (int, error) {
	idVal := c.Locals("id")
	if idVal == nil {
		return 0, errors.New("user ID not found in context")
	}

	idFloat, ok := idVal.(float64)
	if !ok {
		return 0, errors.New("user ID is not a valid float64")
	}

	return int(idFloat), nil
}
