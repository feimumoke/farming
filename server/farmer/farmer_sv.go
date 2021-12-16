package farmer

import (
	"context"
	"feimumoke/farming/v2/api/server"
)

type FarmerSv struct {
}

func (f FarmerSv) SelectFarmer(ctx context.Context, req *server.FarmerReq) (*server.FarmerRsp, error) {
	return &server.FarmerRsp{Result: "Farmer from Server:" + req.Name}, nil
}

func (f FarmerSv) mustEmbedUnimplementedFarmerServiceServer() {
	panic("implement me")
}
