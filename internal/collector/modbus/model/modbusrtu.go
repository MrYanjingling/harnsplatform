package model

import (
	"container/list"
	"go.bug.st/serial"
	"harnsplatform/internal/collector/modbus/runtime"
	"harnsplatform/internal/common"
	"harnsplatform/internal/utils"
	"harnsplatform/internal/utils/binutils"
	"k8s.io/klog/v2"
	"sync"
)

const RtuNonDataLength = 5

type ModbusRtu struct {
}

func (m *ModbusRtu) NewClients(address *runtime.Address, dataFrameCount int) (*runtime.Clients, error) {
	mode := &serial.Mode{
		BaudRate: address.Option.BaudRate,
		Parity:   runtime.ParityToParity[address.Option.Parity],
		DataBits: address.Option.DataBits,
		StopBits: runtime.StopBitsToStopBits[address.Option.StopBits],
	}
	port, err := serial.Open(address.Location, mode)
	if err != nil {
		klog.V(2).InfoS("Failed to connect serial port", "address", address.Location)
		return nil, err
	}

	cs := list.New()
	cs.PushBack(&runtime.SerialClient{
		Timeout: 1,
		Port:    port,
	})

	clients := &runtime.Clients{
		Messengers:   cs,
		Max:          1,
		Idle:         1,
		Mux:          &sync.Mutex{},
		NextRequest:  1,
		ConnRequests: make(map[uint64]chan runtime.Messenger, 0),
		NewMessenger: func() (runtime.Messenger, error) {
			newPort, err := serial.Open(address.Location, mode)
			if err != nil {
				klog.V(2).InfoS("Failed to connect serial port", "address", address.Location)
				return nil, err
			}
			return &runtime.SerialClient{
				Timeout: 1,
				Port:    newPort,
			}, nil
		},
	}
	return clients, nil
}

func (m *ModbusRtu) GenerateReadMessage(slave uint, functionCode uint8, startAddress uint, maxDataSize uint, variables []*runtime.VariableParse, memoryLayout common.MemoryLayout) *runtime.ModBusDataFrame {
	// 01 03 00 00 00 0A C5 CD
	// 01  设备地址
	// 03  功能码
	// 00 00  起始地址
	// 00 0A  寄存器数量(word数量)/线圈数量
	// C5 CD  crc16检验码
	message := make([]byte, 6)
	message[0] = byte(slave)
	message[1] = functionCode
	binutils.WriteUint16BigEndian(message[2:], uint16(startAddress))
	binutils.WriteUint16BigEndian(message[4:], uint16(maxDataSize))
	crc16 := make([]byte, 2)
	binutils.WriteUint16BigEndian(crc16, utils.CheckCrc16sum(message))
	message = append(message, crc16...)

	bytesLength := 0
	switch runtime.FunctionCode(functionCode) {
	case runtime.ReadCoilStatus, runtime.ReadInputStatus:
		if maxDataSize%8 == 0 {
			bytesLength = int(maxDataSize/8 + RtuNonDataLength)
		} else {
			bytesLength = int(maxDataSize/8 + 1 + RtuNonDataLength)
		}
	case runtime.ReadHoldRegister, runtime.ReadInputRegister:
		bytesLength = int(maxDataSize*2 + RtuNonDataLength)
	}

	df := &runtime.ModBusDataFrame{
		Slave:             slave,
		MemoryLayout:      memoryLayout,
		StartAddress:      startAddress,
		FunctionCode:      functionCode,
		MaxDataSize:       maxDataSize,
		TransactionId:     0,
		DataFrame:         message,
		ResponseDataFrame: make([]byte, bytesLength),
		Variables:         make([]*runtime.VariableParse, 0, len(variables)),
	}
	df.Variables = append(df.Variables, variables...)

	return df
}
