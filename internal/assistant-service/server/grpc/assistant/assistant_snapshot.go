package assistant

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	"github.com/UnicomAI/wanwu/api/proto/common"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/config"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) AssistantSnapshotCreate(ctx context.Context, req *assistant_service.AssistantSnapshotReq) (snapshot *assistant_service.AssistantSnapshotResp, err error) {
	assistantId, _ := util.U32(req.AssistantId)
	// 获取assistant详情
	assistant, status := s.cli.GetAssistant(ctx, assistantId, req.Identity.UserId, req.Identity.OrgId)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantErr, status)
	}
	if assistant == nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", "assistant info is nil"))
	}

	// 获取工作流配置详情
	workflows, status := s.cli.GetAssistantWorkflowsByAssistantID(ctx, assistantId)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantErr, status)
	}

	// 获取MCP配置详情
	mcpInfos, status := s.cli.GetAssistantMCPList(ctx, assistantId)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantErr, status)
	}

	// 获取Tool配置详情
	toolInfos, status := s.cli.GetAssistantToolList(ctx, assistantId)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantErr, status)
	}

	// 构造快照
	assistantSnapshot := model.AssistantSnapshot{
		// 基本信息
		AssistantID:  assistantId,
		Version:      req.Version,
		SnapshotDesc: req.Desc,
		// 智能体基本信息
		AssistantInfo: structToJson(assistant),
		// 智能体附表信息
		AssistantToolConfig:     structToJson(toolInfos),
		AssistantMCPConfig:      structToJson(mcpInfos),
		AssistantWorkflowConfig: structToJson(workflows),
		// 身份信息
		UserId: req.Identity.UserId,
		OrgId:  req.Identity.OrgId,
	}

	// 存入数据库
	assistantSnapshotID, status := s.cli.CreateAssistantSnapshot(ctx, &assistantSnapshot)
	if status != nil {
		log.Errorf("CreateAssistantSnapshot failed: %v", status)
		return nil, errStatus(errs.Code_AssistantErr, status)
	}

	return &assistant_service.AssistantSnapshotResp{
		SnapshotId: util.Int2Str(assistantSnapshotID),
	}, nil
}

func (s *Service) AssistantSnapshotUpdate(ctx context.Context, req *assistant_service.AssistantSnapshotUpdateReq) (*emptypb.Empty, error) {
	assistantId, _ := util.U32(req.AssistantId)

	status := s.cli.UpdateAssistantSnapshot(ctx, assistantId, req.Desc, req.Identity.UserId, req.Identity.OrgId)
	if status != nil {
		log.Errorf("UpdateAssistantSnapshot failed: %v", status)
		return nil, errStatus(errs.Code_AssistantErr, status)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) AssistantSnapshotList(ctx context.Context, req *assistant_service.AssistantSnapshotListReq) (*assistant_service.AssistantSnapshotListResp, error) {
	assistantId, _ := util.U32(req.AssistantId)

	assistantSnapshots, status := s.cli.GetAssistantSnapshotList(ctx, assistantId, req.Identity.UserId, req.Identity.OrgId)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantErr, status)
	}

	resp := make([]*assistant_service.AssistantSnapshot, 0, len(assistantSnapshots))
	for _, snapshot := range assistantSnapshots {
		resp = append(resp, &assistant_service.AssistantSnapshot{
			SnapshotId:  util.Int2Str(snapshot.ID),
			AssistantId: util.Int2Str(snapshot.AssistantID),
			Version:     snapshot.Version,
			Desc:        snapshot.SnapshotDesc,
			CreateAt:    snapshot.CreatedAt,
		})

	}
	return &assistant_service.AssistantSnapshotListResp{
		List:  resp,
		Total: int64(len(resp)),
	}, nil
}

