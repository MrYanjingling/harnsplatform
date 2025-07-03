package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"harnsplatform/internal/auth"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/common"
	"harnsplatform/internal/errors"
)

type AgentsRepo struct {
	data *Data
	log  *log.Helper
}

func NewAgentsRepo(data *Data, log *log.Helper) biz.AgentsRepo {
	return &AgentsRepo{
		data: data,
		log:  log,
	}
}

func (s AgentsRepo) Save(ctx context.Context, tt *biz.Agents) (*biz.Agents, error) {
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Create(tt)
	if result.Error != nil {
		s.log.Errorf("failed to save Agents. err:[%v]", result.Error)
		return nil, nil
	}
	return tt, nil
}

func (s AgentsRepo) Update(ctx context.Context, tt *biz.Agents, oldVersion string) (*biz.Agents, error) {
	ctx = context.WithValue(ctx, common.VERSION, oldVersion)

	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Model(tt).Where("version = ?", oldVersion).Updates(tt)
	if result.Error != nil {
		s.log.Errorf("failed to update Agents. err:[%v]", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.GenerateResourceMismatchError(common.AGENTS)
	}
	id, _ := s.FindByID(ctx, tt.Id)
	return id, nil
}

func (s AgentsRepo) FindByID(ctx context.Context, id string) (*biz.Agents, error) {
	tt := biz.Agents{Meta: biz.Meta{}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).First(&tt, "id = ? ", id)
	if result.Error != nil {
		s.log.Errorf("failed to find Agents. err:[%v]", result.Error)
		return nil, errors.GenerateResourceNotFoundError(common.AGENTS)
	}
	// context
	return &tt, nil
}

func (s AgentsRepo) DeleteByID(ctx context.Context, id string, version string) (*biz.Agents, error) {
	tt := biz.Agents{Meta: biz.Meta{Id: id}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Where("version = ?", version).Delete(&tt)
	if result.Error != nil {
		s.log.Errorf("failed to delete Agents. err:[%v]", result.Error)
		return nil, nil
	}
	if result.RowsAffected == 0 {
		return nil, errors.GenerateResourceMismatchError(common.AGENTS)
	}
	return &tt, nil
}

func (s AgentsRepo) DeleteBatch(ctx context.Context, ids []string) error {
	if len(ids) > 0 {
		result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Where("id IN ?", ids).Delete(&biz.Agents{})
		if result.Error != nil {
			s.log.Errorf("failed to delete Agents. err:[%v]", result.Error)
			return result.Error
		}
	}
	return nil
}

func (s AgentsRepo) ListAll(ctx context.Context) ([]*biz.Agents, error) {
	// TODO implement me
	// panic("implement me")
	return nil, nil

}

func (s AgentsRepo) Page(ctx context.Context) ([]*biz.Agents, error) {
	// TODO implement me
	// panic("implement me")
	return nil, nil
}

func (s AgentsRepo) List(ctx context.Context, ttq *biz.AgentsQuery) (*biz.PaginationResponse, error) {
	var data []*biz.Agents
	query := s.data.DB.Model(&biz.Agents{})

	user, _ := ctx.Value(common.USER).(*auth.User)
	query.Where("tenant = ?", user.Tenant)

	if len(ttq.Name) != 0 {
		query.Where("name LIKE ?", "%"+ttq.Name+"%")
	}

	if len(ttq.AgentType) != 0 {
		query.Where("agent_type = ?", ttq.AgentType)
	}

	query, response := ttq.PaginationRequest.List(query, &biz.Agents{})

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	response.Items = data

	return response, nil
}

func (s AgentsRepo) SaveMappings(ctx context.Context, mappings []*biz.Mapping) ([]*biz.Mapping, error) {
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Create(&mappings)
	if result.Error != nil {
		s.log.Errorf("failed to save Mappings. err:[%v]", result.Error)
		return nil, nil
	}
	return mappings, nil
}

func (s AgentsRepo) DeleteMappingByID(ctx context.Context, id string, version string) (*biz.Mapping, error) {
	tt := biz.Mapping{Meta: biz.Meta{Id: id}}
	result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Where("version = ?", version).Delete(&tt)
	if result.Error != nil {
		s.log.Errorf("failed to delete Mappings. err:[%v]", result.Error)
		return nil, nil
	}
	if result.RowsAffected == 0 {
		return nil, errors.GenerateResourceMismatchError(common.MAPPINGS)
	}
	return &tt, nil
}

func (s AgentsRepo) DeleteMappings(ctx context.Context, ids []string) error {
	if len(ids) > 0 {
		result := s.data.DB.WithContext(context.WithoutCancel(ctx)).Where("id IN ?", ids).Delete(&biz.Mapping{})
		if result.Error != nil {
			s.log.Errorf("failed to delete Mappings. err:[%v]", result.Error)
			return result.Error
		}
	}
	return nil
}

func (s AgentsRepo) ListMappings(ctx context.Context, mq *biz.MappingsQuery) (*biz.PaginationResponse, error) {
	var data []*biz.Mapping
	query := s.data.DB.Model(&biz.Mapping{})

	user, _ := ctx.Value(common.USER).(*auth.User)
	query.Where("tenant = ?", user.Tenant)

	if len(mq.AgentId) != 0 {
		query.Where("agent_id = ?", mq.AgentId)
	}

	query, response := mq.PaginationRequest.List(query, &biz.Agents{})

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	response.Items = data

	return response, nil
}
