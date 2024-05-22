package valueObject

import "errors"

type NetworkPortInterval struct {
	Min NetworkPort `json:"min"`
	Max NetworkPort `json:"max"`
}

func NewNetworkPortInterval(
	min NetworkPort,
	max NetworkPort,
) (portInterval NetworkPortInterval, err error) {
	if min.Get() > max.Get() {
		return portInterval, errors.New("MinPortGreaterThanMaxPort")
	}

	return NetworkPortInterval{
		Min: min,
		Max: max,
	}, nil
}
