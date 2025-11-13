package mcp

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/UnicomAI/wanwu/pkg/constant"

	"github.com/UnicomAI/wanwu/api/proto/common"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/mcp-service/client/model"
	"github.com/UnicomAI/wanwu/internal/mcp-service/config"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) CreateCustomTool(ctx context.Context, req *mcp_service.CreateCustomToolReq) (*emptypb.Empty, error) {
	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPCreateCustomToolErr, toErrStatus("mcp_create_custom_tool_err", "identity is empty"))
	}
	if req.ApiAuth == nil {
		return nil, errStatus(errs.Code_MCPCreateCustomToolErr, toErrStatus("mcp_create_custom_tool_err", "apiAuth is empty"))
	}
	apiAuth := &common.ApiAuthWebRequest{
		AuthType:           req.ApiAuth.AuthType,
		ApiKeyHeaderPrefix: req.ApiAuth.ApiKeyHeaderPrefix,
		ApiKeyHeader:       req.ApiAuth.ApiKeyHeader,
		ApiKeyQueryParam:   req.ApiAuth.ApiKeyQueryParam,
		ApiKeyValue:        req.ApiAuth.ApiKeyValue,
	}
	apiAuthBytes, err := json.Marshal(apiAuth)
	if err != nil {
		return nil, errStatus(errs.Code_MCPCreateCustomToolErr, toErrStatus("mcp_create_custom_tool_err", err.Error()))
	}
	if err := s.cli.CreateCustomTool(ctx, &model.CustomTool{
		AvatarPath:    req.AvatarPath,
		Name:          req.Name,
		Description:   req.Description,
		Schema:        req.Schema,
		PrivacyPolicy: req.PrivacyPolicy,
		AuthJSON:      string(apiAuthBytes),
		UserID:        req.Identity.UserId,
		OrgID:         req.Identity.OrgId,
	}); err != nil {
		return nil, errStatus(errs.Code_MCPCreateCustomToolErr, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetCustomToolInfo(ctx context.Context, req *mcp_service.GetCustomToolInfoReq) (*mcp_service.GetCustomToolInfoResp, error) {
	if req.CustomToolId == "" {
		return nil, errStatus(errs.Code_MCPGetCustomToolInfoErr, toErrStatus("mcp_get_custom_tool_info_err", "customToolId is empty"))
	}
	info, err := s.cli.GetCustomTool(ctx, &model.CustomTool{
		ID: util.MustU32(req.CustomToolId),
	})
	if err != nil {
		return nil, grpc_util.ErrorStatus(errs.Code_MCPGetCustomToolInfoErr)
	}
	apiAuthJson := info.AuthJSON
	apiAuth := &common.ApiAuthWebRequest{}
	if err := json.Unmarshal([]byte(apiAuthJson), apiAuth); err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolInfoErr, toErrStatus("mcp_get_custom_tool_info_err", err.Error()))
	}
	return &mcp_service.GetCustomToolInfoResp{
		CustomToolId:  util.Int2Str(info.ID),
		AvatarPath:    info.AvatarPath,
		Name:          info.Name,
		Description:   info.Description,
		Schema:        info.Schema,
		PrivacyPolicy: info.PrivacyPolicy,
		ApiAuth:       apiAuth,
	}, nil
}

func (s *Service) GetCustomToolList(ctx context.Context, req *mcp_service.GetCustomToolListReq) (*mcp_service.GetCustomToolListResp, error) {
	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, toErrStatus("mcp_get_custom_tool_list_err", "identity is empty"))
	}
	infos, err := s.cli.ListCustomTools(ctx, req.Identity.OrgId, req.Identity.UserId, req.Name)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, err)
	}
	list := make([]*mcp_service.GetCustomToolItem, 0)
	for _, info := range infos {
		list = append(list, &mcp_service.GetCustomToolItem{
			CustomToolId: util.Int2Str(info.ID),
			Name:         info.Name,
			Description:  info.Description,
			AvatarPath:   info.AvatarPath,
		})
	}
	return &mcp_service.GetCustomToolListResp{
		List: list,
	}, nil
}

