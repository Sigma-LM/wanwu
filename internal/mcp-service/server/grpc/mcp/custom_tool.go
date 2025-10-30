package mcp

import (
	"context"
	"strings"

	"github.com/UnicomAI/wanwu/pkg/constant"

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
	if err := s.cli.CreateCustomTool(ctx, &model.CustomTool{
		Name:             req.Name,
		Description:      req.Description,
		Schema:           req.Schema,
		PrivacyPolicy:    req.PrivacyPolicy,
		Type:             req.ApiAuth.Type,
		APIKey:           req.ApiAuth.ApiKey,
		AuthType:         req.ApiAuth.AuthType,
		CustomHeaderName: req.ApiAuth.CustomHeaderName,
		UserID:           req.Identity.UserId,
		OrgID:            req.Identity.OrgId,
		ToolSquareId:     req.ToolSquareId,
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
	return &mcp_service.GetCustomToolInfoResp{
		CustomToolId:  util.Int2Str(info.ID),
		Name:          info.Name,
		Description:   info.Description,
		Schema:        info.Schema,
		PrivacyPolicy: info.PrivacyPolicy,
		ApiAuth: &mcp_service.ApiAuthWebRequest{
			Type:             info.Type,
			ApiKey:           info.APIKey,
			AuthType:         info.AuthType,
			CustomHeaderName: info.CustomHeaderName,
		},
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
	if err := s.cli.UpdateCustomTool(ctx, &model.CustomTool{
		ID:               util.MustU32(req.CustomToolId),
		Name:             req.Name,
		Description:      req.Description,
		Schema:           req.Schema,
		PrivacyPolicy:    req.PrivacyPolicy,
		Type:             req.ApiAuth.Type,
		APIKey:           req.ApiAuth.ApiKey,
		AuthType:         req.ApiAuth.AuthType,
		CustomHeaderName: req.ApiAuth.CustomHeaderName,
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
	mcpCfg, exist := config.Cfg().Tool(req.ToolSquareId)
	if !exist {
		return nil, errStatus(errs.Code_MCPGetSquareToolErr, toErrStatus("mcp_get_square_tool_err", "toolSquareId not exist"))
	}

	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPGetSquareToolErr, toErrStatus("mcp_get_square_tool_err", "identity is empty"))
	}

	apiKey := ""
	if mcpCfg.NeedApiKeyInput {
		info, _ := s.cli.GetBuiltinTool(ctx, &model.CustomTool{
			ToolSquareId: req.ToolSquareId,
			OrgID:        req.Identity.OrgId,
			UserID:       req.Identity.UserId,
		})
		if info != nil {
			apiKey = info.APIKey
		}
	}
	return buildSquareToolDetail(mcpCfg, apiKey), nil
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

func buildSquareToolDetail(toolCfg config.ToolConfig, apiKey string) *mcp_service.SquareToolDetail {
	return &mcp_service.SquareToolDetail{
		Info: buildSquareToolInfo(toolCfg),
		BuiltInTools: &mcp_service.BuiltInTools{
			NeedApiKeyInput: toolCfg.NeedApiKeyInput,
			ApiKey:          apiKey,
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
		})
	}
	return &mcp_service.GetToolByToolIdListResp{
		List:               list,
		ToolSquareInfoList: toolSquareInfo,
	}, nil
}

func (s *Service) UpsertBuiltinToolAPIKey(ctx context.Context, req *mcp_service.UpsertBuiltinToolAPIKeyReq) (*emptypb.Empty, error) {

	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, toErrStatus("mcp_update_custom_tool_err", "identity is empty"))
	}
	info, _ := s.cli.GetBuiltinTool(ctx, &model.CustomTool{
		ToolSquareId: req.ToolSquareId,
		OrgID:        req.Identity.OrgId,
		UserID:       req.Identity.UserId,
	})
	if info != nil {
		// update
		if err := s.cli.UpdateCustomTool(ctx, &model.CustomTool{
			ID:     info.ID,
			APIKey: req.ApiKey,
		}); err != nil {
			return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, err)
		}
		return &emptypb.Empty{}, nil
	} else {
		// create
		if err := s.cli.CreateCustomTool(ctx, &model.CustomTool{
			ToolSquareId: req.ToolSquareId,
			APIKey:       req.ApiKey,
			UserID:       req.Identity.UserId,
			OrgID:        req.Identity.OrgId,
		}); err != nil {
			return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, err)
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
		})
	}
	// search builtin tools
	// 先把该用户所有已配置apikey内置工具查出来，构造成map<ToolSquareId, *model.CustomTool>，然后通过ToolSquareId查询apikey
	builtinTools, err := s.cli.ListBuiltinTools(ctx, req.Identity.OrgId, req.Identity.UserId)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, err)
	}
	builtinToolMap := make(map[string]*model.CustomTool)
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
		}
		// 从map中查询内置工具
		if tool, ok := builtinToolMap[toolCfg.ToolSquareId]; ok {
			toolTab.ApiKey = tool.APIKey
		}
		list = append(list, toolTab)
	}
	return &mcp_service.GetToolListResp{
		List:  list,
		Total: int32(len(list)),
	}, nil
}
