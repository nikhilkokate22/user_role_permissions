package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func HashSHA256(value string) string {
	hash := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(value))))
	return hex.EncodeToString(hash[:])
}

