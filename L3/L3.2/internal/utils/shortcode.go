package utils

import (
	"crypto/rand"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateShortCode генерирует случайный короткий код длиной n (по умолчанию 8)
func GenerateShortCode(n int) (string, error) {
	if n <= 0 {
		n = 8
	}
	b := make([]byte, n)
	r := make([]byte, n)
	if _, err := rand.Read(r); err != nil {
		return "", err
	}
	for i := 0; i < n; i++ {
		b[i] = alphabet[int(r[i])%len(alphabet)]
	}
	return string(b), nil
}
