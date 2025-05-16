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

// GetGenres		Get paginated list of genres
// @Summary      	List Genres
// @Description  	Get paginated list of genres
// @Tags         	genres
// @Security     	BearerAuth
// @Produce      	json
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithPagination[[]dto.Genre, dto.Pagination]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/genres [get]
func (h *GenreHandler) GetGenres(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	genres, total, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "genre_handler", "GetGenres", err)
	}

	return c.JSON(dto.ResponseWithPagination[[]dto.Genre, dto.Pagination]{
		Data: genres,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

// GetGenre			Get a Genre by their ID
// @Summary      	Get Genre by ID
// @Description  	Get a Genre by their ID
// @Tags        	genres
// @Security     	BearerAuth
// @Produce      	json
// @Param        	id		path     	int	true  "Genre ID"
// @Success 		200		{object} 	dto.ResponseWithData[dto.Genre]
// @Failure 		404 	{object} 	dto.ErrorResponse "Genre not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router       	/genres/{id} [get]
func (h *GenreHandler) GetGenre(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	genre, err := h.svc.GetGenreById(c.Context(), id)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "genre_handler", "GetGenre", err)
	}

	return c.JSON(dto.ResponseWithData[dto.Genre]{
		Data: genre,
	})
}

// CreateGenre		Create a new genre.
// @Summary 		Create genre
// @Description 	Create a new genre.
// @Tags        	genres
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			genre	 body		dto.CreateGenreRequest true "Genre object that needs to be created"
// @Success 		201 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/genres [post]
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

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully created genre",
	})
}

// UpdateGenre		Update an existing genre.
// @Summary 		Update genre
// @Description 	Update the genre with the specified ID
// @Tags        	genres
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Genre ID"
// @Param 			genre	body		dto.CreateGenreRequest true "Genre object that needs to be updated"
// @Success 		200 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404 	{object} 	dto.ErrorResponse "Genre not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/genres/{id} [put]
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

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully updated genre",
	})
}

// DeleteGenre		Delete an existing genre.
// @Summary 		Delete genre
// @Description 	Delete the genre with the specified ID
// @Tags        	genres
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "Genre ID"
// @Success 		200		{object} 	dto.ResponseMessage
// @Failure 		404 	{object} 	dto.ErrorResponse "Genre not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/genres/{id} [delete]
func (h *GenreHandler) DeleteGenre(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.DeleteGenre(c.Context(), id); err != nil {
		return errs.HandleHTTPError(c, h.log, "genre_handler", "DeleteGenre", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully deleted genre.",
	})
}

// @Summary 		Assign genre to artist
// @Description 	Assign genre to artist
// @Tags        	artists
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Artist ID"
// @Param 			genreId path int true "Genre ID"
// @Success 		201 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404		{object} 	dto.ValidationErrorResponse "Not Found: Genre or artist does not exists."
// @Failure 		409		{object} 	dto.ValidationErrorResponse "Conflict: Genre already exists on artist."
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/artists/{id}/genres/{genreId} [post]
func (h *GenreHandler) CreateArtistGenre(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))
	genreId, _ := strconv.Atoi(c.Params("genreId"))

	if err := h.svc.CreateArtistGenre(c.Context(), artistId, genreId); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "CreateArtistGenre", err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ResponseMessage{
		Message: "Successfully assigned genre to artist.",
	})
}

// @Summary      	List of genres by artist
// @Description  	Get list of genres by artist id
// @Tags         	artists
// @Security     	BearerAuth
// @Produce      	json
// @Param 			id 			path int true "Artist ID"
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithData[[]dto.Genre]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/artists/{id}/genres [get]
func (h *GenreHandler) GetArtistGenres(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))
	_, pageSize, offset := utils.GetPaginationParam(c)

	artistGenres, err := h.svc.GetArtistGenres(c.Context(), artistId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "GetArtistGenres", err)
	}

	return c.JSON(dto.ResponseWithData[[]dto.Genre]{
		Data: artistGenres,
	})
}

