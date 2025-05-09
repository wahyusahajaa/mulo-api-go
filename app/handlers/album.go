package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type AlbumHandler struct {
	svc contracts.AlbumService
	log *logrus.Logger
}

func NewAlbumHandler(svc contracts.AlbumService, log *logrus.Logger) *AlbumHandler {
	return &AlbumHandler{
		svc: svc,
		log: log,
	}
}

func (h *AlbumHandler) GetAlbums(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	albums, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "album_handler", "GetAlbums", err)
	}

	total, err := h.svc.GetCount(c.Context())
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "album_handler", "GetAlbums", err)
	}

	return c.JSON(fiber.Map{
		"data": albums,
		"pagination": dto.Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func (h *AlbumHandler) GetAlbum(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	album, err := h.svc.GetAlbumById(c.Context(), id)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "album_handler", "GetAlbum", err)
	}

	return c.JSON(fiber.Map{
		"data": album,
	})
}

func (h *AlbumHandler) CreateAlbum(c *fiber.Ctx) error {
	var req dto.CreateAlbumRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.CreateAlbum(c.Context(), req); err != nil {
		return utils.HandleHTTPError(c, h.log, "album_handler", "CreateAlbum", err)
	}

	return c.JSON(fiber.Map{
		"message": "Suceessfully created album",
	})
}

func (h *AlbumHandler) UpdateAlbum(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.CreateAlbumRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.UpdateAlbum(c.Context(), req, id); err != nil {
		return utils.HandleHTTPError(c, h.log, "album_handler", "UpdateAlbum", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully updated album",
	})
}

func (h *AlbumHandler) DeleteAlbum(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteAlbum(c.Context(), id); err != nil {
		return utils.HandleHTTPError(c, h.log, "album_handler", "DeleteAlbum", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted album",
	})
}

func (h *AlbumHandler) GetAlbumsByArtistId(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))

	artists, err := h.svc.GetAlbumsByArtistId(c.Context(), artistId)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "album_handler", "GetAlbumsByArtistId", err)
	}

	return c.JSON(fiber.Map{
		"data": artists,
	})
}
