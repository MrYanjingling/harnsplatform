package modbus

import (
	"context"
	"errors"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/collector"
	"harnsplatform/internal/collector/modbus/model"
	"harnsplatform/internal/collector/modbus/runtime"
	"harnsplatform/internal/common"
	"harnsplatform/internal/utils"
	"harnsplatform/internal/utils/binutils"
	"k8s.io/klog/v2"
	"sort"
	"strconv"
	"sync"
	"time"
)

/**
modbus 协议 ADU = 地址(1) + pdu(253) + 16位校验(2) = 256
modbus tcp报文
tcp报文头(6)  +  地址(1)   +   pdu(253) = 260
modbus rtu报文
地址(1) + pdu(253) + 16位校验(2) = 256
modbus rtu over tcp
tcp报文头(6)  +  地址(1)   +   pdu(253)   +  16位校验(2)  = 262
*/

// ModBusDataFrame 报文对应的数据点位
var _ collector.Broker = (*ModbusBroker)(nil)

type ModbusBroker struct {
	NeedCheckTransaction     bool
	NeedCheckCrc16Sum        bool
	ExitCh                   chan struct{}
	Device                   *runtime.ModBusDevice
	Clients                  *runtime.Clients
	FunctionCodeDataFrameMap map[uint8][]*runtime.ModBusDataFrame
	VariableCount            int
	VariableCh               chan *collector.ParseVariableResult
}

func NewBroker(d collector.Device) (collector.Broker, chan *collector.ParseVariableResult, error) {
	device, ok := d.(*runtime.ModBusDevice)
	if !ok {
		klog.V(2).InfoS("Unsupported device,type not Modbus")
		return nil, nil, collector.ErrDeviceType
	}

	needCheckTransaction := false
	needCheckCrc16Sum := false
	switch runtime.StringToModbusModel[device.DeviceModel] {
	case runtime.Tcp:
		needCheckTransaction = true
	case runtime.Rtu:
		needCheckCrc16Sum = true
	case runtime.RtuOverTcp:
		needCheckCrc16Sum = true
	}

	VariableCount := 0
	functionCodeDataFrameMap := make(map[uint8][]*runtime.ModBusDataFrame, 0)
	functionCodeVariableMap := make(map[uint8][]*runtime.Variable, 0)
	for _, variable := range device.Variables {
		functionCodeVariableMap[variable.FunctionCode] = append(functionCodeVariableMap[variable.FunctionCode], variable)
	}
	for code, variables := range functionCodeVariableMap {
		VariableCount = VariableCount + len(variables)
		sort.Sort(runtime.VariableSlice(variables))
		dfs := make([]*runtime.ModBusDataFrame, 0)
		firstVariable := variables[0]
		startOffset := firstVariable.Address - device.PositionAddress
		startAddress := startOffset
		var maxDataSize uint = 0
		vps := make([]*runtime.VariableParse, 0)
		switch runtime.FunctionCode(code) {
		case runtime.ReadCoilStatus, runtime.ReadInputStatus:
			dataFrameDataLength := startAddress + runtime.PerRequestMaxCoil
			for i := 0; i < len(variables); i++ {
				variable := variables[i]
				if variable.Address <= dataFrameDataLength {
					vp := &runtime.VariableParse{
						Variable: variable,
						Start:    variable.Address - startAddress,
					}
					vps = append(vps, vp)
					maxDataSize = variable.Address - startAddress + 1
				} else {
					df := model.ModbusModelers[device.DeviceModel].GenerateReadMessage(device.Slave, code, startAddress, maxDataSize, vps, device.MemoryLayout)
					dfs = append(dfs, df)
					vps = vps[:0:0]
					maxDataSize = 0
					startAddress = variable.Address
					dataFrameDataLength = startAddress + runtime.PerRequestMaxCoil
					i--
				}
			}
		case runtime.ReadHoldRegister, runtime.ReadInputRegister:
			dataFrameDataLength := startAddress + runtime.PerRequestMaxRegister
			for i := 0; i < len(variables); i++ {
				variable := variables[i]
				if variable.Address+common.DataTypeWord[variable.DataType] <= dataFrameDataLength {
					vp := &runtime.VariableParse{
						Variable: variable,
						Start:    (variable.Address - startAddress) * 2,
					}
					vps = append(vps, vp)
					maxDataSize = variable.Address - startAddress + common.DataTypeWord[variable.DataType]
				} else {
					df := model.ModbusModelers[device.DeviceModel].GenerateReadMessage(device.Slave, code, startAddress, maxDataSize, vps, device.MemoryLayout)
					dfs = append(dfs, df)
					vps = vps[:0:0]
					maxDataSize = 0
					startAddress = variable.Address
					dataFrameDataLength = startAddress + runtime.PerRequestMaxRegister
					i--
				}
			}
		}
		if len(vps) > 0 {
			df := model.ModbusModelers[device.DeviceModel].GenerateReadMessage(device.Slave, code, startAddress, maxDataSize, vps, device.MemoryLayout)
			dfs = append(dfs, df)
			vps = vps[:0:0]
		}
		functionCodeDataFrameMap[code] = append(functionCodeDataFrameMap[code], dfs...)
	}

	dataFrameCount := 0
	for _, values := range functionCodeDataFrameMap {
		dataFrameCount += len(values)
	}
	if dataFrameCount == 0 {
		klog.V(2).InfoS("Unnecessary to collect from Modbus device.Because of the variables is empty", "deviceId", device.ID)
		return nil, nil, collector.ErrDeviceEmptyVariable
	}

	clients, err := model.ModbusModelers[device.DeviceModel].NewClients(device.Address, dataFrameCount)
	if err != nil {
		klog.V(2).InfoS("Failed to connect Modbus device", "error", err, "deviceId", device.ID)
		return nil, nil, collector.ErrConnectDevice
	}

	mtc := &ModbusBroker{
		Device:                   device,
		ExitCh:                   make(chan struct{}, 0),
		FunctionCodeDataFrameMap: functionCodeDataFrameMap,
		Clients:                  clients,
		VariableCh:               make(chan *collector.ParseVariableResult, 1),
		VariableCount:            VariableCount,
		NeedCheckCrc16Sum:        needCheckCrc16Sum,
		NeedCheckTransaction:     needCheckTransaction,
	}
	return mtc, mtc.VariableCh, nil
}

