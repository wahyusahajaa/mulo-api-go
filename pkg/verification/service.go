package verification

import (
	"context"
	"fmt"

	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

const maxRetries = 5

type verificationService struct {
	authRepo contracts.AuthRepository
}

func NewVerificationService(authRepo contracts.AuthRepository) VerificationService {
	return &verificationService{
		authRepo: authRepo,
	}
}

func (v *verificationService) GenerateVerificationCode(ctx context.Context) (code string, err error) {
	for i := 0; i < maxRetries; i++ {
		code, err = utils.GenerateUniqueCode()
		if err != nil {
			return "", err
		}

		exists, err := v.authRepo.FindUserVerifiedByCode(ctx, code)
		if err != nil {
			return "", err
		}

		if !exists {
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique code after %d retries", maxRetries)
}
