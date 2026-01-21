package service

import (
	"context"
	"encoding/json"
	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/grpc-consumer/consumer/assistant"
	"github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-processor"
	local_agent "github.com/UnicomAI/wanwu/internal/agent-service/service/local-agent"
	"github.com/UnicomAI/wanwu/internal/agent-service/service/service-model"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/gin-gonic/gin"
)

type SingleAgent struct {
	LocalAgentService local_agent.LocalAgentService
	ChatModelAgent    *adk.ChatModelAgent
	Req               *request.AgentChatParams
	AgentChatInfo     *service_model.AgentChatInfo
	GinContext        *gin.Context
}

// SingleAgentChat 单智能体问答
func SingleAgentChat(ctx *gin.Context, req *request.AgentChatReq) error {
	singleAgentDetail, err := searchSingleAgent(ctx, req)
	if err != nil {
		return err
	}
	agent, err := CreateSingleAgent(ctx, BuildAgentParams(req, singleAgentDetail))
	if err != nil {
		return err
	}
	return agent.Chat(ctx)
}

// CreateSingleAgent 创建单智能体
func CreateSingleAgent(ctx *gin.Context, req *request.AgentChatParams) (*SingleAgent, error) {
	data, _ := json.Marshal(req)
	log.Infof("single agent chat req %s", string(data))
	chatInfo, err := buildAgentChatInfo(ctx, req)
	if err != nil {
		log.Errorf("failed to build chat info: %v", err)
		return nil, err
	}
	localAgentService := local_agent.CreateLocalAgentService(ctx, req, chatInfo)
	//创建模型
	chatModel, err := localAgentService.CreateChatModel(ctx, req, chatInfo)
	if err != nil {
		log.Errorf("failed to create chat model: %v", err)
		return nil, err
	}
	//2.创建智能体
	agent, err := createAgent(ctx, req, chatModel)
	if err != nil {
		log.Errorf("failed to create agent: %v", err)
		return nil, err
	}
	return &SingleAgent{
		LocalAgentService: localAgentService,
		ChatModelAgent:    agent,
		Req:               req,
		AgentChatInfo:     chatInfo,
		GinContext:        ctx,
	}, nil
}

func (s *SingleAgent) Chat(ctx *gin.Context) error {
	//1.执行流式agent问答调用
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           s,
		EnableStreaming: s.Req.Stream,
	})
	iter := runner.Query(ctx, s.Req.Input)

	//2.处理结果
	err := agent_message_processor.AgentMessage(ctx, iter, &request.AgentChatContext{AgentChatReq: s.Req})
	return err
}

func (s *SingleAgent) Name(ctx context.Context) string {
	return s.ChatModelAgent.Name(ctx)
}
func (s *SingleAgent) Description(ctx context.Context) string {
	return s.ChatModelAgent.Description(ctx)
}

func (s *SingleAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	//todo 删除
	log.Infof("single agent run")

	iter, generator := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	var agentInput *adk.AgentInput
	var err error
	go func() {
		defer func() {
			generator.Close()
		}()
		defer util.PrintPanicStack()
		if s.Req.AgentBaseParams.CallDetail { //是否输出调用详情
			agentInput, err = s.LocalAgentService.BuildAgentInput(ctx, s.Req, s.AgentChatInfo, input, generator)
		} else {
			agentInput, err = s.LocalAgentService.BuildAgentInput(ctx, s.Req, s.AgentChatInfo, input, nil)
		}
		if err != nil {
			log.Errorf("failed to build agent input: %v", err)
			generator.Send(&adk.AgentEvent{Err: err})
		}
	}()
	_ = agent_message_processor.AgentMessage(s.GinContext, iter, &request.AgentChatContext{AgentChatReq: s.Req})

	marshal, _ := json.Marshal(agentInput)
	log.Infof("agent input %s", string(marshal))
	return s.ChatModelAgent.Run(ctx, agentInput, options...)
}

// buildAgentChatInfo 构建智能体信息
func buildAgentChatInfo(ctx *gin.Context, req *request.AgentChatParams) (*service_model.AgentChatInfo, error) {
	modelInfo, err := SearchModel(ctx, req.ModelParams.ModelId)
	if err != nil {
		return nil, err
	}
	var functionCall = modelInfo.Config.FunctionCalling != "noSupport"
	var vision = modelInfo.Config.VisionSupport == "support"
	return &service_model.AgentChatInfo{
		FunctionCalling: functionCall,
		VisionSupport:   vision,
		UploadUrl:       len(req.UploadFile) > 0,
		ModelInfo:       modelInfo,
	}, nil
}

// searchSingleAgent 查询智能体详情
func searchSingleAgent(ctx *gin.Context, req *request.AgentChatReq) (*assistant_service.AssistantDetailResp, error) {
	return assistant.GetClient().GetAssistantDetailById(ctx, &assistant_service.GetAssistantDetailByIdReq{
		AssistantId:    req.AssistantId,
		ConversationId: req.ConversationId,
		Draft:          req.Draft,
		Identity: &assistant_service.Identity{
			UserId: req.UserId,
			OrgId:  req.OrgId,
		},
	})
}

// 创建对应智能体
func createAgent(ctx *gin.Context, req *request.AgentChatParams, chatModel model.ToolCallingChatModel) (*adk.ChatModelAgent, error) {
	baseParams := req.AgentBaseParams
	toolsConfig, err := BuildAgentToolsConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Model:       chatModel,
		Name:        baseParams.Name,
		Description: baseParams.Description,
		Instruction: baseParams.Instruction,
		ToolsConfig: toolsConfig,
	})
}