func (s *Service) AssistantSnapshotLatest(ctx context.Context, req *assistant_service.AssistantSnapshotInfoReq) (*assistant_service.AssistantSnapshot, error) {
	assistantId := util.MustU32(req.AssistantId)
	snapshotInfo, status := s.cli.GetAssistantSnapshot(ctx, assistantId, "")
	if status != nil {
		return nil, errStatus(errs.Code_AssistantErr, status)
	}
	if snapshotInfo == nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", "assistant snapshot is nil"))
	}
	return &assistant_service.AssistantSnapshot{
		SnapshotId:  util.Int2Str(snapshotInfo.ID),
		AssistantId: util.Int2Str(snapshotInfo.AssistantID),
		Version:     snapshotInfo.Version,
		Desc:        snapshotInfo.SnapshotDesc,
		CreateAt:    snapshotInfo.CreatedAt,
	}, nil
}
func (s *Service) AssistantSnapshotInfo(ctx context.Context, req *assistant_service.AssistantSnapshotInfoReq) (*assistant_service.AssistantInfo, error) {
	snapshotInfo, status := s.cli.GetAssistantSnapshot(ctx, util.MustU32(req.AssistantId), "")

	if status != nil {
		return nil, errStatus(errs.Code_AssistantErr, status)
	}
	if snapshotInfo == nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", "assistant snapshot is nil"))
	}

	// 解析assistantInfo
	var snapshotAssistant *model.Assistant
	if err := jsonToStruct(snapshotInfo.AssistantInfo, &snapshotAssistant); err != nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", err.Error()))
	}
	if snapshotAssistant == nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", "assistant info is nil"))
	}

	// 解析workflow
	var workFlowConfig []*model.AssistantWorkflow
	if err := jsonToStruct(snapshotInfo.AssistantWorkflowConfig, &workFlowConfig); err != nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", err.Error()))
	}
	// 转换WorkFlows
	var workFlowInfos []*assistant_service.AssistantWorkFlowInfos
	for _, workflow := range workFlowConfig {
		workFlowInfos = append(workFlowInfos, &assistant_service.AssistantWorkFlowInfos{
			Id:         strconv.FormatUint(uint64(workflow.ID), 10),
			WorkFlowId: workflow.WorkflowId,
			Enable:     workflow.Enable,
		})
	}

	// 解析MCP配置详情
	var mcpInfoConfig []*model.AssistantMCP
	if err := jsonToStruct(snapshotInfo.AssistantMCPConfig, &mcpInfoConfig); err != nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", err.Error()))
	}
	// 转换MCP
	var mcpInfos []*assistant_service.AssistantMCPInfos
	for _, mcp := range mcpInfoConfig {
		mcpInfos = append(mcpInfos, &assistant_service.AssistantMCPInfos{
			Id:         strconv.FormatUint(uint64(mcp.ID), 10),
			McpId:      mcp.MCPId,
			McpType:    mcp.MCPType,
			ActionName: mcp.ActionName,
			Enable:     mcp.Enable,
		})
	}

	// 解析Tool配置详情
	var toolInfoConfig []*model.AssistantTool
	if err := jsonToStruct(snapshotInfo.AssistantToolConfig, &toolInfoConfig); err != nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", err.Error()))
	}
	// 转换Tool
	var toolInfos []*assistant_service.AssistantToolInfos
	for _, tool := range toolInfoConfig {
		toolInfos = append(toolInfos, &assistant_service.AssistantToolInfos{
			Id:         strconv.FormatUint(uint64(tool.ID), 10),
			ToolId:     tool.ToolId,
			ToolType:   tool.ToolType,
			ActionName: tool.ActionName,
			Enable:     tool.Enable,
			ToolConfig: tool.ToolConfig,
		})
	}

	// 转换ModelConfig
	var modelConfig *common.AppModelConfig
	if snapshotAssistant.ModelConfig != "" {
		modelConfig = &common.AppModelConfig{}
		if err := json.Unmarshal([]byte(snapshotAssistant.ModelConfig), modelConfig); err != nil {
			return nil, errStatus(errs.Code_AssistantErr, &errs.Status{
				TextKey: "assistant_modelConfig_unmarshal",
				Args:    []string{err.Error()},
			})
		}
	}

	// 转换RerankConfig
	var rerankConfig *common.AppModelConfig
	if snapshotAssistant.RerankConfig != "" {
		rerankConfig = &common.AppModelConfig{}
		if err := json.Unmarshal([]byte(snapshotAssistant.RerankConfig), rerankConfig); err != nil {
			return nil, errStatus(errs.Code_AssistantErr, &errs.Status{
				TextKey: "assistant_rerankConfig_unmarshal",
				Args:    []string{err.Error()},
			})
		}
	}

	// 转换KnowledgeBaseConfig
	var knowledgeBaseConfig *assistant_service.AssistantKnowledgeBaseConfig
	if snapshotAssistant.KnowledgebaseConfig != "" {
		knowledgeBaseConfig = &assistant_service.AssistantKnowledgeBaseConfig{}
		if err := json.Unmarshal([]byte(snapshotAssistant.KnowledgebaseConfig), knowledgeBaseConfig); err != nil {
			return nil, errStatus(errs.Code_AssistantErr, &errs.Status{
				TextKey: "assistant_knowledgeBaseConfig_unmarshal",
				Args:    []string{err.Error()},
			})
		}
	}

	// 转换SafetyConfig
	var safetyConfig *assistant_service.AssistantSafetyConfig
	if snapshotAssistant.SafetyConfig != "" {
		safetyConfig = &assistant_service.AssistantSafetyConfig{}
		if err := json.Unmarshal([]byte(snapshotAssistant.SafetyConfig), safetyConfig); err != nil {
			return nil, errStatus(errs.Code_AssistantErr, &errs.Status{
				TextKey: "assistant_safetyConfig_unmarshal",
				Args:    []string{err.Error()},
			})
		}
	}

	// 转换VisionConfig
	var visionConfig *assistant_service.AssistantVisionConfig
	if snapshotAssistant.VisionConfig != "" {
		visionConfig = &assistant_service.AssistantVisionConfig{}
		if err := json.Unmarshal([]byte(snapshotAssistant.VisionConfig), visionConfig); err != nil {
			return nil, errStatus(errs.Code_AssistantErr, &errs.Status{
				TextKey: "assistant_visionConfig_unmarshal",
				Args:    []string{err.Error()},
			})
		}
		visionConfig.MaxPicNum = config.Cfg().Assistant.MaxPicNum
	}

	return &assistant_service.AssistantInfo{
		AssistantId: util.Int2Str(snapshotAssistant.ID),
		Identity: &assistant_service.Identity{
			UserId: snapshotInfo.UserId,
			OrgId:  snapshotInfo.OrgId,
		},
		AssistantBrief: &common.AppBriefConfig{
			Name:       snapshotAssistant.Name,
			AvatarPath: snapshotAssistant.AvatarPath,
			Desc:       snapshotAssistant.Desc,
		},
		Prologue:            snapshotAssistant.Prologue,
		Instructions:        snapshotAssistant.Instructions,
		RecommendQuestion:   strings.Split(snapshotAssistant.RecommendQuestion, "@#@"),
		ModelConfig:         modelConfig,
		KnowledgeBaseConfig: knowledgeBaseConfig,
		RerankConfig:        rerankConfig,
		SafetyConfig:        safetyConfig,
		VisionConfig:        visionConfig,
		Scope:               int32(snapshotAssistant.Scope),
		WorkFlowInfos:       workFlowInfos,
		McpInfos:            mcpInfos,
		ToolInfos:           toolInfos,
		CreatTime:           snapshotAssistant.CreatedAt,
		UpdateTime:          snapshotAssistant.UpdatedAt,
	}, nil
}

