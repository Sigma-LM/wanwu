package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// CreateCustomTool
//
//	@Tags			tool
//	@Summary		创建自定义工具
//	@Description	创建自定义工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CustomToolCreate	true	"自定义工具信息"
//	@Success		200		{object}	response.Response{}
//	@Router			/tool/custom [post]
func CreateCustomTool(ctx *gin.Context) {
	var req request.CustomToolCreate
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.CreateCustomTool(ctx, getUserID(ctx), getOrgID(ctx), req))
}

// GetCustomTool
//
//	@Tags			tool
//	@Summary		获取自定义工具详情
//	@Description	获取自定义工具详情
//	@Accept			json
//	@Produce		json
//	@Param			customToolId	query		string	true	"customToolId"
//	@Success		200				{object}	response.Response{data=response.CustomToolDetail}
//	@Router			/tool/custom [get]
func GetCustomTool(ctx *gin.Context) {
	resp, err := service.GetCustomTool(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("customToolId"))
	gin_util.Response(ctx, resp, err)
}

// DeleteCustomTool
//
//	@Tags			tool
//	@Summary		删除自定义工具
//	@Description	删除自定义工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CustomToolIDReq	true	"自定义工具ID"
//	@Success		200		{object}	response.Response{}
//	@Router			/tool/custom [delete]
func DeleteCustomTool(ctx *gin.Context) {
	var req request.CustomToolIDReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.DeleteCustomTool(ctx, getUserID(ctx), getOrgID(ctx), req))
}

// UpdateCustomTool
//
//	@Tags			tool
//	@Summary		修改自定义工具
//	@Description	修改自定义工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CustomToolUpdateReq	true	"自定义工具信息"
//	@Success		200		{object}	response.Response{}
//	@Router			/tool/custom [put]
func UpdateCustomTool(ctx *gin.Context) {
	var req request.CustomToolUpdateReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.UpdateCustomTool(ctx, getUserID(ctx), getOrgID(ctx), req))
}

// GetCustomToolList
//
//	@Tags			tool
//	@Summary		获取自定义工具列表
//	@Description	获取自定义工具列表
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"CustomTool名称"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.CustomToolInfo}}
//	@Router			/tool/custom/list [get]
func GetCustomToolList(ctx *gin.Context) {
	resp, err := service.GetCustomToolList(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}

// GetCustomToolSelect
//
//	@Tags			tool
//	@Summary		获取自定义工具列表（用于下拉选择）
//	@Description	获取自定义工具列表（用于下拉选择）
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"CustomTool名称"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.CustomToolSelect}}
//	@Router			/tool/custom/select [get]
func GetCustomToolSelect(ctx *gin.Context) {
	resp, err := service.GetCustomToolSelect(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}

// GetCustomToolActions
//
//	@Tags			tool
//	@Summary		获取可用API列表（根据Schema）
//	@Description	解析自定义工具的Schema转换为API相关数据
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CustomToolSchemaReq	true	"Schema格式数据"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.CustomToolActionInfo}}
//	@Router			/tool/custom/schema [post]
func GetCustomToolActions(ctx *gin.Context) {
	var req request.CustomToolSchemaReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.GetCustomToolActions(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, resp, err)
}

// GetToolSquareDetail
//
//	@Tags			tool.square
//	@Summary		获取内置工具详情
//	@Description	获取内置工具详情
//	@Accept			json
//	@Produce		json
//	@Param			toolSquareId	query		string	true	"toolSquareId"
//	@Success		200				{object}	response.Response{data=response.ToolSquareDetail}
//	@Router			/tool/square [get]
func GetToolSquareDetail(ctx *gin.Context) {
	resp, err := service.GetToolSquareDetail(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("toolSquareId"))
	gin_util.Response(ctx, resp, err)
}

// GetToolSquareList
//
//	@Tags			tool.square
//	@Summary		获取内置工具列表
//	@Description	获取内置工具列表
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"tool名称"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.ToolSquareInfo}}
//	@Router			/tool/square/list [get]
func GetToolSquareList(ctx *gin.Context) {
	resp, err := service.GetToolSquareList(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}

// UpdateToolSquareAPIKey
//
//	@Tags			tool
//	@Summary		修改内置工具
//	@Description	修改内置工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.ToolSquareAPIKeyReq	true	"内置工具信息"
//	@Success		200		{object}	response.Response{}
//	@Router			/tool/builtin [post]
func UpdateToolSquareAPIKey(ctx *gin.Context) {
	var req request.ToolSquareAPIKeyReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.UpsertBuiltinToolAPIKey(ctx, getUserID(ctx), getOrgID(ctx), req))
}

