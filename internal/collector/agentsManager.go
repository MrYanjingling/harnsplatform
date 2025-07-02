package collector

import (
	"context"
	pb "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/collector/modbus"
	"harnsplatform/internal/common"
)

var AgentsManagers = map[common.AgentType]AgentsManager{
	common.AgentTypeModbus: &modbus.AgentsManager{},
}

type AgentsManager interface {
	CreateAgents(ctx context.Context, agents pb.Agents) (*biz.Agents, error)
	// DeleteDevice(device runtime.Device) (runtime.Device, error)
	// UpdateValidation(deviceType v1.DeviceType, device runtime.Device) error
	// UpdateDevice(id string, agents pb.Agents, device runtime.Device) (runtime.Device, error)
}
