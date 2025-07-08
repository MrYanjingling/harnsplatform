package runtime

import (
	"fmt"
	"harnsplatform/internal/collector"
	"harnsplatform/internal/common"
	"harnsplatform/internal/utils/binutils"
	"strconv"
)

var _ collector.Device = (*ModBusDevice)(nil)
var _ collector.VariableValue = (*Variable)(nil)

type Variable struct {
	DataType     common.DataType   `json:"dataType"`     // bool、int16、float32、float64、int32、int64、uint16
	Name         string            `json:"name"`         // 变量名称
	Address      uint              `json:"address"`      // 变量地址
	FunctionCode uint8             `json:"functionCode"` // 功能码 1、2、3、4
	Bits         uint8             `json:"bits"`
	Amount       uint              `json:"amount"`                 // 数量
	Rate         float64           `json:"rate"`                   // 比率
	OffSet       float64           `json:"offset"`                 // 比率
	DefaultValue interface{}       `json:"defaultValue,omitempty"` // 默认值
	Value        interface{}       `json:"value,omitempty"`        // 值
	AccessMode   common.AccessMode `json:"accessMode"`             // 读写属性
}

func (v *Variable) GetVariableAccessMode() common.AccessMode {
	return v.AccessMode
}

func (v *Variable) SetValue(value interface{}) {
	v.Value = value
}

func (v *Variable) GetValue() interface{} {
	return v.Value
}

func (v *Variable) GetVariableName() string {
	return v.Name
}

func (v *Variable) SetVariableName(name string) {
	v.Name = name
}

type ModBusDevice struct {
	collector.DeviceMeta
	CollectorCycle   uint                 `json:"collectorCycle"`                    // 采集周期
	VariableInterval uint                 `json:"variableInterval"`                  // 变量间隔
	Address          *Address             `json:"address"`                           // IP地址\串口地址
	Slave            uint                 `json:"slave"`                             // 下位机号
	MemoryLayout     common.MemoryLayout  `json:"memoryLayout"`                      // 内存布局 DCBA CDAB BADC ABCD
	PositionAddress  uint                 `json:"positionAddress"`                   // 起始地址
	Variables        []*Variable          `json:"variables" binding:"required,dive"` // 自定义变量
	VariablesMap     map[string]*Variable `json:"-"`                                 // 自定义变量Map
}

func (m *ModBusDevice) IndexDevice() {
	m.VariablesMap = make(map[string]*Variable)
	for _, variable := range m.Variables {
		m.VariablesMap[variable.Name] = variable
	}
}

func (m *ModBusDevice) GetVariable(key string) (rv collector.VariableValue, exist bool) {
	if v, isExist := m.VariablesMap[key]; isExist {
		rv = v
		exist = isExist
	}
	return
}

func (m *ModBusDevice) GetVariables() []collector.VariableValue {
	rvs := make([]collector.VariableValue, 0)

	for _, variable := range m.Variables {
		rvs = append(rvs, variable)
	}

	return rvs
}

type Address struct {
	Location string  `json:"location"` // 地址路径
	Option   *Option `json:"option"`   // 地址其他参数
}

type Option struct {
	Port     int             `json:"port,omitempty"`     // 端口号
	BaudRate int             `json:"baudRate,omitempty"` // 波特率
	DataBits int             `json:"dataBits,omitempty"` // 数据位
	Parity   common.Parity   `json:"parity,omitempty"`   // 校验位
	StopBits common.StopBits `json:"stopBits,omitempty"` // 停止位
}

type VariableSlice []*Variable

func (vs VariableSlice) Len() int {
	return len(vs)
}

func (vs VariableSlice) Less(i, j int) bool {
	return vs[i].Address < vs[j].Address
}

func (vs VariableSlice) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

type VariableParse struct {
	Variable *Variable
	Start    uint // 报文中数据[]byte开始位置
}

type ModBusDataFrame struct {
	Slave             uint
	MemoryLayout      common.MemoryLayout
	StartAddress      uint
	FunctionCode      uint8
	MaxDataSize       uint // 最大数量01 代表线圈  03代表word
	TransactionId     uint16
	DataFrame         []byte
	ResponseDataFrame []byte
	Variables         []*VariableParse
}

func (df *ModBusDataFrame) WriteTransactionId() {
	df.TransactionId++
	id := df.TransactionId
	binutils.WriteUint16BigEndian(df.DataFrame, id)
}

