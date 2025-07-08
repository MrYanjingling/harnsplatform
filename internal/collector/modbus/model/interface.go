package model

import (
	"harnsplatform/internal/collector/modbus/runtime"
	"harnsplatform/internal/common"
)

var _ ModbusModeler = (*ModbusTcp)(nil)
var _ ModbusModeler = (*ModbusRtu)(nil)
var _ ModbusModeler = (*ModbusRtuOverTcp)(nil)

var ModbusModelers = map[string]ModbusModeler{
	"modbusTcp":        &ModbusTcp{},
	"modbusRtu":        &ModbusRtu{},
	"modbusRtuOverTcp": &ModbusRtuOverTcp{},
}

type ModbusModeler interface {
	GenerateReadMessage(slave uint, functionCode uint8, startAddress uint, maxDataSize uint, variables []*runtime.VariableParse, memoryLayout common.MemoryLayout) *runtime.ModBusDataFrame
	NewClients(address *runtime.Address, dataFrameCount int) (*runtime.Clients, error)
	// ExecuteAction(messenger modbus.Messenger, variables []*modbus.Variable) error
}
