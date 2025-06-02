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

// @Summary      	List Albums
// @Description  	Get paginated list of albums
// @Tags         	albums
// @Security     	BearerAuth
// @Produce      	json
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithPagination[[]dto.Album, dto.Pagination]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/albums [get]
func (h *AlbumHandler) GetAlbums(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	albums, total, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "album_handler", "GetAlbums", err)
	}

	return c.JSON(dto.ResponseWithPagination[[]dto.AlbumWithArtist, dto.Pagination]{
		Data: albums,
		Pagination: dto.Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// @Summary      	Get Album by ID
// @Description  	Get a Album by their ID
// @Tags        	albums
// @Security     	BearerAuth
// @Produce      	json
// @Param        	id		path     	int	true  "Album ID"
// @Success 		200		{object} 	dto.ResponseWithData[dto.Album]
// @Failure 		404 	{object} 	dto.ErrorResponse "Album not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router       	/albums/{id} [get]
func (h *AlbumHandler) GetAlbum(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	album, err := h.svc.GetAlbumById(c.Context(), id)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "album_handler", "GetAlbum", err)
	}

	return c.JSON(dto.ResponseWithData[dto.AlbumWithArtist]{
		Data: album,
	})
}

// @Summary 		Create album
// @Description 	Create a new album.
// @Tags        	albums
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			album	 body		dto.CreateAlbumRequest true "Album object that needs to be created"
// @Success 		201 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404		{object} 	dto.ErrorResponse "Not Found"
// @Failure 		409		{object} 	dto.ErrorResponse "Conflict album name"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/albums [post]
func (h *AlbumHandler) CreateAlbum(c *fiber.Ctx) error {
	var req dto.CreateAlbumRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Invalid body request.",
		})
	}

	if err := h.svc.CreateAlbum(c.Context(), req); err != nil {
		return errs.HandleHTTPError(c, h.log, "album_handler", "CreateAlbum", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Suceessfully created album.",
	})
}

// @Summary 		Update album
// @Description 	Update the album with the specified ID
// @Tags        	albums
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Album ID"
// @Param 			album	body		dto.CreateAlbumRequest true "album object that needs to be updated"
// @Success 		200 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404 	{object} 	dto.ErrorResponse "Not Found: album or artist not found"
// @Failure 		409		{object} 	dto.ErrorResponse "Conflict: album name"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/albums/{id} [put]
func (h *AlbumHandler) UpdateAlbum(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.CreateAlbumRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Invalid body request.",
		})
	}

	if err := h.svc.UpdateAlbum(c.Context(), req, id); err != nil {
		return errs.HandleHTTPError(c, h.log, "album_handler", "UpdateAlbum", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully updated album.",
	})
}

// @Summary 		Delete album
// @Description 	Delete the album with the specified ID
// @Tags        	albums
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "album ID"
// @Success 		200		{object} 	dto.ResponseMessage
// @Failure 		404 	{object} 	dto.ErrorResponse "Not Found: Album not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/albums/{id} [delete]
func (h *AlbumHandler) DeleteAlbum(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteAlbum(c.Context(), id); err != nil {
		return errs.HandleHTTPError(c, h.log, "album_handler", "DeleteAlbum", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully deleted album.",
	})
}

// @Summary      	List of albums by artist
// @Description  	Get list of albums by artist id
// @Tags         	artists
// @Security     	BearerAuth
// @Produce      	json
// @Param 			id 			path int true "Artist ID"
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithData[[]dto.Album]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/artists/{id}/albums [get]
func (h *AlbumHandler) GetAlbumsByArtistId(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))

	albums, err := h.svc.GetAlbumsByArtistId(c.Context(), artistId)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "album_handler", "GetAlbumsByArtistId", err)
	}

	return c.JSON(dto.ResponseWithData[[]dto.Album]{
		Data: albums,
	})
}
