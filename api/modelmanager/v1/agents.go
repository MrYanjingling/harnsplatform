package v1

import (
	"harnsplatform/internal/biz"
	"harnsplatform/internal/common"
)

var AgentTypeMap = map[string]func() Agents{
	"modbus": func() Agents { return &ModbusAgent{} },
}

type Agents interface {
	GetAgentType() common.AgentType
	biz.ObjectMeta
}

type ModbusAgent struct {
	Name             string                 `json:"name,omitempty"`
	Description      string                 `json:"description,omitempty"`
	AgentType        string                 `json:"agentType,omitempty"`
	CollectorCycle   uint                   `json:"collectorCycle,omitempty"`   // 采集周期毫秒
	VariableInterval uint                   `json:"variableInterval,omitempty"` // 变量间隔
	AgentDetails     biz.ModbusAgentDetails `json:"agentDetails,omitempty"`
	Address          biz.ModbusAgentAddress `json:"address,omitempty"`
	Broker           string                 `json:"broker,omitempty"`
	*biz.Meta        `json:",inline"`
}

func (m *ModbusAgent) GetAgentType() common.AgentType {
	return common.StringToAgentType[m.AgentType]
}

type MQTTAgent struct {
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	AgentType        string `json:"agentType,omitempty"`
	CollectorCycle   uint   `json:"collectorCycle,omitempty"`   // 采集周期毫秒
	VariableInterval uint   `json:"variableInterval,omitempty"` // 变量间隔
	// AgentDetails     JSONMap `json:"agentDetails,omitempty"`
	// Address          JSONMap `json:"address,omitempty"`
	Broker    string `json:"broker,omitempty"`
	*biz.Meta `json:",inline"`
}

// todo
type OpcUaAgent struct {
	Name            string                              `json:"name,omitempty"`
	ThingTypeId     *string                             `json:"parentTypeId,omitempty"`
	Description     string                              `json:"description,omitempty"`
	Characteristics map[string]*biz.Characteristics     `json:"characteristics,omitempty"`
	PropertySets    map[string]map[string]*biz.Property `json:"propertySets,omitempty"`
	Combination     []string                            `json:"combination,omitempty"`
	*biz.Meta       `json:",inline"`
}

// todo
type HttpAgent struct {
	Name            string                              `json:"name,omitempty"`
	ThingTypeId     *string                             `json:"parentTypeId,omitempty"`
	Description     string                              `json:"description,omitempty"`
	Characteristics map[string]*biz.Characteristics     `json:"characteristics,omitempty"`
	PropertySets    map[string]map[string]*biz.Property `json:"propertySets,omitempty"`
	Combination     []string                            `json:"combination,omitempty"`
	*biz.Meta       `json:",inline"`
}

// todo
type SqlAgent struct {
	Name            string                              `json:"name,omitempty"`
	ThingTypeId     *string                             `json:"parentTypeId,omitempty"`
	Description     string                              `json:"description,omitempty"`
	Characteristics map[string]*biz.Characteristics     `json:"characteristics,omitempty"`
	PropertySets    map[string]map[string]*biz.Property `json:"propertySets,omitempty"`
	Combination     []string                            `json:"combination,omitempty"`
	*biz.Meta       `json:",inline"`
}
