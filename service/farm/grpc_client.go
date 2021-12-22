package farm

import (
	"feimumoke/farming/v2/api/server"
	"feimumoke/farming/v2/framework"
	"feimumoke/farming/v2/framework/boot"
	"feimumoke/farming/v2/framework/grpc_util"
	"fmt"
)

var GrpcCli *GrpcClient

type GrpcClient struct {
	FarmerClient server.FarmerServiceClient
	GroundClient server.GroundServiceClient
}

func InitGrpcClient() {
	boot.AddAppGrpcClient([]*grpc_util.ClientMeta{{
		Target:      framework.SERVER_FRAMER_TARGET,
		ProjectName: framework.SERVER_FRAMER_PROJECT,
		ServiceName: framework.SERVER_FRAMER_SERVICE,
		InitFunc:    server.NewFarmerServiceClient,
	}, {
		Target:      framework.SERVER_GROUND_TARGET,
		ProjectName: framework.SERVER_GROUND_PROJECT,
		ServiceName: framework.SERVER_GROUND_SERVICE,
		InitFunc:    server.NewGroundServiceClient,
	}})
	GrpcCli = &GrpcClient{
		FarmerClient: boot.GetAppGrpcClient(fmt.Sprintf("%s-%s", framework.SERVER_FRAMER_PROJECT, framework.SERVER_FRAMER_SERVICE)).(server.FarmerServiceClient),
		GroundClient: boot.GetAppGrpcClient(fmt.Sprintf("%s-%s", framework.SERVER_GROUND_PROJECT, framework.SERVER_GROUND_SERVICE)).(server.GroundServiceClient),
	}
}
