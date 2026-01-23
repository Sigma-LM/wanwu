package sse_util

import (
	"bufio"
	"context"
	"net/http"

	"github.com/UnicomAI/wanwu/pkg/log"
	safe_go_util "github.com/UnicomAI/wanwu/pkg/safe-go-util"
)

type SSEReader struct {
	BusinessKey    string
	Params         string
	StreamReceiver StreamReceiver
}

type StreamReceiver interface {
	Receive() (string, error)
	Close() error
}

type HttpStreamReceiver struct {
	httpResponse *http.Response
	reader       *bufio.Reader
}

func NewHttpStreamReceiver(httpResponse *http.Response) *HttpStreamReceiver {
	return &HttpStreamReceiver{
		httpResponse: httpResponse,
		reader:       bufio.NewReader(httpResponse.Body),
	}
}

func (htr *HttpStreamReceiver) Receive() (string, error) {
	line, err := htr.reader.ReadBytes('\n')
	if err != nil { //异常結束
		return "", err
	}
	return string(line), nil
}

func (htr *HttpStreamReceiver) Close() error {
	return htr.httpResponse.Body.Close()
}

func (sr *SSEReader) ReadStream(ctx context.Context) (<-chan string, error) {
	closer := func(ctx context.Context) {
		err1 := sr.StreamReceiver.Close()
		if err1 != nil {
			log.Errorf("%s close err: %v", sr.BusinessKey, err1)
		}
	}
	rawCh := safe_go_util.SafeChannelReceive[string](ctx, func(ctx context.Context, rawCh chan string) safe_go_util.ChannelReceiveResult[string] {
		content, err := sr.StreamReceiver.Receive()
		return safe_go_util.ChannelResult[string](content, err, sr.BusinessKey, sr.Params)
	}, closer)
	return rawCh, nil
}
