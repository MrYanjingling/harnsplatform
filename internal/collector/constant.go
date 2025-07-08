package collector

import (
	"errors"
	"time"
)

var (
	ErrDeviceType          = errors.New("unsupported device type")
	ErrConnectDevice       = errors.New("unable to connect to device")
	ErrDeviceServerClosed  = errors.New("device server closed")
	ErrDeviceEmptyVariable = errors.New("device variable emptied")
)

var heartBeatTimeInterval = 15 * time.Second