// @Summary 		Delete genre from artist
// @Description 	Delete genre from artist
// @Tags        	artists
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Artist ID"
// @Param 			genreId path int true "Genre ID"
// @Success 		200 	{object} 	dto.ResponseMessage
// @Failure 		404		{object} 	dto.ValidationErrorResponse "Not Found: Genre on artist does not exists."
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/artists/{id}/genres/{genreId} [delete]
func (h *GenreHandler) DeleteArtistGenre(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))
	genreId, _ := strconv.Atoi(c.Params("genreId"))

	if err := h.svc.DeleteArtistGenre(c.Context(), artistId, genreId); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "DeleteArtistGenre", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully deleted genre from artist.",
	})
}

// @Summary 		Assign genre to song
// @Description 	Assign genre to song
// @Tags        	songs
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Song ID"
// @Param 			genreId path int true "Genre ID"
// @Success 		201 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404		{object} 	dto.ValidationErrorResponse "Not Found: Genre or song does not exists."
// @Failure 		409		{object} 	dto.ValidationErrorResponse "Conflict: Genre already exists on song."
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/songs/{id}/genres/{genreId} [post]
func (h *GenreHandler) CreateSongGenre(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("id"))
	genreId, _ := strconv.Atoi(c.Params("genreId"))

	if err := h.svc.CreateSongGenre(c.Context(), songId, genreId); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "CreateSongGenre", err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ResponseMessage{
		Message: "Successfully assigned genre to song.",
	})
}

// @Summary      	List of genres by song
// @Description  	Get paginated list of artists by genre
// @Tags         	songs
// @Security     	BearerAuth
// @Produce      	json
// @Param 			id 			path int true "Song ID"
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithData[[]dto.Genre]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/songs/{id}/genres [get]
func (h *GenreHandler) GetSongGenres(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("id"))
	_, pageSize, offset := utils.GetPaginationParam(c)

	genres, err := h.svc.GetSongGenres(c.Context(), songId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "GetSongGenres", err)
	}

	return c.JSON(dto.ResponseWithData[[]dto.Genre]{
		Data: genres,
	})
}

// @Summary 		Delete genre from song
// @Description 	Delete genre from song
// @Tags        	songs
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Song ID"
// @Param 			genreId path int true "Genre ID"
// @Success 		200 	{object} 	dto.ResponseMessage
// @Failure 		404		{object} 	dto.ValidationErrorResponse "Not Found: Genre on song does not exists."
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/songs/{id}/genres/{genreId} [delete]
func (h *GenreHandler) DeleteSongGenre(c *fiber.Ctx) error {
	songId, _ := strconv.Atoi(c.Params("id"))
	genreId, _ := strconv.Atoi(c.Params("genreId"))

	if err := h.svc.DeleteSongGenre(c.Context(), songId, genreId); err != nil {
		return errs.HandleHTTPError(c, h.log, "artist_handler", "DeleteSongGenre", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully deleted genre from song.",
	})
}

// GetArtists		Get paginated list of artists by genre
// @Summary      	List of artists by genre
// @Description  	Get paginated list of artists by genre
// @Tags         	genres
// @Security     	BearerAuth
// @Produce      	json
// @Param 			id 			path int true "Genre ID"
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithPagination[[]dto.Artist, dto.Pagination]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/genres/{id}/artists [get]
func (h *GenreHandler) GetArtists(c *fiber.Ctx) error {
	genreId, _ := strconv.Atoi(c.Params("id"))
	page, pageSize, offset := utils.GetPaginationParam(c)

	artists, total, err := h.svc.GetAllArtists(c.Context(), genreId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "genre_handler", "GetArtists", err)
	}

	return c.JSON(dto.ResponseWithPagination[[]dto.Artist, dto.Pagination]{
		Data: artists,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

// GetSongs			Get paginated list of songs by genre
// @Summary      	List of songs by genre
// @Description  	Get paginated list of songs by genre
// @Tags         	genres
// @Security     	BearerAuth
// @Produce      	json
// @Param 			id 			path int true "Genre ID"
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithPagination[[]dto.Song, dto.Pagination]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/genres/{id}/songs [get]
func (h *GenreHandler) GetSongs(c *fiber.Ctx) error {
	genreId, _ := strconv.Atoi(c.Params("id"))
	page, pageSize, offset := utils.GetPaginationParam(c)

	songs, total, err := h.svc.GetAllSongs(c.Context(), genreId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "genre_handler", "GetSongs", err)
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
