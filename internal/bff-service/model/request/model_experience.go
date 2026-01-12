package request

import (
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
)

// LlmRequest LLM模型体验请求结构
type LlmRequest struct {
	ModelId           string   `json:"modelId" validate:"required"`   // 模型ID
	SessionId         string   `json:"sessionId" validate:"required"` //会话 ID
	ModelExperienceId string   `json:"modelExperienceId"`             // 体验对话ID
	ParentId          string   `json:"parentId"`                      // 父级ID
	Role              string   `json:"role" validate:"required"`      // 角色(user, assistant, system)
	Content           string   `json:"content" validate:"required"`   // 内容
	FileIdList        []string `json:"fileIdList"`                    // 文件ID列表
	mp_common.LLMParams
}

func (o *LlmRequest) Check() error {
	return nil
}

type ModelExperienceDialogRequest struct {
	ModelId      string      `json:"modelId" validate:"required"`   // 模型
	SessionId    string      `json:"sessionId" validate:"required"` //会话 ID
	Title        string      `json:"title"`                         // 对话标题
	ModelSetting interface{} `json:"modelSetting"`                  //模型参数配置
}

func (o *ModelExperienceDialogRequest) Check() error {
	return nil
}

type ModelExperienceDialogIDReq struct {
	ModelExperienceId string `json:"modelExperienceId" validate:"required"` // 模型体验对话ID
}

func (o *ModelExperienceDialogIDReq) Check() error {
	return nil
}

type FileExtractRequest struct {
	FileName string `json:"fileName" form:"fileName" ` // 文件名
	FilePath string `json:"filePath" form:"filePath" ` // 文件路径
	FileSize int64  `json:"fileSize" form:"fileSize" ` // 文件大小
}

func (o *FileExtractRequest) Check() error {
	return nil
}

type ModelExperienceDialogRecordRequest struct {
	ModelExperienceId string `json:"modelExperienceId" form:"modelExperienceId" validate:"required"` // 体验对话ID
}

func (o *ModelExperienceDialogRecordRequest) Check() error {
	return nil
}
