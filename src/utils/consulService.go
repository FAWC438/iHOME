package utils

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"iHome/src/conf"
)

// GetMicroClientFromConsul
//
//	@Description: 注册 consul 服务发现，从 consul 上找 go-micro 服务
//	@return client.Client
func GetMicroClientFromConsul() client.Client {
	consulReg := consul.NewRegistry(registry.Addrs(conf.ConsulAddr))
	srv := micro.NewService()
	srv.Init(
		micro.Client(client.NewClient()),
		micro.Registry(consulReg), // 将服务注册到 consul
	)
	return srv.Client()
}
