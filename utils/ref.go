package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateReferenceNumber() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "TRX" + hex.EncodeToString(bytes)
}
