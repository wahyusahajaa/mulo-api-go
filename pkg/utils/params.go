package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetPaginationParam(c *fiber.Ctx) (page, pageSize, offset int) {
	page, _ = strconv.Atoi(c.Query("page", "1"))
	pageSize, _ = strconv.Atoi(c.Query("pageSize", "10"))
	offset = (page - 1) * pageSize

	return
}
