package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/pkg/es"
	"github.com/UnicomAI/wanwu/pkg/log"
	sse_util "github.com/UnicomAI/wanwu/pkg/sse-util"
	"github.com/google/uuid"
)

const (
	terminationMessage = "本次回答已被终止"
	esTimeout          = 1 * time.Minute
)

type ConversationParams struct {
	AssistantId    string           `json:"assistantId"`
	ConversationId string           `json:"conversationId"`
	UserId         string           `json:"userId"`
	OrgId          string           `json:"orgId"`
	Query          string           `json:"query"`
	FileInfo       []model.FileInfo `json:"fileInfo"`
}

type AgentChatResp struct {
	Code       int           `json:"code"`
	Message    string        `json:"message"`
	Response   string        `json:"response"`
	SearchList []interface{} `json:"search_list"`
	Finish     int           `json:"finish"`
	EventType  int           `json:"eventType"`
	EventData  interface{}   `json:"eventData"`
}

type ConversationProcessor struct {
	SSEWriter *sse_util.SSEWriter[assistant_service.AssistantConversionStreamResp]
}

func (cp *ConversationProcessor) Process(ctx context.Context, req *ConversationParams, sendRequest func(ctx context.Context) (string, *http.Response, context.CancelFunc, error)) (err error) {
	var fullResponse = &strings.Builder{}
	var searchList *string
	defer func() {
		if err != nil {
			log.Errorf("[Conversation] err: %v", err)
			//错误信息通知
			_ = cp.SSEWriter.WriteLine(assistant_service.AssistantConversionStreamResp{
				Content: buildErrMsg(err),
			}, false, nil, nil)
		}
		if ctx.Err() != nil {
			err = ctx.Err()
			log.Errorf("[Conversation] context err: %v", err)
		}
		//todo delete
		log.Infof("[Conversation] fullResponse: %s", fullResponse.String())
		//保存会话
		saveConversation(ctx, req, buildResponse(fullResponse.String(), err), searchList)
	}()
	//1.执行请求
	businessKey, sseResp, cancel, err := sendRequest(ctx)
	if err != nil {
		return err
	}
	defer cancel()
	//2.读取结果
	SSEReader := &sse_util.SSEReader[string]{
		BusinessKey:    businessKey,
		StreamReceiver: sse_util.NewHttpStreamReceiver(sseResp),
	}
	stream, err := SSEReader.ReadStream(ctx)
	if err != nil {
		return err
	}
	//3.回写结果
	err = cp.SSEWriter.WriteStream(stream, nil, conversationLineBuilder(fullResponse, searchList), nil)
	return err
}

func conversationLineBuilder(fullResponse *strings.Builder, searchList *string) func(s sse_util.SSEWriterClient[assistant_service.AssistantConversionStreamResp], strLine string, streamContextParams interface{}) (assistant_service.AssistantConversionStreamResp, bool, error) {
	return func(s sse_util.SSEWriterClient[assistant_service.AssistantConversionStreamResp], strLine string, streamContextParams interface{}) (assistant_service.AssistantConversionStreamResp, bool, error) {
		conversation, searchResult := processAgentResp(strLine, streamContextParams)
		if searchList == nil && len(searchResult) > 0 {
			searchList = &searchResult
		}
		if len(conversation) > 0 {
			//保存对话
			fullResponse.WriteString(conversation)
		}
		return assistant_service.AssistantConversionStreamResp{
			Content: strLine,
		}, false, nil
	}
}

func processAgentResp(strLine string, streamContextParams interface{}) (string, string) {
	if len(strLine) >= 5 && strLine[:5] == "data:" {
		jsonStrData := strLine[5:]
		// 解析流式数据，提取response字段和search_list
		var agentChatResp = &AgentChatResp{}
		if err1 := json.Unmarshal([]byte(jsonStrData), agentChatResp); err1 == nil {
			var searchList string
			if len(agentChatResp.SearchList) > 0 {
				marshal, err := json.Marshal(agentChatResp.SearchList)
				if err == nil {
					searchList = string(marshal)
				}
			}
			if agentChatResp.EventType != 0 { //非主智能体的先不记录历史
				return "", ""
			}
			return agentChatResp.Response, searchList
		}
	}
	return "", ""
}

func buildResponse(response string, err error) string {
	var conversationResponse = response
	if err != nil {
		if len(conversationResponse) > 0 {
			conversationResponse += "\n"
		}
		conversationResponse += terminationMessage
	}
	return conversationResponse
}

// 构建错误信息,todo 后续考虑创建枚举明细错误信息
func buildErrMsg(err error) string {
	var agentChatResp = &AgentChatResp{
		Code:     1,
		Message:  "智能体处理异常，请稍后重试",
		Response: "智能体处理异常，请稍后重试",
		Finish:   1,
	}
	respString, errR := json.Marshal(agentChatResp)
	if errR != nil {
		log.Errorf("buildErrMsg error: %v", errR)
		return ""
	}
	return string(respString)
}

// 使用独立上下文保存对话的辅助函数
func saveConversation(originalCtx context.Context, req *ConversationParams, response string, searchListPtr *string) {
	if len(req.ConversationId) == 0 {
		return
	}
	var searchList string
	if searchListPtr != nil {
		searchList = *searchListPtr
	}
	// 如果原始上下文已取消，创建一个新的独立上下文
	if originalCtx.Err() != nil {
		ctx, cancel := context.WithTimeout(context.Background(), esTimeout)
		defer cancel()

		if err := saveConversationDetailToES(ctx, req, response, searchList); err != nil {
			log.Errorf("保存聊天记录到ES失败，assistantId: %s, conversationId: %s, error: %v",
				req.AssistantId, req.ConversationId, err)
		}
		return
	}

	// 原始上下文未取消时，继续使用它
	if err := saveConversationDetailToES(originalCtx, req, response, searchList); err != nil {
		log.Errorf("保存聊天记录到ES失败，assistantId: %s, conversationId: %s, error: %v",
			req.AssistantId, req.ConversationId, err)
	}
}

// saveConversationDetailToES 保存聊天记录到ES
func saveConversationDetailToES(ctx context.Context, req *ConversationParams, response, searchList string) error {
	// 根据当前时间生成索引名称，格式为conversation_detail_infos_YYYYMM
	now := time.Now()
	indexName := fmt.Sprintf("conversation_detail_infos_%d%02d", now.Year(), now.Month())

	// 组装ConversationDetails数据
	nowMilli := now.UnixMilli()
	conversationDetail := &model.ConversationDetails{
		Id:             uuid.New().String(),
		AssistantId:    req.AssistantId,
		ConversationId: req.ConversationId,
		Prompt:         req.Query,
		FileInfo:       req.FileInfo,
		Response:       response,
		SearchList:     searchList,
		UserId:         req.UserId,
		OrgId:          req.OrgId,
		CreatedAt:      nowMilli,
		UpdatedAt:      nowMilli,
	}

	// 写入ES
	if err := es.Assistant().IndexDocument(ctx, indexName, conversationDetail); err != nil {
		return fmt.Errorf("写入ES失败: %v", err)
	}

	log.Infof("成功保存聊天记录到ES，索引: %s, assistantId: %s, conversationId: %s",
		indexName, req.AssistantId, req.ConversationId)
	return nil
}
