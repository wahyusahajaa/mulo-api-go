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

type PlaylistHandler struct {
	svc contracts.PlaylistService
	log *logrus.Logger
}

func NewPlaylistHandler(svc contracts.PlaylistService, log *logrus.Logger) *PlaylistHandler {
	return &PlaylistHandler{
		svc: svc,
		log: log,
	}
}

// @Summary      	List of playlists
// @Description  	Get paginated list of playlists
// @Tags         	playlists
// @Security     	BearerAuth
// @Produce      	json
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithPagination[[]dto.Playlist, dto.Pagination]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/playlists [get]
func (h *PlaylistHandler) GetPlaylists(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)
	role := utils.GetRole(c.Context())
	userId := utils.GetUserId(c.Context())

	playlists, total, err := h.svc.GetAll(c.Context(), role, userId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "playlist_handler", "GetPlaylists", err)
	}

	return c.JSON(dto.ResponseWithPagination[[]dto.Playlist, dto.Pagination]{
		Data: playlists,
		Pagination: dto.Pagination{
			PageSize: pageSize,
			Page:     page,
			Total:    total,
		},
	})
}

// @Summary      	Get playlist by ID
// @Description  	Get a playlist by their ID
// @Tags        	playlists
// @Security     	BearerAuth
// @Produce      	json
// @Param        	id		path     	int	true  "Playlist ID"
// @Success 		200		{object} 	dto.ResponseWithData[dto.Playlist]
// @Failure 		404 	{object} 	dto.ErrorResponse "Playlist not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router       	/playlists/{id} [get]
func (h *PlaylistHandler) GetPlaylist(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	role := utils.GetRole(c.Context())
	userId := utils.GetUserId(c.Context())

	playlist, err := h.svc.GetPlaylistById(c.Context(), role, userId, id)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "playlist_handler", "GetPlaylist", err)
	}

	return c.JSON(dto.ResponseWithData[dto.Playlist]{
		Data: playlist,
	})
}

// @Summary 		Create Playlist
// @Description 	Create a new playlist.
// @Tags        	playlists
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			playlist	 body		dto.CreatePlaylistRequest true "Playlist object that needs to be created"
// @Success 		201 		{object} 	dto.ResponseMessage
// @Failure 		400			{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		500 		{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/playlists [post]
func (h *PlaylistHandler) CreatePlaylist(c *fiber.Ctx) error {
	var req dto.CreatePlaylistRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseError{
			Message: "Invalid body request.",
		})
	}

	if err := h.svc.CreatePlaylist(c.Context(), req); err != nil {
		return errs.HandleHTTPError(c, h.log, "playlist_handler", "CreatePlaylist", err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ResponseMessage{
		Message: "Successfully created playlists.",
	})
}

// @Summary 		Update playlist
// @Description 	Update the playlist with the specified ID
// @Tags        	playlists
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Playlist ID"
// @Param 			song	body		dto.CreatePlaylistRequest true "Song object that needs to be updated"
// @Success 		200 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404 	{object} 	dto.ErrorResponse "Playlist not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/playlists/{id} [put]
func (h *PlaylistHandler) UpdatePlaylist(c *fiber.Ctx) error {
	var req dto.CreatePlaylistRequest
	playlistId, _ := strconv.Atoi(c.Params("id"))
	userRole := utils.GetRole(c.Context())
	userId := utils.GetUserId(c.Context())

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseError{
			Message: "Invalid body request.",
		})
	}

	if err := h.svc.UpdatePlaylist(c.Context(), req, userRole, userId, playlistId); err != nil {
		return errs.HandleHTTPError(c, h.log, "playlist_handler", "UpdatePlaylist", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully updated playlist.",
	})
}

// @Summary 		Delete playlist
// @Description 	Delete the playlist with the specified ID
// @Tags        	playlists
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "Song ID"
// @Success 		200		{object} 	dto.ResponseMessage
// @Failure 		404 	{object} 	dto.ErrorResponse "Playlist not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/playlists/{id} [delete]
func (h *PlaylistHandler) DeletePlaylist(c *fiber.Ctx) error {
	playlistId, _ := strconv.Atoi(c.Params("id"))
	userId := utils.GetUserId(c.Context())
	userRole := utils.GetRole(c.Context())

	if err := h.svc.DeletePlaylist(c.Context(), userRole, userId, playlistId); err != nil {
		return errs.HandleHTTPError(c, h.log, "playlist_handler", "DeletePlaylist", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully deleted playlist.",
	})
}

// @Summary      	List of songs by playlist
// @Description  	Get list of songs by playlist id.
// @Tags         	playlists
// @Security     	BearerAuth
// @Produce      	json
// @Param 			id 			path int true "Playlist ID"
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithData[[]dto.Song]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/playlists/{id}/songs [get]
func (h *PlaylistHandler) GetPlaylistSongs(c *fiber.Ctx) error {
	_, pageSize, offset := utils.GetPaginationParam(c)
	playlistId, _ := strconv.Atoi(c.Params("id"))
	role := utils.GetRole(c.Context())
	userId := utils.GetUserId(c.Context())

	songs, err := h.svc.GetPlaylistSongs(c.Context(), role, userId, playlistId, pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "playlist_handler", "GetPlaylistSongs", err)
	}

	return c.JSON(dto.ResponseWithData[[]dto.Song]{
		Data: songs,
	})
}

// @Summary 		Added song to playlist
// @Description 	Added song to playlist
// @Tags        	playlists
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Playlist ID"
// @Param 			songId	path int true "Song ID"
// @Success 		201 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404		{object} 	dto.ValidationErrorResponse "Not Found: Playlist or song does not exists."
// @Failure 		409		{object} 	dto.ValidationErrorResponse "Conflict: Song already exists on playlist."
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/playlists/{id}/songs/{songId} [post]
func (h *PlaylistHandler) CreatePlaylistSong(c *fiber.Ctx) error {
	playlistId, _ := strconv.Atoi(c.Params("id"))
	songId, _ := strconv.Atoi(c.Params("songId"))
	userRole := utils.GetRole(c.Context())
	userId := utils.GetUserId(c.Context())

	if err := h.svc.CreatePlaylistSong(c.Context(), userRole, userId, playlistId, songId); err != nil {
		return errs.HandleHTTPError(c, h.log, "playlist_handler", "CreatePlaylistSong", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully added song to playlists.",
	})
}

// @Summary 		Delete song from playlist
// @Description 	Delete genre from playlist.
// @Tags        	playlists
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path int true "Playlist ID"
// @Param 			songId 	path int true "Song ID"
// @Success 		200 	{object} 	dto.ResponseMessage
// @Failure 		404		{object} 	dto.ValidationErrorResponse "Not Found: Song on playlist does not exists."
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/playlists/{id}/songs/{songId} [delete]
func (h *PlaylistHandler) DeletePlaylistSong(c *fiber.Ctx) error {
	playlistId, _ := strconv.Atoi(c.Params("id"))
	songId, _ := strconv.Atoi(c.Params("songId"))
	userRole := utils.GetRole(c.Context())
	userId := utils.GetUserId(c.Context())

	if err := h.svc.DeletePlaylistSong(c.Context(), userRole, userId, playlistId, songId); err != nil {
		return errs.HandleHTTPError(c, h.log, "playlist_handler", "DeletePlaylist", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully deleted song from playlist.",
	})
}
