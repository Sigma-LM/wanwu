package service

import (
	"context"
	"errors"
	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/grpc-consumer/consumer/assistant"
	"github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-processor"
	agent_preprocessor "github.com/UnicomAI/wanwu/internal/agent-service/service/agent-preprocessor"
	local_agent "github.com/UnicomAI/wanwu/internal/agent-service/service/local-agent"
	service_model "github.com/UnicomAI/wanwu/internal/agent-service/service/service-model"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/supervisor"
	"github.com/gin-gonic/gin"
)

type MultiAgent struct {
	MultiAgent      adk.Agent
	AgentPreprocess *agent_preprocessor.AgentPreprocess
	Req             *request.AgentChatParams
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
	agentChatParams := buildAgentChatParams(multiAgentChatReq)
	chatInfo, err := buildAgentChatInfo(ctx, agentChatParams)
	if err != nil {
		log.Errorf("failed to build chat info: %v", err)
		return nil, err
	}
	//构造智能体服务
	localAgentService := local_agent.CreateLocalAgentService(ctx, agentChatParams, chatInfo)
	//构造supervisor
	sv, err := buildSupervisor(ctx, localAgentService, agentChatParams, chatInfo)
	if err != nil {
		return nil, err
	}
	//构造子智能体
	multiSubAgent, err := buildMultiSubAgent(ctx, multiAgentChatReq)
	if err != nil {
		return nil, err
	}
	//构造多智能体
	agent, err := supervisor.New(ctx, &supervisor.Config{
		Supervisor: sv,
		SubAgents:  multiSubAgent,
	})
	if err != nil {
		return nil, err
	}
	return &MultiAgent{
		MultiAgent: agent,
		AgentPreprocess: &agent_preprocessor.AgentPreprocess{
			LocalAgentService: localAgentService,
			GinContext:        ctx,
			AgentChatInfo:     chatInfo,
		},
		Req: agentChatParams,
	}, nil
}

func (s *MultiAgent) Chat(ctx *gin.Context) error {
	//1.执行流式agent问答调用
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           s,
		EnableStreaming: s.Req.Stream,
	})
	iter := runner.Query(ctx, s.Req.Input)
	//2.处理结果
	err := agent_message_processor.MultiAgentMessage(ctx, iter, &request.AgentChatContext{AgentChatReq: s.Req})
	return err
}

func (s *MultiAgent) Name(ctx context.Context) string {
	return s.MultiAgent.Name(ctx)
}
func (s *MultiAgent) Description(ctx context.Context) string {
	return s.MultiAgent.Description(ctx)
}

func (s *MultiAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	log.Infof("[%s] multi agent run", s.Req.AgentBaseParams.Name)
	//参数预处理
	agentInput := agent_preprocessor.AgentPreProcess(s.AgentPreprocess, input, s.Req)
	return s.MultiAgent.Run(ctx, agentInput, options...)
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
		FilterSubEnable: true,
	})
	if err != nil {
		log.Errorf("failed to get multi assistant by id: %v", err)
		return nil, errors.New("failed to get multi assistant")
	}
	return multiAgent, nil
}

func buildAgentChatParams(multiAgentChatReq *request.MultiAgentChatReq) *request.AgentChatParams {
	var subAgentInfoList []*request.SubAgentInfo
	agentList := multiAgentChatReq.AgentList
	if len(agentList) > 0 {
		for _, subAgent := range agentList {
			subAgentInfoList = append(subAgentInfoList, &request.SubAgentInfo{
				Name:        subAgent.AgentBaseParams.Name,
				Description: subAgent.AgentBaseParams.Description,
			})
		}
	}
	return &request.AgentChatParams{
		MultiAgent:       true,
		SubAgentInfoList: subAgentInfoList,
		Input:            multiAgentChatReq.Input,
		UploadFile:       multiAgentChatReq.UploadFile,
		Stream:           multiAgentChatReq.Stream,
		AgentChatBaseParams: request.AgentChatBaseParams{
			ModelParams:     multiAgentChatReq.ModelParams,
			AgentBaseParams: multiAgentChatReq.AgentBaseParams,
		},
	}
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

func buildSupervisor(ctx *gin.Context, localAgentService local_agent.LocalAgentService, agentChatParams *request.AgentChatParams, chatInfo *service_model.AgentChatInfo) (adk.Agent, error) {
	//创建模型
	chatModel, err := localAgentService.CreateChatModel(ctx, agentChatParams, chatInfo)
	if err != nil {
		log.Errorf("failed to build chat model: %v", err)
		return nil, err
	}

	agentBaseParams := agentChatParams.AgentBaseParams
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        agentBaseParams.Name,
		Description: agentBaseParams.Description,
		Instruction: agentBaseParams.Instruction,
		Model:       chatModel,
		Exit:        &adk.ExitTool{},
	})
}
