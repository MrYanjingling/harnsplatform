package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"harnsplatform/internal/biz"
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

func (s thingTypesRepo) Update(ctx context.Context, tt *biz.ThingTypes) (*biz.ThingTypes, error) {
	t := biz.ThingTypes{Meta: biz.Meta{}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Model(&t).Where("id = ?", tt.Id).Updates(tt)
	if result.Error != nil {
		s.log.Errorf("failed to update thingTypes. err:[%v]", result.Error)
		return nil, nil
	}
	id, _ := s.FindByID(ctx, tt.Id)
	return id, nil
}

func (s thingTypesRepo) FindByID(ctx context.Context, id string) (*biz.ThingTypes, error) {
	tt := biz.ThingTypes{Meta: biz.Meta{}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).First(&tt, "id = ? ", id)
	if result.Error != nil {
		s.log.Errorf("failed to find thingTypes. err:[%v]", result.Error)
		return nil, nil
	}
	// context
	return &tt, nil
}

func (s thingTypesRepo) DeleteByID(ctx context.Context, id string) (*biz.ThingTypes, error) {
	tt := biz.ThingTypes{Meta: biz.Meta{}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Delete(&tt, "id = ?", id)
	if result.Error != nil {
		s.log.Errorf("failed to delete thingTypes. err:[%v]", result.Error)
		return nil, nil
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
