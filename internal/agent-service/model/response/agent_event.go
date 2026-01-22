package response

type AgentEventType int
type SubAgentStatus int

const (
	MainAgentEventType = 0 //单智能体事件/多智能体主智能体
	SubAgentEventType  = 1 //子智能体事件

	SubAgentStartStatus   = 1 //开始事件
	SubAgentProcessStatus = 2 //输出中
	SubAgentEndStatus     = 3 //结束事件
	SubAgentFailStatus    = 4 //子智能体失败
)

type SubAgentEventData struct {
	Status    SubAgentStatus `json:"status"`
	AgentId   string         `json:"agentId"`
	AgentName string         `json:"agentName"`
	TimeCost  string         `json:"timeCost"`
}

func BuildStartSubAgent(agentId, agentName string) *SubAgentEventData {
	return &SubAgentEventData{
		Status:    SubAgentStartStatus,
		AgentId:   agentId,
		AgentName: agentName,
	}
}

func BuildProcessSubAgent(agentId, agentName string) *SubAgentEventData {
	if len(agentId) == 0 || len(agentName) == 0 {
		return nil
	}
	return &SubAgentEventData{
		Status:    SubAgentProcessStatus,
		AgentId:   agentId,
		AgentName: agentName,
	}
}

func BuildEndSubAgent(agentId, agentName, timeCost string) *SubAgentEventData {
	return &SubAgentEventData{
		Status:    SubAgentEndStatus,
		AgentId:   agentId,
		AgentName: agentName,
		TimeCost:  timeCost,
	}
}
