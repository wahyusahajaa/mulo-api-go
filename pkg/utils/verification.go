package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
)

func GenerateUniqueCode() (string, error) {
	max := big.NewInt(100000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate verification code: %w", err)
	}
	return fmt.Sprintf("%05d", n.Int64()), nil
}

type VerificationService interface {
	GenerateVerificationCode(ctx context.Context) (code string, err error)
}

type verificationService struct {
	authRepo contracts.AuthRepository
}

func NewVerification(authRepo contracts.AuthRepository) VerificationService {
	return &verificationService{
		authRepo: authRepo,
	}
}

func (v *verificationService) GenerateVerificationCode(ctx context.Context) (code string, err error) {
	var maxRetries = 5
	for i := 0; i < maxRetries; i++ { // retry up to 5 times
		code, _ = GenerateUniqueCode()
		exists, err := v.authRepo.FindUserVerifiedByCode(ctx, code)

		if err != nil {
			log.Printf("check verify code error: %v", err)
			return "", err
		}

		if !exists {
			return code, nil
		}

		if i == maxRetries-1 {
			// failed If after max retries no unique code found
			return "", fmt.Errorf("failed to generate unique code")
		}
	}

	return "", nil
}
