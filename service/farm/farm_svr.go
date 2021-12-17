package farm

import (
	"context"
	"feimumoke/farming/v2/api/server"
	"feimumoke/farming/v2/api/service"
	"feimumoke/farming/v2/framework/util"
)

type FarmSvr struct {
	*service.UnimplementedFarmServiceServer
	farmerSv server.FarmerServiceClient
	groundSv server.GroundServiceClient
}

func (f FarmSvr) PlantTree(ctx context.Context, req *service.PlantTreeReq) (*service.PlantTreeRsp, error) {
	return &service.PlantTreeRsp{Result: util.StructToString(req) + ":" + util.GetHostString()}, nil
}
