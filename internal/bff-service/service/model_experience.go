package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/minio"
	"github.com/UnicomAI/wanwu/pkg/util"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/pkg/constant"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/log"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/gin-gonic/gin"
)

// LLMModelExperience 处理LLM模型体验请求
func LLMModelExperience(ctx *gin.Context, req *request.LlmRequest, userId, orgId string) {
	log.Infof("开始处理LLM模型体验请求: %+v", req)

	// 获取会话历史
	var messages []mp_common.OpenAIReqMsg

	// 首先检查req是否为nil
	if req == nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, "Request parameters cannot be empty"))
		return
	}

	// modelInfo by modelID
	modelInfo, err := model.GetModel(ctx.Request.Context(), &model_service.GetModelReq{ModelId: req.ModelId})
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	// 确保modelInfo不为nil
	if modelInfo == nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("Model not found: %v", req.ModelId)))
		return
	}

	// 从数据库获取历史对话记录
	dbMessages, err := getHistoricalDialogRecords(ctx.Request.Context(), req.ModelExperienceId, req.SessionId, modelInfo.ModelId, userId, orgId)
	if err != nil {
		log.Errorf("获取历史对话记录失败: %v", err)
	} else if len(dbMessages) > 0 {
		// 使用数据库中的历史记录
		messages = append(messages, dbMessages...)
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(modelInfo.ProviderConfig), &result)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}

	// VisionSupport 支持校验
	visionSupportFlag := false // VisionSupport 校验标识
	vs, ok := result["visionSupport"].(string)
	if ok && mp_common.VSType(vs) == mp_common.VSTypeSupport {
		visionSupportFlag = true
	}

	// 处理文件ID，获取提取的文本并改写提示词
	content := req.Content

	imageContents := []map[string]interface{}{}
	if len(req.FileIdList) > 0 {
		// 将string类型的FileIdList转换为uint32类型
		fileIds := make([]string, 0, len(req.FileIdList))
		for _, idStr := range req.FileIdList {
			fileIds = append(fileIds, idStr)
		}

		// 调用GetModelExperienceFilesByIds获取文件内容
		fileResp, err := model.GetModelExperienceFilesByIds(ctx.Request.Context(), &model_service.GetModelExperienceFilesByIdsReq{
			FileIds: fileIds,
		})
		if err != nil {
			log.Errorf("Failed to get file content: %v", err)
			// Continue using original content on error
		} else {
			// 构建合并后的文档内容
			var context string
			for _, file := range fileResp.Files {
				// 使用strings.Replace替换模板中的变量
				fileContent := strings.Replace(constant.FileItemTemplate, "{{.FileName}}", file.FileName, -1)
				fileContent = strings.Replace(fileContent, "{{.FileContent}}", file.ExtractText, -1)
				context += fileContent
				fileExt := file.FileExt
				// 使用集合存储支持的图片扩展名，提高查找效率
				supportedImageExts := map[string]bool{
					".jpg":  true,
					".jpeg": true,
					".png":  true,
					".gif":  true,
					".bmp":  true,
					".webp": true,
				}
				if supportedImageExts[fileExt] && visionSupportFlag {
					// 从minio下载文件
					fileData, err := minio.DownloadFileToMemory(ctx.Request.Context(), file.FilePath)
					if err != nil {
						log.Errorf("从minio下载文件失败: %v", err)
						continue
					}

					// 转换为base64
					imageBase64 := base64.StdEncoding.EncodeToString(fileData)
					mimeType := "image/" + strings.TrimPrefix(fileExt, ".")
					if fileExt == ".jpg" {
						mimeType = "image/jpeg"
					}
					dataUrl := fmt.Sprintf("data:%s;base64,%s", mimeType, imageBase64)

					// 添加到imageContents
					imageContent := map[string]interface{}{
						"type": "image_url",
						"image_url": map[string]string{
							"url": dataUrl,
						},
					}
					imageContents = append(imageContents, imageContent)

				}
			}

			// 使用常量中的模板
			content = strings.Replace(constant.ModelExperienceTemplate, "{{.Context}}", context, -1)
			content = strings.Replace(content, "{{.Question}}", req.Content, -1)

			var result map[string]interface{}
			err = json.Unmarshal([]byte(modelInfo.ProviderConfig), &result)
			if err != nil {
				gin_util.Response(ctx, nil, err)
				return
			}

		}
	}

	// 如果imageContents不为空，将文本content和imageContents合并
	if len(imageContents) > 0 {
		// 创建文本内容
		textContent := map[string]interface{}{
			"type": "text",
			"text": content,
		}
		// 将文本内容添加到imageContents开头
		imageContents = append(imageContents, textContent)
	}

	// 添加当前用户消息
	currentMsg := mp_common.OpenAIReqMsg{
		Role: mp_common.MsgRole(req.Role),
		Content: func() interface{} {
			if len(imageContents) > 0 {
				return imageContents
			}
			return content
		}(),
	}
	messages = append(messages, currentMsg)

	// 默认使用流式输出
	stream := true

	// 构造LLM请求
	llmReq := &mp_common.LLMReq{
		Model:    modelInfo.Model,
		Messages: messages,
		Stream:   &stream,
	}

	// 根据开关字段设置模型参数
	if req.TemperatureEnable {
		temp := float64(req.Temperature)
		llmReq.Temperature = &temp
	}
	if req.TopPEnable {
		topP := float64(req.TopP)
		llmReq.TopP = &topP
	}
	if req.PresencePenaltyEnable {
		presencePenalty := float64(req.PresencePenalty)
		llmReq.PresencePenalty = &presencePenalty
	}
	if req.FrequencyPenaltyEnable {
		frequencyPenalty := float64(req.FrequencyPenalty)
		llmReq.FrequencyPenalty = &frequencyPenalty
	}
	if req.MaxTokensEnable {
		maxTokens := int(req.MaxTokens)
		llmReq.MaxTokens = &maxTokens
	}

	modelExperienceDialogRecordReq := &model_service.ModelExperienceDialogRecordReq{
		ModelExperienceId: req.ModelExperienceId,
		SessionId:         req.SessionId,
		ModelId:           modelInfo.ModelId,
		OriginalContent:   req.Content,
		HandledContent:    content,
		Role:              req.Role,
		ParentID:          req.ParentId,
		FileIdList:        strings.Join(req.FileIdList, ","),
		UserId:            userId,
		OrgId:             orgId,
	}
	modelExperienceDialogRecord := SaveModelExperienceDialogRecord(ctx, modelExperienceDialogRecordReq, userId, orgId)
	req.ParentId = modelExperienceDialogRecord.Id

	// llm config
	llm, err := mp.ToModelConfig(modelInfo.Provider, modelInfo.ModelType, modelInfo.ProviderConfig)
	if err != nil {
		errMsg := fmt.Sprintf("model %v chat completions err: %v", modelInfo.ModelId, err)
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, errMsg))
		return
	}

	iLLM, ok := llm.(mp.ILLM)
	if !ok {
		errMsg := fmt.Sprintf("model %v chat completions err: invalid provider", modelInfo.ModelId)
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, errMsg))
		return
	}

	// chat completions
	iLLMReq, err := iLLM.NewReq(llmReq)
	if err != nil {
		errMsg := fmt.Sprintf("model %v chat completions NewReq err: %v", modelInfo.ModelId, err)
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, errMsg))
		return
	}
	resp, sseCh, err := iLLM.ChatCompletions(ctx.Request.Context(), iLLMReq)
	if err != nil {
		errMsg := fmt.Sprintf("model %v chat completions err: %v", modelInfo.ModelId, err)
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, errMsg))
		return
	}
	// unary
	if !iLLMReq.Stream() {
		if data, ok := resp.ConvertResp(); ok {
			// 保存模型响应到会话历史
			if ok && data != nil {
				if len(data.Choices) > 0 && data.Choices[0].Message != nil {
					SaveModelExperienceDialogRecord(ctx, &model_service.ModelExperienceDialogRecordReq{
						ModelExperienceId: req.ModelExperienceId,
						SessionId:         req.SessionId,
						ModelId:           modelInfo.ModelId,
						OriginalContent:   data.Choices[0].Message.Content,
						ReasoningContent:  *data.Choices[0].Message.ReasoningContent,
						Role:              string(mp_common.MsgRoleAssistant),
						ParentID:          req.ParentId,
					}, userId, orgId)
				}
			}
			return
		}
		errMsg := fmt.Sprintf("model %v chat completions err: invalid resp", modelInfo.ModelId)
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, errMsg))
		return
	}
	// stream
	var answer string
	var reasonContent string
	var (
		firstFlag = false // 思维链起始标识符，默认思维链未开始
		endFlag   = false // 思维链结束标识符，默认思维链未结束
	)
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Content-Type", "text/event-stream; charset=utf-8")
	var data *mp_common.LLMResp

	// 创建一个通道用于处理goroutine退出信号
	doneChan := make(chan struct{})

	// 启动一个goroutine来监听客户端取消请求
	go func() {
		defer util.PrintPanicStack()
		defer close(doneChan)
		select {
		case <-ctx.Request.Context().Done():
			// 客户端取消了请求，保存部分回答
			log.Infof("客户端取消请求，开始保存部分回答，modelExperienceId: %d, answer length: %d", req.ModelExperienceId, len(answer))
			if answer != "" {
				// 使用一个新的上下文而不是请求上下文，确保即使请求被取消也能保存数据
				saveCtx := context.Background()
				log.Infof("尝试保存部分回答，modelExperienceId: %d, answer长度: %d", req.ModelExperienceId, len(answer))

				// 直接调用model的SaveModelExperienceDialogRecord方法，使用新的上下文
				record, err := model.SaveModelExperienceDialogRecord(saveCtx, &model_service.ModelExperienceDialogRecordReq{
					ModelExperienceId: req.ModelExperienceId,
					SessionId:         req.SessionId,
					ModelId:           modelInfo.ModelId,
					OriginalContent:   answer,
					ReasoningContent:  reasonContent,
					Role:              string(mp_common.MsgRoleAssistant),
					ParentID:          req.ParentId,
					UserId:            userId,
					OrgId:             orgId,
				})

				// 检查保存是否成功
				if err != nil {
					log.Errorf("部分回答保存失败，modelExperienceId: %d, 错误: %v", req.ModelExperienceId, err)
				} else if record != nil {
					log.Infof("部分回答保存成功，modelExperienceId: %d, recordId: %d", req.ModelExperienceId, record.Id)
				} else {
					log.Errorf("部分回答保存失败，modelExperienceId: %d, 未返回记录", req.ModelExperienceId)
				}
			}
		case <-doneChan:
			// 主goroutine完成，退出监听
			return
		}
	}()

	for sseResp := range sseCh {
		data, ok = sseResp.ConvertResp()
		dataStr := ""
		if ok && data != nil && len(data.Choices) > 0 {
			delta := data.Choices[0].Delta
			if delta != nil {
				answer = answer + delta.Content

				if delta.ReasoningContent != nil {
					reasonContent = reasonContent + *delta.ReasoningContent
				}
			}
			if firstFlag && !endFlag && delta.ReasoningContent != nil {
				delta.Content = delta.Content + *delta.ReasoningContent
			}
			if !endFlag && delta.Content != "" && ((delta.ReasoningContent != nil &&
				*delta.ReasoningContent == "") || delta.ReasoningContent == nil) && firstFlag {
				delta.Content = "\n</think>\n" + delta.Content
				endFlag = true
			}
			if !firstFlag && delta.ReasoningContent != nil && *delta.ReasoningContent != "" && delta.Content == "" {
				delta.Content = "<think>\n" +
					delta.Content + *delta.ReasoningContent
				firstFlag = true
			}
			dataByte, _ := json.Marshal(data)
			dataStr = fmt.Sprintf("data: %v\n", string(dataByte))
		} else {
			dataStr = fmt.Sprintf("%v\n", sseResp.String())
		}
		if _, err = ctx.Writer.Write([]byte(dataStr)); err != nil {
			log.Errorf("model %v chat completions sse err: %v", modelInfo.ModelId, err)
			gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v chat completions sse err: %v", modelInfo.ModelId, err)))
			return
		}
		ctx.Writer.Flush()
	}

	// 保存模型响应到会话历史
	SaveModelExperienceDialogRecord(ctx, &model_service.ModelExperienceDialogRecordReq{
		ModelExperienceId: req.ModelExperienceId,
		SessionId:         req.SessionId,
		ModelId:           modelInfo.ModelId,
		OriginalContent:   answer,
		ReasoningContent:  reasonContent,
		Role:              string(mp_common.MsgRoleAssistant),
		ParentID:          req.ParentId,
	}, userId, orgId)

	ctx.Set(gin_util.STATUS, http.StatusOK)
	ctx.Set(gin_util.RESULT, answer)

	log.Infof("LLM model experience request completed: %s", req.ModelId)
	return
}

