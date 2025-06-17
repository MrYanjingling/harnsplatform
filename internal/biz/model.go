package biz

import (
	"gorm.io/gorm"
	"math"
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

type PaginationRequest struct {
	Page     *int `json:"page" form:"page,omitempty"`
	PageSize *int `json:"pageSize" form:"pageSize,omitempty"`
}

func (pr *PaginationRequest) List(db *gorm.DB, desc interface{}) (*gorm.DB, *PaginationResponse) {
	var count int64
	tx := db.Model(desc).Count(&count)
	if pr.Page == nil && pr.PageSize == nil {
		return tx, &PaginationResponse{
			TotalCount: &count,
		}
	}

	totalPages := int(math.Ceil(float64(count) / float64(*pr.PageSize)))

	if *pr.Page > totalPages && totalPages > 0 {
		pr.Page = &totalPages
	}
	// 执行分页查询
	offset := (*pr.Page - 1) * *pr.PageSize
	return db.Offset(offset).Limit(*pr.PageSize), &PaginationResponse{
		Page:       pr.Page,
		PageSize:   pr.PageSize,
		TotalCount: &count,
		TotalPages: &totalPages,
	}
}

type PaginationResponse struct {
	Page       *int        `json:"page,omitempty"`
	PageSize   *int        `json:"pageSize,omitempty"`
	TotalCount *int64      `json:"totalCount,omitempty"`
	TotalPages *int        `json:"totalPages,omitempty"`
	Items      interface{} `json:"items,omitempty"`
}
