package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateUniqueCode() (string, error) {
	max := big.NewInt(100000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate verification code: %w", err)
	}
	return fmt.Sprintf("%05d", n.Int64()), nil
}