// SaveModelExperienceDialog 保存模型体验对话
func SaveModelExperienceDialog(ctx *gin.Context, req *request.ModelExperienceDialogRequest, userId, orgId string) (*response.ModelExperienceDialog, error) {
	log.Infof("开始保存模型体验对话: %+v", req)

	// 将interface{}类型的ModelSetting转换为 json string
	var modelSettingStr string
	if req.ModelSetting != nil {
		jsonBytes, err := json.Marshal(req.ModelSetting)
		if err != nil {
			return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("Model settings serialization error: err: %v", err))
		}
		modelSettingStr = string(jsonBytes)
	}

	res, err := model.SaveModelExperienceDialog(ctx.Request.Context(), &model_service.ModelExperienceDialogReq{
		ModelId:      req.ModelId,
		SessionId:    req.SessionId,
		ModelSetting: modelSettingStr,
		Title:        req.Title,
		UserId:       userId,
		OrgId:        orgId,
	})

	if err != nil {
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("Failed to save model experience dialog record: err: %v", err))
	}
	log.Infof("完成模型体验对话记录保存: ")
	return &response.ModelExperienceDialog{
		ID:           res.Id,
		ModelId:      res.ModelId,
		SessionId:    res.SessionId,
		ModelSetting: res.ModelSetting,
		Title:        res.Title,
	}, nil
}

