package secure

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomToken(n int) (string, error) {
	// n = bytes, 32 → 256-bit
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil { return "", err }
	return base64.RawURLEncoding.EncodeToString(b), nil
}