func (df *ModBusDataFrame) ParseVariableValue(data []byte) []collector.VariableValue {
	vvs := make([]collector.VariableValue, 0, len(df.Variables))
	for _, vp := range df.Variables {
		var value interface{}
		switch FunctionCode(df.FunctionCode) {
		case ReadInputStatus, ReadCoilStatus:
			switch vp.Variable.DataType {
			case common.BOOL:
				v := int(data[vp.Start])
				value = v == 1
			case common.INT16:
				value = int16(data[vp.Start])
			case common.UINT16:
				value = uint16(data[vp.Start])
			case common.INT32:
				value = int32(data[vp.Start])
			case common.INT64:
				value = int64(data[vp.Start])
			case common.FLOAT32:
				value = float32(data[vp.Start])
			case common.FLOAT64:
				value = float64(data[vp.Start])
			}
		case ReadInputRegister, ReadHoldRegister:
			vpData := data[vp.Start:]
			switch vp.Variable.DataType {
			case common.BOOL:
				var v int16
				switch df.MemoryLayout {
				case common.ABCD, common.CDAB:
					v = int16(binutils.ParseUint16BigEndian(vpData))
				case common.BADC, common.DCBA:
					v = int16(binutils.ParseUint16LittleEndian(vpData))
				}
				value = v != 0
			case common.INT16:
				var v interface{}
				switch df.MemoryLayout {
				case common.ABCD, common.CDAB:
					v = int16(binutils.ParseUint16BigEndian(vpData))
				case common.BADC, common.DCBA:
					v = int16(binutils.ParseUint16LittleEndian(vpData))
				}
				if vp.Variable.Rate != 0 && vp.Variable.Rate != 1 {
					value = int16((v.(float64)) * vp.Variable.Rate)
				} else {
					value = v
				}
			case common.UINT16:
				var v interface{}
				switch df.MemoryLayout {
				case common.ABCD, common.CDAB:
					v = binutils.ParseUint16BigEndian(vpData)
				case common.BADC, common.DCBA:
					v = binutils.ParseUint16LittleEndian(vpData)
				}
				if vp.Variable.Rate != 0 && vp.Variable.Rate != 1 {
					value = uint16((v.(float64)) * vp.Variable.Rate)
				} else {
					value = v
				}
			case common.INT32:
				var v interface{}
				switch df.MemoryLayout {
				case common.ABCD:
					v = int32(binutils.ParseUint32BigEndian(vpData))
				case common.BADC:
					// 大端交换
					v = int32(binutils.ParseUint32BigEndianByteSwap(vpData))
				case common.CDAB:
					v = int32(binutils.ParseUint32LittleEndianByteSwap(vpData))
				case common.DCBA:
					v = int32(binutils.ParseUint32LittleEndian(vpData))
				}
				if vp.Variable.Rate != 0 && vp.Variable.Rate != 1 {
					value = int32((v.(float64)) * vp.Variable.Rate)
				} else {
					value = v
				}
			case common.INT64:
				var v interface{}
				switch df.MemoryLayout {
				case common.ABCD:
					v = int64(binutils.ParseUint64BigEndian(vpData))
				case common.BADC:
					v = int64(binutils.ParseUint64BigEndianByteSwap(vpData))
				case common.CDAB:
					v = int64(binutils.ParseUint64LittleEndianByteSwap(vpData))
				case common.DCBA:
					v = int64(binutils.ParseUint64LittleEndian(vpData))
				}
				if vp.Variable.Rate != 0 && vp.Variable.Rate != 1 {
					value = int64((v.(float64)) * vp.Variable.Rate)
				} else {
					value = v
				}
			case common.FLOAT32:
				var v interface{}
				switch df.MemoryLayout {
				case common.ABCD:
					v = binutils.ParseFloat32BigEndian(vpData)
				case common.BADC:
					v = binutils.ParseFloat32BigEndianByteSwap(vpData)
				case common.CDAB:
					v = binutils.ParseFloat32LittleEndianByteSwap(vpData)
				case common.DCBA:
					v = binutils.ParseFloat32LittleEndian(vpData)
				}
				if vp.Variable.Rate != 0 && vp.Variable.Rate != 1 {
					value = v.(float32) * float32(vp.Variable.Rate)
				} else {
					value = v
				}

				if value != float32(int(value.(float32))) {
					value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
				}

			case common.FLOAT64:
				var v interface{}
				switch df.MemoryLayout {
				case common.ABCD:
					v = binutils.ParseFloat64BigEndian(vpData)
				case common.BADC:
					v = binutils.ParseFloat64BigEndianByteSwap(vpData)
				case common.CDAB:
					v = binutils.ParseFloat64LittleEndianByteSwap(vpData)
				case common.DCBA:
					v = binutils.ParseFloat64LittleEndian(vpData)
				}
				if vp.Variable.Rate != 0 && vp.Variable.Rate != 1 {
					value = (v.(float64)) * vp.Variable.Rate
				} else {
					value = v
				}

				if value != float64(int(value.(float64))) {
					value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
				}
			}
		}

		vp.Variable.SetValue(value)
		vvs = append(vvs, &Variable{
			DataType:     vp.Variable.DataType,
			Name:         vp.Variable.Name,
			Address:      vp.Variable.Address,
			FunctionCode: vp.Variable.FunctionCode,
			Rate:         vp.Variable.Rate,
			DefaultValue: vp.Variable.DefaultValue,
			Value:        vp.Variable.Value,
		})
	}
	return vvs
}
