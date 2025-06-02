package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func FiberLogger(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestId := uuid.New().String()
		c.Locals("requestId", requestId)
		err := c.Next()

		logger.WithFields(logrus.Fields{
			"origin":     c.Get("Origin"),
			"method":     c.Method(),
			"url":        c.OriginalURL(),
			"status":     c.Response().StatusCode(),
			"user_agent": c.Get("User-Agent"),
			"requestId":  requestId,
		}).Info("Request processed")

		if err != nil {
			logger.WithError(err).Error("Request error")
		}

		return err
	}
}
