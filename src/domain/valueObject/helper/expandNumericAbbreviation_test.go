package voHelper

import (
	"testing"
)

func TestNewExpandNumericAbbreviation(t *testing.T) {
	t.Run("ValidExpandNumericAbbreviation", func(t *testing.T) {
		validExpandNumericAbbreviations := []string{
			"1",
			"1K",
			"1k",
			"1M",
			"1m",
			"1B",
			"1b",
			"1.5",
			"1.5K",
			"1.5k",
			"1.5M",
			"1.5m",
			"1.5B",
			"1.5b",
			"1.5K+",
			"1.5k+",
			"1.5M+",
			"1.5m+",
			"1.5B+",
			"1.5b+",
		}

		for _, numericAbbreviation := range validExpandNumericAbbreviations {
			_, err := ExpandNumericAbbreviation(numericAbbreviation)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), numericAbbreviation)
			}
		}
	})

	t.Run("InvalidExpandNumericAbbreviation", func(t *testing.T) {
		invalidExpandNumericAbbreviations := []string{
			"",
			"blabla",
			"@testing",
			"UNION SELECT * FROM USERS",
		}

		for _, numericAbbreviation := range invalidExpandNumericAbbreviations {
			_, err := ExpandNumericAbbreviation(numericAbbreviation)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", numericAbbreviation)
			}
		}
	})
}
