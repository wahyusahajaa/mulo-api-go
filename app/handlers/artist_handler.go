package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type ArtistHandler struct {
	svc contracts.ArtistService
	log *logrus.Logger
}

func NewArtistHandler(svc contracts.ArtistService, log *logrus.Logger) *ArtistHandler {
	return &ArtistHandler{
		svc: svc,
		log: log,
	}
}

// GetArtists		Get paginated list of artists
// @Summary      	List of Artists
// @Description  	Get paginated list of artists
// @Tags         	artists
// @Security     	BearerAuth
// @Produce      	json
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithPagination[[]dto.Artist, dto.Pagination]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/artists [get]
func (h *ArtistHandler) GetArtists(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	artists, total, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "GetArtists", err)
	}

	return c.JSON(dto.ResponseWithPagination[[]dto.Artist, dto.Pagination]{
		Data: artists,
		Pagination: dto.Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// @Summary 		Create artist
// @Description 	Create a new artist.
// @Tags        	artists
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			artist	 body		dto.CreateArtistRequest true "artist object that needs to be created"
// @Success 		201 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		409		{object} 	dto.ErrorResponse "Conflict artist name"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/artists [post]
func (h *ArtistHandler) CreateArtist(c *fiber.Ctx) error {
	var req dto.CreateArtistRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Invalid body request.",
		})
	}

	if err := h.svc.CreateArtist(c.Context(), req); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "CreateArtist", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully created artist.",
	})
}

// @Summary      	Get Artist by ID
// @Description  	Get a Artist by their ID
// @Tags        	artists
// @Security     	BearerAuth
// @Produce      	json
// @Param        	id		path     	int	true  "Artist ID"
// @Success 		200		{object} 	dto.ResponseWithData[dto.Artist]
// @Failure 		404 	{object} 	dto.ErrorResponse "Artist not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router       	/artists/{id} [get]
func (h *ArtistHandler) GetArtist(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))

	artist, err := h.svc.GetArtistById(c.Context(), artistId)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "GetArtist", err)
	}

	return c.JSON(dto.ResponseWithData[dto.Artist]{
		Data: artist,
	})
}

// @Summary 		Update artist
// @Description 	Update the artist with the specified ID
// @Tags        	artists
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Artist ID"
// @Param 			artist	body		dto.CreateArtistRequest true "artist object that needs to be updated"
// @Success 		200 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404 	{object} 	dto.ErrorResponse "Not Found: artist not found"
// @Failure 		409		{object} 	dto.ErrorResponse "Conflict: artist name"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/artists/{id} [put]
func (h *ArtistHandler) UpdateArtist(c *fiber.Ctx) error {
	var req dto.CreateArtistRequest
	id, _ := strconv.Atoi(c.Params("id"))

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Invalid body request.",
		})
	}

	if err := h.svc.UpdateArtist(c.Context(), req, id); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "UpdateArtist", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully updated artist.",
	})
}

// @Summary 		Delete artist
// @Description 	Delete the artist with the specified ID
// @Tags        	artists
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "artist ID"
// @Success 		200		{object} 	dto.ResponseMessage
// @Failure 		404 	{object} 	dto.ErrorResponse "Not Found: artist not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/artists/{id} [delete]
func (h *ArtistHandler) DeleteArtist(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteArtist(context.TODO(), id); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "DeleteArtist", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully deleted artist.",
	})
}
