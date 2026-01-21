package assistant

import (
	"context"
	"encoding/json"
	"errors"
	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/config"
	"github.com/UnicomAI/wanwu/internal/assistant-service/service"
	params_process "github.com/UnicomAI/wanwu/internal/assistant-service/service/params-process"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	http_client "github.com/UnicomAI/wanwu/pkg/http-client"
	"github.com/UnicomAI/wanwu/pkg/log"
	sse_util "github.com/UnicomAI/wanwu/pkg/sse-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

func (s *Service) GetMultiAssistantById(ctx context.Context, req *assistant_service.GetMultiAssistantByIdReq) (*assistant_service.MultiAssistantDetailResp, error) {
	agent, agentSnapshot, subAgents, err := s.cli.GetMultiAssistant(ctx, req.AssistantId, req.Identity.GetUserId(), req.Identity.GetOrgId(), req.Draft, req.Version)
	if err != nil {
		return nil, err
	}
	multiAgentParams, err := buildMultiAgentParams(ctx, s.cli, agent, agentSnapshot, req)
	if err != nil {
		return nil, err
	}
	var subParamsList []*assistant_service.AgentDetail
	for _, subAgent := range subAgents {
		params, err := buildSubAgentParams(ctx, s.cli, subAgent)
		if err != nil {
			return nil, err
		}
		subParamsList = append(subParamsList, params)
	}
	return &assistant_service.MultiAssistantDetailResp{
		MultiAgent: multiAgentParams,
		SubAgents:  subParamsList,
	}, nil
}

func (s *Service) MultiAssistantConversionStream(req *assistant_service.MultiAssistantConversionStreamReq, stream grpc.ServerStreamingServer[assistant_service.AssistantConversionStreamResp]) error {
	//会话处理
	conversationProcessor := &service.ConversationProcessor{
		SSEWriter: sse_util.NewGrpcSSEWriter(stream, "MultiAssistantConversionStreamNew", nil),
	}
	err := conversationProcessor.Process(stream.Context(), buildMultiConversationParams(req), buildMultiAgentSendRequest(req))
	if err != nil {
		log.Errorf("Assistant服务处理智能体流式对话失败，assistantId: %s, error: %v", req.AssistantId, err)
		return grpc_util.ErrorStatusWithKey(errs.Code_AssistantConversationErr, "assistant_conversation", "agent服务异常")
	}
	return nil
}

func buildMultiAgentParams(ctx context.Context, cli client.IClient, agent *model.Assistant, snapshot *model.AssistantSnapshot, req *assistant_service.GetMultiAssistantByIdReq) (*assistant_service.AgentDetail, error) {
	clientInfo := &params_process.ClientInfo{
		Cli:       cli,
		Knowledge: Knowledge,
		MCP:       MCP,
	}
	//传入了 ConversationId就会尝试构造历史数据
	userQueryParams := &params_process.UserQueryParams{
		ConversationId: req.ConversationId,
		QueryOrgId:     req.Identity.OrgId,
		QueryUserId:    req.Identity.UserId,
	}
	return service.NewAgentChatParamsBuilder(&params_process.AgentInfo{
		Assistant:         agent,
		AssistantSnapshot: snapshot,
		Draft:             snapshot == nil,
	}, userQueryParams, clientInfo).
		AgentBaseParams().
		ModelParams().
		Build()
}

// buildSubAgentParams 构建子智能体参数，由于子智能体只能选发布后的智能体所以结构体是snapshot
func buildSubAgentParams(ctx context.Context, cli client.IClient, agentSnapshot *model.AssistantSnapshot) (*assistant_service.AgentDetail, error) {
	clientInfo := &params_process.ClientInfo{
		Cli:       cli,
		Knowledge: Knowledge,
		MCP:       MCP,
	}
	var assistant *model.Assistant
	if err := jsonToStruct(agentSnapshot.AssistantInfo, assistant); err != nil {
		log.Errorf("转换智能体信息失败，assistantId: %d, error: %v", agentSnapshot.AssistantID, err)
		return nil, errors.New("build agent info err")
	}
	return service.NewAgentChatParamsBuilder(&params_process.AgentInfo{
		Assistant:         assistant,
		AssistantSnapshot: agentSnapshot,
		Draft:             agentSnapshot == nil,
	}, nil, clientInfo).
		AgentBaseParams().
		ModelParams().
		KnowledgeParams().
		ToolParams().
		Build()
}

func buildMultiConversationParams(req *assistant_service.MultiAssistantConversionStreamReq) *service.ConversationParams {
	return &service.ConversationParams{
		AssistantId:    req.AssistantId,
		ConversationId: req.ConversationId,
		FileInfo:       extractFileInfos(req.FileInfo),
		OrgId:          req.Identity.OrgId,
		Query:          req.Prompt,
		UserId:         req.Identity.UserId,
	}
}

// buildMultiAgentSendRequest 构建底层智能体能力接口请求体
func buildMultiAgentSendRequest(req *assistant_service.MultiAssistantConversionStreamReq) func(ctx context.Context) (string, *http.Response, context.CancelFunc, error) {
	var conversationID string
	// 历史聊天记录配置
	if !req.Trial && req.ConversationId != "" {
		conversationID = req.ConversationId
	}
	// 底层智能体能力接口请求体
	chatReq := service.BuildMultiAgentChatReq(&service.AgentUserInputParams{
		Input:          req.Prompt,
		Stream:         true,
		UploadFile:     extractFileUrls(req.FileInfo),
		ConversationId: conversationID,
		UserId:         req.Identity.UserId,
		OrgId:          req.Identity.OrgId,
		Draft:          req.Draft,
	}, util.MustU32(req.AssistantId))
	var monitorKey = "multi_agent_chat_service"

	return func(ctx context.Context) (string, *http.Response, context.CancelFunc, error) {
		paramsBytes, err := json.Marshal(chatReq)
		if err != nil {
			return monitorKey, nil, nil, err
		}
		// 获取Assistant配置
		assistantConfig := config.Cfg().Assistant
		if assistantConfig.NewSseUrl == "" {
			return monitorKey, nil, nil, errors.New("智能体SSE URL配置错误")
		}
		params := &http_client.HttpRequestParams{
			Body:       paramsBytes,
			Timeout:    5 * time.Minute,
			Url:        assistantConfig.NewSseUrl,
			MonitorKey: monitorKey,
			LogLevel:   http_client.LogAll,
		}
		ctx, cancel := context.WithTimeout(ctx, params.Timeout)
		result, err := http_client.Default().PostJsonOriResp(ctx, params)
		return monitorKey, result, cancel, err
	}
}