// GetToolSelect
//
//	@Tags			tool
//	@Summary		获取工具列表（用于下拉选择）
//	@Description	获取工具列表（用于下拉选择）
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	true	"工具名"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.ToolSelect}}
//	@Router			/tool/select [get]
func GetToolSelect(ctx *gin.Context) {
	resp, err := service.GetToolSelect(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}

// GetToolActionList
//
//	@Tags			tool
//	@Summary		获取工具列表
//	@Description	获取工具列表
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.ToolActionListReq	true	"工具信息"
//	@Success		200		{object}	response.Response{data=response.ToolActionList}
//	@Router			/tool/action/list [get]
func GetToolActionList(ctx *gin.Context) {
	var req request.ToolActionListReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetToolActionList(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, resp, err)
}

// GetToolActionDetail
//
//	@Tags			tool
//	@Summary		获取工具详情
//	@Description	获取工具详情
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.ToolActionReq	true	"工具信息"
//	@Success		200		{object}	response.Response{data=response.ToolActionDetail}
//	@Router			/tool/action/detail [get]
func GetToolActionDetail(ctx *gin.Context) {
	var req request.ToolActionReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetToolActionDetail(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, resp, err)
}

// CreateMCPServer
//
//	@Tags			tool
//	@Summary		创建MCP Server
//	@Description	创建MCP Server
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.MCPServerCreateReq	true	"MCP Server信息"
//	@Success		200		{object}	response.Response{data=response.MCPServerCreateResp}
//	@Router			/mcp/server [post]
func CreateMCPServer(ctx *gin.Context) {
	var req request.MCPServerCreateReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CreateMCPServer(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, resp, err)
}

// UpdateMCPServer
//
//	@Tags			tool
//	@Summary		更新MCP Server
//	@Description	更新MCP Server
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.MCPServerUpdateReq	true	"MCP Server信息"
//	@Success		200		{object}	response.Response{}
//	@Router			/mcp/server [put]
func UpdateMCPServer(ctx *gin.Context) {
	var req request.MCPServerUpdateReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.UpdateMCPServer(ctx, req))
}

// GetMCPServer
//
//	@Tags			tool
//	@Summary		获取MCP Server详情
//	@Description	获取MCP Server详情
//	@Accept			json
//	@Produce		json
//	@Param			mcpServerId	query		string	true	"mcpServerId"
//	@Success		200			{object}	response.Response{data=response.MCPServerDetail}
//	@Router			/mcp/server [get]
func GetMCPServer(ctx *gin.Context) {
	resp, err := service.GetMCPServerDetail(ctx, ctx.Query("mcpServerId"))
	gin_util.Response(ctx, resp, err)
}

// DeleteMCPServer
//
//	@Tags			tool
//	@Summary		删除MCP Server
//	@Description	删除MCP Server
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.MCPServerIDReq	true	"mcpServerId"
//	@Success		200		{object}	response.Response{}
//	@Router			/mcp/server [delete]
func DeleteMCPServer(ctx *gin.Context) {
	var req request.MCPServerIDReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.DeleteMCPServer(ctx, req.MCPServerID)
	gin_util.Response(ctx, nil, err)
}

// GetMCPServerList
//
//	@Tags			tool
//	@Summary		获取MCP Server列表
//	@Description	获取MCP Server列表
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"mcp server名称"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.MCPServerInfo}}
//	@Router			/mcp/server/list [get]
func GetMCPServerList(ctx *gin.Context) {
	resp, err := service.GetMCPServerList(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}

// CreateMCPServerTool
//
//	@Tags			tool
//	@Summary		创建MCP Server工具
//	@Description	创建MCP Server工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.MCPServerToolCreateReq	true	"MCP Server工具信息"
//	@Success		200		{object}	response.Response{}
//	@Router			/mcp/server/tool [post]
func CreateMCPServerTool(ctx *gin.Context) {
	var req request.MCPServerToolCreateReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.CreateMCPServerTool(ctx, req))
}

// UpdateMCPServerTool
//
//	@Tags			tool
//	@Summary		更新MCP Server工具
//	@Description	更新MCP Server工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.MCPServerToolUpdateReq	true	"MCP Server工具信息"
//	@Success		200		{object}	response.Response{}
//	@Router			/mcp/server/tool [put]
func UpdateMCPServerTool(ctx *gin.Context) {
	var req request.MCPServerToolUpdateReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.UpdateMCPServerTool(ctx, req))
}

// DeleteMCPServerTool
//
//	@Tags			tool
//	@Summary		删除MCP Server工具
//	@Description	删除MCP Server工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.MCPServerToolIDReq	true	"mcpServerToolId"
//	@Success		200		{object}	response.Response{}
//	@Router			/mcp/server/tool [delete]
func DeleteMCPServerTool(ctx *gin.Context) {
	var req request.MCPServerToolIDReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.DeleteMCPServerTool(ctx, req.MCPServerToolID)
	gin_util.Response(ctx, nil, err)
}

// CreateMCPServerOpenAPITool
//
//	@Tags			tool
//	@Summary		创建openapi工具
//	@Description	创建openapi工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.MCPServerOpenAPIToolCreate	true	"openapi工具信息"
//	@Success		200		{object}	response.Response{}
//	@Router			/mcp/server/tool/openapi [post]
func CreateMCPServerOpenAPITool(ctx *gin.Context) {
	var req request.MCPServerOpenAPIToolCreate
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.CreateMCPServerOpenAPITool(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, nil, err)
}

// GetMCPServerCustomToolSelect
//
//	@Tags			tool
//	@Summary		获取MCP Server自定义工具列表（用于下拉选择）
//	@Description	获取MCP Server自定义工具列表（用于下拉选择）
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"CustomTool名称"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.MCPServerCustomToolSelect}}
//	@Router			/mcp/server/tool/custom/select [get]
func GetMCPServerCustomToolSelect(ctx *gin.Context) {
	resp, err := service.GetMCPServerCustomToolSelect(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}
