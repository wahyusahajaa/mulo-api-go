package verification

import (
	"context"
	"fmt"

	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

const maxRetries = 5

type verificationService struct {
	userRepo contracts.UserRepository
}

func NewVerificationService(userRepo contracts.UserRepository) VerificationService {
	return &verificationService{
		userRepo: userRepo,
	}
}

func (v *verificationService) GenerateVerificationCode(ctx context.Context) (code string, err error) {
	for i := 0; i < maxRetries; i++ {
		code, err = utils.GenerateUniqueCode()
		if err != nil {
			return "", err
		}

		exists, err := v.userRepo.FindUserVerifiedByCode(ctx, code)
		if err != nil {
			return "", err
		}

		if !exists {
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique code after %d retries", maxRetries)
}
