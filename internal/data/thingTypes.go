package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/common"
	"harnsplatform/internal/errors"
)

type thingTypesRepo struct {
	data *Data
	log  *log.Helper
}

func NewThingTypesRepo(data *Data, log *log.Helper) biz.ThingTypesRepo {
	return &thingTypesRepo{
		data: data,
		log:  log,
	}
}

func (s thingTypesRepo) Save(ctx context.Context, tt *biz.ThingTypes) (*biz.ThingTypes, error) {
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Create(tt)
	if result.Error != nil {
		s.log.Errorf("failed to save thingTypes. err:[%v]", result.Error)
		return nil, nil
	}
	return tt, nil
}

func (s thingTypesRepo) Update(ctx context.Context, tt *biz.ThingTypes, oldVersion string) (*biz.ThingTypes, error) {
	ctx = context.WithValue(ctx, common.VERSION, oldVersion)

	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Model(tt).Where("version = ?", oldVersion).Updates(tt)
	if result.Error != nil {
		s.log.Errorf("failed to update thingTypes. err:[%v]", result.Error)
		return nil, nil
	}
	if result.RowsAffected == 0 {
		return nil, errors.GenerateResourceMismatchError(common.THING_TYPES)
	}
	id, _ := s.FindByID(ctx, tt.Id)
	return id, nil
}

func (s thingTypesRepo) FindByID(ctx context.Context, id string) (*biz.ThingTypes, error) {
	tt := biz.ThingTypes{Meta: biz.Meta{}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).First(&tt, "id = ? ", id)
	if result.Error != nil {
		s.log.Errorf("failed to find thingTypes. err:[%v]", result.Error)
		return nil, errors.GenerateResourceNotFoundError(common.THING_TYPES)
	}
	// context
	return &tt, nil
}

func (s thingTypesRepo) DeleteByID(ctx context.Context, id string, version string) (*biz.ThingTypes, error) {
	tt := biz.ThingTypes{Meta: biz.Meta{Id: id}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Where("version = ?", version).Delete(&tt)
	if result.Error != nil {
		s.log.Errorf("failed to delete thingTypes. err:[%v]", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.GenerateResourceMismatchError(common.THING_TYPES)
	}
	return &tt, nil
}

func (s thingTypesRepo) DeleteBatch(ctx context.Context, ids []string) error {
	if len(ids) > 0 {
		result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Where("id IN ?", ids).Delete(&biz.ThingTypes{})
		if result.Error != nil {
			s.log.Errorf("failed to delete thingTypes. err:[%v]", result.Error)
			return result.Error
		}
	}
	return nil
}

func (s thingTypesRepo) ListAll(ctx context.Context) ([]*biz.ThingTypes, error) {
	// TODO implement me
	panic("implement me")
}

func (s thingTypesRepo) Page(ctx context.Context) ([]*biz.ThingTypes, error) {
	// TODO implement me
	panic("implement me")
}

func (s thingTypesRepo) List(ctx context.Context, ttq *biz.ThingTypesQuery) (*biz.PaginationResponse, error) {
	var data []*biz.ThingTypes
	query := s.data.DB.Model(&biz.ThingTypes{})

	if len(ttq.Name) != 0 {
		query.Where("name LIKE ?", "%"+ttq.Name+"%")
	}

	if len(ttq.ParentTypeId) != 0 {
		query.Where("parent_type_id = ?", ttq.ParentTypeId)
	}

	query, response := ttq.PaginationRequest.List(query, &biz.ThingTypes{})

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	response.Items = data

	return response, nil
}
