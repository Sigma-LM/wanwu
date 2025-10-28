package mcp_util

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
)

var msMgr *mcpServerMgr

type mcpServer struct {
	sseServer           *server.Server
	sseHandler          *transport.SSEHandler
	sseTransport        transport.ServerTransport
	streamableServer    *server.Server
	streamableHandler   *transport.StreamableHTTPHandler
	streamableTransport transport.ServerTransport
}
type mcpServerMgr struct {
	mcpServers map[string]*mcpServer
	mu         sync.RWMutex
}

func Init(ctx context.Context) error {
	if msMgr != nil {
		return fmt.Errorf("mcp server already init")
	}
	msMgr = &mcpServerMgr{
		mcpServers: make(map[string]*mcpServer),
	}
	return nil
}

func GetMCPServerSSEHandler(mcpServerId string) *transport.SSEHandler {
	msMgr.mu.RLock()
	defer msMgr.mu.RUnlock()
	return msMgr.mcpServers[mcpServerId].sseHandler
}

func GetMCPServerStreamableHandler(mcpServerId string) *transport.StreamableHTTPHandler {
	msMgr.mu.RLock()
	defer msMgr.mu.RUnlock()
	return msMgr.mcpServers[mcpServerId].streamableHandler
}

func CheckMCPServerExist(mcpServerId string) bool {
	msMgr.mu.RLock()
	defer msMgr.mu.RUnlock()
	_, exist := msMgr.mcpServers[mcpServerId]
	return exist
}

func StartMCPServer(ctx context.Context, mcpServerId string) error {
	if msMgr == nil {
		return fmt.Errorf("mcp server manager is nil")
	}
	if CheckMCPServerExist(mcpServerId) {
		return fmt.Errorf("mcp server already exist")
	}
	messageUrl, err := url.JoinPath(config.Cfg().Server.ApiBaseUrl, "/openapi/v1/mcp/server/message")
	if err != nil {
		return fmt.Errorf("join message url error: %v", err)
	}
	sseTransport, sseHandler, err := transport.NewSSEServerTransportAndHandler(messageUrl,
		transport.WithSSEServerTransportAndHandlerOptionCopyParamKeys([]string{"key"}))
	if err != nil {
		return fmt.Errorf("new sse transport and hander with error: %v", err)
	}
	sseSrv, err := server.NewServer(sseTransport)
	if err != nil {
		return fmt.Errorf("new server with error: %v", err)
	}
	streamTransport, streamHandler, err := transport.NewStreamableHTTPServerTransportAndHandler(
		transport.WithStreamableHTTPServerTransportAndHandlerOptionStateMode(transport.Stateful))
	if err != nil {
		return fmt.Errorf("new sse transport and hander with error: %v", err)
	}
	streamSrv, err := server.NewServer(streamTransport)
	if err != nil {
		return fmt.Errorf("new server with error: %v", err)
	}
	msMgr.mu.Lock()
	defer msMgr.mu.Unlock()
	msMgr.mcpServers[mcpServerId] = &mcpServer{
		sseServer:           sseSrv,
		sseHandler:          sseHandler,
		sseTransport:        sseTransport,
		streamableServer:    streamSrv,
		streamableHandler:   streamHandler,
		streamableTransport: streamTransport,
	}
	return nil
}

func RegisterMCPServerTools(mcpServerId string, tools []*McpTool) error {
	if !CheckMCPServerExist(mcpServerId) {
		return fmt.Errorf("mcp server doesn't exist")
	}
	msMgr.mu.Lock()
	defer msMgr.mu.Unlock()
	for _, mcpTool := range tools {
		msMgr.mcpServers[mcpServerId].sseServer.RegisterTool(mcpTool.Tool, mcpTool.Handle)
		msMgr.mcpServers[mcpServerId].streamableServer.RegisterTool(mcpTool.Tool, mcpTool.Handle)
	}
	return nil
}

func UnRegisterMCPServerTools(mcpServerId string, tools []string) error {
	if !CheckMCPServerExist(mcpServerId) {
		return fmt.Errorf("mcp server doesn't exist")
	}
	msMgr.mu.Lock()
	defer msMgr.mu.Unlock()
	for _, tool := range tools {
		msMgr.mcpServers[mcpServerId].sseServer.UnregisterTool(tool)
		msMgr.mcpServers[mcpServerId].streamableServer.UnregisterTool(tool)
	}
	return nil
}

func ShutDownMCPServer(ctx context.Context, mcpServerId string) error {
	if !CheckMCPServerExist(mcpServerId) {
		return fmt.Errorf("mcp server doesn't exist")
	}
	msMgr.mu.Lock()
	defer msMgr.mu.Unlock()
	err := msMgr.mcpServers[mcpServerId].sseServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	err = msMgr.mcpServers[mcpServerId].streamableServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	delete(msMgr.mcpServers, mcpServerId)
	return nil
}
