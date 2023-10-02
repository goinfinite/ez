package valueObject

import "errors"

type ScalingPolicy string

const (
	connection ScalingPolicy = "connection"
	cpu        ScalingPolicy = "cpu"
	memory     ScalingPolicy = "memory"
)

func NewScalingPolicy(value string) (ScalingPolicy, error) {
	sp := ScalingPolicy(value)
	if !sp.isValid() {
		return "", errors.New("InvalidScalingPolicy")
	}
	return sp, nil
}

func NewScalingPolicyPanic(value string) ScalingPolicy {
	sp, err := NewScalingPolicy(value)
	if err != nil {
		panic(err)
	}
	return sp
}

func (sp ScalingPolicy) isValid() bool {
	switch sp {
	case connection, cpu, memory:
		return true
	default:
		return false
	}
}

func (sp ScalingPolicy) String() string {
	return string(sp)
}
