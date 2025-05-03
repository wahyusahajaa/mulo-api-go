package utils

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
)

func ParseUserId(c *fiber.Ctx) (int, error) {
	idVal := c.Locals("id")
	if idVal == nil {
		return 0, errors.New("user ID not found in context")
	}

	idFloat, ok := idVal.(float64)
	if !ok {
		return 0, errors.New("user ID is not a valid float64")
	}

	return int(idFloat), nil
}

// Parse image from byte to json
func ParseImageToJSON(img []byte) dto.Image {
	image := dto.Image{}

	if len(img) > 0 {
		_ = json.Unmarshal(img, &image)
	}

	return image
}

// Parse image from json to byte
func ParseImageToByte(image *dto.Image) (imgByte []byte, err error) {
	if image != nil {
		imgByte, err = json.Marshal(image)
		return
	}

	return nil, nil
}
