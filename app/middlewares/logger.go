package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func FiberLogger(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		logger.WithFields(logrus.Fields{
			"method":     c.Method(),
			"url":        c.OriginalURL(),
			"status":     c.Response().StatusCode(),
			"user_agent": c.Get("User-Agent"),
		}).Info("Request processed")

		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Request error")
		}

		return err
	}
}
