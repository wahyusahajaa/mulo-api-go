package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type ArtistHandler struct {
	svc contracts.ArtistService
	log *logrus.Logger
}

func NewArtistHandler(svc contracts.ArtistService, log *logrus.Logger) *ArtistHandler {
	return &ArtistHandler{
		svc: svc,
		log: log,
	}
}

func (h *ArtistHandler) GetArtists(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	artists, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "artist_handler", "GetArtists", err)
	}

	total, err := h.svc.GetCount(c.Context())
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "artist_handler", "GetArtists", err)
	}

	return c.JSON(fiber.Map{
		"data": artists,
		"pagination": dto.Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func (h *ArtistHandler) CreateArtist(c *fiber.Ctx) error {
	var req dto.CreateArtistRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.CreateArtist(c.Context(), req); err != nil {
		return utils.HandleHTTPError(c, h.log, "artist_handler", "CreateArtist", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully created artist",
	})
}

func (h *ArtistHandler) GetArtist(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))

	artist, err := h.svc.GetArtistById(c.Context(), artistId)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "artist_handler", "GetArtist", err)
	}

	return c.JSON(fiber.Map{
		"data": artist,
	})
}

func (h *ArtistHandler) UpdateArtist(c *fiber.Ctx) error {
	var req dto.CreateArtistRequest
	id, _ := strconv.Atoi(c.Params("id"))

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.UpdateArtist(c.Context(), req, id); err != nil {
		return utils.HandleHTTPError(c, h.log, "artist_handler", "UpdateArtist", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully updated artist",
	})
}

func (h *ArtistHandler) DeleteArtist(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteArtist(context.TODO(), id); err != nil {
		return utils.HandleHTTPError(c, h.log, "artist_handler", "DeleteArtist", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted artist",
	})
}

func (h *ArtistHandler) CreateArtistGenre(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))
	genreId, _ := strconv.Atoi(c.Params("genreId"))

	if err := h.svc.CreateArtistGenre(c.Context(), artistId, genreId); err != nil {
		return utils.HandleHTTPError(c, h.log, "artist_handler", "CreateArtistGenre", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully added new artist genre",
	})
}

func (h *ArtistHandler) GetArtistGenres(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))
	_, pageSize, offset := utils.GetPaginationParam(c)

	artistGenres, err := h.svc.GetArtistGenres(c.Context(), artistId, pageSize, offset)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "artist_handler", "GetArtistGenres", err)
	}

	return c.JSON(fiber.Map{
		"data": artistGenres,
	})
}

func (h *ArtistHandler) DeleteArtistGenre(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))
	genreId, _ := strconv.Atoi(c.Params("genreId"))

	if err := h.svc.DeleteArtistGenre(c.Context(), artistId, genreId); err != nil {
		return utils.HandleHTTPError(c, h.log, "artist_handler", "DeleteArtistGenre", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted artist genre",
	})
}
