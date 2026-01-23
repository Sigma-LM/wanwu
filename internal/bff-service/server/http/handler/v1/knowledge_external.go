package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GetKnowledgeExternalAPIList
//
//	@Tags			knowledge
//	@Summary		获取外部知识库API列表
//	@Description	获取外部知识库API列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=response.KnowledgeExternalAPIListResp}
//	@Router			/knowledge/external/api/select [get]
func GetKnowledgeExternalAPIList(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	resp, err := service.GetKnowledgeExternalAPIList(ctx, userId, orgId)
	gin_util.Response(ctx, resp, err)
}

// CreateKnowledgeExternalAPI
//
//	@Tags			knowledge
//	@Summary		创建外部知识库API
//	@Description	创建外部知识库API
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CreateKnowledgeExternalAPIReq	true	"创建外部知识库API请求参数"
//	@Success		200		{object}	response.Response{data=response.CreateKnowledgeExternalAPIResp}
//	@Router			/knowledge/external/api [post]
func CreateKnowledgeExternalAPI(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.CreateKnowledgeExternalAPIReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CreateKnowledgeExternalAPI(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// UpdateKnowledgeExternalAPI
//
//	@Tags			knowledge
//	@Summary		编辑外部知识库API
//	@Description	编辑外部知识库API
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.UpdateKnowledgeExternalAPIReq	true	"修改外部知识库API请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/external/api [put]
func UpdateKnowledgeExternalAPI(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.UpdateKnowledgeExternalAPIReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.UpdateKnowledgeExternalAPI(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// DeleteKnowledgeExternalAPI
//
//	@Tags			knowledge
//	@Summary		删除外部知识库API
//	@Description	删除外部知识库API
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.DeleteKnowledgeExternalAPIReq	true	"删除外部知识库API请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/external/api [delete]
func DeleteKnowledgeExternalAPI(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.DeleteKnowledgeExternalAPIReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.DeleteKnowledgeExternalAPI(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// GetKnowledgeExternalList
//
//	@Tags			knowledge
//	@Summary		获取外部知识库列表
//	@Description	获取外部知识库列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.KnowledgeExternalListReq	true	"获取外部知识库列表请求参数"
//	@Success		200		{object}	response.Response{data=response.KnowledgeExternalListResp}
//	@Router			/knowledge/external/select [get]
func GetKnowledgeExternalList(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.KnowledgeExternalListReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetKnowledgeExternalList(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// CreateKnowledgeExternal
//
//	@Tags			knowledge
//	@Summary		创建外部知识库
//	@Description	创建外部知识库
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CreateKnowledgeExternalReq	true	"创建外部知识库请求参数"
//	@Success		200		{object}	response.Response{data=response.CreateKnowledgeExternalResp}
//	@Router			/knowledge/external [post]
func CreateKnowledgeExternal(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.CreateKnowledgeExternalReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CreateKnowledgeExternal(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// UpdateKnowledgeExternal
//
//	@Tags			knowledge
//	@Summary		编辑外部知识库
//	@Description	编辑外部知识库
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.UpdateKnowledgeExternalReq	true	"编辑外部知识库请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/external [put]
func UpdateKnowledgeExternal(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.UpdateKnowledgeExternalReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.UpdateKnowledgeExternal(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// DeleteKnowledgeExternal
//
//	@Tags			knowledge
//	@Summary		删除外部知识库
//	@Description	删除外部知识库
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.DeleteKnowledgeExternalReq	true	"删除外部知识库请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/external [delete]
func DeleteKnowledgeExternal(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.DeleteKnowledgeExternalReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.DeleteKnowledgeExternal(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}
