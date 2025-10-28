package service

import (
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	mcp_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/mcp-util"
	"github.com/UnicomAI/wanwu/pkg/constant"
	"github.com/gin-gonic/gin"
)

var customToolSchema = &CustomToolSchema{}

func init() {
	AddSchemaBuilder(customToolSchema)
}

type CustomToolSchema struct {
}

func (CustomToolSchema) ToolType() string {
	return constant.MCPServerToolTypeCustomTool
}

func (CustomToolSchema) BuildMcpSchema(ctx *gin.Context, toolInfo *ToolInfo) (*mcp_util.OpenAPISchema, error) {
	customToolInfo, err := mcp.GetCustomToolInfo(ctx, &mcp_service.GetCustomToolInfoReq{
		CustomToolId: toolInfo.Id,
	})
	if err != nil {
		return nil, err
	}
	return &mcp_util.OpenAPISchema{
		Schema: customToolInfo.Schema,
		ApiAuth: &mcp_util.APIAuth{
			Type:  customToolInfo.ApiAuth.Type,
			In:    "header",
			Name:  customToolInfo.ApiAuth.CustomHeaderName,
			Value: customToolInfo.ApiAuth.ApiKey,
		},
		MethodNames: toolInfo.MethodNames,
	}, nil
}
