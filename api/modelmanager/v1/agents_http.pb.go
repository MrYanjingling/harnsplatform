// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             v6.31.1
// source: api/modelmanager/v1/Agents.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/common"
	"harnsplatform/internal/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationAgentsCreateAgents = "/api.modelmanager.v1.Agents/CreateAgents"

type AgentsHTTPServer interface {
	CreateAgents(context.Context, Agents) (*biz.Agents, error)
	UpdateAgentsById(context.Context, Agents) (*biz.Agents, error)
	DeleteAgentsById(context.Context, *biz.Meta) (*biz.Agents, error)
	DeleteAgents(context.Context, *BatchIds) (*BatchIds, error)
	GetAgentsById(context.Context, *biz.Meta) (*biz.Agents, error)
	GetAgents(context.Context, *biz.AgentsQuery) (*biz.PaginationResponse, error)
	CreateAgentsMappings(context.Context, []*biz.Mapping) ([]*biz.Mapping, error)
	DeleteAgentsMappingsById(context.Context, *biz.Meta) (*biz.Mapping, error)
	DeleteAgentsMappings(context.Context, *BatchIds) (*BatchIds, error)
	GetMappingsByAgentsId(context.Context, *biz.MappingsQuery) (*biz.PaginationResponse, error)
}

func RegisterAgentsHTTPServer(s *http.Server, srv AgentsHTTPServer) {
	r := s.Route("/")
	r.POST("/model-manager/v1/agents", CreateAgents(srv))
	r.PUT("/model-manager/v1/agents/{id}", UpdateAgentsById(srv))
	r.GET("/model-manager/v1/agents/{id}", GetAgentsById(srv))
	r.DELETE("/model-manager/v1/agents/{id}", DeleteAgentsById(srv))
	r.POST("/model-manager/v1/deleteAgentsBatch", DeleteAgents(srv))
	r.GET("/model-manager/v1/agents", GetAgents(srv))
	r.POST("/model-manager/v1/agents/{id}/mappings", CreateAgentsMappings(srv))
	r.GET("/model-manager/v1/agents/{id}/mappings", GetMappingsByAgentsId(srv))
	r.DELETE("/model-manager/v1/agents/mappings/{id}", DeleteAgentsMappingsById(srv))
	r.POST("/model-manager/v1/deleteMappingsBatch", DeleteAgentsMappings(srv))
}

func CreateAgents(srv AgentsHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var target struct {
			AgentType string `json:"agentType"`
		}
		if err := ctx.Bind(&target); err != nil {
			return err
		}
		in := AgentTypeMap[target.AgentType]()
		http.SetOperation(ctx, OperationAgentsCreateAgents)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateAgents(ctx, req.(Agents))
		})
		if err := ctx.Bind(in); err != nil {
			return err
		}
		if err := ctx.BindQuery(in); err != nil {
			return err
		}
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		reply := out.(*biz.Agents)
		return ctx.Result(200, reply)
	}
}

func UpdateAgentsById(srv AgentsHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		version := ctx.Header().Get(common.ETAG)
		if len(version) == 0 {
			return errors.GenerateResourcePreconditionRequiredError(common.AGENTS)
		}

		var target struct {
			AgentType string `json:"agentType"`
		}
		if err := ctx.Bind(&target); err != nil {
			return err
		}

		in := AgentTypeMap[target.AgentType]()

		if err := ctx.Bind(in); err != nil {
			return err
		}
		if err := ctx.BindVars(in); err != nil {
			return err
		}

		in.SetVersion(version)

		http.SetOperation(ctx, OperationAgentsCreateAgents)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateAgentsById(ctx, req.(Agents))
		})
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		reply := out.(*biz.Agents)
		return ctx.Result(200, reply)
	}
}

func GetAgentsById(srv AgentsHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in biz.Meta
		if err := ctx.BindVars(&in); err != nil {
			return err
		}

		http.SetOperation(ctx, OperationAgentsCreateAgents)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetAgentsById(ctx, req.(*biz.Meta))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*biz.Agents)
		return ctx.Result(200, reply)
	}
}

