package collector

import (
	"context"
	"fmt"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/common"
	"time"
)

var (
	ErrNotObject = fmt.Errorf("object does not implement the Object interfaces")
)

type CollectStatus byte
type DeviceStatusCh byte

const (
	Collecting CollectStatus = iota
	CollectingError
	Unconnected
	EmptyVariable
	Stopped
	Error
)

var CollectStatusToString = map[CollectStatus]string{
	Collecting:      "collecting",
	CollectingError: "collectingError",
	Unconnected:     "unconnected",
	EmptyVariable:   "emptyVariable",
	Stopped:         "stopped",
	Error:           "error",
}
var StringToCollectStatus = map[string]CollectStatus{
	"collecting":      Collecting,      // 采集中
	"collectingError": CollectingError, // 采集错误
	"unconnected":     Unconnected,     // 未连接
	"emptyVariable":   EmptyVariable,   // 变量为空
	"stopped":         Stopped,         // 停止
	"error":           Error,           // 错误
}

const (
	Restart DeviceStatusCh = iota
	Start
	Stop
)

var DeviceStatusChToString = map[DeviceStatusCh]string{
	Restart: "restart",
	Start:   "start",
	Stop:    "stop",
}
var StringToDeviceStatusCh = map[string]DeviceStatusCh{
	"restart": Restart,
	"start":   Start,
	"stop":    Stop,
}

type Broker interface {
	Collect(ctx context.Context)
	Destroy(ctx context.Context)
	DeliverAction(ctx context.Context, obj map[string]interface{}) error
}

type VariableValue interface {
	SetValue(value interface{})
	GetValue() interface{}
	GetVariableName() string
	SetVariableName(name string)
	GetVariableAccessMode() common.AccessMode
}

type Object interface {
	GetName() string
	SetName(string)
	GetID() string
	SetID(string)
	GetVersion() string
	SetVersion(string)
	GetModTime() time.Time
	SetModTime(time.Time)
}

type GetVariabler interface {
	GetVariable(key string) (VariableValue, bool)
}
type IndexDevice interface {
	IndexDevice()
}

type GetVariables interface {
	GetVariables() []VariableValue
}

type Device interface {
	Object
	GetVariabler
	IndexDevice
	GetVariables
	GetDeviceCode() string
	SetDeviceCode(string)
	GetDeviceType() string
	SetDeviceType(string)
	GetDeviceModel() string
	SetDeviceModel(string)
	GetCollectStatus() string
	SetCollectStatus(string)
}

type Covert interface {
	Convert(agents *biz.Agents) Device
}

var _ Device = (*DeviceMeta)(nil)

type DeviceMeta struct {
	ObjectMeta
	PublishMeta
	DeviceCode    string `json:"deviceCode"`
	DeviceType    string `json:"deviceType"`
	DeviceModel   string `json:"deviceModel"`
	CollectStatus string `json:"collectStatus"`
	// VariablesMap  map[string]VariableValue `json:"-"`
}

func (d *DeviceMeta) GetVariables() []VariableValue {
	return nil
}

func (d *DeviceMeta) IndexDevice() {
}

func (d *DeviceMeta) GetVariable(key string) (rv VariableValue, exist bool) {
	return
}

func (d *DeviceMeta) GetDeviceCode() string {
	return d.DeviceCode
}

func (d *DeviceMeta) SetDeviceCode(s string) {
	d.DeviceCode = s
}

func (d *DeviceMeta) GetDeviceType() string {
	return d.DeviceType
}

func (d *DeviceMeta) SetDeviceType(s string) {
	d.DeviceType = s
}

func (d *DeviceMeta) GetCollectStatus() string {
	return d.CollectStatus
}

func (d *DeviceMeta) SetCollectStatus(collectStatus string) {
	d.CollectStatus = collectStatus
}

func (d *DeviceMeta) GetDeviceModel() string {
	return d.DeviceModel
}

func (d *DeviceMeta) SetDeviceModel(model string) {
	d.DeviceModel = model
}

type ObjectMeta struct {
	Name    string    `json:"name"`
	ID      string    `json:"id"`
	Version string    `json:"eTag"`
	ModTime time.Time `json:"modTime"`
}

func (meta *ObjectMeta) GetName() string              { return meta.Name }
func (meta *ObjectMeta) SetName(name string)          { meta.Name = name }
func (meta *ObjectMeta) GetID() string                { return meta.ID }
func (meta *ObjectMeta) SetID(id string)              { meta.ID = id }
func (meta *ObjectMeta) GetVersion() string           { return meta.Version }
func (meta *ObjectMeta) SetVersion(version string)    { meta.Version = version }
func (meta *ObjectMeta) GetModTime() time.Time        { return meta.ModTime }
func (meta *ObjectMeta) SetModTime(modTime time.Time) { meta.ModTime = modTime }

type PublishMeta struct {
	Topic string `json:"topic,omitempty"`
}

func (pm *PublishMeta) GetTopic() string      { return pm.Topic }
func (pm *PublishMeta) SetTopic(topic string) { pm.Topic = topic }
