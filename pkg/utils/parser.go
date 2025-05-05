package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

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

// Parse request body for check json input
func ParseBody(body []byte, out interface{}) (map[string]string, error) {
	var raw map[string]json.RawMessage

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	t := reflect.TypeOf(out).Elem()
	v := reflect.ValueOf(out).Elem()

	errorMap := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		jsonKey := strings.Split(jsonTag, ",")[0]
		rawValue, ok := raw[jsonKey]
		if !ok {
			continue
		}

		// Try to unmarshal into the field
		if err := json.Unmarshal(rawValue, fieldVal.Addr().Interface()); err != nil {
			errorMap[jsonKey] = fmt.Sprintf("Expected type %s", field.Type.Name())
		}
	}

	if len(errorMap) > 0 {
		return errorMap, nil
	}

	// Full decode if all types are OK
	if err := json.Unmarshal(body, out); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return nil, nil
}
