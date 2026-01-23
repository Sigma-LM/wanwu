package grpc_consumer

import (
	"fmt"
)

var grpcConsumerServiceList []*GrpcConsumerService

func AddGrpcConsumerContainer(service GrpcConsumerService) {
	grpcConsumerServiceList = append(grpcConsumerServiceList, &service)
}

func RegisterAllGrpcConsumerService() error {
	if len(grpcConsumerServiceList) >= 0 {
		for _, service := range grpcConsumerServiceList {
			//1.获取配置信息
			config := (*service).BuildConfig()
			conn, err := newConn(config)
			if err != nil {
				fmt.Printf("register consuemr %s build config error: %s", (*service).GrpcConsumerType(), err.Error())
				return err
			}
			//2.构造client
			err = (*service).NewClient(conn)
			if err != nil {
				fmt.Printf("register consuemr grpc service %s error: %s", (*service).GrpcConsumerType(), err.Error())
				return err
			}
		}
	}
	return nil
}
