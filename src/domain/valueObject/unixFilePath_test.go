package valueObject

import "testing"

func TestUnixFilePath(t *testing.T) {
	ValidRawUnixFilePaths := []interface{}{
		"/", "/root", "/root/", "/home/sandbox/file.php", "/home/sandbox/file",
		"/home/sandbox/file with space.php", "/home/100sandbox/file.php",
		"/home/100sandbox/Imagem - Sem Título.jpg",
		"/home/100sandbox/Imagem - Sem Título & BW.jpg",
		"/home/100sandbox/Imagem - Sem Título # BW.jpg",
		"/home/@directory/file.gif", "/file.php", "/file.tar.br", "/file with space.php",
	}

	t.Run("ValidUnixFilePath", func(t *testing.T) {
		for _, filePath := range ValidRawUnixFilePaths {
			_, err := NewUnixFilePath(filePath)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), filePath)
			}
		}
	})

	t.Run("InvalidUnixFilePath", func(t *testing.T) {
		invalidUnixFilePaths := []interface{}{
			"", "/home/user/file.php?blabla",
			"/home/sandbox/domains/@<php52.sandbox.ntorga.com>",
			"../file.php", "./file.php", "file.php", "/home/../file.php",
		}

		for _, filePath := range invalidUnixFilePaths {
			_, err := NewUnixFilePath(filePath)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", filePath)
			}
		}
	})

	t.Run("ReadFileExtension", func(t *testing.T) {
		for _, rawFilePath := range ValidRawUnixFilePaths {
			filePath := UnixFilePath(rawFilePath.(string))

			if filePath.ReadWithoutExtension() == filePath {
				continue
			}

			_, err := filePath.ReadFileExtension()
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), filePath)
			}

			_, err = filePath.ReadCompoundFileExtension()
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), filePath)
			}
		}
	})
}
