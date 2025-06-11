package biz

import (
	"gorm.io/gorm"
	"harnsplatform/internal/auth"
	"time"
)

type Meta struct {
	Id            string    `gorm:"column:id;primaryKey"`
	Tenant        string    `gorm:"column:tenant;varchar(32);not null;default:main"`
	Version       string    `gorm:"column:version;type:varchar(32);not null"`
	CreatedTime   time.Time `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime   time.Time `gorm:"column:updated_time;autoUpdateTime"`
	CreatedByName string    `gorm:"column:created_by_name;type:varchar(32)"`
	CreatedById   string    `gorm:"column:created_by_id;type:varchar(32)"`
	UpdatedByName string    `gorm:"column:updated_by_name;type:varchar(32)"`
	UpdatedById   string    `gorm:"column:updated_by_id;type:varchar(32)"`
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

func GetCurrentUser(db *gorm.DB) *auth.User {
	cu := &auth.User{
		Id:     "",
		Name:   "",
		Tenant: "",
	}
	ctx := db.Statement.Context
	if user, flag := ctx.Value("user").(*auth.User); flag {
		cu.Name = user.Name
		cu.Id = user.Id
		cu.Tenant = user.Tenant
	}

	return cu
}
