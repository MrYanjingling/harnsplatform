package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	v1 "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/auth"
	"harnsplatform/internal/conf"
	"harnsplatform/internal/service"
	"time"

	// "github.com/go-kratos/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, thingTypes *service.ThingTypesService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			record(),
			idm(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	opts = append(opts, http.Timeout(c.Http.Timeout))
	srv := http.NewServer(opts...)
	v1.RegisterThingTypesHTTPServer(srv, thingTypes)
	return srv
}

func idm() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			c := context.WithValue(ctx, "user", &auth.User{
				Id:     "001",
				Name:   "anonymous",
				Tenant: "main",
			})
			return handler(c, req)
		}
	}
}

func record() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// Start timer
			start := time.Now()
			defer func() {
				// Stop timer
				latency := time.Now().Sub(start).Seconds()
				fmt.Printf("latency [%f]\n", latency)
			}()
			return handler(ctx, req)
		}
	}
}
