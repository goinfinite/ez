package valueObject

import (
	"testing"
)

func TestNewLaunchScript(t *testing.T) {
	t.Run("ValidLaunchScriptFromString", func(t *testing.T) {
		validLaunchScripts := []string{
			`echo "Hello, World!"`,
			`echo "Hello, World!" > /tmp/hello.txt`,
			`echo "Hello, World!" > /tmp/hello.txt; echo "Goodbye, World!" > /tmp/goodbye.txt`,
			`echo "Hello, World!" > /tmp/hello.txt
			echo "Goodbye, World!" > /tmp/goodbye.txt;`,
			`TEXT="Hello, World!"
			echo "${TEXT}" > /tmp/hello.txt;
			TEXT="Goodbye, World!"
			echo "${TEXT}" > /tmp/goodbye.txt`,
		}

		for _, script := range validLaunchScripts {
			_, err := NewLaunchScript(script)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), script)
			}
		}
	})

	t.Run("InvalidLaunchScriptFromString", func(t *testing.T) {
		invalidLaunchScripts := []string{
			"",
			"ls",
		}

		for _, script := range invalidLaunchScripts {
			_, err := NewLaunchScript(script)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", script)
			}
		}
	})
}
