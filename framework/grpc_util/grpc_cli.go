package grpc_util

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"reflect"
	"sync"
	"time"
)

var ClientKeepAlive = keepalive.ClientParameters{
	// send pings every 10 seconds if there is no activity
	Time: 10 * time.Second,
	// wait 1 second for ping ack before considering the connection dead
	Timeout: time.Second,
	// send pings even without active streams
	PermitWithoutStream: true,
}

type GRPCClients struct {
	Clients map[string]interface{}
	Mutex   sync.Mutex
}

type ClientMeta struct {
	// Service target,ip:port,consul address and other service discovery methods
	Target string
	// Project name which the service belongs to
	ProjectName string
	// GRPC service name defined in protobuf
	ServiceName string
	// GRPC client init func
	InitFunc interface{}
}

func (c *GRPCClients) GetClient(serviceName string) interface{} {
	if v, ok := c.Clients[serviceName]; ok {
		return v
	} else {
		return nil
	}
}

func (c *GRPCClients) AddClients(clientConfigs []*ClientMeta) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	defaultInterceptor := func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		msg := fmt.Sprintf("Target: [%v]", cc.Target())
		err := invoker(ctx, method, req, resp, cc, opts...)
		grpclog.Infof("GrpcClientReq [%v] With [%v]", method, msg)
		return err
	}
	opts := []grpc.DialOption{grpc.WithUnaryInterceptor(ChainUnaryClient(defaultInterceptor))}
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithKeepaliveParams(ClientKeepAlive))

	for _, clientConf := range clientConfigs {
		svcKey := fmt.Sprintf("%s-%s", clientConf.ProjectName, clientConf.ServiceName)
		if _, ok := c.Clients[svcKey]; ok {
			continue
		}
		target := clientConf.Target
		conn, err := grpc.DialContext(context.Background(), target, opts...)
		if err != nil {
			grpclog.Error(fmt.Sprintf("Init %s grpc client  error:%s", clientConf.ServiceName, err.Error()))
			continue
		}
		retValue, err := call(clientConf.InitFunc, conn)
		if err != nil {
			grpclog.Error(fmt.Sprintf("call client init func error: %s", err.Error()))
			continue
		}
		client := retValue[0].Interface()
		c.Clients[svcKey] = client
	}
}

func call(f interface{}, params ...interface{}) (result []reflect.Value, err error) {
	fun := reflect.ValueOf(f)
	if fun.Kind() != reflect.Func {
		err = errors.New("non-callable function")
		return
	}
	if len(params) != fun.Type().NumIn() {
		err = errors.New("the number of params is not adapted")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = fun.Call(in)
	return
}