func (broker *ModbusBroker) Destroy(ctx context.Context) {
	broker.ExitCh <- struct{}{}
	broker.Clients.Destroy(ctx)
	close(broker.VariableCh)
}

func (broker *ModbusBroker) Collect(ctx context.Context) {
	go func() {
		for {
			start := time.Now().Unix()
			if !broker.poll(ctx) {
				return
			}
			select {
			case <-broker.ExitCh:
				return
			default:
				end := time.Now().Unix()
				elapsed := end - start
				if elapsed < int64(broker.Device.CollectorCycle) {
					time.Sleep(time.Duration(int64(broker.Device.CollectorCycle)) * time.Second)
				}
			}
		}
	}()
}

func (broker *ModbusBroker) DeliverAction(ctx context.Context, obj map[string]interface{}) error {
	action := make([]*runtime.Variable, 0, len(obj))

	for name, value := range obj {
		vv, _ := broker.Device.GetVariable(name)
		variableValue := vv.(*runtime.Variable)

		v := &runtime.Variable{
			DataType:     variableValue.DataType,
			Name:         variableValue.Name,
			Address:      variableValue.Address,
			Bits:         variableValue.Bits,
			FunctionCode: variableValue.FunctionCode,
			Rate:         variableValue.Rate,
			Amount:       variableValue.Amount,
			AccessMode:   variableValue.AccessMode,
		}
		switch variableValue.DataType {
		case common.BOOL:
			switch value.(type) {
			case bool:
				v.Value = value
			case string:
				b, err := strconv.ParseBool(value.(string))
				if err == nil {
					v.Value = b
				} else {
					return nil
					// return response.ErrBooleanInvalid(name)
				}
			default:
				return nil
				// return response.ErrBooleanInvalid(name)
			}
		case common.INT16:
			switch value.(type) {
			case float64:
				v.Value = int16(value.(float64))
			default:
				return nil
				// return response.ErrInteger16Invalid(name)
			}
		case common.UINT16:
			switch value.(type) {
			case float64:
				v.Value = uint16(value.(float64))
			default:
				return nil
				// return response.ErrInteger16Invalid(name)
			}
		case common.INT32:
			switch value.(type) {
			case float64:
				v.Value = int32(value.(float64))
			default:
				return nil
				// return response.ErrInteger32Invalid(name)
			}
		case common.INT64:
			switch value.(type) {
			case float64:
				v.Value = int64(value.(float64))
			default:
				return nil
				// return response.ErrInteger64Invalid(name)
			}
		case common.FLOAT32:
			switch value.(type) {
			case float64:
				v.Value = float32(value.(float64))
			default:
				return nil
				// return response.ErrFloat32Invalid(name)
			}
		case common.FLOAT64:
			switch value.(type) {
			case float64:
				v.Value = value.(float64)
			default:
				return nil
				// return response.ErrFloat64Invalid(name)
			}
		default:
			klog.V(3).InfoS("Unsupported dataType", "dataType", variableValue.DataType)
		}
		action = append(action, v)
	}

	dataBytes := broker.generateActionBytes(broker.Device.MemoryLayout, action)
	dataFrames := make([][]byte, 0, len(dataBytes))
	for i, dbs := range dataBytes {
		var bytes []byte
		if broker.NeedCheckTransaction {
			bytes = append(bytes, make([]byte, 6)...)
			binutils.WriteUint16BigEndian(bytes[0:], uint16(i))
			binutils.WriteUint16BigEndian(bytes[2:], 0)
			binutils.WriteUint16BigEndian(bytes[4:], uint16(1+len(dbs)))
		}
		bytes = append(bytes, byte(broker.Device.Slave))
		bytes = append(bytes, dbs...)
		if broker.NeedCheckCrc16Sum {
			crc16 := make([]byte, 2)
			binutils.WriteUint16BigEndian(crc16, utils.CheckCrc16sum(bytes))
			bytes = append(bytes, crc16...)
		}
		dataFrames = append(dataFrames, bytes)
	}

	messenger, err := broker.Clients.GetMessenger(ctx)
	if err != nil {
		klog.V(2).InfoS("Failed to get Modbus messenger", "error", err)
		if messenger, err = broker.Clients.NewMessenger(); err != nil {
			return err
		}
	}
	defer broker.Clients.ReleaseMessenger(messenger)

	// errs := &response.MultiError{}
	for _, frame := range dataFrames {
		rp := make([]byte, len(frame))
		_, err = messenger.AskAtLeast(frame, rp, 6)
		if err != nil {
			// errs.Add(modbus.ErrModbusBadConn)
			continue
		}
		if broker.NeedCheckTransaction {
			transactionId := binutils.ParseUint16(rp[:])
			requestTransactionId := binutils.ParseUint16(rp[:])
			if transactionId != requestTransactionId {
				klog.V(2).InfoS("Failed to match Modbus message transaction id", "request transactionId", requestTransactionId, "response transactionId", transactionId)
				// errs.Add(modbus.ErrMessageTransaction)
				continue
			}
			rp = rp[6:]
		}

		slave := rp[0]
		if uint(slave) != broker.Device.Slave {
			klog.V(2).InfoS("Failed to match Modbus slave", "request slave", broker.Device.Slave, "response slave", slave)
			// errs.Add(modbus.ErrMessageSlave)
			continue
		}
		functionCode := rp[1]
		if functionCode&0x80 > 0 {
			klog.V(2).InfoS("Failed to parse Modbus message", "error code", functionCode-128)
			// errs.Add(modbus.ErrMessageFunctionCodeError)
			continue
		}
	}
	//
	return nil

	// if errs.Len() > 0 {
	// 	return errs
	// }

	return nil
}

