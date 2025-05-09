package verification

import "context"

type VerificationService interface {
	GenerateVerificationCode(ctx context.Context) (code string, err error)
}
