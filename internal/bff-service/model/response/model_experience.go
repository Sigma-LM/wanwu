package response

type ModelExperienceDialog struct {
	ID           string `json:"id"`
	ModelId      string `json:"modelId"`
	SessionId    string `json:"sessionId"`
	Title        string `json:"title"`
	ModelSetting string `json:"modelSetting"`
	CreatedAt    string `json:"createdAt"`
}

type ModelExperienceFile struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	FileExt  string `json:"fileExt"`
	FileSize int64  `json:"fileSize"`
}

type ModelExperienceDialogRecord struct {
	Id                string                `json:"id"`
	ModelExperienceId string                `json:"modelExperienceId"` // 模型体验 ID
	ModelId           string                `json:"modelId"`           // 模型 ID
	OriginalContent   string                `json:"originalContent"`   // 原始内容
	ReasoningContent  string                `json:"reasoningContent"`  // 思考过程
	Role              string                `json:"role"`              // 角色
	ParentID          string                `json:"parentID"`          // 父级 ID
	FileList          []ModelExperienceFile `json:"fileList"`          // 文件列表
}
