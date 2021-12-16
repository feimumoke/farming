package identify

import (
	"context"
	"errors"
	"feimumoke/farming/v2/api/server"
	"feimumoke/farming/v2/api/service"
	"feimumoke/farming/v2/framwork/util"
	"sync"
	"time"
)

type Farmer struct {
	name     string
	initTime string
	password string `json:"-"`
}

type IdentifySvr struct {
	farmerSv server.FarmerServiceClient
	groundSv server.GroundServiceClient
	userMap  sync.Map
}

func (i IdentifySvr) Register(ctx context.Context, req *service.RegisterReq) (*service.RegisterRsp, error) {
	if user, ok := i.userMap.Load(req.Name); ok {
		farmer := user.(*Farmer)
		if req.PassWord != farmer.password {
			return nil, errors.New("user name or password not right")
		}
		return &service.RegisterRsp{Result: "ExistUser:" + util.StructToString(user)}, nil
	}
	farmer, err := i.farmerSv.SelectFarmer(ctx, &server.FarmerReq{
		Name:     req.Name,
		PassWord: req.PassWord,
	})
	if err != nil {
		return nil, err
	}
	ground, err := i.groundSv.SelectGround(ctx, &server.GroundReq{
		Kind:  "INIT",
		Owner: req.Name,
	})
	if err != nil {
		return nil, err
	}
	i.userMap.Store(req.Name, &Farmer{
		name:     req.Name,
		initTime: time.Now().Format(util.DefaultTimeLayOut),
		password: req.PassWord,
	})
	return &service.RegisterRsp{Result: util.StructToString(farmer) + "##" + util.StructToString(ground)}, nil

}

func (i IdentifySvr) mustEmbedUnimplementedIdentifyServiceServer() {
	panic("implement me")
}
