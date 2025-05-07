package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type PlaylistHandler struct {
	svc contracts.PlaylistService
	log *logrus.Logger
}

func NewPlaylistHandler(svc contracts.PlaylistService, log *logrus.Logger) *PlaylistHandler {
	return &PlaylistHandler{
		svc: svc,
		log: log,
	}
}

func (h *PlaylistHandler) GetPlaylists(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	playlists, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "playlist_handler", "GetPlaylists", err)
	}

	total, err := h.svc.GetCount(c.Context())
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "playlist_handler", "GetPlaylists", err)
	}

	return c.JSON(fiber.Map{
		"data": playlists,
		"pagination": dto.Pagination{
			PageSize: pageSize,
			Page:     page,
			Total:    total,
		},
	})
}

func (h *PlaylistHandler) GetPlaylist(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	playlist, err := h.svc.GetPlaylistById(c.Context(), id)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "playlist_handler", "GetPlaylist", err)
	}

	return c.JSON(fiber.Map{
		"data": playlist,
	})
}

func (h *PlaylistHandler) CreatePlaylist(c *fiber.Ctx) error {
	var req dto.CreatePlaylistRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.CreatePlaylist(c.Context(), req); err != nil {
		return utils.HandleHTTPError(c, h.log, "playlist_handler", "CreatePlaylist", err)
	}

	return c.JSON(fiber.Map{
		"data": "Successfully created playlists.",
	})
}

func (h *PlaylistHandler) UpdatePlaylist(c *fiber.Ctx) error {
	var req dto.CreatePlaylistRequest
	id, _ := strconv.Atoi(c.Params("id"))

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body.",
		})
	}

	if err := h.svc.UpdatePlaylist(c.Context(), req, id); err != nil {
		return utils.HandleHTTPError(c, h.log, "playlist_handler", "UpdatePlaylist", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully updated playlist",
	})
}

func (h *PlaylistHandler) DeletePlaylist(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeletePlaylist(c.Context(), id); err != nil {
		return utils.HandleHTTPError(c, h.log, "playlist_handler", "DeletePlaylist", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted playlist",
	})
}
