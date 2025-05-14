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

type FavoriteHandler struct {
	svc contracts.FavoriteService
	log *logrus.Logger
}

func NewFavoriteHandler(svc contracts.FavoriteService, log *logrus.Logger) *FavoriteHandler {
	return &FavoriteHandler{
		svc: svc,
		log: log,
	}
}

func (h *FavoriteHandler) GetFavoriteSongsByUserId(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)
	userId := utils.GetUserId(c.Context())

	songs, err := h.svc.GetAllFavoriteSongsByUserId(c.Context(), userId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "favorite_handler", "GetFavoriteSongsByUserId", err)
	}

	total, err := h.svc.GetCountFavoriteSongsByUserId(c.Context(), userId)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "favorite_handler", "GetFavoriteSongsByUserId", err)
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

func (h *FavoriteHandler) CreateFavoriteSong(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("songId"))
	userId := utils.GetUserId(c.Context())

	if err := h.svc.CreateFavoriteSong(c.Context(), userId, songId); err != nil {
		return errs.HandleHTTPError(c, h.log, "favorite_handler", "CreateFavoriteSong", err)
	}

	return c.JSON(fiber.Map{
		"message": "Succesfully added song to favorite",
	})
}

func (h *FavoriteHandler) DeleteFavoriteSong(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("songId"))
	userId := utils.GetUserId(c.Context())

	if err := h.svc.DeleteFavoriteSong(c.Context(), userId, songId); err != nil {
		return errs.HandleHTTPError(c, h.log, "favorite_handler", "DeleteFavoriteSong", err)
	}

	return c.JSON(fiber.Map{
		"message": "Succesfully deleted song from favorite",
	})
}
