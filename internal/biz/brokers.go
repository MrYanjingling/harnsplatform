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

type Brokers struct {
	Name                  string  `gorm:"column:name;type:varchar(64)" json:"name"`
	Description           string  `gorm:"column:description;type:varchar(256)" json:"description"`
	DeployDetails         JSONMap `gorm:"column:deploy_details;type:json" json:"deployDetails"`                  // 部署相关信息 IP node
	RuntimeType           string  `gorm:"column:description;type:varchar(256)" json:"runtimeType"`               // single单一架构  redundancy冗余架构
	TimeSeriesStorePeriod JSONMap `gorm:"column:timeSeries_store_period;type:json" json:"TimeSeriesStorePeriod"` // 时序数据存储周期
	Sink                  JSONMap `gorm:"column:sink;type:json" json:"sink"`                                     // 配置了ThingId相关参数才能sink
	OnBoard               bool    `gorm:"column:on_board;type:TINYINT(1)" json:"onBoard"`                        // 是否注册在线
	Secret                []byte  `gorm:"column:secret;type:varbinary(32)" json:"secret"`                        // 密钥
	OnLine                bool    `gorm:"-" json:"onLine"`                                                       // 上线 在线
	Meta                  `gorm:"embedded"`
}
type BrokersQuery struct {
	Name               string `json:"name,omitempty"`
	*PaginationRequest `json:",inline"`
}

type DeployDetails struct {
	Ip string `json:"ip,omitempty"`
}

type TimeSeriesStorePeriod struct {
	Flag     bool   `json:"flag,omitempty"`
	TimeType string `json:"timeType,omitempty"`
	Period   int    `json:"period,omitempty"`
}

type Sink struct {
	Flag          bool                   `json:"flag,omitempty"`
	SinkMQ        string                 `json:"sinkMQ,omitempty"` // kafka  mqtt
	MQConfig      map[string]interface{} `json:"mqConfig,omitempty"`
	PushFrequency int                    `json:"pushFrequency,omitempty"` // 单位 秒
}

func (t *Brokers) BeforeSave(db *gorm.DB) error {
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

func (t *Brokers) BeforeUpdate(db *gorm.DB) error {
	user := auth.GetCurrentUser(db)
	if user.Name != "" {
		t.Meta.UpdatedByName = user.Name
		t.Meta.UpdatedById = user.Id
	}
	return nil
}

func (t *Brokers) AfterUpdate(db *gorm.DB) error {
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

func (t *Brokers) BeforeDelete(db *gorm.DB) error {
	return nil
}

type BrokersRepo interface {
	Save(context.Context, *Brokers) (*Brokers, error)
	Update(context.Context, *Brokers, string) (*Brokers, error)
	FindByID(context.Context, string) (*Brokers, error)
	DeleteByID(context.Context, string, string) (*Brokers, error)
	DeleteBatch(context.Context, []string) error
	ListAll(context.Context) ([]*Brokers, error)
	List(ctx context.Context, query *BrokersQuery) (*PaginationResponse, error)
}

type BrokersUsecase struct {
	repo BrokersRepo
	log  *log.Helper
}

func NewBrokersUsecase(repo BrokersRepo, logger *log.Helper) *BrokersUsecase {
	return &BrokersUsecase{repo: repo, log: logger}
}

func (ttu *BrokersUsecase) CreateBrokers(ctx context.Context, tt *Brokers) (*Brokers, error) {
	return ttu.repo.Save(ctx, tt)
}

func (ttu *BrokersUsecase) GetBrokersById(ctx context.Context, id string) (*Brokers, error) {
	return ttu.repo.FindByID(ctx, id)
}

func (ttu *BrokersUsecase) DeleteBrokersById(ctx context.Context, id string, version string) (*Brokers, error) {
	byID, err := ttu.repo.FindByID(ctx, id)
	if err != nil {
		// 404
		return nil, err
	}
	ctx = context.WithValue(ctx, common.BROKERS, byID)
	_, err = ttu.repo.DeleteByID(ctx, id, version)
	if err != nil {
		// 428
		return nil, err
	}
	return byID, nil
}

func (ttu *BrokersUsecase) UpdateBrokersById(ctx context.Context, tt *Brokers, oldVersion string) (*Brokers, error) {
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

func (ttu *BrokersUsecase) DeleteBrokers(ctx context.Context, ids []string) error {
	err := ttu.repo.DeleteBatch(ctx, ids)
	if err != nil {
		return err
	}
	return nil
}

func (ttu *BrokersUsecase) GetBrokers(ctx context.Context, ttq *BrokersQuery) (*PaginationResponse, error) {
	pr, err := ttu.repo.List(ctx, ttq)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
