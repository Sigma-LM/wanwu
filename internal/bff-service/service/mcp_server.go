package service

import (
	"bytes"
	"context"
	"io"

	app_service "github.com/UnicomAI/wanwu/api/proto/app-service"
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	mcp_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/mcp-util"
	"github.com/UnicomAI/wanwu/pkg/constant"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	openapi3_util "github.com/UnicomAI/wanwu/pkg/openapi3-util"
	"github.com/gin-gonic/gin"
)

func CreateMCPServer(ctx *gin.Context, userID, orgID string, req request.MCPServerCreateReq) error {
	resp, err := mcp.CreateMCPServer(ctx.Request.Context(), &mcp_service.CreateMCPServerReq{
		Name:       req.Name,
		Desc:       req.Desc,
		AvatarPath: req.Avatar.Key,
		Identity: &mcp_service.Identity{
			OrgId:  orgID,
			UserId: userID,
		},
	})
	if err != nil {
		return err
	}
	_, err = app.GenApiKey(ctx.Request.Context(), &app_service.GenApiKeyReq{
		AppId:   resp.McpServerId,
		AppType: constant.AppTypeMCPServer,
		UserId:  userID,
		OrgId:   orgID,
	})
	if err != nil {
		return err
	}
	err = mcp_util.StartMCPServer(ctx, resp.McpServerId)
	if err != nil {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_start_err", err.Error())
	}
	return err
}

