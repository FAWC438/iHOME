package main

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4/registry"
	"user/conf"
	"user/handler"
	"user/model"
	pb "user/proto"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "user"
	version = "latest"
)

func main() {
	// 初始化 MySQL 表
	_, err := model.InitMySQL()
	if err != nil {
		return
	}

	consulReg := consul.NewRegistry(registry.Addrs(conf.ConsulAddr))
	// Create service
	srv := micro.NewService()
	srv.Init(
		micro.Address("127.0.0.1:9081"), // 指定微服务端口
		micro.Name(service),
		micro.Registry(consulReg),
		micro.Version(version),
	)

	// Register handler
	if err := pb.RegisterUserHandler(srv.Server(), new(handler.User)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
