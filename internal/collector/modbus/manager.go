package modbus

import (
	"context"
	"github.com/imdario/mergo"
	pb "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/common"
	"harnsplatform/internal/errors"
)

type AgentsManager struct {
}

func (m *AgentsManager) ValidateMappings(ctx context.Context, mappings []*biz.Mapping) error {
	return nil
}

func (m *AgentsManager) CreateAgents(ctx context.Context, agents pb.Agents) (*biz.Agents, error) {
	modbusAgents, ok := agents.(*pb.ModbusAgent)
	if !ok {
		return nil, errors.GenerateAgentsUnsupportedError(common.AgentTypeToString[agents.GetAgentType()])
	}
	bz := &biz.Agents{
		Name:             modbusAgents.Name,
		AgentType:        common.MODBUS,
		Description:      modbusAgents.Description,
		CollectorCycle:   modbusAgents.CollectorCycle,
		VariableInterval: modbusAgents.VariableInterval,
		Broker:           modbusAgents.Broker,
	}

	adv := map[string]interface{}{}
	if err := mergo.Map(&adv, modbusAgents.AgentDetails); err != nil {
		return nil, err
	}
	bz.AgentDetails = adv

	av := map[string]interface{}{}
	if err := mergo.Map(&av, modbusAgents.Address); err != nil {
		return nil, err
	}
	bz.Address = av

	return bz, nil
}
