package mp_jina

import (
	"context"
	"net/url"

	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
)

type MultiModalRerank struct {
	ApiKey      string `json:"apiKey"`      // ApiKey
	EndpointUrl string `json:"endpointUrl"` // 推理url
	ContextSize *int   `json:"contextSize"` // 上下文长度
}

func (cfg *MultiModalRerank) Tags() []mp_common.Tag {
	tags := []mp_common.Tag{
		{
			Text: mp_common.TagMultiModalRerank,
		},
	}
	tags = append(tags, mp_common.GetTagsByContentSize(cfg.ContextSize)...)
	return tags
}

func (cfg *MultiModalRerank) NewReq(req *mp_common.MultiModalRerankReq) (mp_common.IMultiModalRerankReq, error) {
	m, err := req.Data()
	if err != nil {
		return nil, err
	}
	return mp_common.NewRerankReq(m), nil
}

func (cfg *MultiModalRerank) MultiModalRerank(ctx context.Context, req mp_common.IMultiModalRerankReq, headers ...mp_common.Header) (mp_common.IMultiModalRerankResp, error) {
	b, err := mp_common.MultiModalRerank(ctx, "jina", cfg.ApiKey, cfg.rerankUrl(), req.Data(), headers...)
	if err != nil {
		return nil, err
	}
	return mp_common.NewMultiModalRerankResp(string(b)), nil
}

func (cfg *MultiModalRerank) rerankUrl() string {
	ret, _ := url.JoinPath(cfg.EndpointUrl, "/rerank")
	return ret
}
