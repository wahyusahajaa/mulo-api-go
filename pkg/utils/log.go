package utils

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Get requestId from context
func GetRequestId(ctx context.Context) string {
	v := ctx.Value("requestId")
	if id, ok := v.(string); ok {
		return id
	}
	return ""
}

// Set log error
func LogError(log *logrus.Logger, ctx context.Context, layer, operation string, err error) {
	log.WithFields(logrus.Fields{
		"layer":     layer,
		"operation": operation,
		"error":     err.Error(),
		"requestId": GetRequestId(ctx),
	}).Error("operation failed")
}

func LogWarn(log *logrus.Logger, ctx context.Context, layer, operation string, err error) {
	log.WithFields(logrus.Fields{
		"layer":     layer,
		"operation": operation,
		"error":     err.Error(),
		"requestId": GetRequestId(ctx),
	}).Warn("operation completed with warnings")
}
