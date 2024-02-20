package infraHelper

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

func GenHash(value string) string {
	hash := sha3.New256()
	hash.Write([]byte(value))
	hashStr := hex.EncodeToString(hash.Sum(nil))
	return hashStr
}

func GenShortHash(value string) string {
	return GenHash(string(value))[:12]
}
