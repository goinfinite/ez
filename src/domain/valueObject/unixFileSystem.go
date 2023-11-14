package valueObject

import (
	"errors"
	"regexp"
)

const fileSystemRegexp = `^[\w-]{3,24}$`

type UnixFileSystem string

func NewUnixFileSystem(value string) (UnixFileSystem, error) {
	fs := UnixFileSystem(value)
	if !fs.isValid() {
		return "", errors.New("InvalidUnixFileSystem")
	}
	return fs, nil
}

func NewUnixFileSystemPanic(value string) UnixFileSystem {
	fs, err := NewUnixFileSystem(value)
	if err != nil {
		panic(err)
	}
	return fs
}

func (fs UnixFileSystem) isValid() bool {
	fsRegex := regexp.MustCompile(fileSystemRegexp)
	return fsRegex.MatchString(string(fs))
}

func (fs UnixFileSystem) String() string {
	return string(fs)
}
