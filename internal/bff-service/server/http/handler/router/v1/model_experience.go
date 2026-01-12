package v1

import (
	"net/http"

	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func registerModelExperience(apiV1 *gin.RouterGroup) {
	mid.Sub("model").Reg(apiV1, "/model/experience/llm", http.MethodPost, v1.ModelExperienceLLM, "LLM模型体验")
	mid.Sub("model").Reg(apiV1, "/model/experience/dialog", http.MethodPost, v1.ModelExperienceSaveDialog, "保存模型体验对话")
	mid.Sub("model").Reg(apiV1, "/model/experience/file/extract", http.MethodPost, v1.ModelExperienceFileExtract, "文本提取")
	mid.Sub("model").Reg(apiV1, "/model/experience/dialogs", http.MethodGet, v1.GetDialogs, "获取模型体验对话列表")
	mid.Sub("model").Reg(apiV1, "/model/experience/dialog", http.MethodDelete, v1.DeleteDialog, "删除模型体验对话")
	mid.Sub("model").Reg(apiV1, "/model/experience/dialog/records", http.MethodGet, v1.GetDialogRecords, "查询模型体验对话历史记录")

}
