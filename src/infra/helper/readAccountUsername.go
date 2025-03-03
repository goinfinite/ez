package infraHelper

import (
	"os/user"

	"github.com/goinfinite/ez/src/domain/valueObject"
)

func ReadAccountUsername(accountId valueObject.AccountId) valueObject.UnixUsername {
	unknownUsername := valueObject.UnixUsername("unknown")

	userInfo, err := user.LookupId(accountId.String())
	if err != nil {
		return unknownUsername
	}

	accountUsername, err := valueObject.NewUnixUsername(userInfo.Username)
	if err != nil {
		return unknownUsername
	}

	return accountUsername
}
