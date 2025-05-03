package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
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

func (h *ArtistHandler) GetArtists(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	artists, err := h.svc.GetAll(c.Context(), pageSize, offset)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	results := []dto.Artist{}

	for _, v := range artists {
		result := dto.Artist{
			Id:    v.Id,
			Name:  v.Name,
			Slug:  v.Slug,
			Image: utils.ParseImageToJSON(v.Image),
		}
		results = append(results, result)
	}

	total, err := h.svc.GetCount(c.Context())

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

func (h *ArtistHandler) CreateArtist(c *fiber.Ctx) error {
	var req dto.CreateArtistRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if errorMaps, err := utils.RequestValidate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errorMaps,
		})
	}

	slug := utils.MakeSlug(req.Name)
	exists, err := h.svc.CheckDuplicateArtistBySlug(c.Context(), slug)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Artist name already exists",
		})
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid image object",
		})
	}

	if err := h.svc.Create(c.Context(), req.Name, slug, imgByte); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully create artist",
	})
}

func (h *ArtistHandler) GetArtist(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))

	artist, err := h.svc.GetArtistById(c.Context(), artistId)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if artist == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Artist not found",
		})
	}

	result := dto.Artist{
		Id:    artist.Id,
		Name:  artist.Name,
		Slug:  artist.Slug,
		Image: utils.ParseImageToJSON(artist.Image),
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *ArtistHandler) UpdateArtist(c *fiber.Ctx) error {
	var req dto.CreateArtistRequest
	artistId, _ := strconv.Atoi(c.Params("id"))

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if errorMaps, err := utils.RequestValidate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation errors",
			"errors":  errorMaps,
		})
	}

	artist, err := h.svc.GetArtistById(c.Context(), artistId)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if artist == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Artist not found",
		})
	}

	slug := utils.MakeSlug(req.Name)

	if artist.Name != req.Name {
		exists, err := h.svc.CheckDuplicateArtistBySlug(c.Context(), slug)

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if exists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Artist name already exists",
			})
		}
	}

	imgByte, err := utils.ParseImageToByte(req.Image)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := h.svc.Update(c.Context(), req.Name, slug, imgByte, artistId); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully update artist",
	})
}

func (h *ArtistHandler) DeleteArtist(c *fiber.Ctx) error {
	artistId, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.Delete(c.Context(), artistId); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully delete artist",
	})
}
