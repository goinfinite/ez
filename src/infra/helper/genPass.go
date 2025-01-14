package infraHelper

import (
	"math/rand"
	"strings"
)

func oneCharCharsetGuarantor(
	originalString []byte,
	charset string,
) []byte {
	if strings.ContainsAny(string(originalString), charset) {
		return originalString
	}

	randomStringIndex := rand.Intn(len(originalString))
	isFirstChar := randomStringIndex == 0
	if isFirstChar {
		randomStringIndex++
	}
	isLastChar := randomStringIndex == len(originalString)-1
	if isLastChar {
		randomStringIndex--
	}
	if randomStringIndex >= len(originalString) {
		randomStringIndex = len(originalString) - 1
	}

	randomCharsetIndex := rand.Intn(len(charset))
	originalString[randomStringIndex] = charset[randomCharsetIndex]

	return originalString
}

func GenPass(length int) string {
	lowercaseAlphabetCharset := "abcdefghijklmnopqrstuvwxyz"
	uppercaseAlphabetCharset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numericCharset := "0123456789"
	alphanumericCharset := lowercaseAlphabetCharset + uppercaseAlphabetCharset + numericCharset
	symbolCharset := "!@#$%^&*()_+"

	pass := make([]byte, length)
	previousCharset := alphanumericCharset
	currentCharset := alphanumericCharset

	for i := 0; i < length; i++ {
		currentCharset = alphanumericCharset

		if previousCharset != symbolCharset {
			nextCharsetShouldUseSymbol := rand.Float32() < 0.1
			isTipChar := i == 0 || i == length-1
			if nextCharsetShouldUseSymbol && !isTipChar {
				currentCharset = symbolCharset
			}
		}
		currentCharsetLen := len(currentCharset)

		randomCharsetIndex := rand.Intn(currentCharsetLen)
		pass[i] = currentCharset[randomCharsetIndex]
		previousCharset = currentCharset
	}

	if length > 4 {
		pass = oneCharCharsetGuarantor(pass, lowercaseAlphabetCharset)
		pass = oneCharCharsetGuarantor(pass, uppercaseAlphabetCharset)
		pass = oneCharCharsetGuarantor(pass, numericCharset)
		pass = oneCharCharsetGuarantor(pass, "!@#$%^&*()_+")
	}

	return string(pass)
}
