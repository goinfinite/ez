package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type InstructionSetArchitecture string

var ValidInstructionSetArchitectures = []string{
	"amd64",
	"arm",
	"arm64",
	"i386",
	"riscv64",
}

func NewInstructionSetArchitecture(value string) (InstructionSetArchitecture, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	if !slices.Contains(ValidInstructionSetArchitectures, value) {
		return "", errors.New("InvalidInstructionSetArchitecture")
	}
	return InstructionSetArchitecture(value), nil
}

func NewInstructionSetArchitecturePanic(value string) InstructionSetArchitecture {
	isa, err := NewInstructionSetArchitecture(value)
	if err != nil {
		panic(err)
	}
	return isa
}

func (isa InstructionSetArchitecture) String() string {
	return string(isa)
}
