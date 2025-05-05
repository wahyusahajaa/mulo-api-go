package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
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
		utils.LogError(h.log, c.Context(), "song_handler", "GetSongs", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}

	total, err := h.svc.GetCount(c.Context())
	if err != nil {
		utils.LogError(h.log, c.Context(), "song_handler", "GetSongs", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": utils.GetRequestId(c.Context()),
		})
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
		var notFoundErr utils.NotFoundError
		if errors.As(err, &notFoundErr) {
			utils.LogWarn(h.log, c.Context(), "song_handler", "GetSong", notFoundErr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": notFoundErr.Error(),
			})
		}

		utils.LogError(h.log, c.Context(), "song_handler", "GetSong", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}

	return c.JSON(fiber.Map{
		"data": song,
	})
}

func (h *SongHandler) CreateSong(c *fiber.Ctx) error {
	var req dto.CreateSongRequest

	errorsMap, err := utils.ParseBody(c.Body(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if errorsMap != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid data type",
			"errors":  errorsMap,
		})
	}

	if err := h.svc.CreateSong(c.Context(), req); err != nil {
		valErr := utils.BadReqError{}
		if errors.As(err, &valErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Validation failed",
				"errors":  valErr.Errors,
			})
		}

		notFoundErr := utils.NotFoundError{}
		if errors.As(err, &notFoundErr) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": notFoundErr.Error(),
			})
		}

		conflictErr := utils.ConflictError{}
		if errors.As(err, &conflictErr) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": conflictErr.Error(),
			})
		}

		utils.LogError(h.log, c.Context(), "song_handler", "CreateSong", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}

	return c.JSON(fiber.Map{
		"data": "Successfully created song",
	})
}

func (h *SongHandler) UpdateSong(c *fiber.Ctx) error {
	var req dto.CreateSongRequest
	id, _ := strconv.Atoi(c.Params("id"))

	errorsMap, err := utils.ParseBody(c.Body(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if errorsMap != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid data type",
			"errors":  errorsMap,
		})
	}

	if err := h.svc.UpdateSong(c.Context(), req, id); err != nil {

		valErr := utils.BadReqError{}
		if errors.As(err, &valErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Validation failed",
				"errors":  valErr.Errors,
			})
		}

		notFoundErr := utils.NotFoundError{}
		if errors.As(err, &notFoundErr) {
			utils.LogWarn(h.log, c.Context(), "song_handler", "UpdateSong", notFoundErr)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": notFoundErr.Error(),
			})
		}

		utils.LogError(h.log, c.Context(), "song_handler", "UpdateSong", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully updated song",
	})
}

func (h *SongHandler) DeleteSong(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteSong(c.Context(), id); err != nil {

		notFoundErr := utils.NotFoundError{}
		if errors.As(err, &notFoundErr) {
			utils.LogWarn(h.log, c.Context(), "song_handler", "DeleteSong", notFoundErr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": notFoundErr.Error(),
			})
		}

		utils.LogError(h.log, c.Context(), "song_handler", "DeleteSong", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted song",
	})
}
