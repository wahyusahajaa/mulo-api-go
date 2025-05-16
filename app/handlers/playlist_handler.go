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