func (s *Service) GetCustomToolByCustomToolIdList(ctx context.Context, req *mcp_service.GetCustomToolByCustomToolIdListReq) (*mcp_service.GetCustomToolListResp, error) {
	if len(req.CustomToolIdList) == 0 {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, toErrStatus("mcp_get_custom_tool_list_err", "customToolIdList is empty"))
	}

	// 批量转换 string 为 uint32
	var ids []uint32
	for _, idStr := range req.CustomToolIdList {
		id := util.MustU32(idStr)
		ids = append(ids, id)
	}

	infos, err := s.cli.ListCustomToolsByCustomToolIDs(ctx, ids)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, err)
	}
	list := make([]*mcp_service.GetCustomToolItem, 0)
	for _, info := range infos {
		list = append(list, &mcp_service.GetCustomToolItem{
			CustomToolId: util.Int2Str(info.ID),
			Name:         info.Name,
			Description:  info.Description,
		})
	}
	return &mcp_service.GetCustomToolListResp{
		List: list,
	}, nil
}

func (s *Service) UpdateCustomTool(ctx context.Context, req *mcp_service.UpdateCustomToolReq) (*emptypb.Empty, error) {
	if req.CustomToolId == "" {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, toErrStatus("mcp_update_custom_tool_err", "customToolId is empty"))
	}
	if req.ApiAuth == nil {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, toErrStatus("mcp_update_custom_tool_err", "apiAuth is empty"))
	}
	apiAuth := &common.ApiAuthWebRequest{
		AuthType:           req.ApiAuth.AuthType,
		ApiKeyHeaderPrefix: req.ApiAuth.ApiKeyHeaderPrefix,
		ApiKeyHeader:       req.ApiAuth.ApiKeyHeader,
		ApiKeyQueryParam:   req.ApiAuth.ApiKeyQueryParam,
		ApiKeyValue:        req.ApiAuth.ApiKeyValue,
	}
	apiAuthBytes, err := json.Marshal(apiAuth)
	if err != nil {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, toErrStatus("mcp_update_custom_tool_err", err.Error()))
	}
	if err := s.cli.UpdateCustomTool(ctx, &model.CustomTool{
		ID:            util.MustU32(req.CustomToolId),
		AvatarPath:    req.AvatarPath,
		Name:          req.Name,
		Description:   req.Description,
		Schema:        req.Schema,
		PrivacyPolicy: req.PrivacyPolicy,
		AuthJSON:      string(apiAuthBytes),
	}); err != nil {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteCustomTool(ctx context.Context, req *mcp_service.DeleteCustomToolReq) (*emptypb.Empty, error) {
	if req.CustomToolId == "" {
		return nil, errStatus(errs.Code_MCPDeleteCustomToolErr, toErrStatus("mcp_delete_custom_tool_err", "customToolId is empty"))
	}
	if err := s.cli.DeleteCustomTool(ctx, util.MustU32(req.CustomToolId)); err != nil {
		return nil, errStatus(errs.Code_MCPDeleteCustomToolErr, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetSquareTool(ctx context.Context, req *mcp_service.GetSquareToolReq) (*mcp_service.SquareToolDetail, error) {
	toolCfg, exist := config.Cfg().Tool(req.ToolSquareId)
	if !exist {
		return nil, errStatus(errs.Code_MCPGetSquareToolErr, toErrStatus("mcp_get_square_tool_err", "toolSquareId not exist"))
	}

	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPGetSquareToolErr, toErrStatus("mcp_get_square_tool_err", "identity is empty"))
	}

	apiAuth := &common.ApiAuthWebRequest{}
	if toolCfg.NeedApiKeyInput {
		info, _ := s.cli.GetBuiltinTool(ctx, &model.BuiltinTool{
			ToolSquareId: req.ToolSquareId,
			OrgID:        req.Identity.OrgId,
			UserID:       req.Identity.UserId,
		})
		if info != nil {
			apiAuthJson := info.AuthJSON
			if err := json.Unmarshal([]byte(apiAuthJson), apiAuth); err != nil {
				return nil, errStatus(errs.Code_MCPGetSquareToolErr, toErrStatus("mcp_get_square_tool_err", err.Error()))
			}
			apiAuth.AuthType = toolCfg.AuthType
			apiAuth.ApiKeyHeaderPrefix = toolCfg.ApiKeyHeaderPrefix
			apiAuth.ApiKeyHeader = toolCfg.ApiKeyHeader
			apiAuth.ApiKeyQueryParam = toolCfg.ApiKeyQueryParam
		}
	}
	return buildSquareToolDetail(toolCfg, apiAuth), nil
}

func (s *Service) GetSquareToolList(ctx context.Context, req *mcp_service.GetSquareToolListReq) (*mcp_service.SquareToolList, error) {
	var toolSquareInfo []*mcp_service.ToolSquareInfo
	for _, toolCfg := range config.Cfg().Tools {
		if req.Name != "" && !strings.Contains(toolCfg.Name, req.Name) {
			continue
		}
		toolSquareInfo = append(toolSquareInfo, buildSquareToolInfo(*toolCfg))
	}
	return &mcp_service.SquareToolList{Infos: toolSquareInfo}, nil
}

func buildSquareToolInfo(toolCfg config.ToolConfig) *mcp_service.ToolSquareInfo {
	return &mcp_service.ToolSquareInfo{
		ToolSquareId: toolCfg.ToolSquareId,
		AvatarPath:   toolCfg.AvatarPath,
		Name:         toolCfg.Name,
		Desc:         toolCfg.Desc,
		Tags:         toolCfg.Tags,
	}
}

func buildSquareToolDetail(toolCfg config.ToolConfig, apiAuth *common.ApiAuthWebRequest) *mcp_service.SquareToolDetail {
	return &mcp_service.SquareToolDetail{
		Info: buildSquareToolInfo(toolCfg),
		BuiltInTools: &mcp_service.BuiltInTools{
			NeedApiKeyInput: toolCfg.NeedApiKeyInput,
			ApiAuth:         apiAuth,
			Detail:          toolCfg.Detail,
			ActionSum:       int32(len(toolCfg.Tools)),
			Tools:           convertMCPTools(toolCfg.Tools),
		},
		Schema: toolCfg.Schema,
	}
}

func (s *Service) GetToolByIdList(ctx context.Context, req *mcp_service.GetToolByToolIdListReq) (*mcp_service.GetToolByToolIdListResp, error) {
	// 内置工具
	var toolSquareInfo []*mcp_service.ToolSquareInfo
	for _, toolCfg := range config.Cfg().Tools {
		for _, builtInTool := range req.BuiltInToolIdList {
			if builtInTool != "" && !strings.Contains(toolCfg.Name, builtInTool) {
				toolSquareInfo = append(toolSquareInfo, buildSquareToolInfo(*toolCfg))
			}
		}
	}
	// 自定义工具
	var ids []uint32
	for _, idStr := range req.CustomToolIdList {
		id := util.MustU32(idStr)
		ids = append(ids, id)
	}

	infos, err := s.cli.ListCustomToolsByCustomToolIDs(ctx, ids)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, err)
	}
	list := make([]*mcp_service.GetCustomToolItem, 0)
	for _, info := range infos {
		list = append(list, &mcp_service.GetCustomToolItem{
			CustomToolId: util.Int2Str(info.ID),
			Name:         info.Name,
			Description:  info.Description,
			AvatarPath:   info.AvatarPath,
		})
	}
	return &mcp_service.GetToolByToolIdListResp{
		List:               list,
		ToolSquareInfoList: toolSquareInfo,
	}, nil
}

