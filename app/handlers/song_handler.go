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

// @Summary      	List Songs
// @Description  	Get paginated list of songs
// @Tags         	songs
// @Security     	BearerAuth
// @Produce      	json
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithPagination[[]dto.Song, dto.Pagination]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/songs [get]
func (h *SongHandler) GetSongs(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	songs, total, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "GetSongs", err)
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

// @Summary      	Get song by ID
// @Description  	Get a song by their ID
// @Tags        	songs
// @Security     	BearerAuth
// @Produce      	json
// @Param        	id		path     	int	true  "Song ID"
// @Success 		200		{object} 	dto.ResponseWithData[dto.Song]
// @Failure 		404 	{object} 	dto.ErrorResponse "Song not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router       	/songs/{id} [get]
func (h *SongHandler) GetSong(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	song, err := h.svc.GetSongById(c.Context(), id)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "GetSong", err)
	}

	return c.JSON(dto.ResponseWithData[dto.Song]{
		Data: song,
	})
}

// @Summary 		Create song
// @Description 	Create a new song.
// @Tags        	songs
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			song	 body		dto.CreateSongRequest true "Song object that needs to be created"
// @Success 		201 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/songs [post]
func (h *SongHandler) CreateSong(c *fiber.Ctx) error {
	var req dto.CreateSongRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Invalid body request.",
		})
	}

	if err := h.svc.CreateSong(c.Context(), req); err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "CreateSong", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully created song.",
	})
}

// @Summary 		Update song
// @Description 	Update the song with the specified ID
// @Tags        	songs
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Song ID"
// @Param 			song	body		dto.CreateSongRequest true "Song object that needs to be updated"
// @Success 		200 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404 	{object} 	dto.ErrorResponse "Song not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/songs/{id} [put]
func (h *SongHandler) UpdateSong(c *fiber.Ctx) error {
	var req dto.CreateSongRequest
	id, _ := strconv.Atoi(c.Params("id"))

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Invalid body request.",
		})
	}

	if err := h.svc.UpdateSong(c.Context(), req, id); err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "UpdateSong", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully updated song.",
	})
}

// @Summary 		Delete song
// @Description 	Delete the song with the specified ID
// @Tags        	songs
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "Song ID"
// @Success 		200		{object} 	dto.ResponseMessage
// @Failure 		404 	{object} 	dto.ErrorResponse "Song not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteSong(c.Context(), id); err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "DeleteSong", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully deleted song.",
	})
}

// @Summary      	List of Songs by album
// @Description  	Get paginated list of songs by album
// @Tags         	albums
// @Security     	BearerAuth
// @Produce      	json
// @Param        	id			path     	int	true  "Album ID"
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithPagination[[]dto.Song, dto.Pagination]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/albums/{id}/songs [get]
func (h *SongHandler) GetSongsByAlbumId(c *fiber.Ctx) error {
	albumId, _ := strconv.Atoi(c.Params("id"))
	page, pageSize, offset := utils.GetPaginationParam(c)

	songs, total, err := h.svc.GetSongsByAlbumId(c.Context(), albumId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "song_handler", "GetSongsByAlbumId", err)
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
