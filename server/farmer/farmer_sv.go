package farmer

import (
	"context"
	"encoding/json"
	"feimumoke/farming/v2/api/server"
	"feimumoke/farming/v2/framework/util"
	"fmt"
)

func NewFarmerSv() *FarmerSv {
	return &FarmerSv{}
}

type FarmerSv struct {
	*server.UnimplementedFarmerServiceServer
}

func (f *FarmerSv) SelectFarmer(ctx context.Context, req *server.FarmerReq) (*server.FarmerRsp, error) {
	//grpclog.Infof("SelectFarmer req: [%v]", util.StructToString(req))
	return &server.FarmerRsp{Result: "Farmer from Server:" + req.Name + ":" + util.GetHostString()}, nil
}

func (f *FarmerSv) HelpFarmer(ctx context.Context, req *server.HelpFarmerReq) (*server.HelpFarmerRsp, error) {
	fmt.Println(ctx.Value("a"))
	return &server.HelpFarmerRsp{Result: "Farmer Help Server:" + ":" + util.GetHostString()}, nil
}

func Test() {
	mi := make([]*MethodInfo, 0)
	mi = append(mi, &MethodInfo{
		Name: "SelectFarmer",
		Data: `{"Name":"Zhangsan"}`,
	})

	mi = append(mi, &MethodInfo{
		Name: "HelpFarmer",
		Data: `{"Name":"Zhangsan"}`,
	})

	ctx := context.Background()
	ctx = context.WithValue(ctx, "a", "ffvvv")
	s := &FarmerSv{}
	for _, m := range mi {
		ApplyFunc(m, ctx, s)
	}

}

type MethodInfo struct {
	Name string
	Data string
}

func ApplyFunc(ms *MethodInfo, ctx context.Context, s *FarmerSv) {
	if ms.Name == "SelectFarmer" {
		var data server.FarmerReq
		json.Unmarshal([]byte(ms.Data), &data)
		farmer, _ := s.SelectFarmer(ctx, &data)
		fmt.Println(farmer)
	}
	if ms.Name == "HelpFarmer" {
		var data server.HelpFarmerReq
		json.Unmarshal([]byte(ms.Data), &data)
		farmer, _ := s.HelpFarmer(ctx, &data)
		fmt.Println(farmer)
	}
}
