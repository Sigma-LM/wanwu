package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// ModelExperienceLLM
//
//	@Tags			model experience
//	@Summary		模型体验
//	@Description	LLM模型体验
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.LlmRequest	true	"LLM模型体验"
//	@Success		200		{object}	response.Response
//	@Router			/model/experience/llm [post]
func ModelExperienceLLM(ctx *gin.Context) {
	var req request.LlmRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}

	userId := ctx.GetString(gin_util.USER_ID)
	orgId := ctx.GetHeader(gin_util.X_ORG_ID)

	service.LLMModelExperience(ctx, &req, userId, orgId)
}

// ModelExperienceSaveDialog
//
//	@Tags			model experience
//	@Summary		新建/保存对话
//	@Description	新建/保存对话
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.ModelExperienceDialogRequest	true	"模型体验对话"
//	@Success		200		{object}	response.Response{data=response.ModelExperienceDialog}
//	@Router			/model/experience/dialog [post]
func ModelExperienceSaveDialog(ctx *gin.Context) {
	var req request.ModelExperienceDialogRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}

	userId := ctx.GetString(gin_util.USER_ID)
	orgId := ctx.GetHeader(gin_util.X_ORG_ID)

	// 调用服务层处理LLM请求
	res, err := service.SaveModelExperienceDialog(ctx, &req, userId, orgId)
	gin_util.Response(ctx, res, err)
}

// ModelExperienceFileExtract
//
//	@Tags			model experience
//	@Summary		文本提取
//	@Description	文本提取
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.FileExtractRequest	true	"文件提取请求"
//	@Success		200		{object}	response.Response{data=response.ModelExperienceFile}
//	@Router			/model/experience/file/extract [post]
func ModelExperienceFileExtract(ctx *gin.Context) {
	var req request.FileExtractRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}

	userId := ctx.GetString(gin_util.USER_ID)
	orgId := ctx.GetHeader(gin_util.X_ORG_ID)

	res, err := service.ModelExperienceFileExtract(ctx, &req, userId, orgId)
	gin_util.Response(ctx, res, err)
}

// GetDialogs
//
//	@Tags			model experience
//	@Summary		获取模型体验对话列表
//	@Description	获取模型体验对话列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=response.ListResult{list=model_service.ModelExperienceDialog}}
//	@Router			/model/experience/dialogs [get]
func GetDialogs(ctx *gin.Context) {

	userId := ctx.GetString(gin_util.USER_ID)
	orgId := ctx.GetHeader(gin_util.X_ORG_ID)

	// 调用服务层处理LLM请求
	res, err := service.GetModelExperienceDialogs(ctx, userId, orgId)
	gin_util.Response(ctx, res, err)
}

// DeleteDialog
//
//	@Tags			model experience
//	@Summary		删除模型体验对话
//	@Description	删除模型体验对话
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data						body		request.ModelExperienceDialogIDReq	true	"模型体验对话ID"
//	@Success		200							{object}	response.Response
//	@Router			/model/experience/dialog	 [delete]
func DeleteDialog(ctx *gin.Context) {
	var req request.ModelExperienceDialogIDReq
	if !gin_util.Bind(ctx, &req) {
		return
	}

	userId := ctx.GetString(gin_util.USER_ID)
	orgId := ctx.GetHeader(gin_util.X_ORG_ID)

	// 调用服务层处理LLM请求
	err := service.DeleteDialog(ctx, &req, userId, orgId)
	gin_util.Response(ctx, nil, err)
}

// GetDialogRecords
//
//	@Tags			model experience
//	@Summary		获取模型体验对话记录列表
//	@Description	获取模型体验对话记录列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			modelExperienceId	query		uint32	true	"模型体验对话ID"
//	@Success		200					{object}	response.Response{data=response.ListResult{list=response.ModelExperienceDialogRecord}}
//	@Router			/model/experience/dialog/records [get]
func GetDialogRecords(ctx *gin.Context) {
	var req request.ModelExperienceDialogRecordRequest
	if !gin_util.BindQuery(ctx, &req) {
		return
	}

	userId := ctx.GetString(gin_util.USER_ID)
	orgId := ctx.GetHeader(gin_util.X_ORG_ID)

	// 调用服务层处理LLM请求
	res, err := service.GetModelExperienceDialogRecords(ctx, &req, userId, orgId)
	gin_util.Response(ctx, res, err)
}
