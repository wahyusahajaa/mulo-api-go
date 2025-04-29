package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var customMessage = map[string]string{
	"required": "This field is required",
	"min":      "Value is too short",
	"email":    "Invalid email format",
}

func GetJSONFieldName(structType any, fieldName string) string {
	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if field, ok := t.FieldByName(fieldName); ok {
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			return strings.Split(tag, ",")[0]
		}
	}
	return fieldName
}

func RequestValidate(input any) (map[string]string, error) {
	err := validate.Struct(input)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		errorsMap := map[string]string{}

		for _, e := range errs {
			field := GetJSONFieldName(input, e.StructField())
			tag := e.Tag()

			msg, ok := customMessage[tag]
			if !ok {
				msg = "Validation failed on " + tag
			}

			errorsMap[field] = msg
		}

		return errorsMap, fmt.Errorf("validation failed")
	}

	return nil, nil
}