func UpdateMCPServer(ctx *gin.Context, req request.MCPServerUpdateReq) error {
	_, err := mcp.UpdateMCPServer(ctx.Request.Context(), &mcp_service.UpdateMCPServerReq{
		McpServerId: req.MCPServerID,
		Name:        req.Name,
		Desc:        req.Desc,
		AvatarPath:  req.Avatar.Key,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetMCPServerDetail(ctx *gin.Context, mcpServerId string) (*response.MCPServerDetail, error) {
	mcpServerInfo, err := mcp.GetMCPServer(ctx.Request.Context(), &mcp_service.GetMCPServerReq{
		McpServerId: mcpServerId,
	})
	if err != nil {
		return nil, err
	}
	mcpServerTools, err := mcp.GetMCPServerToolList(ctx.Request.Context(), &mcp_service.GetMCPServerToolListReq{
		McpServerId: mcpServerId,
	})
	if err != nil {
		return nil, err
	}
	return toMCPServerDetail(ctx, mcpServerInfo, mcpServerTools.List), nil
}

func DeleteMCPServer(ctx *gin.Context, mcpServerId string) error {
	_, err := app.DeleteApp(ctx.Request.Context(), &app_service.DeleteAppReq{
		AppId:   mcpServerId,
		AppType: constant.AppTypeMCPServer,
	})
	if err != nil {
		return err
	}
	_, err = mcp.DeleteMCPServer(ctx.Request.Context(), &mcp_service.DeleteMCPServerReq{
		McpServerId: mcpServerId,
	})
	if err != nil {
		return err
	}
	// 关闭 mcp server
	err = mcp_util.ShutDownMCPServer(ctx, mcpServerId)
	if err != nil {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_shutdown_err", err.Error())
	}
	return nil
}

func GetMCPServerList(ctx *gin.Context, userID, orgID, name string) (*response.ListResult, error) {
	resp, err := mcp.GetMCPServerList(ctx.Request.Context(), &mcp_service.GetMCPServerListReq{
		Name: name,
		Identity: &mcp_service.Identity{
			OrgId:  orgID,
			UserId: userID,
		},
	})
	if err != nil {
		return nil, err
	}
	var list []response.MCPServerInfo
	for _, mcpServerInfo := range resp.List {
		list = append(list, toMCPServerInfo(ctx, mcpServerInfo))
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func CreateMCPServerTool(ctx *gin.Context, req request.MCPServerToolCreateReq) error {
	toolInfo := &ToolInfo{
		Id:          req.Id,
		ToolType:    req.Type,
		MethodNames: []string{req.MethodName},
	}
	if req.Type == constant.MCPServerToolTypeCustomTool {
		info, err := mcp.GetCustomToolInfo(ctx.Request.Context(), &mcp_service.GetCustomToolInfoReq{
			CustomToolId: req.Id,
		})
		if err != nil {
			return err
		}
		toolInfo.Name = info.Name
	}
	schema, err := CreateMcpSchema(ctx, req.MCPServerID, toolInfo)
	if err != nil || schema == nil {
		return err
	}
	_, err = mcp.CreateMCPServerTool(ctx.Request.Context(), &mcp_service.CreateMCPServerToolReq{
		McpServerId:         req.MCPServerID,
		McpServiceToolInfos: schema.McpServerTool,
	})
	if err != nil {
		return err
	}
	err = mcp_util.RegisterMCPServerTools(req.MCPServerID, schema.McpTool)
	if err != nil {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_register_tool_err", err.Error())
	}
	return err
}

func UpdateMCPServerTool(ctx *gin.Context, req request.MCPServerToolUpdateReq) error {
	tool, err := mcp.GetMCPServerTool(ctx.Request.Context(), &mcp_service.GetMCPServerToolReq{
		McpServerToolId: req.MCPServerToolID,
	})
	if err != nil {
		return err
	}
	if tool.Name == req.MethodName && tool.Desc == req.Desc {
		return nil
	}
	schema, err := UpdateMcpSchema(req, tool)
	if err != nil || schema == nil || len(schema.McpServerTool) < 1 {
		return err
	}
	mcpServerTool := schema.McpServerTool[0]
	_, err = mcp.UpdateMCPServerTool(ctx.Request.Context(), &mcp_service.UpdateMCPServerToolReq{
		McpServerToolId: mcpServerTool.McpServerToolId,
		Name:            mcpServerTool.Name,
		Desc:            mcpServerTool.Desc,
		Schema:          mcpServerTool.Schema,
	})
	if err != nil {
		return err
	}
	err = mcp_util.UnRegisterMCPServerTools(tool.McpServerId, []string{tool.Name})
	if err != nil {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_unregister_tool_err", err.Error())
	}
	err = mcp_util.RegisterMCPServerTools(tool.McpServerId, schema.McpTool)
	if err != nil {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_register_tool_err", err.Error())
	}
	return err
}

func DeleteMCPServerTool(ctx *gin.Context, mcpServerToolId string) error {
	info, err := mcp.GetMCPServerTool(ctx.Request.Context(), &mcp_service.GetMCPServerToolReq{
		McpServerToolId: mcpServerToolId,
	})
	if err != nil {
		return err
	}
	_, err = mcp.DeleteMCPServerTool(ctx.Request.Context(), &mcp_service.DeleteMCPServerToolReq{
		McpServerToolId: mcpServerToolId,
	})
	if err != nil {
		return err
	}
	err = mcp_util.UnRegisterMCPServerTools(info.McpServerId, []string{info.Name})
	if err != nil {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_unregister_tool_err", err.Error())
	}
	return nil
}

func CreateMCPServerOpenAPITool(ctx *gin.Context, userID, orgID string, req request.MCPServerOpenAPIToolCreate) error {
	schema, err := CreateMcpSchema(ctx, req.MCPServerID, &ToolInfo{
		Name:        req.Name,
		ToolType:    constant.MCPServerToolTypeOpenAPI,
		MethodNames: req.MethodNames,
		Schema:      req.Schema,
		ApiAuth:     req.ApiAuth,
	})
	if err != nil || schema == nil {
		return err
	}
	_, err = mcp.CreateMCPServerTool(ctx.Request.Context(), &mcp_service.CreateMCPServerToolReq{
		McpServerId:         req.MCPServerID,
		McpServiceToolInfos: schema.McpServerTool,
	})
	if err != nil {
		return err
	}
	err = mcp_util.RegisterMCPServerTools(req.MCPServerID, schema.McpTool)
	if err != nil {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_register_tool_err", err.Error())
	}
	return err
}

func GetMCPServerCustomToolSelect(ctx *gin.Context, userID, orgID, name string) (*response.ListResult, error) {
	resp, err := mcp.GetCustomToolList(ctx.Request.Context(), &mcp_service.GetCustomToolListReq{
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	var list []response.MCPServerCustomToolSelect
	for _, item := range resp.List {
		info, err := mcp.GetCustomToolInfo(ctx.Request.Context(), &mcp_service.GetCustomToolInfoReq{
			CustomToolId: item.CustomToolId,
		})
		if err != nil {
			return nil, err
		}
		apis, err := GetSchemaActions(ctx, info.Schema)
		if err != nil {
			continue
		}
		toolList := toMCPServerCustomToolSelect(item, apis)
		list = append(list, toolList...)
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func GetMCPServerSSE(ctx *gin.Context, mcpServerId string, key string) error {
	if !mcp_util.CheckMCPServerExist(mcpServerId) {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_not_exist")
	}
	queryParams := ctx.Request.URL.Query()
	queryParams.Set("key", key)
	ctx.Request.URL.RawQuery = queryParams.Encode()
	mcp_util.GetMCPServerSSEHandler(mcpServerId).HandleSSE().ServeHTTP(ctx.Writer, ctx.Request)
	return nil
}

func GetMCPServerMessage(ctx *gin.Context, mcpServerId string) error {
	if !mcp_util.CheckMCPServerExist(mcpServerId) {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_not_exist")
	}
	var body []byte
	if cb, ok := ctx.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	if body != nil {
		// 调用前再次确保Body可用（防止中间件已读取）
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	mcp_util.GetMCPServerSSEHandler(mcpServerId).HandleMessage().ServeHTTP(ctx.Writer, ctx.Request)
	return nil
}

func GetMCPServerStreamable(ctx *gin.Context, mcpServerId string) error {
	if !mcp_util.CheckMCPServerExist(mcpServerId) {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_not_exist")
	}
	var body []byte
	if cb, ok := ctx.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	if body != nil {
		// 调用前再次确保Body可用（防止中间件已读取）
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	mcp_util.GetMCPServerStreamableHandler(mcpServerId).HandleMCP().ServeHTTP(ctx.Writer, ctx.Request)
	return nil
}

func StartMCPServer(ctx context.Context) error {
	err := mcp_util.Init(ctx)
	if err != nil {
		return err
	}
	mcpServerList, err := mcp.GetMCPServerList(ctx, &mcp_service.GetMCPServerListReq{
		Identity: &mcp_service.Identity{
			OrgId:  "",
			UserId: "",
		},
	})
	if err != nil {
		return err
	}
	for _, mcpServerInfo := range mcpServerList.List {
		mcpServerToolList, err := mcp.GetMCPServerToolList(ctx, &mcp_service.GetMCPServerToolListReq{
			McpServerId: mcpServerInfo.McpServerId,
		})
		if err != nil {
			return err
		}
		mcpTools := make([]*mcp_util.McpTool, 0)
		for _, tool := range mcpServerToolList.List {
			doc, err := mcp_util.ParseOpenApiContent(tool.Schema)
			if err != nil {
				return err
			}
			mcpTool := &mcp_util.McpTool{}
			mcpTool.Tool = mcp_util.ConvertMcpTool(doc, tool.Name)
			mcpTool.Handle = mcp_util.ConvertMcpHandler(doc, tool.Name, &mcp_util.APIAuth{
				Type:  tool.ApiAuth.AuthType,
				In:    tool.ApiAuth.AuthIn,
				Name:  tool.ApiAuth.AuthName,
				Value: tool.ApiAuth.AuthValue,
			})
			mcpTools = append(mcpTools, mcpTool)
		}
		err = mcp_util.StartMCPServer(ctx, mcpServerInfo.McpServerId)
		if err != nil {
			return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_start_err", err.Error())
		}
		err = mcp_util.RegisterMCPServerTools(mcpServerInfo.McpServerId, mcpTools)
		if err != nil {
			return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_mcp_server_register_tool_err", err.Error())
		}
	}
	return nil
}

// internal

func GetSchemaActions(ctx *gin.Context, schema string) ([]response.CustomToolActionInfo, error) {
	doc, err := openapi3_util.LoadFromData([]byte(schema))
	if err != nil {
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFInvalidArg, err.Error())
	}
	if err := openapi3_util.ValidateDoc(ctx.Request.Context(), doc); err != nil {
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFInvalidArg, err.Error())
	}
	list := openapiSchema2ToolList(doc)
	return list, nil
}

func toMCPServerCustomToolSelect(item *mcp_service.GetCustomToolItem, apis []response.CustomToolActionInfo) []response.MCPServerCustomToolSelect {
	var list []response.MCPServerCustomToolSelect
	var methods []response.MCPServerCustomToolApi
	for _, api := range apis {
		methods = append(methods, response.MCPServerCustomToolApi{
			MethodName:  api.Name,
			Description: api.Desc,
		})
	}
	list = append(list, response.MCPServerCustomToolSelect{
		UniqueId:     constant.MCPServerToolTypeCustomTool + "-" + item.CustomToolId,
		CustomToolId: item.CustomToolId,
		Name:         item.Name,
		Description:  item.Description,
		Methods:      methods,
	})
	return list
}

func toMCPServerInfo(ctx *gin.Context, mcpServerInfo *mcp_service.MCPServerInfo) response.MCPServerInfo {
	return response.MCPServerInfo{
		MCPServerID: mcpServerInfo.McpServerId,
		Avatar:      CacheAvatar(ctx, mcpServerInfo.AvatarPath, true),
		Name:        mcpServerInfo.Name,
		Desc:        mcpServerInfo.Desc,
		ToolNum:     mcpServerInfo.ToolNum,
	}
}

func toMCPServerDetail(ctx *gin.Context, mcpServerInfo *mcp_service.MCPServerInfo, mcpServerToolInfos []*mcp_service.MCPServerToolInfo) *response.MCPServerDetail {
	var mcpServerTools []response.MCPServerToolInfo
	for _, mcpServerToolInfo := range mcpServerToolInfos {
		mcpServerTools = append(mcpServerTools, response.MCPServerToolInfo{
			MCPServerToolID: mcpServerToolInfo.McpServerToolId,
			MethodName:      mcpServerToolInfo.Name,
			Type:            mcpServerToolInfo.Type,
			Name:            mcpServerToolInfo.AppToolName,
			Desc:            mcpServerToolInfo.Desc,
		})
	}
	return &response.MCPServerDetail{
		MCPServerID:       mcpServerInfo.McpServerId,
		Avatar:            CacheAvatar(ctx, mcpServerInfo.AvatarPath, true),
		Name:              mcpServerInfo.Name,
		Desc:              mcpServerInfo.Desc,
		SSEURL:            mcpServerInfo.SseUrl,
		SSEExample:        mcpServerInfo.SseExample,
		StreamableURL:     mcpServerInfo.StreamableUrl,
		StreamableExample: mcpServerInfo.StreamableExample,
		Tools:             mcpServerTools,
	}
}
