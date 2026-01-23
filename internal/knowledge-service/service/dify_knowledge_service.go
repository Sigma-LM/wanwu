package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/config"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/http"
	http_client "github.com/UnicomAI/wanwu/pkg/http-client"
	"github.com/UnicomAI/wanwu/pkg/log"
)

type DifyGetDatasetsParams struct {
	Keyword    string   `json:"keyword"`    // 关键词
	TagIds     []string `json:"tagIds"`     // 标签
	Page       int      `json:"page"`       // 分页序号
	Limit      int      `json:"limit"`      // 分页大小
	IncludeAll bool     `json:"includeAll"` // 是否包含所有数据集
}

type DifyGetDatasetsResp struct {
	Data    []*DifyDatasetData `json:"data"`
	HasMore bool               `json:"has_more"`
	Page    int                `json:"page"`
	Limit   int                `json:"limit"`
	Total   int                `json:"total"`
}

type DifyDatasetData struct {
	Id                      string                      `json:"id"`
	Name                    string                      `json:"name"`
	Description             string                      `json:"description"`
	Provider                string                      `json:"provider"`
	Permission              string                      `json:"permission"`
	DataSourceType          string                      `json:"data_source_type"`
	IndexingTechnique       string                      `json:"indexing_technique"`
	AppCount                int64                       `json:"app_count"`
	DocumentCount           int64                       `json:"document_count"`
	WordCount               int64                       `json:"word_count"`
	CreatedBy               string                      `json:"created_by"`
	AuthorName              string                      `json:"author_name"`
	CreatedAt               int64                       `json:"created_at"`
	UpdatedBy               string                      `json:"updated_by"`
	UpdatedAt               int64                       `json:"updated_at"`
	EmbeddingModel          string                      `json:"embedding_model"`
	EmbeddingModelProvider  string                      `json:"embedding_model_provider"`
	EmbeddingAvailable      bool                        `json:"embedding_available"`
	RetrievalModelDict      *DifyRetrievalModelDict     `json:"retrieval_model_dict"`
	Tags                    []*DifyTag                  `json:"tags"`
	DocForm                 string                      `json:"doc_form"`
	ExternalKnowledgeInfo   *DifyExternalKnowledgeInfo  `json:"external_knowledge_info"`
	ExternalRetrievalModel  *DifyExternalRetrievalModel `json:"external_retrieval_model"`
	DocMetaData             []*DifyDocMetaData          `json:"doc_meta_data"`
	BuiltInFieldEnabled     bool                        `json:"built_in_field_enabled"`
	PipelineId              string                      `json:"pipeline_id"`
	RuntimeMode             string                      `json:"runtime_mode"`
	ChunkStructure          string                      `json:"chunk_structure"`
	IconInfo                *DifyIconInfo               `json:"icon_info"`
	IsPublished             bool                        `json:"is_published"`
	TotalDocuments          int64                       `json:"total_documents"`
	TotalAvailableDocuments int64                       `json:"total_available_documents"`
	EnableAPI               bool                        `json:"enable_api"`
	IsMultimodal            bool                        `json:"is_multimodal"`
}

type DifyRetrievalModelDict struct {
	SearchMethod          string              `json:"search_method"`
	RerankingEnable       bool                `json:"reranking_enable"`
	RerankingMode         string              `json:"reranking_mode"`
	RerankingModel        *DifyRerankingModel `json:"reranking_model"`
	TopK                  int64               `json:"top_k"`
	ScoreThresholdEnabled bool                `json:"score_threshold_enabled"`
	ScoreThreshold        float32             `json:"score_threshold"`
	Weights               *DifyWeights        `json:"weights"`
}

type DifyWeights struct {
	WeightType     string              `json:"weight_type"`
	KeywordSetting *DifyKeywordSetting `json:"keyword_setting"`
	VectorSetting  *DifyVectorSetting  `json:"vector_setting"`
}

type DifyKeywordSetting struct {
	KeywordWeight float32 `json:"keyword_weight"`
}

type DifyVectorSetting struct {
	VectorWeight          float32 `json:"vector_weight"`
	EmbeddingModelName    string  `json:"embedding_model_name"`
	EmbeddingProviderName string  `json:"embedding_provider_name"`
}

type DifyRerankingModel struct {
	RerankingProviderName string `json:"reranking_provider_name"`
	RerankingModelName    string `json:"reranking_model_name"`
}

type DifyTag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type DifyExternalKnowledgeInfo struct {
	ExternalKnowledgeId          string `json:"external_knowledge_id"`
	ExternalKnowledgeAPIId       string `json:"external_knowledge_api_id"`
	ExternalKnowledgeAPIName     string `json:"external_knowledge_api_name"`
	ExternalKnowledgeAPIEndpoint string `json:"external_knowledge_api_endpoint"`
}

type DifyExternalRetrievalModel struct {
	TopK                  int32   `json:"top_k"`
	ScoreThresholdEnabled bool    `json:"score_threshold_enabled"`
	ScoreThreshold        float32 `json:"score_threshold"`
}

type DifyDocMetaData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type DifyIconInfo struct {
	IconType       string `json:"icon_type"`
	Icon           string `json:"icon"`
	IconBackground string `json:"icon_background"`
	IconUrl        string `json:"icon_url"`
}

// DifyGetDatasets dify获取知识库列表
func DifyGetDatasets(ctx context.Context, externalAPI *model.KnowledgeExternalAPI, difyGetDatasetsParams *DifyGetDatasetsParams) (*DifyGetDatasetsResp, error) {
	difyServer := config.GetConfig().DifyServer
	difyUrl := externalAPI.BaseUrl + difyServer.GetDatasetsUri
	params := map[string]string{}
	params["keyword"] = difyGetDatasetsParams.Keyword
	params["page"] = strconv.Itoa(difyGetDatasetsParams.Page)
	params["limit"] = strconv.Itoa(difyGetDatasetsParams.Limit)
	params["include_all"] = strconv.FormatBool(difyGetDatasetsParams.IncludeAll)
	headers := map[string]string{}
	headers["Authorization"] = fmt.Sprintf("Bearer %s", externalAPI.APIKey)
	result, err := http.GetClient().Get(ctx, &http_client.HttpRequestParams{
		Headers:    headers,
		Params:     params,
		Url:        difyUrl,
		Timeout:    time.Duration(difyServer.Timeout) * time.Second,
		MonitorKey: "dify_get_datasets",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return nil, err
	}
	var resp DifyGetDatasetsResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf("knowledge get datasets unmarshal err: %v", err.Error())
		return nil, err
	}
	return &resp, nil
}

// DifyGetDataset dify获取知识库详情
func DifyGetDataset(ctx context.Context, externalAPI *model.KnowledgeExternalAPI, externalKnowledgeId string) (*DifyDatasetData, error) {
	difyServer := config.GetConfig().DifyServer
	difyUrl := externalAPI.BaseUrl + strings.ReplaceAll(difyServer.GetDatasetUri, "{dataset_id}", externalKnowledgeId)
	headers := map[string]string{}
	headers["Authorization"] = fmt.Sprintf("Bearer %s", externalAPI.APIKey)
	result, err := http.GetClient().Get(ctx, &http_client.HttpRequestParams{
		Headers:    headers,
		Url:        difyUrl,
		Timeout:    time.Duration(difyServer.Timeout) * time.Second,
		MonitorKey: "dify_get_dataset",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return nil, err
	}
	var resp DifyDatasetData
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf("knowledge get dataset unmarshal err: %v", err.Error())
		return nil, err
	}
	return &resp, nil
}
