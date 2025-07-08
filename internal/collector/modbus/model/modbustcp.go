package model

import (
	"container/list"
	"fmt"
	"harnsplatform/internal/collector/modbus/runtime"
	"harnsplatform/internal/common"
	"harnsplatform/internal/utils/binutils"
	"k8s.io/klog/v2"
	"net"
	"sync"
)

const TcpNonDataLength = 9

type ModbusTcp struct {
}

func (m *ModbusTcp) NewClients(address *runtime.Address, dataFrameCount int) (*runtime.Clients, error) {
	tcpChannel := dataFrameCount/5 + 1
	addr := fmt.Sprintf("%s:%d", address.Location, address.Option.Port)
	cs := list.New()
	for i := 0; i < tcpChannel; i++ {
		tunnel, err := net.Dial("tcp", addr)
		if err != nil {
			klog.V(2).InfoS("Failed to connect Modbus server", "error", err)
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

func (m *ModbusTcp) GenerateReadMessage(slave uint, functionCode uint8, startAddress uint, maxDataSize uint, variables []*runtime.VariableParse, memoryLayout common.MemoryLayout) *runtime.ModBusDataFrame {
	// 00 01 00 00 00 06 18 03 00 02 00 02
	// 00 01  此次通信事务处理标识符，一般每次通信之后将被要求加1以区别不同的通信数据报文
	// 00 00  表示协议标识符，00 00为modbus协议
	// 00 06  数据长度，用来指示接下来数据的长度，单位字节
	// 18  设备地址，用以标识连接在串行线或者网络上的远程服务端的地址。以上七个字节也被称为modbus报文头
	// 03  功能码，此时代码03为读取保持寄存器数据
	// 00 02  起始地址
	// 00 02  寄存器数量(word数量)/线圈数量
	message := make([]byte, 12)

	binutils.WriteUint16BigEndian(message[2:], 0) // 协议版本
	binutils.WriteUint16BigEndian(message[4:], 6) // 剩余长度
	message[6] = byte(slave)
	message[7] = functionCode
	binutils.WriteUint16BigEndian(message[8:], uint16(startAddress))
	binutils.WriteUint16BigEndian(message[10:], uint16(maxDataSize))

	bytesLength := 0
	switch runtime.FunctionCode(functionCode) {
	case runtime.ReadCoilStatus, runtime.ReadInputStatus:
		if maxDataSize%8 == 0 {
			bytesLength = int(maxDataSize/8 + TcpNonDataLength)
		} else {
			bytesLength = int(maxDataSize/8 + 1 + TcpNonDataLength)
		}
	case runtime.ReadHoldRegister, runtime.ReadInputRegister:
		bytesLength = int(maxDataSize*2 + TcpNonDataLength)
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
