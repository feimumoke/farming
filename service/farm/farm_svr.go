package farm

import (
	"context"
	"feimumoke/farming/v2/api/server"
	"feimumoke/farming/v2/api/service"
	"feimumoke/farming/v2/framwork/util"
)

type FarmSvr struct {
	farmerSv server.FarmerServiceClient
	groundSv server.GroundServiceClient
}

func (f FarmSvr) PlantTree(ctx context.Context, req *service.PlantTreeReq) (*service.PlantTreeRsp, error) {
	return &service.PlantTreeRsp{Result: util.StructToString(req)}, nil
}

func (f FarmSvr) mustEmbedUnimplementedFarmServiceServer() {
	panic("implement me")
}