func (s *Service) AssistantSnapshotRollback(ctx context.Context, req *assistant_service.AssistantSnapshotRollbackReq) (*emptypb.Empty, error) {
	assistantId := util.MustU32(req.AssistantId)
	version := req.Version
	if version == "" {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", "version is empty"))
	}

	// 获取指定版本的快照信息
	assistantSnapshot, status := s.cli.GetAssistantSnapshot(ctx, assistantId, version)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantErr, status)
	}

	// --- AssistantInfo ---
	var assistantInfo *model.Assistant
	if err := jsonToStruct(assistantSnapshot.AssistantInfo, &assistantInfo); err != nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", err.Error()))
	}
	if assistantInfo == nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", "assistant info is nil"))
	}

	// --- AssistantToolConfig ---
	var assistantToolConfig []*model.AssistantTool
	if err := jsonToStruct(assistantSnapshot.AssistantToolConfig, &assistantToolConfig); err != nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", err.Error()))
	}

	// --- AssistantMCPConfig ---
	var assistantMCPConfig []*model.AssistantMCP
	if err := jsonToStruct(assistantSnapshot.AssistantMCPConfig, &assistantMCPConfig); err != nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", err.Error()))
	}

	// --- AssistantWorkflowConfig ---
	var assistantWorkflowConfig []*model.AssistantWorkflow
	if err := jsonToStruct(assistantSnapshot.AssistantWorkflowConfig, &assistantWorkflowConfig); err != nil {
		return nil, errStatus(errs.Code_AssistantErr, toErrStatus("assistant_snapshot", err.Error()))
	}

	// 执行回滚事务
	status = s.cli.RollbackAssistantSnapshot(ctx, assistantInfo, assistantToolConfig, assistantMCPConfig, assistantWorkflowConfig, req.Identity.UserId, req.Identity.OrgId)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantErr, status)
	}

	return nil, nil
}

// --------------------- internal methods ---------------------

// json 字符串转结构体
func jsonToStruct(jsonStr string, v interface{}) error {
	if jsonStr == "" {
		return errors.New("json string is empty")
	}
	if err := json.Unmarshal([]byte(jsonStr), v); err != nil {
		log.Errorf("json unmarshal failed: %v", err)
		return err
	}
	return nil
}

// 结构体转json
func structToJson(v interface{}) string {
	if v == nil {
		return ""
	}
	// 即使结构体为空，也进行序列化，确保返回"{}"而不是空字符串
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		log.Errorf("json marshal failed: %v", err)
		return ""
	}
	return string(jsonBytes)
}
