package cliMiddleware

import (
	"fmt"
	"os"
	"os/user"
)

func PreventRootless() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get current user:", err)
		os.Exit(1)
	}

	if currentUser.Username != "root" {
		fmt.Println("Only root can run Ez.")
		os.Exit(1)
	}
}
