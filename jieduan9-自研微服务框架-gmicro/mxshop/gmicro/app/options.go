package app

import (
	"mxshop/gmicro/registry"
	"mxshop/gmicro/server/restserver"
	"mxshop/gmicro/server/rpcserver"
	"net/url"
	"os"
	"time"
)

type Option func(o *options)

type options struct {
	id        string
	endpoints []*url.URL // url类型
	name      string

	sigs []os.Signal

	//允许用户传入自己的注册的服务实现
	// gmicro 不把注册中心写死成 consul，
	// 它只规定：你给我一个注册器，我在启动和停止时调用它。
	registrar        registry.Registrar // 注册器
	registrarTimeout time.Duration      // 注册器的超时时间

	//stop超时时间
	stopTimeout time.Duration

	restServer *restserver.Server
	rpcServer  *rpcserver.Server
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithRPCServer(server *rpcserver.Server) Option {
	return func(o *options) {
		o.rpcServer = server
	}
}

func WithRestServer(server *restserver.Server) Option {
	return func(o *options) {
		o.restServer = server
	}
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithSigs(sigs []os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}