func (broker *ModbusBroker) poll(ctx context.Context) bool {
	select {
	case <-broker.ExitCh:
		return false
	default:
		sw := &sync.WaitGroup{}
		dfvCh := make(chan *collector.ParseVariableResult, 0)
		for _, DataFrames := range broker.FunctionCodeDataFrameMap {
			for _, frame := range DataFrames {
				sw.Add(1)
				go broker.message(ctx, frame, dfvCh, sw, broker.Clients)
			}
		}
		go broker.rollVariable(ctx, dfvCh)
		sw.Wait()
		close(dfvCh)
		return true
	}
}

func (broker *ModbusBroker) message(ctx context.Context, dataFrame *runtime.ModBusDataFrame, pvrCh chan<- *collector.ParseVariableResult, sw *sync.WaitGroup, clients *runtime.Clients) {
	defer sw.Done()
	defer func() {
		if err := recover(); err != nil {
			klog.V(2).InfoS("Failed to ask Modbus server message", "error", err)
		}
	}()
	messenger, err := clients.GetMessenger(ctx)
	defer broker.Clients.ReleaseMessenger(messenger)
	if err != nil {
		klog.V(2).InfoS("Failed to get messenger", "error", err)
		if messenger, err = broker.Clients.NewMessenger(); err != nil {
			return
		}
	}

	var buf []byte

	if err := broker.retry(func(messenger runtime.Messenger, dataFrame *runtime.ModBusDataFrame) error {
		if broker.NeedCheckTransaction {
			dataFrame.WriteTransactionId()
		}
		_, err := messenger.AskAtLeast(dataFrame.DataFrame, dataFrame.ResponseDataFrame, 6)
		if err != nil {
			return runtime.ErrModbusBadConn
		}
		buf, err = broker.ValidateAndExtractMessage(dataFrame)
		if err != nil {
			return runtime.ErrModbusServerBadResp
		}
		return nil
	}, messenger, dataFrame); err != nil {
		klog.V(2).InfoS("Failed to connect modbus server", "error", err)
		pvrCh <- &collector.ParseVariableResult{Err: []error{err}}
		return
	}

	pvrCh <- &collector.ParseVariableResult{Err: nil, VariableSlice: dataFrame.ParseVariableValue(buf)}
}

