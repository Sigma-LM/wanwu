package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// ModelExperienceLLM
//
//	@Tags			model.experience
//	@Summary		模型体验
//	@Description	LLM模型体验
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.ModelExperienceLlmRequest	true	"LLM模型体验"
//	@Success		200		{object}	response.Response
//	@Router			/model/experience/llm [post]
func ModelExperienceLLM(ctx *gin.Context) {
	var req request.ModelExperienceLlmRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	service.ModelExperienceLLM(ctx, getUserID(ctx), getOrgID(ctx), &req)
}

// ModelExperienceSaveDialog
//
//	@Tags			model.experience
//	@Summary		新建/保存对话
//	@Description	新建/保存对话
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.ModelExperienceDialogRequest	true	"模型体验对话"
//	@Success		200		{object}	response.Response{}
//	@Router			/model/experience/dialog [post]
func ModelExperienceSaveDialog(ctx *gin.Context) {
	var req request.ModelExperienceDialogRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.SaveModelExperienceDialog(ctx, getUserID(ctx), getOrgID(ctx), &req)
	gin_util.Response(ctx, resp, err)
}

// ModelExperienceListDialogs
//
//	@Tags			model.experience
//	@Summary		获取模型体验对话列表
//	@Description	获取模型体验对话列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=response.ListResult{list=model_service.ModelExperienceDialog}}
//	@Router			/model/experience/dialogs [get]
func ModelExperienceListDialogs(ctx *gin.Context) {
	resp, err := service.ListModelExperienceDialogs(ctx, getUserID(ctx), getOrgID(ctx))
	gin_util.Response(ctx, resp, err)
}

// ModelExperienceDeleteDialog
//
//	@Tags			model.experience
//	@Summary		删除模型体验对话
//	@Description	删除模型体验对话
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data						body		request.ModelExperienceDialogIDRequest	true	"模型体验对话ID"
//	@Success		200							{object}	response.Response
//	@Router			/model/experience/dialog	 [delete]
func ModelExperienceDeleteDialog(ctx *gin.Context) {
	var req request.ModelExperienceDialogIDRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.DeleteModelExperienceDialog(ctx, getUserID(ctx), getOrgID(ctx), req.ModelExperienceId)
	gin_util.Response(ctx, nil, err)
}

// ModelExperienceListDialogRecords
//
//	@Tags			model.experience
//	@Summary		获取模型体验对话记录列表
//	@Description	获取模型体验对话记录列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			modelExperienceId	query		uint32	true	"模型体验对话ID"
//	@Success		200					{object}	response.Response{data=response.ListResult{list=response.ModelExperienceDialogRecord}}
//	@Router			/model/experience/dialog/records [get]
func ModelExperienceListDialogRecords(ctx *gin.Context) {
	var req request.ModelExperienceDialogRecordRequest
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.ListModelExperienceDialogRecords(ctx, getUserID(ctx), getOrgID(ctx), &req)
	gin_util.Response(ctx, resp, err)
}
