package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/oklog/ulid/v2"
	pb "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/common"
	randutil "harnsplatform/internal/utils"
	"strconv"
	"time"
)

type ThingsService struct {
	// pb.UnimplementedThingsServer

	tu  *biz.ThingsUsecase
	ttu *biz.ThingTypesUsecase
	log *log.Helper
}

func NewThingsService(ttu *biz.ThingTypesUsecase, tu *biz.ThingsUsecase, logger *log.Helper) *ThingsService {
	return &ThingsService{
		ttu: ttu,
		tu:  tu,
		log: logger,
	}
}

// CreateThings Validate in this
func (s *ThingsService) CreateThings(ctx context.Context, req *pb.Things) (*biz.Things, error) {
	id := ulid.MustNewDefault(time.Now()).String()

	tt := &biz.Things{
		Name:            req.Name,
		ThingTypeId:     req.ThingTypeId,
		Description:     req.Description,
		Characteristics: map[string]interface{}{},
		PropertySets:    map[string]interface{}{},
		Combination:     []string{},
		Meta: biz.Meta{
			Id:      id,
			Version: strconv.FormatUint(randutil.Uint64n(), 10),
		},
	}
	for key, characteristic := range req.Characteristics {
		tt.Characteristics[key] = &biz.Characteristics{
			Name:  characteristic.Name,
			Unit:  characteristic.Unit,
			Value: randutil.Ternary(len(characteristic.Value) == 0, characteristic.DefaultValue, characteristic.Value),
		}
	}
	if req.ThingTypeId != nil {
		thingType, err := s.ttu.GetThingTypesById(ctx, *req.ThingTypeId)
		if err != nil {
			return nil, err
		} else {
			for key, ps := range thingType.PropertySets {
				ps := ps.(map[string]*biz.Property)
				for k, property := range ps {
					tt.PropertySets[key].(map[string]*biz.Property)[k] = &biz.Property{
						Name:       property.Name,
						Unit:       property.Unit,
						Value:      property.Value,
						DataType:   property.DataType,
						AccessMode: property.AccessMode,
						Min:        property.Min,
						Max:        property.Max,
					}
				}
			}
		}
	} else {
		tt.Combination = append(tt.Combination, req.Combination...)
	}

	Things, err := s.tu.CreateThings(ctx, tt)
	if err != nil {
		return nil, err
	}

	return Things, nil
}

func (s *ThingsService) GetThingsById(ctx context.Context, req *pb.Things) (*biz.Things, error) {
	id, err := s.tu.GetThingsById(ctx, req.Meta.Id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (s *ThingsService) UpdateThingsById(ctx context.Context, req *pb.Things) (*biz.Things, error) {
	c := context.WithValue(ctx, common.META, req.Meta)

	tt := &biz.Things{
		Name:            req.Name,
		ThingTypeId:     req.ThingTypeId,
		Description:     req.Description,
		Characteristics: map[string]interface{}{},
		PropertySets:    map[string]interface{}{},
		Meta: biz.Meta{
			Id: req.Id,
		},
	}
	for key, characteristic := range req.Characteristics {
		tt.Characteristics[key] = &biz.Characteristics{
			Name:  characteristic.Name,
			Unit:  characteristic.Unit,
			Value: randutil.Ternary(len(characteristic.Value) == 0, characteristic.DefaultValue, characteristic.Value),
		}
	}
	if req.ThingTypeId != nil {
		thingType, err := s.ttu.GetThingTypesById(ctx, *req.ThingTypeId)
		if err != nil {
			return nil, err
		} else {
			for key, ps := range thingType.PropertySets {
				ps := ps.(map[string]*biz.Property)
				for k, property := range ps {
					tt.PropertySets[key].(map[string]*biz.Property)[k] = &biz.Property{
						Name:       property.Name,
						Unit:       property.Unit,
						Value:      property.Value,
						DataType:   property.DataType,
						AccessMode: property.AccessMode,
						Min:        property.Min,
						Max:        property.Max,
					}
				}
			}
		}
	} else {
		tt.Combination = append(tt.Combination, req.Combination...)
	}

	id, err := s.tu.UpdateThingsById(c, tt)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *ThingsService) DeleteThingsById(ctx context.Context, req *pb.Things) (*biz.Things, error) {
	c := context.WithValue(ctx, common.META, req.Meta)
	id, err := s.tu.DeleteThingsById(c, req.Meta.Id)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *ThingsService) DeleteThings(ctx context.Context, req *pb.BatchIds) (*pb.BatchIds, error) {
	err := s.tu.DeleteThings(ctx, req.Ids)
	if err != nil {
		return req, err
	}
	return req, nil
}

func (s *ThingsService) GetThings(ctx context.Context, req *biz.ThingsQuery) (*biz.PaginationResponse, error) {
	pr, err := s.tu.GetThings(ctx, req)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
