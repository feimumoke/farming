package boot

import "feimumoke/farming/v2/framework/grpc_util"

func GetAppGrpcClient(serviceName string) interface{} {
	return AppCtx.grpcClients.GetClient(serviceName)
}

func AddAppGrpcClient(clientConfigs []*grpc_util.ClientMeta) {
	AppCtx.grpcClients.AddClients(clientConfigs)
}
