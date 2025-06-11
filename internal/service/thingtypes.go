package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/biz"
	randutil "harnsplatform/internal/utils"
	"strconv"
	"time"
)

type ThingTypesService struct {
	pb.UnimplementedThingTypesServer

	ttu *biz.ThingTypesUsecase
	log *log.Helper
}

func NewThingTypesService(ttu *biz.ThingTypesUsecase, logger log.Logger) *ThingTypesService {
	return &ThingTypesService{
		ttu: ttu,
		log: log.NewHelper(logger),
	}
}

// CreateThingTypes Validate in this
func (s *ThingTypesService) CreateThingTypes(ctx context.Context, req *pb.CreateThingTypesRequest) (*pb.CreateThingTypesReply, error) {
	id := ulid.MustNewDefault(time.Now()).String()

	tt := &biz.ThingTypes{
		Name:            req.GetName(),
		ParentTypeId:    req.GetParentTypeId(),
		Description:     req.GetDescription(),
		Characteristics: map[string]interface{}{},
		PropertySets:    map[string]interface{}{},
		Meta: biz.Meta{
			Id:      id,
			Version: strconv.FormatUint(randutil.Uint64n(), 10),
		},
	}
	for key, characteristic := range req.GetCharacteristics() {
		tt.Characteristics[key] = &biz.Characteristics{
			Name:         characteristic.GetName(),
			Unit:         characteristic.GetUnit(),
			Length:       characteristic.GetLength(),
			DataType:     characteristic.GetDataType(),
			DefaultValue: characteristic.GetDefaultValue(),
		}
	}
	for key, ps := range req.GetPropertySets() {
		tt.PropertySets[key] = make(map[string]*biz.Property, 0)
		for k, property := range ps.GetProperties() {
			tt.PropertySets[key].(map[string]*biz.Property)[k] = &biz.Property{
				Name:       property.GetName(),
				Unit:       property.GetUnit(),
				Value:      property.GetValue(),
				DataType:   pb.DataType_name[int32(property.GetDataType())],
				AccessMode: pb.AccessMode_name[int32(property.GetAccessMode())],
				Min:        property.GetMin(),
				Max:        property.GetMax(),
			}
		}
	}

	thingTypes, err := s.ttu.CreateThingTypes(ctx, tt)
	if err != nil {
		return nil, err
	}

	return &pb.CreateThingTypesReply{
		Name:            thingTypes.Name,
		ParentTypeId:    thingTypes.ParentTypeId,
		Description:     thingTypes.Description,
		Characteristics: req.Characteristics,
		PropertySets:    req.PropertySets,
		Meta: &pb.Meta{
			Id:            thingTypes.Meta.Id,
			Version:       thingTypes.Meta.Version,
			Tenant:        thingTypes.Meta.Tenant,
			CreatedById:   thingTypes.Meta.CreatedById,
			UpdatedById:   thingTypes.Meta.UpdatedById,
			CreatedByName: thingTypes.Meta.CreatedByName,
			UpdatedByName: thingTypes.Meta.UpdatedByName,
			CreatedTime:   timestamppb.New(thingTypes.Meta.CreatedTime),
			UpdatedTime:   timestamppb.New(thingTypes.Meta.UpdatedTime),
		},
	}, nil
}
func (s *ThingTypesService) UpdateThingTypes(ctx context.Context, req *pb.UpdateThingTypesRequest) (*pb.UpdateThingTypesReply, error) {
	return &pb.UpdateThingTypesReply{}, nil
}
func (s *ThingTypesService) DeleteThingTypes(ctx context.Context, req *pb.DeleteThingTypesRequest) (*pb.DeleteThingTypesReply, error) {
	return &pb.DeleteThingTypesReply{}, nil
}
func (s *ThingTypesService) GetThingTypes(ctx context.Context, req *pb.GetThingTypesRequest) (*pb.GetThingTypesReply, error) {
	return &pb.GetThingTypesReply{}, nil
}
func (s *ThingTypesService) ListThingTypes(ctx context.Context, req *pb.ListThingTypesRequest) (*pb.ListThingTypesReply, error) {
	return &pb.ListThingTypesReply{}, nil
}
