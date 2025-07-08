package model

import (
	"container/list"
	"fmt"
	"harnsplatform/internal/collector/modbus/runtime"
	"harnsplatform/internal/common"
	"harnsplatform/internal/utils"
	"harnsplatform/internal/utils/binutils"
	"k8s.io/klog/v2"
	"net"
	"sync"
)

const RtuOverTcpNonDataLength = 5

type ModbusRtuOverTcp struct {
}

func (m *ModbusRtuOverTcp) NewClients(address *runtime.Address, dataFrameCount int) (*runtime.Clients, error) {
	tcpChannel := dataFrameCount/5 + 1
	addr := fmt.Sprintf("%s:%d", address.Location, address.Option.Port)
	cs := list.New()
	for i := 0; i < tcpChannel; i++ {
		tunnel, err := net.Dial("tcp", addr)
		if err != nil {
			klog.V(2).InfoS("Failed to connect modbus server", "error", err)
			return nil, err
		}
		c := &runtime.TcpClient{
			Tunnel:  tunnel,
			Timeout: 1,
		}
		cs.PushBack(c)
	}

	clients := &runtime.Clients{
		Messengers:   cs,
		Max:          tcpChannel,
		Idle:         tcpChannel,
		Mux:          &sync.Mutex{},
		NextRequest:  1,
		ConnRequests: make(map[uint64]chan runtime.Messenger, 0),
		NewMessenger: func() (runtime.Messenger, error) {
			tunnel, err := net.Dial("tcp", addr)
			if err != nil {
				klog.V(2).InfoS("Failed to connect modbus server", "error", err)
				return nil, err
			}
			return &runtime.TcpClient{
				Tunnel:  tunnel,
				Timeout: 1,
			}, nil
		},
	}
	return clients, nil
}

func (m *ModbusRtuOverTcp) GenerateReadMessage(slave uint, functionCode uint8, startAddress uint, maxDataSize uint, variables []*runtime.VariableParse, memoryLayout common.MemoryLayout) *runtime.ModBusDataFrame {
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
			bytesLength = int(maxDataSize/8 + RtuOverTcpNonDataLength)
		} else {
			bytesLength = int(maxDataSize/8 + 1 + RtuOverTcpNonDataLength)
		}
	case runtime.ReadHoldRegister, runtime.ReadInputRegister:
		bytesLength = int(maxDataSize*2 + RtuOverTcpNonDataLength)
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
