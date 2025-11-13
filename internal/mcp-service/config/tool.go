package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	openapi3_util "github.com/UnicomAI/wanwu/pkg/openapi3-util"
)

type ToolConfig struct {
	ToolSquareId       string          `json:"tool_square_id" mapstructure:"tool_square_id"`
	Name               string          `json:"name" mapstructure:"name"`
	Desc               string          `json:"desc" mapstructure:"desc"`
	AvatarPath         string          `json:"avatar_path" mapstructure:"avatar_path"`
	Detail             string          `json:"detail" mapstructure:"detail"`
	Tags               string          `json:"tags" mapstructure:"tags"`
	Tools              []McpToolConfig `json:"tools" mapstructure:"tools"`
	AuthType           string          `json:"auth_type" mapstructure:"auth_type"`
	ApiKeyHeaderPrefix string          `json:"api_key_header_prefix" mapstructure:"api_key_header_prefix"`
	ApiKeyHeader       string          `json:"api_key_header" mapstructure:"api_key_header"`
	ApiKeyQueryParam   string          `json:"api_key_query_param" mapstructure:"api_key_query_param"`
	ApiKeyValue        string          `json:"api_key_value" mapstructure:"api_key_value"`
	Schema             string          `json:"schema" mapstructure:"-"`
	SchemaPath         string          `json:"schema_path" mapstructure:"schema_path"`
	NeedApiKeyInput    bool            `json:"need_api_key_input" mapstructure:"need_api_key_input"`
}

func (tool *ToolConfig) load() error {
	avatarPath := filepath.Join(ConfigDir, tool.AvatarPath)
	if _, err := os.ReadFile(avatarPath); err != nil {
		return fmt.Errorf("load tool %v avatar path %v err: %v", tool.ToolSquareId, avatarPath, err)
	}
	schemaPath := filepath.Join(ConfigDir, tool.SchemaPath)
	schemaOpenAPI, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("load tool %v schema path %v err: %v", tool.ToolSquareId, schemaPath, err)
	}
	if err := openapi3_util.ValidateSchema(context.Background(), schemaOpenAPI); err != nil {
		return grpc_util.ErrorStatus(errs.Code_MCPGeneral, err.Error())
	}
	tool.Schema = string(schemaOpenAPI)
	return nil
}
