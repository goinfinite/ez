package infraHelper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/goinfinite/ez/src/domain/valueObject"
	"golang.org/x/crypto/sha3"
)

func GenWeakHash(value string) string {
	hash := md5.Sum([]byte(value))
	return hex.EncodeToString(hash[:])
}

func GenStrongHash(value string) string {
	hash := sha3.New256()
	hash.Write([]byte(value))
	return hex.EncodeToString(hash.Sum(nil))
}

func GenStrongShortHash(value string) string {
	return GenStrongHash(string(value))[:12]
}

func GenFileHash(filePath valueObject.UnixFilePath) (valueObject.Hash, error) {
	rawHash, err := RunCmd("sha256sum", filePath.String())
	if err != nil {
		return "", err
	}
	if len(rawHash) < 64 {
		return "", errors.New("InvalidHashLength")
	}

	return valueObject.NewHash(rawHash[:64])
}
