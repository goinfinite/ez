package valueObject

import (
	"testing"
)

func TestNewContainerProfileName(t *testing.T) {
	t.Run("ValidContainerProfileName", func(t *testing.T) {
		validContainerProfileName := []any{
			"névoa", "cirrus", "cirrostratus", "stratocumulus", "cumulus",
			"altostratus", "testContainer", "production",
		}

		for _, name := range validContainerProfileName {
			_, err := NewContainerProfileName(name)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", name, err.Error())
			}
		}
	})

	t.Run("InvalidContainerProfileName", func(t *testing.T) {
		invalidContainerProfileName := []any{
			"", 10,
		}

		for _, name := range invalidContainerProfileName {
			_, err := NewContainerProfileName(name)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", name)
			}
		}
	})
}
