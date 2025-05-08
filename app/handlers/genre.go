package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type GenreHandler struct {
	svc contracts.GenreService
	log *logrus.Logger
}

func NewGenreHandler(svc contracts.GenreService, log *logrus.Logger) *GenreHandler {
	return &GenreHandler{
		svc: svc,
		log: log,
	}
}

func (h *GenreHandler) GetGenres(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)
	genres, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "genre_handler", "GetGenres", err)
	}

	total, err := h.svc.GetCount(c.Context())
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "genre_handler", "GetGenres", err)
	}

	return c.JSON(fiber.Map{
		"data": genres,
		"pagination": dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

func (h *GenreHandler) GetGenre(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	genre, err := h.svc.GetGenreById(c.Context(), id)

	if err != nil {
		return utils.HandleHTTPError(c, h.log, "genre_handler", "GetGenre", err)
	}

	return c.JSON(fiber.Map{
		"data": genre,
	})
}

func (h *GenreHandler) CreateGenre(c *fiber.Ctx) error {
	var req dto.CreateGenreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.CreateGenre(c.Context(), req); err != nil {
		return utils.HandleHTTPError(c, h.log, "genre_handler", "CreateGenre", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully created genre",
	})
}

func (h *GenreHandler) UpdateGenre(c *fiber.Ctx) error {
	var req dto.CreateGenreRequest
	id, _ := strconv.Atoi(c.Params("id"))

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.UpdateGenre(c.Context(), req, id); err != nil {
		return utils.HandleHTTPError(c, h.log, "genre_handler", "UpdateGenre", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully updated genre",
	})
}

func (h *GenreHandler) DeleteGenre(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteGenre(c.Context(), id); err != nil {
		return utils.HandleHTTPError(c, h.log, "genre_handler", "DeleteGenre", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully Deleted genre",
	})
}
