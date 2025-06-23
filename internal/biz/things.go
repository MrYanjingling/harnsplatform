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

type Things struct {
	Name            string      `gorm:"column:name;type:varchar(64)"`
	ThingTypeId     *string     `gorm:"column:thing_type_id;type:varchar(32)"`
	Description     string      `gorm:"column:description;type:varchar(256)"`
	Characteristics JSONMap     `gorm:"column:characteristics;type:json"`
	PropertySets    JSONMap     `gorm:"column:property_sets;type:json"`
	Combination     StringSlice `gorm:"column:combination;type:text"`
	Meta            `gorm:"embedded"`
}
type ThingsQuery struct {
	Name               string `json:"name,omitempty"`
	ThingTypeId        string `json:"thingTypeId,omitempty"`
	*PaginationRequest `json:",inline"`
}

func (t *Things) BeforeSave(db *gorm.DB) error {
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

func (t *Things) BeforeUpdate(db *gorm.DB) error {
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
		if tt, ok := ctx.Value(common.THINGS).(*Things); ok {
			if readyUpdate.GetVersion() != tt.GetVersion() {
				return errors2.GenerateResourceMismatchError(common.THINGS)
			}
		}
	}
	return nil
}

func (t *Things) AfterUpdate(db *gorm.DB) error {

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

		result := db.WithContext(c).Model(&Things{}).Where("id = ?", t.Id).Update(common.VERSION, version)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (t *Things) BeforeDelete(db *gorm.DB) error {
	ctx := db.Statement.Context
	meta := ctx.Value(common.META)
	if readyDelete, ok := meta.(*Meta); ok {
		if tt, ok := ctx.Value(common.THINGS).(*Things); ok {
			if readyDelete.GetVersion() != tt.GetVersion() {
				return errors2.GenerateResourceMismatchError(common.THINGS)
			}
		}
		return nil
	}
	return nil
}

type ThingsRepo interface {
	Save(context.Context, *Things) (*Things, error)
	Update(context.Context, *Things) (*Things, error)
	FindByID(context.Context, string) (*Things, error)
	DeleteByID(context.Context, string) (*Things, error)
	DeleteBatch(context.Context, []string) error
	ListAll(context.Context) ([]*Things, error)
	List(ctx context.Context, query *ThingsQuery) (*PaginationResponse, error)
}

type ThingsUsecase struct {
	repo ThingsRepo
	log  *log.Helper
}

func NewThingsUsecase(repo ThingsRepo, logger *log.Helper) *ThingsUsecase {
	return &ThingsUsecase{repo: repo, log: logger}
}

func (ttu *ThingsUsecase) CreateThings(ctx context.Context, tt *Things) (*Things, error) {
	return ttu.repo.Save(ctx, tt)
}

func (ttu *ThingsUsecase) GetThingsById(ctx context.Context, id string) (*Things, error) {
	return ttu.repo.FindByID(ctx, id)
}

func (ttu *ThingsUsecase) DeleteThingsById(ctx context.Context, id string) (*Things, error) {
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

func (ttu *ThingsUsecase) UpdateThingsById(ctx context.Context, tt *Things) (*Things, error) {
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

func (ttu *ThingsUsecase) DeleteThings(ctx context.Context, ids []string) error {
	err := ttu.repo.DeleteBatch(ctx, ids)
	if err != nil {
		return err
	}
	return nil
}

func (ttu *ThingsUsecase) GetThings(ctx context.Context, ttq *ThingsQuery) (*PaginationResponse, error) {
	pr, err := ttu.repo.List(ctx, ttq)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
