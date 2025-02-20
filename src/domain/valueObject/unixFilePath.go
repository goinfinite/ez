package valueObject

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const unixFilePathRegexExpression = `^\/?[^\n\r\t\f\0\?\[\]\<\>]+$`
const unixFileRelativePathRegexExpression = `\.\.\/|^\.\/|^\/\.\/`

type UnixFilePath string

func NewUnixFilePath(value interface{}) (filePath UnixFilePath, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return filePath, errors.New("UnixFilePathValueMustBeString")
	}

	unixFilePathRegex := regexp.MustCompile(unixFilePathRegexExpression)
	if !unixFilePathRegex.MatchString(stringValue) {
		return filePath, errors.New("InvalidUnixFilePath")
	}

	isOnlyFileName := !strings.Contains(stringValue, "/")
	if isOnlyFileName {
		return filePath, errors.New("PathIsFileNameOnly")
	}

	unixFileRelativePathRegex := regexp.MustCompile(unixFileRelativePathRegexExpression)
	if unixFileRelativePathRegex.MatchString(stringValue) {
		return filePath, errors.New("RelativePathNotAllowed")
	}

	return UnixFilePath(stringValue), nil
}

func (vo UnixFilePath) ReadWithoutExtension() UnixFilePath {
	unixFilePathExtStr := filepath.Ext(string(vo))
	if unixFilePathExtStr == "" {
		return vo
	}

	unixFilePathWithoutExtStr := strings.TrimSuffix(string(vo), unixFilePathExtStr)
	unixFilePathWithoutExt, _ := NewUnixFilePath(unixFilePathWithoutExtStr)
	return unixFilePathWithoutExt
}

func (vo UnixFilePath) ReadFileName() UnixFileName {
	unixFileBase := filepath.Base(string(vo))
	unixFileName, _ := NewUnixFileName(unixFileBase)
	return unixFileName
}

func (vo UnixFilePath) ReadFileExtension() (UnixFileExtension, error) {
	unixFileExtensionStr := filepath.Ext(string(vo))
	return NewUnixFileExtension(unixFileExtensionStr)
}

func (vo UnixFilePath) ReadCompoundFileExtension() (UnixFileExtension, error) {
	fileNameParts := strings.Split(vo.ReadFileName().String(), ".")
	if len(fileNameParts) < 3 {
		return vo.ReadFileExtension()
	}
	extensionsOnly := fileNameParts[1:]
	return NewUnixFileExtension(strings.Join(extensionsOnly, "."))
}

func (vo UnixFilePath) ReadFileNameWithoutExtension() UnixFileName {
	fileBase := filepath.Base(string(vo))
	fileExt, err := vo.ReadCompoundFileExtension()
	if err != nil {
		return vo.ReadFileName()
	}
	fileBaseWithoutExtStr := strings.TrimSuffix(fileBase, "."+fileExt.String())
	fileNameWithoutExt, _ := NewUnixFileName(fileBaseWithoutExtStr)
	return fileNameWithoutExt
}

func (vo UnixFilePath) ReadFileDir() UnixFilePath {
	unixFileDirPath, _ := NewUnixFilePath(filepath.Dir(string(vo)))
	return unixFileDirPath
}

func (vo UnixFilePath) String() string {
	return string(vo)
}
