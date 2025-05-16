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

func (h *FavoriteHandler) GetFavoriteSongsByUserID(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)
	userID := utils.GetUserId(c.Context())

	songs, total, err := h.svc.GetFavoriteSongsByUserID(c.Context(), userID, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "favorite_handler", "GetFavoriteSongsByUserID", err)
	}

	return c.JSON(dto.ResponseWithPagination[[]dto.Song, dto.Pagination]{
		Data: songs,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

func (h *FavoriteHandler) AddFavoriteSong(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("songId"))
	userId := utils.GetUserId(c.Context())

	if err := h.svc.AddFavoriteSong(c.Context(), userId, songId); err != nil {
		return errs.HandleHTTPError(c, h.log, "favorite_handler", "AddFavoriteSong", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Succesfully added song to favorite.",
	})
}

func (h *FavoriteHandler) RemoveFavoriteSong(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("songId"))
	userId := utils.GetUserId(c.Context())

	if err := h.svc.RemoveFavoriteSong(c.Context(), userId, songId); err != nil {
		return errs.HandleHTTPError(c, h.log, "favorite_handler", "RemoveFavoriteSong", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Succesfully removed song from favorite.",
	})
}
