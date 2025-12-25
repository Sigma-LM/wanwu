package app

import (
	app_service "github.com/UnicomAI/wanwu/api/proto/app-service"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/app-service/client"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
)

type Service struct {
	app_service.UnimplementedAppServiceServer
	cli client.IClient
}

func NewService(cli client.IClient) *Service {
	return &Service{
		cli: cli,
	}
}

func errStatus(code errs.Code, status *errs.Status) error {
	return grpc_util.ErrorStatusWithKey(code, status.TextKey, status.Args...)
}

type iReq interface {
	GetPageNo() int32 // 从1开始
	GetPageSize() int32
}

func toOffset(req iReq) int32 {
	if req.GetPageNo() < 1 || req.GetPageSize() < 0 {
		return -1
	}
	return (req.GetPageNo() - 1) * req.GetPageSize()
}
