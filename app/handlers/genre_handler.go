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

	genres, total, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "genre_handler", "GetGenres", err)
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
		return errs.HandleHTTPError(c, h.log, "genre_handler", "GetGenre", err)
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
		return errs.HandleHTTPError(c, h.log, "genre_handler", "CreateGenre", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
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
		return errs.HandleHTTPError(c, h.log, "genre_handler", "UpdateGenre", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully updated genre",
	})
}

func (h *GenreHandler) DeleteGenre(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteGenre(c.Context(), id); err != nil {
		return errs.HandleHTTPError(c, h.log, "genre_handler", "DeleteGenre", err)
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Successfully Deleted genre",
	})
}

func (h *GenreHandler) CreateArtistGenre(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))
	genreId, _ := strconv.Atoi(c.Params("genreId"))

	if err := h.svc.CreateArtistGenre(c.Context(), artistId, genreId); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "CreateArtistGenre", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully added new artist genre",
	})
}

func (h *GenreHandler) GetArtistGenres(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))
	_, pageSize, offset := utils.GetPaginationParam(c)

	artistGenres, err := h.svc.GetArtistGenres(c.Context(), artistId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "GetArtistGenres", err)
	}

	return c.JSON(fiber.Map{
		"data": artistGenres,
	})
}

func (h *GenreHandler) DeleteArtistGenre(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))
	genreId, _ := strconv.Atoi(c.Params("genreId"))

	if err := h.svc.DeleteArtistGenre(c.Context(), artistId, genreId); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "DeleteArtistGenre", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted artist genre",
	})
}

func (h *GenreHandler) CreateSongGenre(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("id"))
	genreId, _ := strconv.Atoi(c.Params("genreId"))

	if err := h.svc.CreateSongGenre(c.Context(), songId, genreId); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "CreateSongGenre", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully added new song genre",
	})
}

func (h *GenreHandler) GetSongGenres(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("id"))
	_, pageSize, offset := utils.GetPaginationParam(c)

	artistGenres, err := h.svc.GetSongGenres(c.Context(), songId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "GetSongGenres", err)
	}

	return c.JSON(fiber.Map{
		"data": artistGenres,
	})
}

func (h *GenreHandler) DeleteSongGenre(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("id"))
	genreId, _ := strconv.Atoi(c.Params("genreId"))

	if err := h.svc.DeleteSongGenre(c.Context(), songId, genreId); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "DeleteSongGenre", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted song genre",
	})
}

func (h *GenreHandler) GetArtists(c *fiber.Ctx) error {
	genreId, _ := strconv.Atoi(c.Params("id"))
	page, pageSize, offset := utils.GetPaginationParam(c)

	artists, total, err := h.svc.GetAllArtists(c.Context(), genreId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "genre_handler", "GetArtists", err)
	}

	return c.JSON(fiber.Map{
		"data": artists,
		"pagination": dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

func (h *GenreHandler) GetSongs(c *fiber.Ctx) error {
	genreId, _ := strconv.Atoi(c.Params("id"))
	page, pageSize, offset := utils.GetPaginationParam(c)

	songs, total, err := h.svc.GetAllSongs(c.Context(), genreId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "genre_handler", "GetSongs", err)
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