func (broker *ModbusBroker) retry(fun func(messenger runtime.Messenger, dataFrame *runtime.ModBusDataFrame) error, messenger runtime.Messenger, dataFrame *runtime.ModBusDataFrame) error {
	for i := 0; i < 3; i++ {
		err := fun(messenger, dataFrame)
		if err == nil {
			return nil
		} else if errors.Is(err, runtime.ErrModbusBadConn) {
			messenger.Close()
			newMessenger, err := broker.Clients.NewMessenger()
			if err != nil {
				return err
			}
			messenger.Reset(newMessenger)
		} else {
			klog.V(2).InfoS("Failed to connect Modbus server", "error", err)
		}
	}
	return runtime.ErrManyRetry
}

func (broker *ModbusBroker) ValidateAndExtractMessage(df *runtime.ModBusDataFrame) ([]byte, error) {
	buf := df.ResponseDataFrame[:]

	if broker.NeedCheckTransaction {
		transactionId := binutils.ParseUint16(buf[:])
		if transactionId != df.TransactionId {
			klog.V(2).InfoS("Failed to match message transaction id", "request transactionId", df.TransactionId, "response transactionId", transactionId)
			return nil, runtime.ErrMessageTransaction
		}
		buf = buf[6:]
	}

	slave := buf[0]
	if uint(slave) != df.Slave {
		klog.V(2).InfoS("Failed to match modbus slave", "request slave", df.Slave, "response slave", slave)
		return nil, runtime.ErrMessageSlave
	}
	functionCode := buf[1]
	if functionCode&0x80 > 0 {
		klog.V(2).InfoS("Failed to parse modbus tcp message", "error code", functionCode-128)
		return nil, runtime.ErrMessageFunctionCodeError
	}

	byteDataLength := buf[2]
	if broker.NeedCheckCrc16Sum {
		if int(byteDataLength)+5 != len(buf) {
			klog.V(2).InfoS("Failed to get message enough length")
			return nil, runtime.ErrMessageDataLengthNotEnough
		}
		checkBufData := buf[:byteDataLength+3]
		sum := utils.CheckCrc16sum(checkBufData)
		crc := binutils.ParseUint16BigEndian(buf[byteDataLength+3 : byteDataLength+5])
		if sum != crc {
			klog.V(2).InfoS("Failed to check CRC16")
			return nil, runtime.ErrCRC16Error
		}
	} else {
		if int(byteDataLength)+3 != len(buf) {
			klog.V(2).InfoS("Failed to get message enough length")
			return nil, runtime.ErrMessageDataLengthNotEnough
		}
	}

	var bb []byte
	switch runtime.FunctionCode(buf[1]) {
	case runtime.ReadCoilStatus, runtime.ReadInputStatus:
		// 数组解压
		bb = binutils.ExpandBool(buf[3:], int(byteDataLength))
	case runtime.ReadHoldRegister, runtime.ReadInputRegister:
		bb = binutils.Dup(buf[3:])
	case runtime.WriteSingleCoil, runtime.WriteSingleRegister, runtime.WriteMultipleCoil, runtime.WriteMultipleRegister:
	default:
		klog.V(2).InfoS("Unsupported function code", "functionCode", buf[1])
	}

	return bb, nil
}

