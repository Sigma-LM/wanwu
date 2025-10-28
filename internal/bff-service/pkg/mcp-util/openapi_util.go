package mcp_util

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/getkin/kin-openapi/openapi3"
)

const timeout = 60 * time.Second

var urlRegex = regexp.MustCompile(`(https|http)://.*$`)

type APIAuth struct {
	Type  string `json:"type"`
	In    string `json:"in"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type McpTool struct {
	ToolSchema string
	Tool       *protocol.Tool
	Handle     server.ToolHandlerFunc
}

type OpenAPISchema struct {
	Schema      string
	ApiAuth     *APIAuth
	MethodNames []string
}

func OpenApiToMcpToolList(content *OpenAPISchema) ([]*McpTool, error) {
	//将openapi数据转换成openapi结构体
	parseContent, err := ParseOpenApiContent(content.Schema)
	if err != nil {
		return nil, err
	}
	var mcpTools []*McpTool
	for _, methodName := range content.MethodNames {
		mcpTools = append(mcpTools, &McpTool{
			ToolSchema: BuildOpenApiContent(parseContent, methodName),
			Tool:       ConvertMcpTool(parseContent, methodName),
			Handle:     ConvertMcpHandler(parseContent, methodName, content.ApiAuth),
		})
	}
	return mcpTools, nil
}

func ParseOpenApiContent(content string) (*openapi3.T, error) {
	// 创建新的 OpenAPI 加载器
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	doc, err := loader.LoadFromData([]byte(content))
	if err != nil {
		return nil, err
	}
	// 验证文档
	if err = doc.Validate(loader.Context); err != nil {
		return nil, err
	}
	// 校验字段
	if !checkFields(doc) {
		return nil, errors.New("check field error")
	}
	return doc, nil
}

// 生成只包含单个方法的openapi schema
func BuildOpenApiContent(doc *openapi3.T, operationID string) string {
	operation, path, method := ParseOpenApiOperation(doc, operationID)
	pathItem := &openapi3.PathItem{}
	pathItem.SetOperation(method, operation)
	doc.Paths = openapi3.NewPaths(openapi3.WithPath(path, pathItem))
	content, _ := doc.MarshalJSON()
	return string(content)
}

func ParseOpenApiOperation(doc *openapi3.T, operationID string) (operation *openapi3.Operation, path, method string) {
	for pathName, pathItem := range doc.Paths.Map() {
		for methodName, op := range pathItem.Operations() {
			if op.OperationID == operationID {
				operation = op
				method = methodName
				path = pathName
				break
			}
		}
	}
	return
}

func checkFields(doc *openapi3.T) bool {
	if doc.Servers != nil {
		for _, server := range doc.Servers {
			if !urlRegex.MatchString(server.URL) {
				return false
			}
		}
	}
	return true
}

// 转换到 mcpTool
func ConvertMcpTool(doc *openapi3.T, operationID string) *protocol.Tool {
	operation, _, _ := ParseOpenApiOperation(doc, operationID)
	inputSchema := BuildInputSchema(operation.Parameters, operation.RequestBody)
	inputSchemaJSON, _ := json.MarshalIndent(inputSchema, "", "  ")
	tool := protocol.NewToolWithRawSchema(operationID, operation.Description, inputSchemaJSON)
	return tool
}

func ConvertMcpHandler(doc *openapi3.T, operationID string, apiAuth *APIAuth) server.ToolHandlerFunc {
	return func(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		operation, path, method := ParseOpenApiOperation(doc, operationID)
		baseURL := doc.Servers[rand.Intn(len(doc.Servers))].URL
		finalURL, err := url.JoinPath(baseURL, path)
		if err != nil {
			return nil, err
		}

		params := req.Arguments

		pathParams := make(map[string]interface{})
		queryParams := make(map[string]interface{})
		bodyParams := make(map[string]interface{})
		headers := make(map[string]string)

		for _, param := range operation.Parameters {
			if param.Value.In == "path" {
				pathParams[param.Value.Name] = params[param.Value.Name]
			} else if param.Value.In == "query" {
				queryParams[param.Value.Name] = params[param.Value.Name]
			} else if param.Value.In == "header" {
				headers[param.Value.Name] = params[param.Value.Name].(string)
			}
		}
		if requestBodyMap, ok := params["requestBody"].(map[string]interface{}); ok {
			bodyParams = requestBodyMap
		}

		// path param
		if len(pathParams) > 0 {
			for paramName, paramValue := range pathParams {
				placeholder := fmt.Sprintf("{%s}", paramName)
				if strings.Contains(finalURL, placeholder) {
					var strValue string
					switch v := paramValue.(type) {
					case string:
						strValue = v
					case nil:
						strValue = ""
					default:
						strValue = fmt.Sprintf("%v", v)
					}
					finalURL = strings.ReplaceAll(finalURL, placeholder, strValue)
				}
			}
		}
		// query param
		if len(queryParams) > 0 {
			parsedURL, err := url.Parse(finalURL)
			if err != nil {
				return newCallToolResultText(fmt.Sprintf("parser url: %v", err), true), err
			}
			q := parsedURL.Query()
			for paramName, paramValue := range queryParams {
				var strValue string
				switch v := paramValue.(type) {
				case string:
					strValue = v
				case nil:
					continue
				default:
					strValue = fmt.Sprintf("%v", v)
				}

				q.Add(paramName, strValue)
			}
			parsedURL.RawQuery = q.Encode()
			finalURL = parsedURL.String()
		}
		// apiAuth
		if apiAuth != nil && apiAuth.Type == "API Key" {
			if apiAuth.In == "query" {
				parsedURL, err := url.Parse(finalURL)
				if err != nil {
					return newCallToolResultText(fmt.Sprintf("parser url: %v", err), true), err
				}

				q := parsedURL.Query()
				q.Add(apiAuth.Name, apiAuth.Value)
				parsedURL.RawQuery = q.Encode()
				finalURL = parsedURL.String()
			} else if apiAuth.In == "header" {
				headers[apiAuth.Name] = apiAuth.Value
			}
		}
		// requestBody
		var reqBody io.Reader = nil
		if len(bodyParams) > 0 {
			jsonParams, err := json.Marshal(bodyParams)
			if err != nil {
				return newCallToolResultText(fmt.Sprintf("Error json marshal: %v", err), true), err
			}
			reqBody = bytes.NewBuffer(jsonParams)
		}

		request, err := http.NewRequestWithContext(ctx, method, finalURL, reqBody)
		if err != nil {
			return newCallToolResultText(fmt.Sprintf("Error creating request: %v", err), true), err
		}

		if reqBody != nil {
			request.Header.Set("Content-Type", "application/json")
		}
		// header
		for key, value := range headers {
			request.Header.Set(key, value)
		}
		client := &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		resp, err := client.Do(request)
		if err != nil {
			return newCallToolResultText(fmt.Sprintf("Error executing request: %v", err), true), err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return newCallToolResultText(fmt.Sprintf("Error reading response: %v", err), true), err
		}
		return newCallToolResultText(string(body), false), nil
	}
}

func BuildInputSchema(params openapi3.Parameters, requestBody *openapi3.RequestBodyRef) map[string]any {
	inputSchema := map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
	properties := inputSchema["properties"].(map[string]any)
	var required []string
	// param
	for _, p := range params {
		if p.Value.Schema != nil {
			prop := extractProperty(p.Value.Schema.Value)
			if p.Value.Description != "" {
				prop["description"] = p.Value.Description
			}
			properties[p.Value.Name] = prop
			if p.Value.Required {
				required = append(required, p.Value.Name)
			}
		}
	}
	// requestBody
	if requestBody != nil {
		if content := requestBody.Value.Content.Get("application/json"); content != nil {
			properties["requestBody"] = extractProperty(content.Schema.Value)
			required = append(required, "requestBody")
		}
	}

	if len(required) > 0 {
		inputSchema["required"] = required
	}
	return inputSchema
}

func extractProperty(s *openapi3.Schema) map[string]any {
	if s == nil {
		return nil
	}
	sType := getType(s)
	prop := map[string]any{}
	if sType != "" {
		prop["type"] = s.Type
	}
	if s.Format != "" {
		prop["format"] = s.Format
	}
	// Object properties
	if sType == "object" {
		objProps := map[string]any{}
		for propName, propSchema := range s.Properties {
			objProps[propName] = extractProperty(propSchema.Value)
		}
		prop["properties"] = objProps
	}
	// Array items
	if sType == "array" {
		prop["items"] = extractProperty(s.Items.Value)
	}
	return prop
}

func getType(types *openapi3.Schema) string {
	if types == nil || types.Type == nil {
		return "unknown"
	}
	slice := types.Type.Slice()
	if len(slice) == 0 {
		return "unknown"
	}
	return slice[0]
}

func newCallToolResultText(text string, isError bool) *protocol.CallToolResult {
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: text,
			},
		},
		IsError: isError,
	}
}
