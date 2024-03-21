package pkg

import (
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"os"
	"os/signal"
	"qq-guild-bot/internal/pkg/config"
	"syscall"
	"time"
)

func RegisterConsul(s *grpc.Server) error {
	consulConfig := consulApi.DefaultConfig()
	consul, _ := consulApi.NewClient(consulConfig)
	serviceId := config.AppConf.ServiceId
	registration := &consulApi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    config.AppConf.ServiceName,
		Port:    config.AppConf.GrpcPort,
		Address: config.AppConf.ConsulHost,
		Check: &consulApi.AgentServiceCheck{
			GRPC:     fmt.Sprintf("%s:%d", "host.docker.internal", config.AppConf.GrpcPort),
			Interval: "10s",
			Timeout:  "30s",
		},
	}

	// 创建健康检查服务的实例
	healthServer := health.NewServer()

	// 注册健康检查服务到 gRPC 服务器
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	// 可以设置服务的健康状态
	healthServer.SetServingStatus(config.AppConf.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	if err := consul.Agent().ServiceRegister(registration); err != nil {
		return err
	}

	// 监听系统信号
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		fmt.Println("注销 Consul 服务并退出...")
		_ = consul.Agent().ServiceDeregister(serviceId)

		// 延时以完成任何正在进行的操作
		time.Sleep(time.Second * 5)
		os.Exit(0)
	}()

	return nil
}
