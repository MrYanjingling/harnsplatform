package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"harnsplatform/internal/auth"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/common"
	"harnsplatform/internal/errors"
)

type BrokersRepo struct {
	data *Data
	log  *log.Helper
}

func NewBrokersRepo(data *Data, log *log.Helper) biz.BrokersRepo {
	return &BrokersRepo{
		data: data,
		log:  log,
	}
}

func (s BrokersRepo) Save(ctx context.Context, tt *biz.Brokers) (*biz.Brokers, error) {
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Create(tt)
	if result.Error != nil {
		s.log.Errorf("failed to save Brokers. err:[%v]", result.Error)
		return nil, nil
	}
	return tt, nil
}

func (s BrokersRepo) Update(ctx context.Context, tt *biz.Brokers, oldVersion string) (*biz.Brokers, error) {
	ctx = context.WithValue(ctx, common.VERSION, oldVersion)

	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Model(tt).Where("version = ?", oldVersion).Updates(tt)
	if result.Error != nil {
		s.log.Errorf("failed to update Brokers. err:[%v]", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.GenerateResourceMismatchError(common.BROKERS)
	}
	id, _ := s.FindByID(ctx, tt.Id)
	return id, nil
}

func (s BrokersRepo) FindByID(ctx context.Context, id string) (*biz.Brokers, error) {
	tt := biz.Brokers{Meta: biz.Meta{}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).First(&tt, "id = ? ", id)
	if result.Error != nil {
		s.log.Errorf("failed to find Brokers. err:[%v]", result.Error)
		return nil, errors.GenerateResourceNotFoundError(common.BROKERS)
	}
	// context
	return &tt, nil
}

func (s BrokersRepo) DeleteByID(ctx context.Context, id string, version string) (*biz.Brokers, error) {
	tt := biz.Brokers{Meta: biz.Meta{Id: id}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Where("version = ?", version).Delete(&tt)
	if result.Error != nil {
		s.log.Errorf("failed to delete Brokers. err:[%v]", result.Error)
		return nil, nil
	}
	if result.RowsAffected == 0 {
		return nil, errors.GenerateResourceMismatchError(common.BROKERS)
	}
	return &tt, nil
}

func (s BrokersRepo) DeleteBatch(ctx context.Context, ids []string) error {
	if len(ids) > 0 {
		result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Where("id IN ?", ids).Delete(&biz.Brokers{})
		if result.Error != nil {
			s.log.Errorf("failed to delete Brokers. err:[%v]", result.Error)
			return result.Error
		}
	}
	return nil
}

func (s BrokersRepo) ListAll(ctx context.Context) ([]*biz.Brokers, error) {
	// TODO implement me
	// panic("implement me")
	return nil, nil

}

func (s BrokersRepo) Page(ctx context.Context) ([]*biz.Brokers, error) {
	// TODO implement me
	// panic("implement me")
	return nil, nil
}

func (s BrokersRepo) List(ctx context.Context, ttq *biz.BrokersQuery) (*biz.PaginationResponse, error) {
	var data []*biz.Brokers
	query := s.data.DB.Model(&biz.Brokers{})

	user, _ := ctx.Value(common.USER).(*auth.User)
	query.Where("tenant = ?", user.Tenant)

	if len(ttq.Name) != 0 {
		query.Where("name LIKE ?", "%"+ttq.Name+"%")
	}

	query, response := ttq.PaginationRequest.List(query, &biz.Brokers{})

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	response.Items = data

	return response, nil
}
