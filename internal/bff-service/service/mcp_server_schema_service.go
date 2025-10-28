package service

import (
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	mcp_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/mcp-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/gin-gonic/gin"
)

var mcpSchemaBuilderMap = make(map[string]McpSchemaBuilder)

func AddSchemaBuilder(builder McpSchemaBuilder) {
	mcpSchemaBuilderMap[builder.ToolType()] = builder
}

type McpSchema struct {
	McpTool       []*mcp_util.McpTool
	McpServerTool []*mcp_service.MCPServerToolInfo
}

type ToolInfo struct {
	ToolType    string
	Id          string
	Name        string
	MethodNames []string
	Schema      string
	ApiAuth     request.CustomToolApiAuthWebRequest
}

type McpSchemaBuilder interface {
	ToolType() string
	BuildMcpSchema(ctx *gin.Context, toolInfo *ToolInfo) (*mcp_util.OpenAPISchema, error)
}

// CreateMcpSchema 创建mcp schema
func CreateMcpSchema(ctx *gin.Context, mcpServerId string, toolInfo *ToolInfo) (*McpSchema, error) {
	builder := mcpSchemaBuilderMap[toolInfo.ToolType]
	//构造schema
	mcpSchema, err := builder.BuildMcpSchema(ctx, toolInfo)
	if err != nil {
		return nil, err
	}
	//构造mcp tool
	mcpTools, err := mcp_util.OpenApiToMcpToolList(mcpSchema)
	if err != nil {
		return nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_convert_tool_err", err.Error())
	}

	//调用底层的mcpTool
	mcpServerTools, err := buildMcpTools(mcpTools, mcpServerId, toolInfo, mcpSchema)
	if err != nil {
		return nil, err
	}

	return &McpSchema{
		McpTool:       mcpTools,
		McpServerTool: mcpServerTools,
	}, nil
}

// UpdateMcpSchema 更新mcp schema
func UpdateMcpSchema(req request.MCPServerToolUpdateReq, tool *mcp_service.MCPServerToolInfo) (*McpSchema, error) {
	mcpTool := &mcp_util.McpTool{}
	doc, err := mcp_util.ParseOpenApiContent(tool.Schema)
	if err != nil {
		return nil, err
	}
	op, _, _ := mcp_util.ParseOpenApiOperation(doc, tool.Name)
	op.OperationID = req.MethodName
	op.Description = req.Desc
	mcpTool.Tool = mcp_util.ConvertMcpTool(doc, req.MethodName)
	mcpTool.Handle = mcp_util.ConvertMcpHandler(doc, req.MethodName, buildAPIAuth(tool.ApiAuth))
	content := mcp_util.BuildOpenApiContent(doc, op.OperationID)
	mcpServerTool := &mcp_service.MCPServerToolInfo{
		McpServerToolId: req.MCPServerToolID,
		Name:            req.MethodName,
		Desc:            req.Desc,
		Schema:          content,
	}
	return &McpSchema{
		McpTool:       []*mcp_util.McpTool{mcpTool},
		McpServerTool: []*mcp_service.MCPServerToolInfo{mcpServerTool},
	}, nil
}

func buildMcpTools(mcpTools []*mcp_util.McpTool, mcpServerId string, toolInfo *ToolInfo, schema *mcp_util.OpenAPISchema) ([]*mcp_service.MCPServerToolInfo, error) {
	var mcpServerTools []*mcp_service.MCPServerToolInfo
	for _, tool := range mcpTools {
		mcpServerTools = append(mcpServerTools, &mcp_service.MCPServerToolInfo{
			McpServerId: mcpServerId,
			AppToolId:   toolInfo.Id,
			Type:        toolInfo.ToolType,
			AppToolName: toolInfo.Name,
			Name:        tool.Tool.Name,
			Desc:        tool.Tool.Description,
			Schema:      tool.ToolSchema,
			ApiAuth: &mcp_service.ApiAuth{
				AuthType:  schema.ApiAuth.Type,
				AuthIn:    schema.ApiAuth.In,
				AuthName:  schema.ApiAuth.Name,
				AuthValue: schema.ApiAuth.Value,
			},
		})
	}
	return mcpServerTools, nil
}

func buildAPIAuth(apiAuth *mcp_service.ApiAuth) *mcp_util.APIAuth {
	return &mcp_util.APIAuth{
		Type:  apiAuth.AuthType,
		In:    apiAuth.AuthIn,
		Name:  apiAuth.AuthName,
		Value: apiAuth.AuthValue,
	}
}
