package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type AlbumHandler struct {
	albumService  contracts.AlbumService
	artistService contracts.ArtistService
	log           *logrus.Logger
}

func NewAlbumHandler(albumService contracts.AlbumService, artistService contracts.ArtistService, log *logrus.Logger) *AlbumHandler {
	return &AlbumHandler{
		albumService:  albumService,
		artistService: artistService,
		log:           log,
	}
}

func (h *AlbumHandler) GetAlbums(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	albums, err := h.albumService.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Deduplicate artist IDs
	artistIdMap := make(map[int]struct{})
	for _, album := range albums {
		artistIdMap[album.ArtistId] = struct{}{}
	}

	artistIds := make([]any, 0, len(artistIdMap))
	for id := range artistIdMap {
		artistIds = append(artistIds, id)
	}

	inClause, args := utils.BuildInClause(1, artistIds)

	artists, err := h.artistService.GetArtistByIds(c.Context(), inClause, args)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	//  Build artist lookup map
	artistMap := make(map[int]models.Artist)
	for _, artist := range artists {
		artistMap[artist.Id] = artist
	}

	results := make([]dto.Album, 0, len(albums))
	for _, album := range albums {
		result := dto.Album{
			Id:    album.Id,
			Name:  album.Name,
			Slug:  album.Slug,
			Image: utils.ParseImageToJSON(album.Image),
		}

		if artist, ok := artistMap[album.ArtistId]; ok {
			result.Artist.Id = artist.Id
			result.Artist.Name = artist.Name
			result.Artist.Slug = artist.Slug
			result.Artist.Image = utils.ParseImageToJSON(artist.Image)
		}

		results = append(results, result)
	}

	total, err := h.albumService.GetCount(c.Context())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	pagination := dto.Pagination{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	return c.JSON(fiber.Map{
		"data":       results,
		"pagination": pagination,
	})
}

func (h *AlbumHandler) GetAlbum(c *fiber.Ctx) error {
	albumId, _ := strconv.Atoi(c.Params("id"))

	album, err := h.albumService.GetAlbumById(c.Context(), albumId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if album == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Album not found",
		})
	}

	result := dto.Album{
		Id:    album.Id,
		Name:  album.Name,
		Slug:  album.Slug,
		Image: utils.ParseImageToJSON(album.Image),
		Artist: dto.Artist{
			Id:    album.Artist.Id,
			Name:  album.Artist.Name,
			Slug:  album.Artist.Slug,
			Image: utils.ParseImageToJSON(album.Artist.Image),
		},
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *AlbumHandler) CreateAlbum(c *fiber.Ctx) error {
	var req dto.CreateAlbumRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation errors",
			"errors":  errorsMap,
		})
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid image object",
		})
	}

	artist, err := h.artistService.GetArtistById(c.Context(), req.ArtistId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if artist == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Artist with id: %d not found", req.ArtistId),
		})
	}

	slug := utils.MakeSlug(req.Name)
	exists, err := h.albumService.CheckDuplicateAlbumBySlug(c.Context(), slug)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Album name already exists",
		})
	}

	if err := h.albumService.Create(c.Context(), req.ArtistId, req.Name, slug, imgByte); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Suceessfully create album",
	})
}

func (h *AlbumHandler) UpdateAlbum(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.CreateAlbumRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation errors",
			"errors":  errorsMap,
		})
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid image object",
		})
	}

	album, err := h.albumService.GetAlbumById(c.Context(), id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if album == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Album not found",
		})
	}

	artist, err := h.artistService.GetArtistById(c.Context(), req.ArtistId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if artist == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Artist with id: %d not found", req.ArtistId),
		})
	}

	slug := utils.MakeSlug(req.Name)
	if album.Slug != slug {
		exists, err := h.albumService.CheckDuplicateAlbumBySlug(c.Context(), slug)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if exists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Album name already exitsts",
			})
		}
	}

	if err := h.albumService.UpdateAlbum(c.Context(), req.ArtistId, req.Name, slug, imgByte, id); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully update album",
	})
}

func (h *AlbumHandler) DeleteAlbum(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	album, err := h.albumService.GetAlbumById(c.Context(), id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if album == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Album not found",
		})
	}

	if err := h.albumService.DeleteAlbum(c.Context(), id); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully delete album",
	})
}
