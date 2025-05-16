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

// @Summary      	List of favorite songs
// @Description  	Get paginated list of favorite songs
// @Tags         	favorites
// @Security     	BearerAuth
// @Produce      	json
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithPagination[[]dto.Song, dto.Pagination]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/favorites/songs [get]
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

// @Summary 		Add song to favorite
// @Description 	Add song to favorite
// @Tags        	favorites
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			songId path int true "Song ID"
// @Success 		201 	{object} 	dto.ResponseMessage
// @Failure 		404		{object} 	dto.ValidationErrorResponse "Not Found: Song does not exists."
// @Failure 		409		{object} 	dto.ValidationErrorResponse "Conflict: Song already exists on favorites."
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/favorites/songs/{songId} [post]
func (h *FavoriteHandler) AddFavoriteSong(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("songId"))
	userId := utils.GetUserId(c.Context())

	if err := h.svc.AddFavoriteSong(c.Context(), userId, songId); err != nil {
		return errs.HandleHTTPError(c, h.log, "favorite_handler", "AddFavoriteSong", err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ResponseMessage{
		Message: "Succesfully added song to favorite.",
	})
}

// @Summary 		Remove song from favorite
// @Description 	Remove song from favorite
// @Tags        	favorites
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			songId path int true "Song ID"
// @Success 		200 	{object} 	dto.ResponseMessage
// @Failure 		404		{object} 	dto.ValidationErrorResponse "Not Found: Song on favorites does not exists."
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/favorites/songs/{songId} [delete]
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
