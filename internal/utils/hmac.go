package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func VerifySignature(timestamp string, payload []byte, secret string, claimedSignature string) bool {
	signedPayload := timestamp + "." + string(payload)

	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(signedPayload)); err != nil {
		return false
	}
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	matched := hmac.Equal([]byte(expectedSignature), []byte(claimedSignature))
	return matched
}
