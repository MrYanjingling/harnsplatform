package collector

import (
	"harnsplatform/internal/collector/modbus"
	"harnsplatform/internal/common"
)

var AgentsManagers = map[common.AgentType]AgentsManager{
	common.AgentTypeModbus: &modbus.AgentsManager{},
}


