package dto

import "github.com/speedianet/control/src/domain/valueObject"

type UpdateFailedAttempt struct {
	IpAddress valueObject.IpAddress `json:"ipAddress"`
	Count     int                   `json:"count"`
}
