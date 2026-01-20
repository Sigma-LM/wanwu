package service

import (
	"errors"
	"fmt"
	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/grpc-consumer/consumer/assistant"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/util"
	"github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-flow/prompt"
	"github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-processor"
	local_agent "github.com/UnicomAI/wanwu/internal/agent-service/service/local-agent"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/supervisor"
	"github.com/cloudwego/eino/components/model"
	"github.com/gin-gonic/gin"
	"strings"
)

type MultiAgent struct {
	MultiAgent  adk.Agent
	Stream      bool
	Input       string
	UploadFiles []string
}

// MultiAgentChat 多智能体问答
func MultiAgentChat(ctx *gin.Context, req *request.MultiAgentChatParams) error {
	multiAgentConfig, err := searchMultiAgentConfig(ctx, req)
	if err != nil {
		return err
	}
	agent, err := CreateSupervisorMultiAgent(ctx, req, multiAgentConfig)
	if err != nil {
		return err
	}
	return agent.Chat(ctx)
}

func CreateSupervisorMultiAgent(ctx *gin.Context, multiAgentChatParams *request.MultiAgentChatParams, multiAgentConfig *assistant_service.MultiAssistantDetailResp) (*MultiAgent, error) {
	var multiAgentChatReq = BuildMultiAgentParams(multiAgentChatParams, multiAgentConfig)
	chatModel, err := buildSupervisorChatModel(ctx, multiAgentChatReq)
	if err != nil {
		log.Errorf("failed to build chat model: %v", err)
		return nil, err
	}

	sv, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "supervisor",
		Description: "the agent responsible to supervise tasks",
		Instruction: buildMultiAgentPrompt(multiAgentChatReq),
		Model:       chatModel,
		Exit:        &adk.ExitTool{},
	})
	if err != nil {
		return nil, err
	}

	multiSubAgent, err := buildMultiSubAgent(ctx, multiAgentChatReq)
	if err != nil {
		return nil, err
	}
	agent, err := supervisor.New(ctx, &supervisor.Config{
		Supervisor: sv,
		SubAgents:  multiSubAgent,
	})
	if err != nil {
		return nil, err
	}
	return &MultiAgent{
		MultiAgent:  agent,
		Stream:      multiAgentChatReq.Stream,
		Input:       multiAgentChatReq.Input,
		UploadFiles: multiAgentChatReq.UploadFile,
	}, nil
}

func (s *MultiAgent) Chat(ctx *gin.Context) error {
	//1.执行流式agent问答调用
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           s.MultiAgent,
		EnableStreaming: s.Stream,
	})
	iter := runner.Query(ctx, s.Input)

	//2.处理结果,todo 还需要详细调试处理结果
	err := agent_message_processor.AgentMessage(ctx, iter, &request.AgentChatContext{AgentChatReq: &request.AgentChatParams{
		Input:      s.Input,
		Stream:     s.Stream,
		UploadFile: s.UploadFiles,
	}})
	return err
}

// searchMultiAgentConfig 查询多智能体配置信息
func searchMultiAgentConfig(ctx *gin.Context, req *request.MultiAgentChatParams) (*assistant_service.MultiAssistantDetailResp, error) {
	multiAgent, err := assistant.GetClient().GetMultiAssistantById(ctx, &assistant_service.GetMultiAssistantByIdReq{
		AssistantId:    req.MultiAgentId,
		ConversationId: req.ConversationId,
		Draft:          req.Draft,
		Identity: &assistant_service.Identity{
			UserId: req.UserId,
			OrgId:  req.OrgId,
		},
	})
	if err != nil {
		log.Errorf("failed to get multi assistant by id: %v", err)
		return nil, errors.New("failed to get multi assistant")
	}
	return multiAgent, nil
}

func buildSupervisorChatModel(ctx *gin.Context, multiAgentChatReq *request.MultiAgentChatReq) (model.ToolCallingChatModel, error) {
	req := &request.AgentChatParams{
		AgentChatBaseParams: request.AgentChatBaseParams{
			ModelParams: multiAgentChatReq.ModelParams,
		},
	}
	chatInfo, err := buildAgentChatInfo(ctx, req)
	if err != nil {
		log.Errorf("failed to build chat info: %v", err)
		return nil, err
	}
	return local_agent.CreateChatModel(ctx, chatInfo, req)
}

func buildMultiAgentPrompt(multiAgentChatReq *request.MultiAgentChatReq) string {
	return fmt.Sprintf(prompt.SupervisorPrompt, util.EnglishNumber(len(multiAgentChatReq.AgentList)), buildAgentDesc(multiAgentChatReq))
}

func buildAgentDesc(multiAgentChatReq *request.MultiAgentChatReq) string {
	builder := strings.Builder{}
	for _, params := range multiAgentChatReq.AgentList {
		builder.WriteString(fmt.Sprintf(prompt.SupervisorAgentTemplate, params.AgentBaseParams.Name, params.AgentBaseParams.Description))
	}
	return builder.String()
}

func buildMultiSubAgent(ctx *gin.Context, multiAgentChatReq *request.MultiAgentChatReq) ([]adk.Agent, error) {
	var subAgents []adk.Agent
	for _, agentParams := range multiAgentChatReq.AgentList {
		subAgent, err := CreateSingleAgent(ctx, &request.AgentChatParams{
			AgentChatBaseParams: *agentParams,
			Stream:              multiAgentChatReq.Stream,
			UploadFile:          multiAgentChatReq.UploadFile,
		})
		if err != nil {
			return nil, err
		}
		subAgents = append(subAgents, subAgent)
	}
	return subAgents, nil
}
