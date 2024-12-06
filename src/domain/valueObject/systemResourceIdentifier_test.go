package valueObject

import "testing"

func TestSystemResourceIdentifier(t *testing.T) {
	t.Run("ValidSystemResourceIdentifier", func(t *testing.T) {
		validSystemResourceIdentifier := []interface{}{
			"sri://0:account/120", "sri://0:cron/1", "sri://1:mapping/2",
		}

		for _, rawIdentifier := range validSystemResourceIdentifier {
			_, err := NewSystemResourceIdentifier(rawIdentifier)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), rawIdentifier)
			}
		}
	})

	t.Run("InvalidSystemResourceIdentifier", func(t *testing.T) {
		invalidSystemResourceIdentifier := []interface{}{
			"", "sri://0:/",
		}

		for _, rawIdentifier := range invalidSystemResourceIdentifier {
			_, err := NewSystemResourceIdentifier(rawIdentifier)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", rawIdentifier)
			}
		}
	})
}
