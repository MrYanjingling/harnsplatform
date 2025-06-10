package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	v1 "harnsplatform/api/modelmanager/v1"
)

var (
	// ErrUserNotFound is user not found.
	ErrStudentNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
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
		t.Meta.updatedByName = user.Name
		t.Meta.CreatedById = user.Id
		t.Meta.updatedById = user.Id
		t.Meta.Tenant = user.Tenant
	}
	return nil
}

func (t *ThingTypes) BeforeUpdate(db *gorm.DB) error {
	user := GetCurrentUser(db)
	if user.Name != "" {
		t.Meta.updatedByName = user.Name
		t.Meta.updatedById = user.Id
	}
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

// GreeterRepo is a Greater repo.
type ThingTypesRepo interface {
	Save(context.Context, *ThingTypes) (*ThingTypes, error)
	Update(context.Context, *ThingTypes) (*ThingTypes, error)
	FindByID(context.Context, int64) (*ThingTypes, error)
	ListAll(context.Context) ([]*ThingTypes, error)
}

// GreeterUsecase is a Greeter usecase.
type ThingTypesUsecase struct {
	repo ThingTypesRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewThingTypesUsecase(repo ThingTypesRepo, logger log.Logger) *ThingTypesUsecase {
	return &ThingTypesUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (ttu *ThingTypesUsecase) CreateThingTypes(ctx context.Context, tt *ThingTypes) (*ThingTypes, error) {
	ttu.log.WithContext(ctx).Debug("CreateThingTypes: %v", tt)
	return ttu.repo.Save(ctx, tt)
}
