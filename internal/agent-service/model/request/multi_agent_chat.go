package request

type MultiAgentChatParams struct {
	MultiAgentId   uint32   `json:"multiAgentId"  validate:"required"` //多智能体ID
	Input          string   `json:"input" validate:"required"`
	UserId         string   `json:"userId"  validate:"required"`
	OrgId          string   `json:"orgId"  validate:"required"`
	UploadFile     []string `json:"uploadFile"`
	Stream         bool     `json:"stream"`
	Draft          bool     `json:"draft"`
	ConversationId string   `json:"conversationId"` //会话ID
}

type MultiAgentChatReq struct {
	Input       string                 `json:"input"`
	UploadFile  []string               `json:"uploadFile"`
	Stream      bool                   `json:"stream"`
	ModelParams *ModelParams           `json:"modelParams"` // 模型参数
	AgentList   []*AgentChatBaseParams `json:"agentList"`
}

func (c *MultiAgentChatParams) Check() error {
	return nil
}
