package valueObject

import "errors"

type NetworkPortInterval struct {
	Min NetworkPort  `json:"min"`
	Max *NetworkPort `json:"max"`
}

func NewNetworkPortInterval(
	min NetworkPort,
	max *NetworkPort,
) (NetworkPortInterval, error) {
	if max != nil && min.Get() > max.Get() {
		return NetworkPortInterval{}, errors.New("MinPortGreaterThanMaxPort")
	}

	return NetworkPortInterval{
		Min: min,
		Max: max,
	}, nil
}
