package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/oklog/ulid/v2"
	pb "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/biz"
	randutil "harnsplatform/internal/utils"
	"strconv"
	"time"
)

type ThingTypesService struct {
	// pb.UnimplementedThingTypesServer

	ttu *biz.ThingTypesUsecase
	log *log.Helper
}

func NewThingTypesService(ttu *biz.ThingTypesUsecase, logger *log.Helper) *ThingTypesService {
	return &ThingTypesService{
		ttu: ttu,
		log: logger,
	}
}

// CreateThingTypes Validate in this
func (s *ThingTypesService) CreateThingTypes(ctx context.Context, req *pb.ThingTypes) (*biz.ThingTypes, error) {
	id := ulid.MustNewDefault(time.Now()).String()

	tt := &biz.ThingTypes{
		Name:            req.Name,
		ParentTypeId:    req.ParentTypeId,
		Description:     req.Description,
		Characteristics: map[string]interface{}{},
		PropertySets:    map[string]interface{}{},
		Meta: biz.Meta{
			Id:      id,
			Version: strconv.FormatUint(randutil.Uint64n(), 10),
		},
	}
	for key, characteristic := range req.Characteristics {
		tt.Characteristics[key] = &biz.Characteristics{
			Name:         characteristic.Name,
			Unit:         characteristic.Unit,
			Length:       characteristic.Length,
			DataType:     characteristic.DataType,
			DefaultValue: characteristic.DefaultValue,
		}
	}
	for key, ps := range req.PropertySets {
		tt.PropertySets[key] = make(map[string]*biz.Property, 0)
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

	thingTypes, err := s.ttu.CreateThingTypes(ctx, tt)
	if err != nil {
		return nil, err
	}

	return thingTypes, nil
}

func (s *ThingTypesService) GetThingTypesById(ctx context.Context, req *pb.ThingTypes) (*biz.ThingTypes, error) {
	id, err := s.ttu.GetThingTypesById(ctx, req.Meta.Id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (s *ThingTypesService) UpdateThingTypesById(ctx context.Context, req *pb.ThingTypes) (*biz.ThingTypes, error) {
	tt := &biz.ThingTypes{
		Name:            req.Name,
		ParentTypeId:    req.ParentTypeId,
		Description:     req.Description,
		Characteristics: map[string]interface{}{},
		PropertySets:    map[string]interface{}{},
		Meta: biz.Meta{
			Id: req.Id,
		},
	}
	for key, characteristic := range req.Characteristics {
		tt.Characteristics[key] = &biz.Characteristics{
			Name:         characteristic.Name,
			Unit:         characteristic.Unit,
			Length:       characteristic.Length,
			DataType:     characteristic.DataType,
			DefaultValue: characteristic.DefaultValue,
		}
	}
	for key, ps := range req.PropertySets {
		tt.PropertySets[key] = make(map[string]*biz.Property, 0)
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

	id, err := s.ttu.UpdateThingTypesById(ctx, tt, req.Version)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *ThingTypesService) DeleteThingTypesById(ctx context.Context, req *pb.ThingTypes) (*biz.ThingTypes, error) {
	id, err := s.ttu.DeleteThingTypesById(ctx, req.Meta.Id, req.Version)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *ThingTypesService) DeleteThingTypes(ctx context.Context, req *pb.BatchIds) (*pb.BatchIds, error) {
	err := s.ttu.DeleteThingTypes(ctx, req.Ids)
	if err != nil {
		return req, err
	}
	return req, nil
}

func (s *ThingTypesService) GetThingTypes(ctx context.Context, req *biz.ThingTypesQuery) (*biz.PaginationResponse, error) {
	pr, err := s.ttu.GetThingTypes(ctx, req)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
