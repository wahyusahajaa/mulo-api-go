package errs

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

func HandleHTTPError(c *fiber.Ctx, log *logrus.Logger, layer, operation string, err error) error {
	switch e := err.(type) {
	case *BadRequestError:
		resBody := make(fiber.Map)
		resBody["message"] = e.Message
		if e.Errors != nil {
			resBody["errors"] = e.Errors
		}
		return c.Status(fiber.StatusBadRequest).JSON(resBody)
	case *ConflictError:
		utils.LogWarn(log, c.Context(), layer, operation, err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": e.Message,
		})
	case *NotFoundError:
		utils.LogWarn(log, c.Context(), layer, operation, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": e.Message,
		})
	case *GoneError:
		utils.LogWarn(log, c.Context(), layer, operation, err)
		return c.Status(fiber.StatusGone).JSON(fiber.Map{
			"message": e.Message,
		})
	default:
		utils.LogError(log, c.Context(), layer, operation, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}
}
