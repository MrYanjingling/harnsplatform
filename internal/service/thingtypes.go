package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "harnsplatform/api/modelmanager/v1"
	"harnsplatform/internal/biz"
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

	tt := &biz.ThingTypes{
		Name:            req.GetName(),
		ParentTypeId:    req.GetParentTypeId(),
		Description:     req.GetDescription(),
		Characteristics: map[string]interface{}{},
		PropertySets:    map[string]interface{}{},
		Meta: biz.Meta{
			Id: "11121",
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

	s.ttu.CreateThingTypes(ctx, tt)

	return &pb.CreateThingTypesReply{}, nil
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
