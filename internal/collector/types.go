package collector

import (
	"harnsplatform/internal/biz"
	"harnsplatform/internal/collector/modbus"
)

type ParseVariableResult struct {
	VariableSlice []VariableValue
	Err           []error
}

type NewBroker func(object Device) (Broker, chan *ParseVariableResult, error)

var DeviceTypeBrokerMap = map[string]NewBroker{
	"modbus": modbus.NewBroker,
	// "opcUa":  opcua.NewBroker,
	// "s7":     s7.NewBroker,
}

type ConvertDevice func(object *biz.Agents) Device

var ConvertDeviceMap = map[string]ConvertDevice{
	"modbus": modbus.ConvertDevice,
	// "opcUa":  opcua.NewBroker,
	// "s7":     s7.NewBroker,
}
