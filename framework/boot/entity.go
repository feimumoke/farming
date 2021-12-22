package boot

import (
	"context"
	"feimumoke/farming/v2/framework/grpc_util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

var AppCtx *AppContext

type GrpcRegFunc func(server *grpc.Server)
type GateWayRegFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

type AppContext struct {
	grpcClients *grpc_util.GRPCClients
	boot        *Bootstrapper
}

type grpcEntity struct {
	Port     string
	Opts     []grpc.ServerOption
	RegFuncs []GrpcRegFunc
	server   *grpc.Server
	listener net.Listener
}

type httpEntity struct {
	HttpPort        string
	GrpcPort        string
	GateWayRegFuncs []GateWayRegFunc
	Opts            []grpc.DialOption
	MuxOpts         []runtime.ServeMuxOption
	server          *http.Server
}
