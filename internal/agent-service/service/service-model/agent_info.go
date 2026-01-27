package service_model

type AgentChatInfo struct {
	ModelInfo       *ModelInfo
	FunctionCalling bool `json:"functionCalling"` // 是否支持functionCall
	VisionSupport   bool `json:"visionSupport"`   // 是否支持多模态
	UploadUrl       bool `json:"uploadUrl"`       // 是否上传文件
}
