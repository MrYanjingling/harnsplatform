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

func NewThingTypesRepo(data *Data, logger log.Logger) biz.ThingTypesRepo {
	return &thingTypesRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (s thingTypesRepo) Save(ctx context.Context, tt *biz.ThingTypes) (*biz.ThingTypes, error) {
	// c, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	// defer cancelFunc()

	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Create(tt)
	if result.Error != nil {
		s.log.Errorf("failed to save thingTypes. err:[%v]", result.Error)
		return nil, biz.ErrStudentNotFound
	}
	return tt, nil
}

func (s thingTypesRepo) Update(ctx context.Context, tt *biz.ThingTypes) (*biz.ThingTypes, error) {
	// TODO implement me
	panic("implement me")
}

func (s thingTypesRepo) FindByID(ctx context.Context, i int64) (*biz.ThingTypes, error) {
	// TODO implement me
	panic("implement me")
}

func (s thingTypesRepo) ListAll(ctx context.Context) ([]*biz.ThingTypes, error) {
	// TODO implement me
	panic("implement me")
}
