package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type SongHandler struct {
	svc contracts.SongService
	log *logrus.Logger
}

func NewSongHandler(svc contracts.SongService, log *logrus.Logger) *SongHandler {
	return &SongHandler{
		svc: svc,
		log: log,
	}
}

func (h *SongHandler) GetSongs(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	songs, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "GetSongs", err)
	}

	total, err := h.svc.GetCount(c.Context())
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "GetSongs", err)
	}

	return c.JSON(fiber.Map{
		"data": songs,
		"pagination": dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

func (h *SongHandler) GetSong(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	song, err := h.svc.GetSongById(c.Context(), id)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "GetSong", err)
	}

	return c.JSON(fiber.Map{
		"data": song,
	})
}

func (h *SongHandler) CreateSong(c *fiber.Ctx) error {
	var req dto.CreateSongRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.CreateSong(c.Context(), req); err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "CreateSong", err)
	}

	return c.JSON(fiber.Map{
		"data": "Successfully created song",
	})
}

func (h *SongHandler) UpdateSong(c *fiber.Ctx) error {
	var req dto.CreateSongRequest
	id, _ := strconv.Atoi(c.Params("id"))

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.UpdateSong(c.Context(), req, id); err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "UpdateSong", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully updated song",
	})
}

func (h *SongHandler) DeleteSong(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteSong(c.Context(), id); err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "DeleteSong", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted song",
	})
}
