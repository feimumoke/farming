package farmer

import (
	"context"
	"feimumoke/farming/v2/api/server"
	"feimumoke/farming/v2/framework/util"
)

func NewFarmerSv() *FarmerSv {
	return &FarmerSv{}
}

type FarmerSv struct {
	*server.UnimplementedFarmerServiceServer
}

func (f *FarmerSv) SelectFarmer(ctx context.Context, req *server.FarmerReq) (*server.FarmerRsp, error) {
	return &server.FarmerRsp{Result: "Farmer from Server:" + req.Name + ":" + util.GetHostString()}, nil
}
