package server

import (
	v1 "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/conf"
	"harnsplatform/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, thingTypes *service.ThingTypesService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	opts = append(opts, grpc.Timeout(c.Grpc.Timeout))
	srv := grpc.NewServer(opts...)
	v1.RegisterThingTypesServer(srv, thingTypes)
	return srv
}
