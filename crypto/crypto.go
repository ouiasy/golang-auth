package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("error generating random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
