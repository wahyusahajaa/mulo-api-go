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
		utils.LogError(h.log, c.Context(), "genre_handler", "GetGenres", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Error,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}

	total, err := h.svc.GetCount(c.Context())
	if err != nil {
		utils.LogError(h.log, c.Context(), "genre_handler", "GetGenres", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Error,
			"requestId": utils.GetRequestId(c.Context()),
		})
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

		notFoundErr := utils.NotFoundError{}
		if errors.As(err, &notFoundErr) {
			utils.LogWarn(h.log, c.Context(), "genre_handler", "GetGenre", notFoundErr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": notFoundErr.Error(),
			})
		}

		utils.LogError(h.log, c.Context(), "genre_handler", "GetGenre", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}

	return c.JSON(fiber.Map{
		"data": genre,
	})
}

func (h *GenreHandler) CreateGenre(c *fiber.Ctx) error {
	var req dto.CreateGenreRequest

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

	if err := h.svc.CreateGenre(c.Context(), req); err != nil {

		valErr := utils.BadReqError{}
		if errors.As(err, &valErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Validation failed",
				"errors":  valErr.Errors,
			})
		}

		utils.LogError(h.log, c.Context(), "genre_handler", "CreateGenre", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully created genre",
	})
}

func (h *GenreHandler) UpdateGenre(c *fiber.Ctx) error {
	var req dto.CreateGenreRequest
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

	if err := h.svc.UpdateGenre(c.Context(), req, id); err != nil {

		badReqErr := utils.BadReqError{}
		if errors.As(err, &badReqErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Validation failed",
				"errors":  badReqErr.Errors,
			})
		}

		notFoundErr := utils.NotFoundError{}
		if errors.As(err, &notFoundErr) {
			utils.LogWarn(h.log, c.Context(), "genre_handler", "UpdateGenre", notFoundErr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": notFoundErr.Error(),
			})
		}

		utils.LogError(h.log, c.Context(), "genre_handler", "UpdateGenre", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Error,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully updated genre",
	})
}

func (h *GenreHandler) DeleteGenre(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteGenre(c.Context(), id); err != nil {

		errNotFound := utils.NotFoundError{}
		if errors.As(err, &errNotFound) {
			utils.LogWarn(h.log, c.Context(), "genre_handler", "DeleteGenre", errNotFound)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": errNotFound.Error(),
			})
		}

		utils.LogError(h.log, c.Context(), "genre_handler", "DeleteGenre", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":   fiber.ErrInternalServerError.Message,
			"requestId": utils.GetRequestId(c.Context()),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully Deleted genre",
	})
}
