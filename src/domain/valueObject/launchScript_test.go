package valueObject

import (
	"testing"
)

func TestNewLaunchScript(t *testing.T) {
	t.Run("ValidLaunchScript", func(t *testing.T) {
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

	t.Run("ValidEncodedLaunchScript", func(t *testing.T) {
		validEncodedLaunchScripts := []string{
			`ZWNobyAiSGVsbG8sIFdvcmxkISI=`,
			`cHJpbnRlbnYgPiAvdG1wL2hlbGxvLnR4dA==`,
			`VEVYVD0iSGVsbG8sIFdvcmxkISIKZWNobyAiJHtURVhUfSIgPiAvdG1wL2hlbGxvLnR4dDsKVEVYVD0iR29vZGJ5ZSwgV29ybGQhIgplY2hvICIke1RFWFR9IiA+IC90bXAvZ29vZGJ5ZS50eHQ=`,
		}

		for _, script := range validEncodedLaunchScripts {
			encodedScript, err := NewEncodedContent(script)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), script)
			}

			_, err = NewLaunchScriptFromEncodedContent(encodedScript)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), script)
			}
		}
	})

	t.Run("InvalidLaunchScript", func(t *testing.T) {
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
