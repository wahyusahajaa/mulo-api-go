package utils

import (
	"regexp"
	"strings"
)

func MakeSlug(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Remove non-alphanumeric characters except spaces
	re := regexp.MustCompile(`[^\w\s-]`)
	text = re.ReplaceAllString(text, "")

	// Replace multiple spaces or underscores with single hyphen
	re = regexp.MustCompile(`[\s_]+`)
	text = re.ReplaceAllString(text, "-")

	return text
}
