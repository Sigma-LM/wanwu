package mp_yuanjing

import (
	"context"
	"encoding/json"
	"github.com/UnicomAI/wanwu/pkg/util"
	"net/url"

	"github.com/UnicomAI/wanwu/pkg/log"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
)

type Rerank struct {
	ApiKey      string `json:"apiKey"`      // ApiKey
	EndpointUrl string `json:"endpointUrl"` // 推理url
	ContextSize *int   `json:"contextSize"` // 上下文长度
}

func (cfg *Rerank) Tags() []mp_common.Tag {
	tags := []mp_common.Tag{
		{
			Text: mp_common.TagRerank,
		},
	}
	tags = append(tags, mp_common.GetTagsByContentSize(cfg.ContextSize)...)
	return tags
}

func (cfg *Rerank) NewReq(req *mp_common.RerankReq) (mp_common.IRerankReq, error) {
	instruction := "Given a web search query, retrieve relevant passages that answer the query"
	if req.Instruction == nil {
		req.Instruction = &instruction
	}
	m := map[string]interface{}{
		"instruction": req.Instruction,
		"documents":   req.Documents,
		"query":       req.Query,
	}
	return mp_common.NewRerankReq(m), nil
}

func (cfg *Rerank) Rerank(ctx context.Context, req mp_common.IRerankReq, headers ...mp_common.Header) (mp_common.IRerankResp, error) {
	b, err := mp_common.Rerank(ctx, "yuanjing", cfg.ApiKey, cfg.rerankUrl(), req.Data(), headers...)
	if err != nil {
		return nil, err
	}
	return &rerankResp{raw: string(b)}, nil
}

func (cfg *Rerank) rerankUrl() string {
	ret, _ := url.JoinPath(cfg.EndpointUrl, "/rerank")
	return ret
}

// --- rerankResp ---

type rerankResp struct {
	raw     string
	Results []mp_common.Result `json:"results"`
}

func (resp *rerankResp) String() string {
	return resp.raw
}

func (resp *rerankResp) Data() (interface{}, bool) {
	ret := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(resp.raw), &ret); err != nil {
		log.Errorf("yuanjing rerank resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}
	return ret, true
}

func (resp *rerankResp) ConvertResp() (*mp_common.RerankResp, bool) {
	if err := json.Unmarshal([]byte(resp.raw), resp); err != nil {
		log.Errorf("yuanjing rerank resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}
	if err := util.Validate(resp); err != nil {
		log.Errorf("yuanjing rerank resp validate err: %v", err)
		return nil, false
	}
	return &mp_common.RerankResp{
		Results: resp.Results,
	}, true
}
