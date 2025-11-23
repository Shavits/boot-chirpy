package auth

import (
	"crypto/rand"
	"encoding/hex"
)



func MakeRefreshToken() (string, error) {
    randData := make([]byte, 32)
    _, err := rand.Read(randData)
    if err != nil {
        return "", err
    }
    encodedString := hex.EncodeToString(randData)
    return encodedString, nil
}