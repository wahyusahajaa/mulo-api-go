package utils

import (
	"context"
)

// Get requestId from context
func GetRequestId(ctx context.Context) string {
	v := ctx.Value("requestId")
	if id, ok := v.(string); ok {
		return id
	}
	return ""
}

// Get user_id from context
func GetUserId(ctx context.Context) (id int) {
	v := ctx.Value("id")
	if id, ok := v.(float64); ok {
		return int(id)
	}
	return
}

// Get role from context
func GetRole(ctx context.Context) string {
	return ctx.Value("role").(string)
}