// SaveModelExperienceDialogRecord 保存模型体验对话记录
func SaveModelExperienceDialogRecord(ctx *gin.Context, req *model_service.ModelExperienceDialogRecordReq, userId, orgId string) *model_service.ModelExperienceDialogRecord {
	log.Infof("开始保存模型体验对话记录: %+v", req)

	// 设置用户ID和组织ID
	req.UserId = userId
	req.OrgId = orgId

	modelExperienceDialogRecord, err := model.SaveModelExperienceDialogRecord(ctx.Request.Context(), req)

	if err != nil {
		gin_util.Response(ctx, nil, err)
		return nil
	}
	log.Infof("完成模型体验对话记录保存: %s", req)
	return modelExperienceDialogRecord
}

// getHistoricalDialogRecords 从数据库获取历史对话记录
func getHistoricalDialogRecords(ctx context.Context, modelExperienceId, sessionId, modelId, userId, orgId string) ([]mp_common.OpenAIReqMsg, error) {
	log.Infof("开始从数据库获取历史对话记录，modelExperienceId: %d, sessionId: %s, model: %s, userId: %s, orgId: %s", modelExperienceId, sessionId, modelId, userId, orgId)

	// 调用model-service获取对话记录
	recordsResp, err := model.GetModelExperienceDialogRecords(ctx, &model_service.GetModelExperienceDialogRecordReq{
		ModelExperienceId: modelExperienceId,
		SessionId:         sessionId,
		ModelId:           modelId,
		UserId:            userId,
		OrgId:             orgId,
	})
	if err != nil {
		log.Errorf("获取模型体验对话记录失败: %v", err)
		return nil, fmt.Errorf("获取模型体验对话记录失败: %v", err)
	}

	// 将数据库记录转换为OpenAIReqMsg格式
	var messages []mp_common.OpenAIReqMsg
	for _, record := range recordsResp.Record {
		content := record.HandledContent
		if content == "" {
			content = record.OriginalContent
		}
		messages = append(messages, mp_common.OpenAIReqMsg{
			Role:    mp_common.MsgRole(record.Role),
			Content: content,
		})
	}

	log.Infof("成功获取历史对话记录，共%d条", len(messages))
	return messages, nil
}

