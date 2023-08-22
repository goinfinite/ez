package valueObject

import "errors"

type SvcCmd string

func NewSvcCmd(value string) (SvcCmd, error) {
	cmd := SvcCmd(value)
	if !cmd.isValid() {
		return "", errors.New("InvalidSvcCmd")
	}
	return cmd, nil
}

func NewSvcCmdPanic(value string) SvcCmd {
	cmd, err := NewSvcCmd(value)
	if err != nil {
		panic(err)
	}
	return cmd
}

func (cmd SvcCmd) isValid() bool {
	isTooShort := len(string(cmd)) < 3
	isTooLong := len(string(cmd)) > 1024
	return !isTooShort && !isTooLong
}

func (cmd SvcCmd) String() string {
	return string(cmd)
}
