package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/imdario/mergo"
	"github.com/oklog/ulid/v2"
	pb "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/biz"
	randutil "harnsplatform/internal/utils"
	"strconv"
	"time"
)

type BrokersService struct {
	// pb.UnimplementedBrokersServer

	tu  *biz.BrokersUsecase
	ttu *biz.ThingTypesUsecase
	log *log.Helper
}

func NewBrokersService(tu *biz.BrokersUsecase, ttu *biz.ThingTypesUsecase, logger *log.Helper) *BrokersService {
	return &BrokersService{
		ttu: ttu,
		tu:  tu,
		log: logger,
	}
}

// CreateBrokers Validate in this
func (s *BrokersService) CreateBrokers(ctx context.Context, req *pb.Brokers) (*biz.Brokers, error) {
	id := ulid.MustNewDefault(time.Now()).String()
	tt := &biz.Brokers{
		Name:        req.Name,
		Description: req.Description,
		RuntimeType: req.RuntimeType,
		OnBoard:     false,
		OnLine:      false,
		Meta: biz.Meta{
			Id:      id,
			Version: strconv.FormatUint(randutil.Uint64n(), 10),
		},
	}

	ddv := map[string]interface{}{}
	if err := mergo.Map(&ddv, req.DeployDetails); err != nil {
		return nil, err
	}
	tt.DeployDetails = ddv

	sv := map[string]interface{}{}
	if err := mergo.Map(&sv, req.Sink); err != nil {
		return nil, err
	}
	tt.Sink = sv

	tsv := map[string]interface{}{}
	if err := mergo.Map(&tsv, req.TimeSeriesStorePeriod); err != nil {
		return nil, err
	}
	tt.TimeSeriesStorePeriod = tsv

	Brokers, err := s.tu.CreateBrokers(ctx, tt)
	if err != nil {
		return nil, err
	}

	return Brokers, nil
}

func (s *BrokersService) GetBrokersById(ctx context.Context, req *biz.Meta) (*biz.Brokers, error) {
	id, err := s.tu.GetBrokersById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (s *BrokersService) UpdateBrokersById(ctx context.Context, req *pb.Brokers) (*biz.Brokers, error) {
	tt := &biz.Brokers{
		Name:        req.Name,
		Description: req.Description,
		RuntimeType: req.RuntimeType,
		OnBoard:     false,
		OnLine:      false,
		Meta: biz.Meta{
			Id:      req.GetId(),
		},
	}

	ddv := map[string]interface{}{}
	if err := mergo.Map(&ddv, req.DeployDetails); err != nil {
		return nil, err
	}
	tt.DeployDetails = ddv

	sv := map[string]interface{}{}
	if err := mergo.Map(&sv, req.Sink); err != nil {
		return nil, err
	}
	tt.Sink = sv

	tsv := map[string]interface{}{}
	if err := mergo.Map(&tsv, req.TimeSeriesStorePeriod); err != nil {
		return nil, err
	}
	tt.TimeSeriesStorePeriod = tsv

	id, err := s.tu.UpdateBrokersById(ctx, tt, req.Version)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *BrokersService) DeleteBrokersById(ctx context.Context, req *biz.Meta) (*biz.Brokers, error) {
	id, err := s.tu.DeleteBrokersById(ctx, req.GetId(), req.GetVersion())
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *BrokersService) DeleteBrokers(ctx context.Context, req *pb.BatchIds) (*pb.BatchIds, error) {
	err := s.tu.DeleteBrokers(ctx, req.Ids)
	if err != nil {
		return req, err
	}
	return req, nil
}

func (s *BrokersService) GetBrokers(ctx context.Context, req *biz.BrokersQuery) (*biz.PaginationResponse, error) {
	pr, err := s.tu.GetBrokers(ctx, req)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
