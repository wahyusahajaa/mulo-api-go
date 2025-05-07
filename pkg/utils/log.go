package utils

import (
	"context"

	"github.com/sirupsen/logrus"
)

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