// 辅助函数：将map转换为JSON字符串
func convertMapToJson(data map[string]interface{}) string {
	jsonBytes, _ := json.Marshal(data)
	return string(jsonBytes)
}

// ModelExperienceFileExtract 文本提取
func ModelExperienceFileExtract(ctx *gin.Context, req *request.FileExtractRequest, userId, orgId string) (*response.ModelExperienceFile, error) {
	log.Infof("开始提取文件文本: %+v", req)

	// 构造请求数据
	requestData := map[string]interface{}{
		"upload_file_url": req.FilePath,
	}

	// 发送POST请求到文档解析服务
	resp, err := http.Post("http://agent-wanwu:15003/doc_pra", "application/json",
		strings.NewReader(convertMapToJson(requestData)))
	if err != nil {
		log.Errorf("调用文档解析服务失败: %v", err)
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("Document parsing service call failed: err: %v", err))
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("读取文档解析服务响应失败: %v", err)
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, "Failed to parse document parsing service response")
	}

	// 解析响应
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Errorf("解析文档解析服务响应失败: %v", err)
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, "Failed to parse document parsing response")
	}

	// 保存解析结果到数据库
	fileReq := &model_service.ModelExperienceFileReq{
		FileName:    req.FileName,
		FilePath:    req.FilePath,
		FileExt:     util.FileExt(req.FilePath),
		ExtractText: string(body),
		FileSize:    req.FileSize,
		UserId:      userId,
		OrgId:       orgId,
	}

	modelExperienceFile, err := model.SaveModelExperienceFile(ctx.Request.Context(), fileReq)
	if err != nil {
		log.Errorf("保存文件解析结果到数据库失败: %v", err)
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, "Failed to save file parsing result")
	}

	log.Infof("完成文件文本提取: %s", req.FilePath)
	res := &response.ModelExperienceFile{
		ID:       modelExperienceFile.Id,
		FileName: modelExperienceFile.FileName,
		FilePath: modelExperienceFile.FilePath,
		FileExt:  modelExperienceFile.FileExt,
		FileSize: modelExperienceFile.FileSize,
	}

	return res, nil
}

