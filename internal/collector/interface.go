package collector

import (
	"context"
	pb "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/biz"
)

type AgentsManager interface {
	CreateAgents(ctx context.Context, agents pb.Agents) (*biz.Agents, error)
	ValidateMappings(ctx context.Context, mappings []*biz.Mapping) error
}
