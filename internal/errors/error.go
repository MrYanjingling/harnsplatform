package errors

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
)

type ErrorReason int32

const (
	ErrorReason_GREETER_UNSPECIFIED ErrorReason = 0
	ErrorReason_USER_NOT_FOUND      ErrorReason = 1
	ErrorReason_RESOURCE_MISMATCH   ErrorReason = 2
	ErrorReason_RESOURCE_NOT_FOUND  ErrorReason = 4
)

// Enum value maps for ErrorReason.
var (
	ErrorReasonName = map[int32]string{
		0: "GREETER_UNSPECIFIED",
		1: "USER_NOT_FOUND",
		2: "RESOURCE_MISMATCH",
		3: "RESOURCE_PRECONDITION_REQUIRED",
		4: "RESOURCE_NOT_FOUND",
	}
	ErrorReasonValue = map[string]int32{
		"GREETER_UNSPECIFIED":            0,
		"USER_NOT_FOUND":                 1,
		"RESOURCE_MISMATCH":              2,
		"RESOURCE_PRECONDITION_REQUIRED": 3,
		"RESOURCE_NOT_FOUND":             4,
	}
)

func (x ErrorReason) Enum() *ErrorReason {
	p := new(ErrorReason)
	*p = x
	return p
}

func (x ErrorReason) String() string {
	return ErrorReasonName[int32(x)]
}

func GenerateResourceMismatchError(resourceName string) error {
	return errors.New(412, ErrorReason_RESOURCE_MISMATCH.String(), fmt.Sprintf("%s resource mismatch.", resourceName))
}

func GenerateResourcePreconditionRequiredError(resourceName string) error {
	return errors.New(428, ErrorReason_RESOURCE_MISMATCH.String(), fmt.Sprintf("%s resource precondition required.", resourceName))
}

func GenerateResourceNotFoundError(resourceName string) error {
	return errors.New(404, ErrorReason_RESOURCE_NOT_FOUND.String(), fmt.Sprintf("%s resource not found.", resourceName))
}
