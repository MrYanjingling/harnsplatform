package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"harnsplatform/internal/auth"
	"harnsplatform/internal/common"
	errors2 "harnsplatform/internal/errors"
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
	if db.Statement.Context.Value(common.SKIP_UPDATE) == true {
		return nil
	}

	user := auth.GetCurrentUser(db)
	if user.Name != "" {
		t.Meta.UpdatedByName = user.Name
		t.Meta.UpdatedById = user.Id
	}

	ctx := db.Statement.Context
	meta := ctx.Value(common.META)
	if readyUpdate, ok := meta.(*Meta); ok {
		if tt, ok := ctx.Value(common.THING_TYPES).(*ThingTypes); ok {
			if readyUpdate.GetVersion() != tt.GetVersion() {
				return errors2.GenerateResourceMismatchError(common.THING_TYPES)
			}
		}
		// ver, _ := strconv.ParseUint(readyUpdate.GetVersion(), 10, 64)
		// t.Version = strconv.FormatUint(ver+uint64(rand.Intn(100)), 10)
		// result := db.Model(&ThingTypes{}).Where("id = ?", t.Id).Update("version", readyUpdate.GetVersion())
	}
	return nil
}

func (t *ThingTypes) AfterUpdate(db *gorm.DB) error {

	if db.Statement.Context.Value(common.SKIP_UPDATE) == true {
		return nil
	}

	ctx := db.Statement.Context
	meta := ctx.Value(common.META)
	if readyUpdate, ok := meta.(*Meta); ok {
		// set version
		ver, _ := strconv.ParseUint(readyUpdate.GetVersion(), 10, 64)
		version := strconv.FormatUint(ver+uint64(rand.Intn(100)), 10)

		c := context.WithValue(ctx, common.SKIP_UPDATE, true)

		result := db.WithContext(c).Model(&ThingTypes{}).Where("id = ?", t.Id).Update(common.VERSION, version)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (t *ThingTypes) BeforeDelete(db *gorm.DB) error {
	ctx := db.Statement.Context
	meta := ctx.Value(common.META)
	if readyDelete, ok := meta.(*Meta); ok {
		if tt, ok := ctx.Value(common.THING_TYPES).(*ThingTypes); ok {
			if readyDelete.GetVersion() != tt.GetVersion() {
				return errors2.GenerateResourceMismatchError(common.THING_TYPES)
			}
		}
		return nil
	}
	return nil
}

type ThingTypesRepo interface {
	Save(context.Context, *ThingTypes) (*ThingTypes, error)
	Update(context.Context, *ThingTypes) (*ThingTypes, error)
	FindByID(context.Context, string) (*ThingTypes, error)
	DeleteByID(context.Context, string) (*ThingTypes, error)
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

func (ttu *ThingTypesUsecase) DeleteThingTypesById(ctx context.Context, id string) (*ThingTypes, error) {
	byID, err := ttu.repo.FindByID(ctx, id)
	if err != nil {
		// 404
		return nil, err
	}
	ctx = context.WithValue(ctx, common.THING_TYPES, byID)
	_, err = ttu.repo.DeleteByID(ctx, id)
	if err != nil {
		// 428
		return nil, err
	}
	return byID, nil
}

func (ttu *ThingTypesUsecase) UpdateThingTypesById(ctx context.Context, tt *ThingTypes) (*ThingTypes, error) {
	byID, err := ttu.repo.FindByID(ctx, tt.Id)
	if err != nil {
		// 404
		return nil, err
	}
	ctx = context.WithValue(ctx, common.THING_TYPES, byID)
	updateID, err := ttu.repo.Update(ctx, tt)
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
