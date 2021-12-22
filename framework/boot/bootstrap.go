package boot

import (
	"context"
	"feimumoke/farming/v2/framework"
	"feimumoke/farming/v2/framework/grpc_util"
	"feimumoke/farming/v2/framework/util"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ServerKeepAlivePolicy = keepalive.EnforcementPolicy{
		// If a client pings more than once every 5 seconds, terminate the connection
		MinTime: 5 * time.Second,
		// Allow pings even when there are no active streams
		PermitWithoutStream: true,
	}
	ServerKeepAlive = keepalive.ServerParameters{
		// If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionIdle: 15 * time.Second,
		// If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAge: 30 * time.Second,
		// Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		MaxConnectionAgeGrace: 5 * time.Second,
		// Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Time: 5 * time.Second,
		// Wait 1 second for the ping ack before assuming the connection is dead
		Timeout: 1 * time.Second,
	}
)

type Bootstrapper struct {
	grpcEntities  map[string]*grpcEntity
	httpEntities  map[string]*httpEntity
	shutDownHooks map[string]func()
	startTime     time.Time
	mutex         sync.Mutex
}

func NewBootstrapper() *Bootstrapper {
	b := &Bootstrapper{mutex: sync.Mutex{}}
	b.grpcEntities = make(map[string]*grpcEntity)
	b.httpEntities = make(map[string]*httpEntity)
	return b
}

func (b *Bootstrapper) Bootstrap() func(duration time.Duration) {
	b.startTime = time.Now()
	b.initLog()
	grpclog.Infof("start to init grpc")
	b.bootGrpc()
	grpclog.Infof("start to init http")
	b.bootGateway()
	grpclog.Infof("Bootstrap success")
	AppCtx = &AppContext{
		grpcClients: &grpc_util.GRPCClients{Clients: make(map[string]interface{}), Mutex: sync.Mutex{}},
		boot:        b,
	}
	return b.shutdown
}

func (b *Bootstrapper) RegisterGrpcFunc(port string, funcs ...GrpcRegFunc) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	entity, ok := b.grpcEntities[port]
	if !ok {
		entity = &grpcEntity{
			Port:     port,
			Opts:     nil,
			RegFuncs: []GrpcRegFunc{},
		}
	}
	entity.RegFuncs = append(entity.RegFuncs, funcs...)
	b.grpcEntities[port] = entity
}

func (b *Bootstrapper) RegisterGateWayFunc(port, grpcPort string, funcs ...GateWayRegFunc) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	entity, ok := b.httpEntities[port]
	if !ok {
		entity = &httpEntity{
			HttpPort:        port,
			GrpcPort:        grpcPort,
			Opts:            nil,
			GateWayRegFuncs: []GateWayRegFunc{},
		}
	}
	entity.GateWayRegFuncs = append(entity.GateWayRegFuncs, funcs...)
	b.httpEntities[port] = entity
}

func (b *Bootstrapper) initLog() {
	proj := os.Getenv("PROJECT_NAME")
	if proj == "" {
		proj = "app"
	}
	file, err := os.OpenFile(framework.LogPath+proj+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
	grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(file, file, file, 99))
}