// GetModelExperienceDialogs 获取模型体验对话列表
func GetModelExperienceDialogs(ctx *gin.Context, userId, orgId string) (*response.ListResult, error) {
	modelExperienceDialogs, err := model.GetModelExperienceDialogs(ctx, &model_service.GetModelExperienceDialogReq{
		UserId: userId,
		OrgId:  orgId,
	})
	if err != nil {
		log.Errorf("获取模型体验对话列表失败: %v", err)
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, "Failed to get model experience dialog list")
	}

	return &response.ListResult{
		List:  modelExperienceDialogs.ModelExperienceDialog,
		Total: modelExperienceDialogs.Total,
	}, nil
}

// DeleteDialog 删除模型体验对话
func DeleteDialog(ctx *gin.Context, req *request.ModelExperienceDialogIDReq, userId, orgId string) error {
	_, err := model.DeleteModelExperienceDialog(ctx, &model_service.ModelExperienceDialogIdReq{
		ModelExperienceId: req.ModelExperienceId,
		UserId:            userId,
		OrgId:             orgId,
	})
	if err != nil {
		log.Errorf("删除模型体验对话失败: %v", err)
		return grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("Failed to delete model experience dialog: %v", err))
	}

	return nil
}

// GetModelExperienceDialogRecords 获取模型体验对话记录
func GetModelExperienceDialogRecords(ctx *gin.Context, req *request.ModelExperienceDialogRecordRequest, userId, orgId string) (*response.ListResult, error) {
	modelExperienceDialog, err := model.GetModelExperienceDialog(ctx, &model_service.ModelExperienceDialogIdReq{
		ModelExperienceId: req.ModelExperienceId,
		UserId:            userId,
		OrgId:             orgId,
	})
	if err != nil {
		log.Errorf("获取模型体验对话失败: %v", err)
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("Failed to get model experience dialog: %v", err))
	}
	modelExperienceDialogRecords, err := model.GetModelExperienceDialogRecords(ctx, &model_service.GetModelExperienceDialogRecordReq{
		ModelExperienceId: req.ModelExperienceId,
		ModelId:           modelExperienceDialog.ModelId,
		UserId:            userId,
		OrgId:             orgId,
	})
	if err != nil {
		log.Errorf("获取模型体验对话记录失败: %v", err)
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("Failed to get model experience dialog records: %v", err))
	}
	var list []*response.ModelExperienceDialogRecord
	for _, record := range modelExperienceDialogRecords.Record {

		var respFileList []response.ModelExperienceFile
		if record.FileIdList != "" {
			// 将字符串文件ID列表转换为uint32切片
			fileIdStrs := strings.Split(record.FileIdList, ",")
			var fileIds []string
			for _, idStr := range fileIdStrs {
				if idStr != "" {
					fileIds = append(fileIds, idStr)
				}
			}
			fileResp, err := model.GetModelExperienceFilesByIds(ctx, &model_service.GetModelExperienceFilesByIdsReq{
				FileIds: fileIds,
				UserId:  userId,
				OrgId:   orgId,
			})
			if err != nil {
				log.Errorf("获取模型体验文件列表失败: %v", err)
				return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("Failed to get model experience file list: %v", err))
			}
			// 将proto文件列表转换为response文件列表

			for _, file := range fileResp.Files {
				respFileList = append(respFileList, response.ModelExperienceFile{
					FileName: file.FileName,
					FilePath: file.FilePath,
					FileExt:  file.FileExt,
					FileSize: file.FileSize,
					ID:       file.Id,
				})
			}
		}

		list = append(list, &response.ModelExperienceDialogRecord{
			Id:                record.Id,
			ModelId:           record.ModelId,
			ModelExperienceId: record.ModelExperienceId,
			OriginalContent:   record.OriginalContent,
			ReasoningContent:  record.ReasoningContent,
			Role:              record.Role,
			ParentID:          record.ParentID,
			FileList:          respFileList,
		})
	}

	return &response.ListResult{
		List:  list,
		Total: modelExperienceDialogRecords.Total,
	}, nil
}
