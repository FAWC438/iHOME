package main

import (
	"getCaptcha/conf"
	"getCaptcha/handler"
	pb "getCaptcha/proto"
	"go-micro.dev/v4/registry"

	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "getCaptcha"
	version = "latest"
)

func main() {
	consulReg := consul.NewRegistry(registry.Addrs(conf.ConsulAddr))
	// Create service
	srv := micro.NewService()
	srv.Init(
		micro.Address("127.0.0.1:9080"), // 指定微服务端口
		micro.Name(service),
		micro.Registry(consulReg),
		micro.Version(version),
	)

	// Register handler
	if err := pb.RegisterGetCaptchaHandler(srv.Server(), new(handler.GetCaptcha)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
