package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	v1 "harnsplatform/api/modelmanager/v1"
	"math/rand"
	"strconv"
)

var (
	// ErrUserNotFound is user not found.
	ErrStudentNotFound = errors.NotFound(v1.ErrorReason_RESOURCE_MISMATCH.String(), "thingTypes resource mismatch")
)

type ThingTypes struct {
	Name            string  `gorm:"column:name;type:varchar(64)"`
	ParentTypeId    string  `gorm:"column:parent_type_id;type:varchar(32)"`
	Description     string  `gorm:"column:description;type:varchar(256)"`
	Characteristics JSONMap `gorm:"column:characteristics;type:json"`
	PropertySets    JSONMap `gorm:"column:property_sets;type:json"`
	Meta            Meta    `gorm:"embedded"`
}

func (t *ThingTypes) BeforeSave(db *gorm.DB) error {
	user := GetCurrentUser(db)
	if user.Name != "" {
		t.Meta.CreatedByName = user.Name
		t.Meta.UpdatedByName = user.Name
		t.Meta.CreatedById = user.Id
		t.Meta.UpdatedById = user.Id
		t.Meta.Tenant = user.Tenant
	}
	return nil
}

func (t *ThingTypes) BeforeUpdate(db *gorm.DB) error {
	user := GetCurrentUser(db)
	if user.Name != "" {
		t.Meta.UpdatedByName = user.Name
		t.Meta.UpdatedById = user.Id
	}

	// 从上下文中获取是否已经查询过最新版本
	if latest, ok := db.Get("tt_l_v"); ok {
		if latestThingTypes, ok := latest.(ThingTypes); ok {
			if t.Meta.GetVersion() != latestThingTypes.Meta.GetVersion() {
				return GenerateResourceMismatchError("thingTypes")
			}

			// set version
			ver, _ := strconv.ParseUint(t.Meta.GetVersion(), 10, 64)
			t.Meta.SetVersion(strconv.FormatUint(ver+uint64(rand.Intn(100)), 10))
			return nil
		}
	}

	// 未查询过则执行查询
	var latest ThingTypes
	if err := db.First(&latest, t.Meta.Id).Error; err != nil {
		return err
	}
	// 保存到上下文，避免重复查询
	db.Set("tt_l_v", latest)

	if t.Meta.GetVersion() != latest.Meta.GetVersion() {
		return errors.New(421, "数据已被修改，乐观锁检测失败", "")
	}

	// set version
	ver, _ := strconv.ParseUint(t.Meta.GetVersion(), 10, 64)
	t.Meta.SetVersion(strconv.FormatUint(ver+uint64(rand.Intn(100)), 10))
	return nil
}

type Characteristics struct {
	Name         string `json:"name,omitempty"`
	Unit         string `json:"unit,omitempty"`
	Length       string `json:"length,omitempty"`
	DataType     string `json:"dataType,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty"`
}

type PropertySet struct {
	properties map[string]*Property
}

type Property struct {
	Name       string `json:"name,omitempty"`
	Unit       string `json:"unit,omitempty"`
	Value      string `json:"value,omitempty"`
	DataType   string `json:"dataType,omitempty"`
	AccessMode string `json:"accessMode,omitempty"`
	Min        string `json:"min,omitempty"`
	Max        string `json:"max,omitempty"`
}

type ThingTypesRepo interface {
	Save(context.Context, *ThingTypes) (*ThingTypes, error)
	Update(context.Context, *ThingTypes) (*ThingTypes, error)
	FindByID(context.Context, int64) (*ThingTypes, error)
	ListAll(context.Context) ([]*ThingTypes, error)
}

type ThingTypesUsecase struct {
	repo ThingTypesRepo
	log  *log.Helper
}

func NewThingTypesUsecase(repo ThingTypesRepo, logger log.Logger) *ThingTypesUsecase {
	return &ThingTypesUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (ttu *ThingTypesUsecase) CreateThingTypes(ctx context.Context, tt *ThingTypes) (*ThingTypes, error) {
	ttu.log.WithContext(ctx).Debug("CreateThingTypes: %v", tt)
	return ttu.repo.Save(ctx, tt)
}
