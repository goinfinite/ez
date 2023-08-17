package valueObject

import "errors"

type Hash string

func NewHash(value string) (Hash, error) {
	hash := Hash(value)
	if !hash.isValid() {
		return "", errors.New("InvalidHash")
	}
	return hash, nil
}

func NewHashPanic(value string) Hash {
	hash := Hash(value)
	if !hash.isValid() {
		panic("InvalidHash")
	}
	return hash
}

func (hash Hash) isValid() bool {
	isTooShort := len(string(hash)) < 6
	isTooLong := len(string(hash)) > 64
	return !isTooShort && !isTooLong
}

func (hash Hash) String() string {
	return string(hash)
}
