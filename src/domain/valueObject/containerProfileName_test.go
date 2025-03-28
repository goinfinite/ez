package valueObject

import (
	"testing"
)

func TestNewContainerProfileName(t *testing.T) {
	t.Run("ValidContainerProfileName", func(t *testing.T) {
		validContainerProfileNames := []any{
			"névoa", "cirrus", "cirrostratus", "stratocumulus", "cumulus",
			"altostratus", "testContainer", "production",
		}

		for _, name := range validContainerProfileNames {
			_, err := NewContainerProfileName(name)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", name, err.Error())
			}
		}
	})

	t.Run("InvalidContainerProfileName", func(t *testing.T) {
		invalidContainerProfileNames := []any{
			"", 10,
		}

		for _, name := range invalidContainerProfileNames {
			_, err := NewContainerProfileName(name)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", name)
			}
		}
	})
}
