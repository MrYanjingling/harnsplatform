package collector

import (
	"context"
	"errors"
	"harnsplatform/internal/biz"
	"k8s.io/klog/v2"
	"os"
	"strings"
	"sync"
	"time"
)

type Option func(*Manager)

type Manager struct {
	mm      *ModelManager
	ts      *TimeSeriesManager
	tsStore bool
	// agents           *sync.Map
	devices          *sync.Map
	heartBeatDevices *sync.Map
	brokers          map[string]Broker
	brokerReturnCh   map[string]chan *ParseVariableResult
	stopCh           <-chan struct{}
	deviceStatusCh   chan string
	mux              *sync.Mutex
}

func NewManager(mm *ModelManager, ts *TimeSeriesManager, tsStore bool, stop <-chan struct{}, opts ...Option) *Manager {
	m := &Manager{
		devices:          &sync.Map{},
		heartBeatDevices: &sync.Map{},
		mm:               mm,
		mux:              &sync.Mutex{},
		ts:               ts,
		tsStore:          tsStore,
		brokers:          make(map[string]Broker, 0),
		brokerReturnCh:   make(map[string]chan *ParseVariableResult, 0),
		stopCh:           stop,
		deviceStatusCh:   make(chan string, 0),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Manager) Init() {
	// devices, _ := m.store.LoadResource()
	m.mm.Init()
	if m.tsStore {
		m.ts.Init()
	}

	// m.agents = m.mm.GetAgents()

	m.mm.GetAgents().Range(func(key, value any) bool {
		agents := value.(*biz.Agents)
		device := ConvertDeviceMap[agents.AgentType](agents)
		device.IndexDevice()
		m.devices.Store(device.GetID(), device)

		if err := m.readyCollect(device); err != nil {
			if errors.Is(err, ErrConnectDevice) {
				// 开启探测协程 15S一次
				m.heartBeatDevices.Store(device.GetID(), device)
			} else {
				klog.V(2).InfoS("Failed to start process collect device data", "deviceId", device.GetID())
			}
		}

		return true
	})

	go m.heartBeatDetection()
	go m.listeningDeviceStatusCh()
}

// func (m *Manager) ListDevices(filter * DeviceFilter, exploded bool) ([] Device, error) {
// 	rds := make([] Device, 0)
// 	predicates :=  ParseTypeFilter(filter)
//
// 	// descend
// 	byModTime := func(d1, d2  Device) bool { return d1.GetModTime().Before(d2.GetModTime()) }
// 	sorter :=  ByDevice(byModTime)
//
// 	m.devices.Range(func(key, value interface{}) bool {
// 		isMatch := true
// 		v := value.( Device)
// 		for _, p := range predicates {
// 			if !p(v) {
// 				isMatch = false
// 				break
// 			}
// 		}
// 		if isMatch {
// 			rds = sorter.Insert(rds, v)
// 		}
// 		return true
// 	})
//
// 	if !exploded {
// 		for i := range rds {
// 			rds[i] = m.foldDevice(rds[i])
// 		}
// 	}
//
// 	return rds, nil
// }

func (m *Manager) GetDeviceById(id string, exploded bool) ( Device, error) {
	d, isExist := m.devices.Load(id)
	if !isExist {
		return nil, os.ErrNotExist
	}
	device, _ := d.( Device)
	return device, nil
}

func (m *Manager) SwitchDeviceStatus(id string, status string) error {
	if _, err := m.GetDeviceById(id, true); err != nil {
		klog.V(2).InfoS("Failed to find device", "deviceId", id)
		return err
	}
	if _, ok :=  StringToDeviceStatusCh[status]; !ok {
		klog.V(2).InfoS("Unsupported device status", "status", status)
		return nil
		// return response.ErrDeviceOperatorUnSupported(status)
	}
	dsc := id + "-" + status
	m.deviceStatusCh <- dsc
	return nil
}

// func (m *Manager) DeliverAction(id string, actions []map[string]interface{}) error {
// 	device, err := m.GetDeviceById(id, true)
// 	if err != nil {
// 		klog.V(2).InfoS("Failed to find device", "deviceId", id)
// 		return response.NewMultiError(response.ErrDeviceNotFound(id))
// 	}
//
// 	errs := &response.MultiError{}
// 	legalActions := make(map[string]interface{}, 0)
// 	for _, item := range actions {
// 		for k, v := range item {
// 			if _, exist := legalActions[k]; exist {
// 				errs.Add(response.ErrResourceExists(k))
// 				continue
// 			}
// 			if v, ok := device.GetVariable(k); !ok {
// 				errs.Add(response.ErrResourceNotFound(k))
// 				continue
// 			} else if v.GetVariableAccessMode() != constant.AccessModeReadWrite {
// 				errs.Add(response.ErrResourceNotFound(k))
// 				continue
// 			}
// 			legalActions[k] = v
// 		}
// 	}
//
// 	if errs.Len() > 0 {
// 		return errs
// 	}
//
// 	if len(legalActions) == 0 {
// 		return response.NewMultiError(response.ErrLegalActionNotFound)
// 	}
//
// 	if device.GetCollectStatus() ==  CollectStatusToString[ Unconnected] {
// 		klog.V(2).InfoS("Failed to connect device", "deviceId", id)
// 		return response.NewMultiError(response.ErrDeviceNotConnect(id))
// 	}
//
// 	return m.brokers[id].DeliverAction(context.Background(), legalActions)
// }

func (m Manager) cancelCollect(obj Device) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	// switch status
	obj.SetCollectStatus(CollectStatusToString[Stopped])
	// delete heartBeat devices if exist
	if _, exist := m.heartBeatDevices.Load(obj.GetID()); exist {
		m.heartBeatDevices.Delete(obj.GetID())
	}
	if v, ok := m.brokers[obj.GetID()]; ok {
		v.Destroy(context.Background())
		delete(m.brokers, obj.GetID())
		delete(m.brokerReturnCh, obj.GetID())
	}
	return nil
}