func (b *Bootstrapper) bootGrpc() {
	unaries := make([]grpc.UnaryServerInterceptor, 0)
	serverUnaryLogInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		now := time.Now().UnixNano()
		msg := ""
		if peer, ok := peer.FromContext(ctx); ok {
			network := peer.Addr.Network()
			msg += fmt.Sprintf("remoteNetwork:[%v],Address [%v]", network, peer.Addr.String())
		}
		if headers, ok := metadata.FromIncomingContext(ctx); ok {
			msg += fmt.Sprintf("Headers: [%v]", util.StructToString(headers))
		}
		resp, err = handler(ctx, req)
		grpclog.Infof("FarmingGrpc Res [%v]; Resp [%v] err [%v] with msg [%v],cost [%v] nanoseconds", util.StructToString(req), util.StructToString(resp), err, msg, time.Now().UnixNano()-now)
		return resp, err
	}
	unaries = append(unaries, serverUnaryLogInterceptor)
	streams := make([]grpc.StreamServerInterceptor, 0)
	serverStreamLogInterceptor := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		msg := ""
		if peer, ok := peer.FromContext(stream.Context()); ok {
			network := peer.Addr.Network()
			msg += fmt.Sprintf("remoteNetwork:[%v],Address [%v]", network, peer.Addr.String())
		}
		if headers, ok := metadata.FromIncomingContext(stream.Context()); ok {
			msg += fmt.Sprintf("Headers: [%v]", util.StructToString(headers))
		}
		err := handler(srv, stream)
		grpclog.Infof("FarmingGrpcStream Res [%v]; err [%v],with msg [%v]", info.FullMethod, err, msg)
		return err
	}
	streams = append(streams, serverStreamLogInterceptor)

	for _, entity := range b.grpcEntities {
		listener, err := net.Listen("tcp4", "0.0.0.0:"+entity.Port)
		if err != nil {
			log.Fatal(err)
		}
		entity.listener = listener
		options := []grpc.ServerOption{grpc.UnaryInterceptor(grpc_util.ChainUnaryServer(unaries...))}
		options = append(options, grpc.StreamInterceptor(grpc_util.ChainStreamServer(streams...)))
		options = append(options, grpc.KeepaliveEnforcementPolicy(ServerKeepAlivePolicy), grpc.KeepaliveParams(ServerKeepAlive))
		entity.Opts = options
		entity.server = grpc.NewServer(entity.Opts...)
		for _, reg := range entity.RegFuncs {
			reg(entity.server)
		}
		go func(entity *grpcEntity) {
			if err := entity.server.Serve(entity.listener); err != nil {
				log.Fatal(err)
			}
		}(entity)
		grpclog.Infof("start grpc for %v success", entity.Port)
	}
}

func (b *Bootstrapper) bootGateway() {

	for _, entity := range b.httpEntities {
		gRPCEndpoint := "0.0.0.0:" + entity.GrpcPort
		httpEndpoint := "0.0.0.0:" + entity.HttpPort
		if _, ok := b.grpcEntities[entity.GrpcPort]; !ok {
			grpclog.Errorf("not have grpc port %v", entity.GrpcPort)
			continue
		}
		ctx := context.Background()
		gwMux := runtime.NewServeMux(DefaultHTTPServeMuxOpt()...)
		entity.Opts = []grpc.DialOption{grpc.WithInsecure()}
		entity.server = &http.Server{
			Addr:    httpEndpoint,
			Handler: gwMux,
		}
		for _, reg := range entity.GateWayRegFuncs {
			if err := reg(ctx, gwMux, gRPCEndpoint, entity.Opts); err != nil {
				log.Fatal(err)
			}
		}

		go func(entity *httpEntity) {
			if err := entity.server.ListenAndServe(); err != nil {
				log.Fatal(err)
			}
		}(entity)
		grpclog.Infof("start http for %v success", entity.HttpPort)
	}
}

func DefaultHTTPServeMuxOpt() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithOutgoingHeaderMatcher(func(s string) (string, bool) {
			return s, true
		}),
		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			return s, true
		}),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			meta := make(map[string]string)
			meta["lang"] = "zh"
			return metadata.New(meta)
		}),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions:   protojson.MarshalOptions{},
			UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
		}),
	}
}

func (b *Bootstrapper) shutdown(duration time.Duration) {
	sig := <-b.registerShutdownHook()
	grpclog.Infof("shutdown by %v", sig.String())
	for name, f := range b.shutDownHooks {
		f()
		grpclog.Infof("execute shutDownHook [%v]", name)
	}
	for _, entity := range b.httpEntities {
		entity.server.Shutdown(context.Background())
	}
	for _, entity := range b.grpcEntities {
		entity.server.GracefulStop()
	}
	grpclog.Infof("all shutDownHook executed")
	time.Sleep(duration)
	os.Exit(0)
}

func (b *Bootstrapper) registerShutdownHook() chan os.Signal {
	shutdownHook := make(chan os.Signal)
	signal.Notify(shutdownHook,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		os.Interrupt)
	return shutdownHook
}
