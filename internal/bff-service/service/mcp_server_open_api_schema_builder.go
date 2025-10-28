package service

import (
	mcp_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/mcp-util"
	"github.com/UnicomAI/wanwu/pkg/constant"
	"github.com/gin-gonic/gin"
)

var openApiSchema = &OpenApiSchema{}

func init() {
	AddSchemaBuilder(openApiSchema)
}

type OpenApiSchema struct {
}

func (OpenApiSchema) ToolType() string {
	return constant.MCPServerToolTypeOpenAPI
}

func (OpenApiSchema) BuildMcpSchema(ctx *gin.Context, toolInfo *ToolInfo) (*mcp_util.OpenAPISchema, error) {
	return &mcp_util.OpenAPISchema{
		Schema: toolInfo.Schema,
		ApiAuth: &mcp_util.APIAuth{
			Type:  toolInfo.ApiAuth.Type,
			In:    "header",
			Name:  toolInfo.ApiAuth.CustomHeaderName,
			Value: toolInfo.ApiAuth.APIKey,
		},
		MethodNames: toolInfo.MethodNames,
	}, nil
}
