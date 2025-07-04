package model_manager

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
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
func NewHTTPServer(c *conf.Server, thingTypes *service.ThingTypesService, things *service.ThingsService, agents *service.AgentsService, logger *log.Helper) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			record(logger),
			selector.Server(idm(logger)).Match(NewWhiteListMatcher()).Build(),
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

	// auth

	v1.RegisterThingTypesHTTPServer(srv, thingTypes)
	v1.RegisterThingsHTTPServer(srv, things)
	v1.RegisterAgentsHTTPServer(srv, agents)
	return srv
}

func idm(logger *log.Helper) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			hc, _ := ctx.(http.Context)
			token := hc.Header().Get("Authorization")
			logger.Infof("Authorization:[%s]", token)
			c := context.WithValue(hc, "user", &auth.User{
				Id:     "001",
				Name:   "anonymous",
				Tenant: "main",
			})
			return handler(c, req)
		}
	}
}

func record(log *log.Helper) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// Start timer
			start := time.Now()
			defer func() {
				// Stop timer
				latency := time.Now().Sub(start).Seconds()
				log.Debugf("latency [%f]", latency)
			}()
			return handler(ctx, req)
		}
	}
}

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/login"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
