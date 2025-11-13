package service

import (
	"fmt"
	net_url "net/url"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/constant"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func CreateChatflow(ctx *gin.Context, orgID, name, desc, iconUri string) (*response.CozeWorkflowIDData, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.CreateUri)
	ret := &response.CozeWorkflowIDResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"space_id":  orgID,
			"name":      name,
			"desc":      desc,
			"icon_uri":  iconUri,
			"flow_mode": "3",
		}).
		SetResult(ret).
		Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_create", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_create", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_create", fmt.Sprintf("code %v msg %v", ret.Code, ret.Msg))
	}
	return ret.Data, nil
}

func cozeChatflowInfo2Model(chatflowInfo *response.CozeWorkflowListDataWorkflow) response.AppBriefInfo {
	return response.AppBriefInfo{
		AppId:     chatflowInfo.WorkflowId,
		AppType:   constant.AppTypeChatflow,
		Name:      chatflowInfo.Name,
		Desc:      chatflowInfo.Desc,
		Avatar:    cacheWorkflowAvatar(chatflowInfo.URL, constant.AppTypeChatflow),
		CreatedAt: util.Time2Str(chatflowInfo.CreateTime * 1000),
		UpdatedAt: util.Time2Str(chatflowInfo.UpdateTime * 1000),
	}
}
