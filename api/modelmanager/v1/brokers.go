package v1

import (
	"harnsplatform/internal/biz"
)

type Brokers struct {
	Name                  string                     `json:"name"`
	Description           string                     `json:"description"`
	DeployDetails         *biz.DeployDetails         `json:"deployDetails"`         // 部署相关信息 IP node
	RuntimeType           string                     `json:"runtimeType"`           // single单一架构  redundancy冗余架构
	TimeSeriesStorePeriod *biz.TimeSeriesStorePeriod `json:"TimeSeriesStorePeriod"` // 时序数据存储周期
	Sink                  *biz.Sink                  `json:"sink"`                  // 配置了ThingId相关参数才能sink
	*biz.Meta             `gorm:"embedded"`
}
