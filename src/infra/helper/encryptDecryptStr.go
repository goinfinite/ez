package infraHelper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
)

func EncryptStr(secretKey string, plainText string) (string, error) {
	secretKeyBytes, err := base64.RawURLEncoding.DecodeString(secretKey)
	if err != nil {
		log.Printf("EncryptSecretKeyError: %s", err)
		return "", errors.New("EncryptSecretKeyError")
	}

	block, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		log.Printf("EncryptCipherError: %s", err)
		return "", errors.New("EncryptCipherError")
	}

	plainTextBytes := []byte(plainText)
	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Printf("EncryptIvGenerationError: %s", err)
		return "", errors.New("EncryptIvGenerationError")
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptStr(secretKey string, encryptedText string) (string, error) {
	apiKeyDecoded, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", errors.New("DecryptDecodingError")
	}
	if len(apiKeyDecoded) < aes.BlockSize {
		return "", errors.New("DecryptDecodedTooShort")
	}

	secretKeyBytes, err := base64.RawURLEncoding.DecodeString(secretKey)
	if err != nil {
		return "", errors.New("DecryptSecretDecodingError")
	}

	block, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		return "", errors.New("DecryptCipherError")
	}

	apiKeyDecryptedBinary := make([]byte, len(apiKeyDecoded)-aes.BlockSize)
	iv := apiKeyDecoded[:aes.BlockSize]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(apiKeyDecryptedBinary, apiKeyDecoded[aes.BlockSize:])

	return string(apiKeyDecryptedBinary), nil
}