func (broker *ModbusBroker) rollVariable(ctx context.Context, ch chan *collector.ParseVariableResult) {
	rvs := make([]collector.VariableValue, 0, broker.VariableCount)
	errs := make([]error, 0)
	for {
		select {
		case pvr, ok := <-ch:
			if !ok {
				broker.VariableCh <- &collector.ParseVariableResult{Err: errs, VariableSlice: rvs}
				return
			} else if pvr.Err != nil {
				errs = append(errs, pvr.Err...)
			} else {
				for _, variable := range pvr.VariableSlice {
					rvs = append(rvs, variable)
				}
			}
		}
	}
}

func (broker *ModbusBroker) generateActionBytes(memoryLayout common.MemoryLayout, action []*runtime.Variable) [][]byte {
	dataBytes := make([][]byte, 0, len(action))

	for _, variable := range action {
		// functioncode + startAddress
		pduByte := make([]byte, 3)
		var fc byte
		dataByte := make([]byte, 0)
		switch runtime.FunctionCode(variable.FunctionCode) {
		case runtime.ReadCoilStatus, runtime.ReadInputStatus:
			switch variable.DataType {
			// 65280
			case common.BOOL:
				fc = byte(runtime.WriteSingleCoil)
				if variable.Value.(bool) {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(65280))...)
				} else {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(0))...)
				}
			case common.INT16:
				fc = byte(runtime.WriteSingleCoil)
				if variable.Value.(int16) > 0 {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(65280))...)
				} else {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(0))...)
				}
			case common.UINT16:
				fc = byte(runtime.WriteSingleCoil)
				if variable.Value.(uint16) > 0 {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(65280))...)
				} else {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(0))...)
				}
			case common.INT32:
				fc = byte(runtime.WriteSingleCoil)
				if variable.Value.(int32) > 0 {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(65280))...)
				} else {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(0))...)
				}
			case common.INT64:
				fc = byte(runtime.WriteSingleCoil)
				if variable.Value.(int64) > 0 {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(65280))...)
				} else {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(0))...)
				}
			case common.FLOAT32:
				fc = byte(runtime.WriteSingleCoil)
				if variable.Value.(float32) > 0 {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(65280))...)
				} else {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(0))...)
				}
			case common.FLOAT64:
				fc = byte(runtime.WriteSingleCoil)
				if variable.Value.(float64) > 0 {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(65280))...)
				} else {
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(0))...)
				}
			}
		case runtime.ReadHoldRegister, runtime.ReadInputRegister:
			switch variable.DataType {
			case common.BOOL:
				klog.V(2).InfoS("Unsupported bool variable with read hold register", "variableName", variable.Name)
				// todo 需要先读取数据
			case common.INT16:
				fc = byte(runtime.WriteSingleRegister)

				var value int16
				if variable.Rate != 0 && variable.Rate != 1 {
					value = int16((variable.Value.(float64)) * variable.Rate)
				} else {
					value = variable.Value.(int16)
				}
				switch memoryLayout {
				case common.ABCD, common.CDAB:
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(value))...)
				case common.BADC, common.DCBA:
					dataByte = append(dataByte, binutils.Uint16ToBytesLittleEndian(uint16(value))...)
				}
			case common.UINT16:
				fc = byte(runtime.WriteSingleRegister)

				var value uint16
				if variable.Rate != 0 && variable.Rate != 1 {
					value = uint16((variable.Value.(float64)) * variable.Rate)
				} else {
					value = variable.Value.(uint16)
				}
				switch memoryLayout {
				case common.ABCD, common.CDAB:
					dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(value)...)
				case common.BADC, common.DCBA:
					dataByte = append(dataByte, binutils.Uint16ToBytesLittleEndian(value)...)
				}
			case common.INT32:
				fc = byte(runtime.WriteMultipleRegister)
				registerAmount := 2
				dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(registerAmount))...)
				dataByte = append(dataByte, byte(2*registerAmount))

				var value int32
				if variable.Rate != 0 && variable.Rate != 1 {
					value = int32((variable.Value.(float64)) * variable.Rate)
				} else {
					value = variable.Value.(int32)
				}
				switch memoryLayout {
				case common.ABCD:
					dataByte = append(dataByte, binutils.Uint32ToBytesBigEndian(uint32(value))...)
				case common.BADC:
					// 大端交换
					dataByte = append(dataByte, binutils.Uint32ToBytesBigEndianByteSwap(uint32(value))...)
				case common.CDAB:
					dataByte = append(dataByte, binutils.Uint32ToBytesLittleEndianByteSwap(uint32(value))...)
				case common.DCBA:
					dataByte = append(dataByte, binutils.Uint32ToBytesLittleEndian(uint32(value))...)
				}
			case common.INT64:
				fc = byte(runtime.WriteMultipleRegister)
				registerAmount := 4
				dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(registerAmount))...)
				dataByte = append(dataByte, byte(2*registerAmount))

				var value int64
				if variable.Rate != 0 && variable.Rate != 1 {
					value = int64((variable.Value.(float64)) * variable.Rate)
				} else {
					value = variable.Value.(int64)
				}
				switch memoryLayout {
				case common.ABCD:
					dataByte = append(dataByte, binutils.Uint64ToBytesBigEndian(uint64(value))...)
				case common.BADC:
					// 大端交换
					dataByte = append(dataByte, binutils.Uint64ToBytesBigEndianByteSwap(uint64(value))...)
				case common.CDAB:
					dataByte = append(dataByte, binutils.Uint64ToBytesLittleEndianByteSwap(uint64(value))...)
				case common.DCBA:
					dataByte = append(dataByte, binutils.Uint64ToBytesLittleEndian(uint64(value))...)
				}
			case common.FLOAT32:
				fc = byte(runtime.WriteMultipleRegister)
				registerAmount := 2
				dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(registerAmount))...)
				dataByte = append(dataByte, byte(2*registerAmount))

				var value float32
				if variable.Rate != 0 && variable.Rate != 1 {
					value = variable.Value.(float32) * float32(16000)
				} else {
					value = variable.Value.(float32)
				}
				switch memoryLayout {
				case common.ABCD:
					dataByte = append(dataByte, binutils.Float32ToBytesBigEndian(value)...)
				case common.BADC:
					// 大端交换
					dataByte = append(dataByte, binutils.Float32ToBytesBigEndianByteSwap(value)...)
				case common.CDAB:
					dataByte = append(dataByte, binutils.Float32ToBytesLittleEndianByteSwap(value)...)
				case common.DCBA:
					dataByte = append(dataByte, binutils.Float32ToBytesLittleEndian(value)...)
				}
			case common.FLOAT64:
				fc = byte(runtime.WriteMultipleRegister)
				registerAmount := 4
				dataByte = append(dataByte, binutils.Uint16ToBytesBigEndian(uint16(registerAmount))...)
				dataByte = append(dataByte, byte(2*registerAmount))

				var value float64
				if variable.Rate != 0 && variable.Rate != 1 {
					value = (variable.Value.(float64)) * variable.Rate
				} else {
					value = variable.Value.(float64)
				}
				switch memoryLayout {
				case common.ABCD:
					dataByte = append(dataByte, binutils.Float64ToBytesBigEndian(value)...)
				case common.BADC:
					// 大端交换
					dataByte = append(dataByte, binutils.Float64ToBytesBigEndianByteSwap(value)...)
				case common.CDAB:
					dataByte = append(dataByte, binutils.Float64ToBytesLittleEndianByteSwap(value)...)
				case common.DCBA:
					dataByte = append(dataByte, binutils.Float64ToBytesLittleEndian(value)...)
				}
			}
		}
		pduByte[0] = fc
		binutils.WriteUint16BigEndian(pduByte[1:], uint16(variable.Address))
		pduByte = append(pduByte, dataByte...)

		dataBytes = append(dataBytes, pduByte)
	}
	return dataBytes
}

func ConvertDevice(agents *biz.Agents) collector.Device {
	return nil
}
