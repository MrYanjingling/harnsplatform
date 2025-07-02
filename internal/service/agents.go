package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/oklog/ulid/v2"
	pb "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/collector"
	randutil "harnsplatform/internal/utils"
	"strconv"
	"time"
)

type AgentsService struct {
	// pb.UnimplementedAgentsServer

	au  *biz.AgentsUsecase
	log *log.Helper
}

func NewAgentsService(au *biz.AgentsUsecase, logger *log.Helper) *AgentsService {
	return &AgentsService{
		au:  au,
		log: logger,
	}
}

// CreateAgents Validate in this
func (s *AgentsService) CreateAgents(ctx context.Context, req pb.Agents) (*biz.Agents, error) {
	manager := collector.AgentsManagers[req.GetAgentType()]
	agents, err := manager.CreateAgents(ctx, req)
	if err != nil {
		return nil, err
	}

	agents.Id = ulid.MustNewDefault(time.Now()).String()
	agents.Version = strconv.FormatUint(randutil.Uint64n(), 10)

	createAgents, err := s.au.CreateAgents(ctx, agents)
	if err != nil {
		return nil, err
	}
	return createAgents, nil
}

func (s *AgentsService) GetAgentsById(ctx context.Context, req *biz.Meta) (*biz.Agents, error) {
	id, err := s.au.GetAgentsById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (s *AgentsService) UpdateAgentsById(ctx context.Context, req pb.Agents) (*biz.Agents, error) {
	manager := collector.AgentsManagers[req.GetAgentType()]
	agents, err := manager.CreateAgents(ctx, req)
	if err != nil {
		return nil, err
	}

	agents.Id = req.GetId()

	updateAgents, err := s.au.UpdateAgentsById(ctx, agents, req.GetVersion())
	if err != nil {
		return nil, err
	}

	return updateAgents, nil
}

func (s *AgentsService) DeleteAgentsById(ctx context.Context, req *biz.Meta) (*biz.Agents, error) {
	id, err := s.au.DeleteAgentsById(ctx, req.GetId(), req.GetVersion())
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *AgentsService) DeleteAgents(ctx context.Context, req *pb.BatchIds) (*pb.BatchIds, error) {
	err := s.au.DeleteAgents(ctx, req.Ids)
	if err != nil {
		return req, err
	}
	return req, nil
}

func (s *AgentsService) GetAgents(ctx context.Context, req *biz.AgentsQuery) (*biz.PaginationResponse, error) {
	pr, err := s.au.GetAgents(ctx, req)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