func (m *Manager) readyCollect(obj Device) error {
	broker, results, err := DeviceTypeBrokerMap[obj.GetDeviceType()](obj)
	if err != nil {
		switch {
		case errors.Is(err, ErrConnectDevice):
			obj.SetCollectStatus(CollectStatusToString[Unconnected])
			return err
		case errors.Is(err, ErrDeviceEmptyVariable):
			obj.SetCollectStatus(CollectStatusToString[EmptyVariable])
			return nil
		default:
			return err
		}
	}
	obj.SetCollectStatus(CollectStatusToString[Collecting])
	klog.V(2).InfoS("Succeed to collect data", "deviceId", obj.GetID())
	m.mux.Lock()
	defer m.mux.Unlock()
	m.brokers[obj.GetID()] = broker
	m.brokerReturnCh[obj.GetID()] = results

	// topic := obj.GetTopic()
	// if len(topic) == 0 {
	// 	topic = fmt.Sprintf("data/%s/v1/%s", m.gatewayMeta.ID, obj.GetID())
	// 	obj.SetTopic(topic)
	// }

	broker.Collect(context.Background())
	go func(deviceId string, ch chan *ParseVariableResult) {
		for {
			select {
			case _, ok := <-m.stopCh:
				if !ok {
					return
				}
			case pvr, ok := <-results:
				if ok {
					if v, ok := m.devices.Load(deviceId); ok {
						if len(pvr.Err) == 0 {
							if v.(Device).GetCollectStatus() !=  CollectStatusToString[ Collecting] {
								v.(Device).SetCollectStatus( CollectStatusToString[ Collecting])
							}
							// pds := make([]PointData, 0, len(pvr.VariableSlice))
							// for _, value := range pvr.VariableSlice {
							// 	pd := PointData{
							// 		DataPointId: value.GetVariableName(),
							// 		Value:       value.GetValue(),
							// 	}
							// 	pds = append(pds, pd)
							// }

							// m.processData(pds)
							// publishData :=  PublishData{Payload:  Payload{Data: [] TimeSeriesData{{
							// 	Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
							// 	Values:    pds,
							// }}}}
							//
							// marshal, _ := json.Marshal(publishData)
							// token := m.mqttClient.Publish(topic, 1, false, marshal)
							// if token.WaitTimeout(mqttTimeout) && token.Error() == nil {
							// 	klog.V(5).InfoS("Succeed to publish MQTT", "topic", topic, "data", publishData)
							// } else {
							// 	klog.V(1).InfoS("Failed to publish MQTT", "topic", topic, "err", token.Error())
							// }
						} else {
							v.( Device).SetCollectStatus( CollectStatusToString[ CollectingError])
						}
					} else {
						klog.V(2).InfoS("Failed to load device", "deviceId", deviceId)
					}
				} else {
					klog.V(2).InfoS("Stopped to collect data", "deviceId", deviceId)
					return
				}
			}
		}
	}(obj.GetID(), results)
	return nil
}

