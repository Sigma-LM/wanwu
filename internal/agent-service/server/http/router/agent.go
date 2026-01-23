package handler

import (
	"net/http"

	http_server "github.com/UnicomAI/wanwu/internal/agent-service/pkg/http-server"
	"github.com/UnicomAI/wanwu/internal/agent-service/server/http/handler"
)

func init() {
	group := http_server.Group("/agent")
	group.Register("/chat", http.MethodPost, handler.AgentChat, "智能体流式问答")

	multiGroup := http_server.Group("/multi_agent")
	multiGroup.Register("/chat", http.MethodPost, handler.MultiAgentChat, "多智能体流式问答")
}
