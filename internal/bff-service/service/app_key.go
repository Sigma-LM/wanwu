package service

import (
	"net/url"

	app_service "github.com/UnicomAI/wanwu/api/proto/app-service"
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/constant"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

func GetAppBaseUrl(ctx *gin.Context, req request.GetAppBaseUrlRequest) (string, error) {
	if req.AppType == constant.AppTypeWorkflow {
		apiBaseUrl, err := url.JoinPath(config.Cfg().Server.ApiBaseUrl, "/openapi/v1", req.AppType, "/run")
		if err != nil {
			return "", grpc_util.ErrorStatus(err_code.Code_BFFGeneral, err.Error())
		}
		return apiBaseUrl, nil
	}
	apiBaseUrl, err := url.JoinPath(config.Cfg().Server.ApiBaseUrl, "/openapi/v1", req.AppType, "/chat")
	if err != nil {
		return "", grpc_util.ErrorStatus(err_code.Code_BFFGeneral, err.Error())
	}
	return apiBaseUrl, nil
}

func GetAppKeyByKey(ctx *gin.Context, appKey string) (*app_service.AppKeyInfo, error) {
	return app.GetAppKeyByKey(ctx.Request.Context(), &app_service.GetAppKeyByKeyReq{AppKeyId: appKey})
}

func GenAppKey(ctx *gin.Context, userId, orgId string, req request.GenAppKeyRequest) (*response.AppResponse, error) {
	key, err := app.GenAppKey(ctx.Request.Context(), &app_service.GenAppKeyReq{
		AppId:   req.AppId,
		AppType: req.AppType,
		UserId:  userId,
		OrgId:   orgId,
	})
	if err != nil {
		return nil, err
	}
	return &response.AppResponse{
		ApiID:     key.AppKeyId,
		ApiKey:    key.AppKey,
		CreatedAt: util.Time2Str(key.CreatedAt),
	}, nil
}

func DelAppKey(ctx *gin.Context, req request.DelAppKeyRequest) error {
	_, err := app.DelAppKey(ctx.Request.Context(), &app_service.DelAppKeyReq{
		AppKeyId: req.ApiId,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetAppKeyList(ctx *gin.Context, userId string, req request.GetAppKeyListRequest) ([]*response.AppResponse, error) {
	apiKeyList, err := app.GetAppKeyList(ctx.Request.Context(), &app_service.GetAppKeyListReq{
		AppId:   req.AppId,
		AppType: req.AppType,
		UserId:  userId,
	})
	if err != nil {
		return nil, err
	}
	var apiRes []*response.AppResponse
	for _, apiKeyInfo := range apiKeyList.Info {
		apiRes = append(apiRes, toAppResp(apiKeyInfo))
	}
	return apiRes, nil
}

func toAppResp(apiKeyInfo *app_service.AppKeyInfo) *response.AppResponse {
	return &response.AppResponse{
		ApiID:     apiKeyInfo.AppKeyId,
		ApiKey:    apiKeyInfo.AppKey,
		CreatedAt: util.Time2Str(apiKeyInfo.CreatedAt),
	}
}
