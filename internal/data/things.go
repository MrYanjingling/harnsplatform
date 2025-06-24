package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"harnsplatform/internal/biz"
)

type ThingsRepo struct {
	data *Data
	log  *log.Helper
}

func NewThingsRepo(data *Data, log *log.Helper) biz.ThingsRepo {
	return &ThingsRepo{
		data: data,
		log:  log,
	}
}

func (s ThingsRepo) Save(ctx context.Context, tt *biz.Things) (*biz.Things, error) {
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Create(tt)
	if result.Error != nil {
		s.log.Errorf("failed to save Things. err:[%v]", result.Error)
		return nil, nil
	}
	return tt, nil
}

func (s ThingsRepo) Update(ctx context.Context, tt *biz.Things) (*biz.Things, error) {
	t := biz.Things{Meta: biz.Meta{}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Model(&t).Where("id = ?", tt.Id).Updates(tt)
	if result.Error != nil {
		s.log.Errorf("failed to update Things. err:[%v]", result.Error)
		return nil, nil
	}
	id, _ := s.FindByID(ctx, tt.Id)
	return id, nil
}

func (s ThingsRepo) FindByID(ctx context.Context, id string) (*biz.Things, error) {
	tt := biz.Things{Meta: biz.Meta{}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).First(&tt, "id = ? ", id)
	if result.Error != nil {
		s.log.Errorf("failed to find Things. err:[%v]", result.Error)
		return nil, nil
	}
	// context
	return &tt, nil
}

func (s ThingsRepo) DeleteByID(ctx context.Context, id string) (*biz.Things, error) {
	tt := biz.Things{Meta: biz.Meta{}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Delete(&tt, "id = ?", id)
	if result.Error != nil {
		s.log.Errorf("failed to delete Things. err:[%v]", result.Error)
		return nil, nil
	}
	return &tt, nil
}

func (s ThingsRepo) DeleteBatch(ctx context.Context, ids []string) error {
	if len(ids) > 0 {
		result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Where("id IN ?", ids).Delete(&biz.Things{})
		if result.Error != nil {
			s.log.Errorf("failed to delete Things. err:[%v]", result.Error)
			return result.Error
		}
	}
	return nil
}

func (s ThingsRepo) ListAll(ctx context.Context) ([]*biz.Things, error) {
	// TODO implement me
	// panic("implement me")
	return nil, nil

}

func (s ThingsRepo) Page(ctx context.Context) ([]*biz.Things, error) {
	// TODO implement me
	// panic("implement me")
	return nil, nil
}

func (s ThingsRepo) List(ctx context.Context, ttq *biz.ThingsQuery) (*biz.PaginationResponse, error) {
	var data []*biz.Things
	query := s.data.DB.Model(&biz.Things{})

	if len(ttq.Name) != 0 {
		query.Where("name LIKE ?", "%"+ttq.Name+"%")
	}

	if len(ttq.ThingTypeId) != 0 {
		query.Where("parent_type_id = ?", ttq.ThingTypeId)
	}

	query, response := ttq.PaginationRequest.List(query, &biz.Things{})

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	response.Items = data

	return response, nil
}
