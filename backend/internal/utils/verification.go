package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)

	for i := range code {
		n, err := rand.Int(rand.Reader, bigInt(len(digits)))
		if err != nil {
			return "", fmt.Errorf("error generating random number: %w", err)
		}
		code[i] = digits[n.Int64()]
	}

	return string(code), nil
}

func bigInt(n int) *big.Int {
	return big.NewInt(int64(n))
}
