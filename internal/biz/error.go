package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "harnsplatform/api/modelmanager/v1"
)

func GenerateResourceMismatchError(resourceName string) error {
	return errors.New(412, v1.ErrorReason_RESOURCE_MISMATCH.String(), fmt.Sprintf("%s resource mismatch.", resourceName))
}