func DeleteAgentsById(srv AgentsHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in biz.Meta
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		version := ctx.Header().Get(common.ETAG)
		if len(version) == 0 {
			return errors.GenerateResourcePreconditionRequiredError(common.AGENTS)
		}

		in.SetVersion(version)

		http.SetOperation(ctx, OperationAgentsCreateAgents)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteAgentsById(ctx, req.(*biz.Meta))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*biz.Agents)
		return ctx.Result(200, reply)
	}
}

func DeleteAgents(srv AgentsHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var bis BatchIds
		if err := ctx.Bind(&bis); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAgentsCreateAgents)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteAgents(ctx, req.(*BatchIds))
		})
		out, err := h(ctx, &bis)
		if err != nil {
			return err
		}
		reply := out.(*BatchIds)
		return ctx.Result(200, reply)
	}
}

func GetAgents(srv AgentsHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		ttq := biz.AgentsQuery{
			PaginationRequest: &biz.PaginationRequest{},
		}
		if err := ctx.BindQuery(&ttq); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAgentsCreateAgents)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetAgents(ctx, req.(*biz.AgentsQuery))
		})
		out, err := h(ctx, &ttq)
		if err != nil {
			return err
		}
		reply := out.(*biz.PaginationResponse)
		return ctx.Result(200, reply)
	}
}

func CreateAgentsMappings(srv AgentsHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in []*biz.Mapping
		http.SetOperation(ctx, OperationAgentsCreateAgents)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateAgentsMappings(ctx, req.([]*biz.Mapping))
		})
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		reply := out.([]*biz.Mapping)
		return ctx.Result(200, reply)
	}
}

func DeleteAgentsMappingsById(srv AgentsHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in biz.Meta
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		version := ctx.Header().Get(common.ETAG)
		if len(version) == 0 {
			return errors.GenerateResourcePreconditionRequiredError(common.MAPPINGS)
		}

		in.SetVersion(version)

		http.SetOperation(ctx, OperationAgentsCreateAgents)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteAgentsMappingsById(ctx, req.(*biz.Meta))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*biz.Mapping)
		return ctx.Result(200, reply)
	}
}

func DeleteAgentsMappings(srv AgentsHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var bis BatchIds
		if err := ctx.Bind(&bis); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAgentsCreateAgents)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteAgentsMappings(ctx, req.(*BatchIds))
		})
		out, err := h(ctx, &bis)
		if err != nil {
			return err
		}
		reply := out.(*BatchIds)
		return ctx.Result(200, reply)
	}
}

func GetMappingsByAgentsId(srv AgentsHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		ttq := biz.MappingsQuery{
			PaginationRequest: &biz.PaginationRequest{},
		}
		if err := ctx.BindQuery(&ttq); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAgentsCreateAgents)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetMappingsByAgentsId(ctx, req.(*biz.MappingsQuery))
		})
		out, err := h(ctx, &ttq)
		if err != nil {
			return err
		}
		reply := out.(*biz.PaginationResponse)
		return ctx.Result(200, reply)
	}
}

type AgentsHTTPClient interface {
	// CreateAgents(ctx context.Context, req *biz.Agents, opts ...http.CallOption) (rsp *biz.Agents, err error)
	GetAgentsByBrokerId(ctx context.Context, req *biz.AgentsQuery, opts ...http.CallOption) ([]*biz.Agents, error)
}

type AgentsHTTPClientImpl struct {
	cc *http.Client
}

func NewAgentsHTTPClient(client *http.Client) AgentsHTTPClient {
	return &AgentsHTTPClientImpl{client}
}

func (c *AgentsHTTPClientImpl) CreateAgents(ctx context.Context, in *Agents, opts ...http.CallOption) (*biz.Agents, error) {
	var out biz.Agents
	pattern := "/model-manager/v1/Agents"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAgentsCreateAgents))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AgentsHTTPClientImpl) GetAgentsByBrokerId(ctx context.Context, in *biz.AgentsQuery, opts ...http.CallOption) ([]*biz.Agents, error) {
	var out []*biz.Agents
	pattern := "/model-manager/v1/agents"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAgentsCreateAgents))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
