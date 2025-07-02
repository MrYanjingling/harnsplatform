package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"harnsplatform/internal/auth"
	"harnsplatform/internal/common"
	"math/rand"
	"strconv"
)

type Agents struct {
	Name             string     `gorm:"column:name;type:varchar(64)" json:"name"`
	AgentType        string     `gorm:"column:agent_type;type:varchar(32)" json:"agentType"`
	Description      string     `gorm:"column:description;type:varchar(256)" json:"dataType"`
	CollectorCycle   uint       `gorm:"column:collector_cycle;type:int" json:"collectorCycle"`     // 采集周期毫秒
	VariableInterval uint       `gorm:"column:variable_interval;type:int" json:"variableInterval"` // 变量间隔
	AgentDetails     JSONMap    `gorm:"column:agent_details;type:json" json:"agentDetails"`
	Address          JSONMap    `gorm:"column:address;type:json" json:"address"`
	Mappings         []*Mapping `gorm:"foreignKey:agent_id;references:id" json:"mappings"`
	Broker           string     `gorm:"column:broker;type:varchar(32);default:main;index:base_index,priority:2" json:"broker"`
	Meta             `gorm:"embedded"`
}

type Mapping struct {
	AgentId      string            `gorm:"column:agent_id;type:varchar(32);index:idx_agent_id" json:"agentId"`
	DataType     common.DataType   `gorm:"column:data_type;type:varchar(32)" json:"dataType"`                     // bool、int16、float32、float64、int32、int64、uint16
	Name         string            `gorm:"column:name;type:varchar(32)"  json:"name"`                             // 变量名称
	Variable     string            `gorm:"column:variable;type:varchar(32)"  json:"variable"`                     // 变量地址 4655536 functionCode = 4
	Rate         string            `gorm:"column:rate;type:varchar(32)"  json:"rate"`                             // 比率
	Offset       string            `gorm:"column:offset;type:varchar(32)"  json:"offset"`                         // 数量
	DefaultValue string            `gorm:"column:default_value;type:varchar(256)"  json:"defaultValue,omitempty"` // 默认值
	Value        interface{}       `gorm:"-" json:"value,omitempty"`                                              // 值
	AccessMode   common.AccessMode `gorm:"column:access_mode;type:varchar(2)"  json:"accessMode"`                 // 读写属性
	Target       `gorm:"embedded"`
}

type Target struct {
	ThingId         string `gorm:"column:thing_id;type:varchar(32)"  json:"thingId"`
	ThingName       string `gorm:"column:thing_name;type:varchar(64)"  json:"thingName"`
	ThingTypeName   string `gorm:"column:thing_type_name;type:varchar(64)"  json:"thingTypeName"`
	ThingTypeId     string `gorm:"column:thing_type_id;type:varchar(32)"  json:"thingTypeId"`
	PropertySetName string `gorm:"column:propertySet_name;type:varchar(64)"  json:"propertySetName"`
	Property        string `gorm:"column:property;type:varchar(64)"  json:"property"`
}

type AgentsQuery struct {
	Name               string `json:"name,omitempty"`
	AgentType          string `json:"agentType,omitempty"`
	*PaginationRequest `json:",inline"`
}

type ModbusAgentDetails struct {
	Slave           uint   `json:"slave" binding:"required"`                                  // 下位机号
	MemoryLayout    string `json:"memoryLayout" binding:"required,oneof=ABCD BADC CDAB DCBA"` // 内存布局 DCBA CDAB BADC ABCD
	PositionAddress uint   `json:"positionAddress,omitempty"`                                 // 起始地址
}

type ModbusAgentAddress struct {
	Location string                    `json:"location"` // 地址路径
	Option   *ModbusAgentAddressOption `json:"option"`   // 地址其他参数
}

