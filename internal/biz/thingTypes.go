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

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrStudentNotFound = errors.NotFound(v1.ErrorReason_RESOURCE_MISMATCH.String(), "thingTypes resource mismatch")
// )

type ThingTypes struct {
	Name            string  `gorm:"column:name;type:varchar(64)"`
	ParentTypeId    string  `gorm:"column:parent_type_id;type:varchar(32)"`
	Description     string  `gorm:"column:description;type:varchar(256)"`
	Characteristics JSONMap `gorm:"column:characteristics;type:json"`
	PropertySets    JSONMap `gorm:"column:property_sets;type:json"`
	Meta            `gorm:"embedded"`
}

type Characteristics struct {
	Name         string `json:"name,omitempty"`
	Unit         string `json:"unit,omitempty"`
	Length       string `json:"length,omitempty"`
	DataType     string `json:"dataType,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty"`
	Value        string `json:"value,omitempty"`
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

type ThingTypesQuery struct {
	Name               string `json:"name,omitempty"`
	ParentTypeId       string `json:"parentTypeId,omitempty"`
	*PaginationRequest `json:",inline"`
}

func (t *ThingTypes) BeforeSave(db *gorm.DB) error {
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

func (t *ThingTypes) BeforeUpdate(db *gorm.DB) error {
	user := auth.GetCurrentUser(db)
	if user.Name != "" {
		t.Meta.UpdatedByName = user.Name
		t.Meta.UpdatedById = user.Id
	}

	return nil
}

func (t *ThingTypes) AfterUpdate(db *gorm.DB) error {
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

func (t *ThingTypes) BeforeDelete(db *gorm.DB) error {
	return nil
}

type ThingTypesRepo interface {
	Save(context.Context, *ThingTypes) (*ThingTypes, error)
	Update(context.Context, *ThingTypes, string) (*ThingTypes, error)
	FindByID(context.Context, string) (*ThingTypes, error)
	DeleteByID(context.Context, string, string) (*ThingTypes, error)
	DeleteBatch(context.Context, []string) error
	ListAll(context.Context) ([]*ThingTypes, error)
	List(ctx context.Context, query *ThingTypesQuery) (*PaginationResponse, error)
}

type ThingTypesUsecase struct {
	repo ThingTypesRepo
	log  *log.Helper
}

func NewThingTypesUsecase(repo ThingTypesRepo, logger *log.Helper) *ThingTypesUsecase {
	return &ThingTypesUsecase{repo: repo, log: logger}
}

func (ttu *ThingTypesUsecase) CreateThingTypes(ctx context.Context, tt *ThingTypes) (*ThingTypes, error) {
	return ttu.repo.Save(ctx, tt)
}

func (ttu *ThingTypesUsecase) GetThingTypesById(ctx context.Context, id string) (*ThingTypes, error) {
	return ttu.repo.FindByID(ctx, id)
}

func (ttu *ThingTypesUsecase) DeleteThingTypesById(ctx context.Context, id string, version string) (*ThingTypes, error) {
	byID, err := ttu.repo.FindByID(ctx, id)
	if err != nil {
		// 404
		return nil, err
	}
	ctx = context.WithValue(ctx, common.THING_TYPES, byID)
	_, err = ttu.repo.DeleteByID(ctx, id, version)
	if err != nil {
		// 428
		return nil, err
	}
	return byID, nil
}

func (ttu *ThingTypesUsecase) UpdateThingTypesById(ctx context.Context, tt *ThingTypes, oldVersion string) (*ThingTypes, error) {
	_, err := ttu.repo.FindByID(ctx, tt.Id)
	if err != nil {
		// 404
		return nil, err
	}
	// ctx = context.WithValue(ctx, common.THING_TYPES, byID)
	updateID, err := ttu.repo.Update(ctx, tt, oldVersion)
	if err != nil {
		// 428
		return nil, err
	}
	return updateID, nil
}

func (ttu *ThingTypesUsecase) DeleteThingTypes(ctx context.Context, ids []string) error {
	err := ttu.repo.DeleteBatch(ctx, ids)
	if err != nil {
		return err
	}
	return nil
}

func (ttu *ThingTypesUsecase) GetThingTypes(ctx context.Context, ttq *ThingTypesQuery) (*PaginationResponse, error) {
	pr, err := ttu.repo.List(ctx, ttq)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