func (m *Manager) Shutdown(context context.Context) error {
	for _, c := range m.brokers {
		c.Destroy(context)
	}
	return nil
}

// func (m *Manager) foldDevice(device  Device)  Device {
// 	return & DeviceMeta{
// 		ObjectMeta:  ObjectMeta{
// 			Name:    device.GetName(),
// 			ID:      device.GetID(),
// 			Version: device.GetVersion(),
// 			ModTime: device.GetModTime(),
// 		},
// 		DeviceModel:   device.GetDeviceModel(),
// 		DeviceCode:    device.GetDeviceCode(),
// 		DeviceType:    device.GetDeviceType(),
// 		CollectStatus: device.GetCollectStatus(),
// 	}
// }

func (m *Manager) heartBeatDetection() {
	tick := time.Tick(heartBeatTimeInterval)
	for {
		select {
		case _, ok := <-m.stopCh:
			if !ok {
				return
			}
		case <-tick:
			resumeDevices := make([]string, 0, 0)
			m.heartBeatDevices.Range(func(key, value any) bool {
				d := value.(Device)
				if err := m.readyCollect(d); err == nil {
					resumeDevices = append(resumeDevices, key.(string))
					return true
				}
				return false
			})
			if len(resumeDevices) > 0 {
				for _, deviceId := range resumeDevices {
					m.heartBeatDevices.Delete(deviceId)
				}
			}
		}
	}
}

func (m *Manager) listeningDeviceStatusCh() {
	for {
		select {
		case _, ok := <-m.stopCh:
			if !ok {
				return
			}
		case statusCh, ok := <-m.deviceStatusCh:
			if !ok {
				return
			}
			split := strings.Split(statusCh, "-")
			deviceId := split[0]
			status := split[1]
			d, exist := m.devices.Load(deviceId)
			if !exist {
				klog.V(2).InfoS("Failed to find device", "deviceId", deviceId)
			}
			m.switchDeviceStatus(d.(Device), status)
		}
	}
}

func (m *Manager) switchDeviceStatus(device Device, status string) {
	cs := device.GetCollectStatus()
	switch StringToCollectStatus[cs] {
	case Collecting:
		switch StringToDeviceStatusCh[status] {
		case Start:
			return
		case Restart:
			_ = m.cancelCollect(device)
			if err := m.readyCollect(device); err != nil {
				if errors.Is(err, ErrConnectDevice) {
					m.heartBeatDevices.Store(device.GetID(), device)
				} else {
					klog.V(2).InfoS("Failed to start process collect device data", "deviceId", device.GetID())
				}
			}
			return
		case Stop:
			_ = m.cancelCollect(device)
			return
		}
	case CollectingError, Error:
		switch StringToDeviceStatusCh[status] {
		case Restart, Start:
			_ = m.cancelCollect(device)
			if err := m.readyCollect(device); err != nil {
				if errors.Is(err, ErrConnectDevice) {
					m.heartBeatDevices.Store(device.GetID(), device)
				} else {
					klog.V(2).InfoS("Failed to start process collect device data", "deviceId", device.GetID())
				}
			}
			return
		case Stop:
			_ = m.cancelCollect(device)
			return
		}
	case EmptyVariable, Unconnected:
		switch StringToDeviceStatusCh[status] {
		case Restart, Start:
			_ = m.cancelCollect(device)
			if err := m.readyCollect(device); err != nil {
				if errors.Is(err, ErrConnectDevice) {
					m.heartBeatDevices.Store(device.GetID(), device)
				} else {
					klog.V(2).InfoS("Failed to start process collect device data", "deviceId", device.GetID())
				}
			}
			return
		case Stop:
			_ = m.cancelCollect(device)
			return
		}
	case Stopped:
		switch StringToDeviceStatusCh[status] {
		case Restart, Start:
			if err := m.readyCollect(device); err != nil {
				if errors.Is(err, ErrConnectDevice) {
					m.heartBeatDevices.Store(device.GetID(), device)
				} else {
					klog.V(2).InfoS("Failed to start process collect device data", "deviceId", device.GetID())
				}
			}
			return
		case Stop:
			return
		}
	}
}