func (s *Service) UpsertBuiltinToolAPIKey(ctx context.Context, req *mcp_service.UpsertBuiltinToolAPIKeyReq) (*emptypb.Empty, error) {
	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPUpdateBuiltinToolErr, toErrStatus("mcp_update_builtin_tool_err", "identity is empty"))
	}
	apiAuth := &common.ApiAuthWebRequest{}
	builtToolCfg, exist := config.Cfg().Tool(req.ToolSquareId)
	if !exist {
		return nil, errStatus(errs.Code_MCPUpdateBuiltinToolErr, toErrStatus("mcp_update_builtin_tool_err", "toolSquareId not exist"))
	}

	info, _ := s.cli.GetBuiltinTool(ctx, &model.BuiltinTool{
		ToolSquareId: req.ToolSquareId,
		OrgID:        req.Identity.OrgId,
		UserID:       req.Identity.UserId,
	})
	if info != nil {
		// update
		apiAuthJson := info.AuthJSON
		if err := json.Unmarshal([]byte(apiAuthJson), apiAuth); err != nil {
			return nil, errStatus(errs.Code_MCPUpdateBuiltinToolErr, toErrStatus("mcp_update_builtin_tool_err", err.Error()))
		}
		apiAuth.AuthType = builtToolCfg.AuthType
		apiAuth.ApiKeyHeaderPrefix = builtToolCfg.ApiKeyHeaderPrefix
		apiAuth.ApiKeyHeader = builtToolCfg.ApiKeyHeader
		apiAuth.ApiKeyQueryParam = builtToolCfg.ApiKeyQueryParam
		apiAuth.ApiKeyValue = req.ApiKey
		apiAuthBytes, _ := json.Marshal(apiAuth)
		info.AuthJSON = string(apiAuthBytes)
		if err := s.cli.UpdateBuiltinTool(ctx, info); err != nil {
			return nil, errStatus(errs.Code_MCPUpdateBuiltinToolErr, err)
		}
		return &emptypb.Empty{}, nil
	} else {
		// create
		apiAuth.AuthType = builtToolCfg.AuthType
		apiAuth.ApiKeyHeaderPrefix = builtToolCfg.ApiKeyHeaderPrefix
		apiAuth.ApiKeyHeader = builtToolCfg.ApiKeyHeader
		apiAuth.ApiKeyQueryParam = builtToolCfg.ApiKeyQueryParam
		apiAuth.ApiKeyValue = req.ApiKey
		apiAuthBytes, _ := json.Marshal(apiAuth)

		if err := s.cli.CreateBuiltinTool(ctx, &model.BuiltinTool{
			ToolSquareId: req.ToolSquareId,
			AuthJSON:     string(apiAuthBytes),
			UserID:       req.Identity.UserId,
			OrgID:        req.Identity.OrgId,
		}); err != nil {
			return nil, errStatus(errs.Code_MCPUpdateBuiltinToolErr, err)
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetToolSelect(ctx context.Context, req *mcp_service.GetToolSelectReq) (*mcp_service.GetToolListResp, error) {
	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, toErrStatus("mcp_get_custom_tool_list_err", "identity is empty"))
	}
	// search custom tools
	infos, err := s.cli.ListCustomTools(ctx, req.Identity.OrgId, req.Identity.UserId, req.Name)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, err)
	}
	list := make([]*mcp_service.GetToolItem, 0)
	for _, info := range infos {
		needApiKeyInput := true
		if info.Type == model.ApiAuthNone {
			needApiKeyInput = false
		}
		list = append(list, &mcp_service.GetToolItem{
			ToolId:          util.Int2Str(info.ID),
			ToolName:        info.Name,
			Desc:            info.Description,
			ToolType:        constant.ToolTypeCustom,
			ApiKey:          info.APIKey,
			NeedApiKeyInput: needApiKeyInput,
			AvatarPath:      info.AvatarPath,
		})
	}
	// search builtin tools
	// 先把该用户所有已配置apikey内置工具查出来，构造成map<ToolSquareId, *model.CustomTool>，然后通过ToolSquareId查询apikey
	builtinTools, err := s.cli.ListBuiltinTools(ctx, req.Identity.OrgId, req.Identity.UserId)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, err)
	}
	builtinToolMap := make(map[string]*model.BuiltinTool)
	for _, tool := range builtinTools {
		builtinToolMap[tool.ToolSquareId] = tool
	}

	for _, toolCfg := range config.Cfg().Tools {
		if req.Name != "" && !strings.Contains(toolCfg.Name, req.Name) {
			continue
		}
		toolTab := &mcp_service.GetToolItem{
			ToolId:          toolCfg.ToolSquareId,
			ToolName:        toolCfg.Name,
			Desc:            toolCfg.Desc,
			ToolType:        constant.ToolTypeBuiltIn,
			NeedApiKeyInput: toolCfg.NeedApiKeyInput,
			AvatarPath:      toolCfg.AvatarPath,
		}
		// 从map中查询内置工具
		if tool, ok := builtinToolMap[toolCfg.ToolSquareId]; ok {
			apiAuth := &common.ApiAuthWebRequest{}
			if err := json.Unmarshal([]byte(tool.AuthJSON), apiAuth); err != nil {
				return nil, errStatus(errs.Code_MCPGetCustomToolListErr, toErrStatus("mcp_get_custom_tool_list_err", err.Error()))
			}
			toolTab.ApiKey = apiAuth.ApiKeyValue

		}
		list = append(list, toolTab)
	}
	return &mcp_service.GetToolListResp{
		List:  list,
		Total: int32(len(list)),
	}, nil
}
