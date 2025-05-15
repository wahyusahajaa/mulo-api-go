package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func RequestValidate(req any) (map[string]string, error) {
	if err := validate.Struct(req); err != nil {
		errorMap := make(map[string]string)

		for _, fe := range err.(validator.ValidationErrors) {
			fieldName := getJSONFieldName(req, fe.StructField())
			message := getErrorMessage(fe)

			errorMap[fieldName] = message
		}

		return errorMap, fmt.Errorf("validation failed")
	}

	return nil, nil
}

func getJSONFieldName(obj any, field string) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if f, ok := t.FieldByName(field); ok {
		tag := f.Tag.Get("json")
		if tag != "" {
			return strings.Split(tag, ",")[0]
		}
	}

	return strings.ToLower(field)
}

func getErrorMessage(fe validator.FieldError) string {
	customErrorMessage := map[string]string{
		"required": "Field is required",
		"min":      fmt.Sprintf("Minimum value is %s", fe.Param()),
		"max":      fmt.Sprintf("Maximum value is %s", fe.Param()),
		"email":    "Must be a valid email",
		"len":      fmt.Sprintf("Length must be %s characters", fe.Param()),
		"gt":       "Field must be Greater than 0",
	}

	if result, ok := customErrorMessage[fe.Tag()]; ok {
		return result
	}

	return "Invalid value"
}
