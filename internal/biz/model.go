package biz

import (
	"time"
)

type Meta struct {
	Id            string    `json:"id" gorm:"column:id;primaryKey"`
	Tenant        string    `json:"tenant" gorm:"column:tenant;varchar(32);not null;default:main"`
	Version       string    `json:"version" gorm:"column:version;type:varchar(32);not null"`
	CreatedTime   time.Time `json:"createdTime" gorm:"column:created_time;autoCreateTime"`
	UpdatedTime   time.Time `json:"updatedTime" gorm:"column:updated_time;autoUpdateTime"`
	CreatedByName string    `json:"createdByName" gorm:"column:created_by_name;type:varchar(32)"`
	CreatedById   string    `json:"createdById" gorm:"column:created_by_id;type:varchar(32)"`
	UpdatedByName string    `json:"updatedByName" gorm:"column:updated_by_name;type:varchar(32)"`
	UpdatedById   string    `json:"updatedById" gorm:"column:updated_by_id;type:varchar(32)"`
}

func (m *Meta) GetVersion() string {
	return m.Version
}

func (m *Meta) SetVersion(version string) {
	m.Version = version
}

type OptimisticLock interface {
	GetVersion() string
	SetVersion(version string)
}