type ModbusAgentAddressOption struct {
	Port     int    `json:"port,omitempty"`     // 端口号
	BaudRate int    `json:"baudRate,omitempty"` // 波特率
	DataBits int    `json:"dataBits,omitempty"` // 数据位
	Parity   string `json:"parity,omitempty"`   // 校验位
	StopBits string `json:"stopBits,omitempty"` // 停止位
}

func (t *Agents) BeforeSave(db *gorm.DB) error {
	user := auth.GetCurrentUser(db)
	if user.Name != "" {
		t.Meta.CreatedByName = user.Name
		t.Meta.UpdatedByName = user.Name
		t.Meta.CreatedById = user.Id
		t.Meta.UpdatedById = user.Id
		t.Meta.Tenant = user.Tenant
	}
	return nil
}

func (t *Agents) BeforeUpdate(db *gorm.DB) error {
	user := auth.GetCurrentUser(db)
	if user.Name != "" {
		t.Meta.UpdatedByName = user.Name
		t.Meta.UpdatedById = user.Id
	}
	return nil
}

func (t *Agents) AfterUpdate(db *gorm.DB) error {
	if db.Statement.RowsAffected != 0 {
		ver := db.Statement.Context.Value(common.VERSION)
		if v, ok := ver.(string); ok {
			ov, _ := strconv.ParseUint(v, 10, 64)
			version := strconv.FormatUint(ov+uint64(rand.Intn(100)), 10)
			err := db.Model(t).UpdateColumn(common.VERSION, version).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Agents) BeforeDelete(db *gorm.DB) error {
	return nil
}

type AgentsRepo interface {
	Save(context.Context, *Agents) (*Agents, error)
	Update(context.Context, *Agents, string) (*Agents, error)
	FindByID(context.Context, string) (*Agents, error)
	DeleteByID(context.Context, string, string) (*Agents, error)
	DeleteBatch(context.Context, []string) error
	ListAll(context.Context) ([]*Agents, error)
	List(ctx context.Context, query *AgentsQuery) (*PaginationResponse, error)
}

type AgentsUsecase struct {
	repo AgentsRepo
	log  *log.Helper
}

func NewAgentsUsecase(repo AgentsRepo, logger *log.Helper) *AgentsUsecase {
	return &AgentsUsecase{repo: repo, log: logger}
}

func (ttu *AgentsUsecase) CreateAgents(ctx context.Context, tt *Agents) (*Agents, error) {
	return ttu.repo.Save(ctx, tt)
}

func (ttu *AgentsUsecase) GetAgentsById(ctx context.Context, id string) (*Agents, error) {
	return ttu.repo.FindByID(ctx, id)
}

func (ttu *AgentsUsecase) DeleteAgentsById(ctx context.Context, id string, version string) (*Agents, error) {
	byID, err := ttu.repo.FindByID(ctx, id)
	if err != nil {
		// 404
		return nil, err
	}
	ctx = context.WithValue(ctx, common.THINGS, byID)
	_, err = ttu.repo.DeleteByID(ctx, id, version)
	if err != nil {
		// 428
		return nil, err
	}
	return byID, nil
}

func (ttu *AgentsUsecase) UpdateAgentsById(ctx context.Context, tt *Agents, oldVersion string) (*Agents, error) {
	// todo consider deleted
	_, err := ttu.repo.FindByID(ctx, tt.Id)
	if err != nil {
		// 404
		return nil, err
	}

	updateID, err := ttu.repo.Update(ctx, tt, oldVersion)
	if err != nil {
		// 428
		return nil, err
	}
	return updateID, nil
}

func (ttu *AgentsUsecase) DeleteAgents(ctx context.Context, ids []string) error {
	err := ttu.repo.DeleteBatch(ctx, ids)
	if err != nil {
		return err
	}
	return nil
}

func (ttu *AgentsUsecase) GetAgents(ctx context.Context, ttq *AgentsQuery) (*PaginationResponse, error) {
	pr, err := ttu.repo.List(ctx, ttq)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
