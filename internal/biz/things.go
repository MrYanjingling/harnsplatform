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
	user := auth.GetCurrentUser(db)
	if user.Name != "" {
		t.Meta.UpdatedByName = user.Name
		t.Meta.UpdatedById = user.Id
	}
	return nil
}

func (t *Things) AfterUpdate(db *gorm.DB) error {
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

func (t *Things) BeforeDelete(db *gorm.DB) error {
	return nil
}

type ThingsRepo interface {
	Save(context.Context, *Things) (*Things, error)
	Update(context.Context, *Things, string) (*Things, error)
	FindByID(context.Context, string) (*Things, error)
	DeleteByID(context.Context, string, string) (*Things, error)
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

func (ttu *ThingsUsecase) DeleteThingsById(ctx context.Context, id string, version string) (*Things, error) {
	byID, err := ttu.repo.FindByID(ctx, id)
	if err != nil {
		// 404
		return nil, err
	}
	ctx = context.WithValue(ctx, common.THINGS, byID)
	_, err = ttu.repo.DeleteByID(ctx, id, version)
	if err != nil {
		// 428
		return nil, err
	}
	return byID, nil
}

func (ttu *ThingsUsecase) UpdateThingsById(ctx context.Context, tt *Things, oldVersion string) (*Things, error) {
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
