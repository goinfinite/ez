package infraHelper

import "github.com/goinfinite/ez/src/domain/valueObject"

// In order to run containers, the user must have a systemd session.
// Since the user is not actually logged in the system, the session
// has not been created yet. This is why the containers fail to start.
// The solution is to enable lingering for the user when adding the
// first container and disable it when removing the account.

func DisableLingering(accountId valueObject.AccountId) error {
	_, err := RunCmd("loginctl", "disable-linger", accountId.String())
	return err
}

func EnableLingering(accountId valueObject.AccountId) error {
	_ = DisableLingering(accountId)
	_, err := RunCmd("loginctl", "enable-linger", accountId.String())
	return err
}
