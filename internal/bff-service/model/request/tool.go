package request

type CustomToolCreate struct {
	Name          string                      `json:"name" validate:"required"`        // 名称
	Description   string                      `json:"description" validate:"required"` // 描述
	ApiAuth       CustomToolApiAuthWebRequest `json:"apiAuth" validate:"required"`     // api身份认证
	Schema        string                      `json:"schema"  validate:"required"`     // schema
	PrivacyPolicy string                      `json:"privacyPolicy"`                   // 隐私政策
}

func (req *CustomToolCreate) Check() error { return nil }

type CustomToolApiAuthWebRequest struct {
	Type             string `json:"type" validate:"required,oneof='None' 'API Key'"` // 认证类型 None 或 APIKey
	APIKey           string `json:"apiKey"`                                          // apiKey 仅当认证类型为API Key时必填
	CustomHeaderName string `json:"customHeaderName"`                                // Custom Header Name 仅当认证类型为API Key时必填
	AuthType         string `json:"authType" validate:"omitempty,oneof=Custom"`      // Auth类型 仅当认证类型为API Key时必填，也可以为空
}

func (req *CustomToolApiAuthWebRequest) Check() error { return nil }

type CustomToolUpdateReq struct {
	CustomToolID  string                      `json:"customToolId" validate:"required"` // 自定义工具ID
	Name          string                      `json:"name" validate:"required"`         // 名称
	Description   string                      `json:"description" validate:"required"`  // 描述
	ApiAuth       CustomToolApiAuthWebRequest `json:"apiAuth" validate:"required"`      // api身份认证
	Schema        string                      `json:"schema"  validate:"required"`      // schema
	PrivacyPolicy string                      `json:"privacyPolicy"`                    // 隐私政策
}

func (req *CustomToolUpdateReq) Check() error { return nil }

type CustomToolIDReq struct {
	CustomToolID string `json:"customToolId" validate:"required"` // 自定义工具id
}

func (req *CustomToolIDReq) Check() error { return nil }

type CustomToolSchemaReq struct {
	Schema string `json:"schema" validate:"required"` // schema
}

func (req *CustomToolSchemaReq) Check() error { return nil }

type ToolSquareAPIKeyReq struct {
	ToolSquareID string `json:"toolSquareId" validate:"required"` // 广场toolId
	APIKey       string `json:"apiKey"`                           // apiKey
}

func (req *ToolSquareAPIKeyReq) Check() error { return nil }

type ToolActionListReq struct {
	ToolId   string `json:"toolId" validate:"required"`                        // 工具id
	ToolType string `json:"toolType" validate:"required,oneof=builtin custom"` // 工具类型
}

func (req *ToolActionListReq) Check() error { return nil }

type ToolActionReq struct {
	ToolId     string `json:"toolId" validate:"required"`                        // 工具id
	ToolType   string `json:"toolType" validate:"required,oneof=builtin custom"` // 工具类型
	ActionName string `json:"actionName" validate:"required"`                    // action名称
}

func (req *ToolActionReq) Check() error { return nil }

type MCPServerCreateReq struct {
	Avatar Avatar `json:"avatar"`                   // 图标
	Name   string `json:"name" validate:"required"` // 名称
	Desc   string `json:"desc" validate:"required"` // 描述
}

func (req *MCPServerCreateReq) Check() error { return nil }

type MCPServerUpdateReq struct {
	MCPServerID string `json:"mcpServerId" validate:"required"` // mcp server Id
	Avatar      Avatar `json:"avatar"`                          // 图标
	Name        string `json:"name" validate:"required"`        // 名称
	Desc        string `json:"desc" validate:"required"`        // 描述
}

func (req *MCPServerUpdateReq) Check() error { return nil }

type MCPServerIDReq struct {
	MCPServerID string `json:"mcpServerId" validate:"required"`
}

func (req *MCPServerIDReq) Check() error {
	return nil
}

type MCPServerToolCreateReq struct {
	MCPServerID string `json:"mcpServerId" validate:"required"` // mcp server Id
	Id          string `json:"id" validate:"required"`          // 应用或工具id
	Type        string `json:"type" validate:"required"`        // mcp server tool类型
	MethodName  string `json:"methodName" validate:"required"`  // 显示名称
}

func (req *MCPServerToolCreateReq) Check() error { return nil }

type MCPServerToolUpdateReq struct {
	MCPServerToolID string `json:"mcpServerToolId" validate:"required"` //mcp server tool id
	MethodName      string `json:"methodName" validate:"required"`      // 显示名称
	Desc            string `json:"desc" validate:"required"`            // 描述
}

func (req *MCPServerToolUpdateReq) Check() error { return nil }

type MCPServerToolIDReq struct {
	MCPServerToolID string `json:"mcpServerToolId" validate:"required"` //mcp server tool id
}

func (req *MCPServerToolIDReq) Check() error { return nil }

type MCPServerOpenAPIToolCreate struct {
	MCPServerID   string                      `json:"mcpServerId" validate:"required"` // mcp server Id
	Name          string                      `json:"name" validate:"required"`        // 名称
	ApiAuth       CustomToolApiAuthWebRequest `json:"apiAuth" validate:"required"`     // api身份认证
	Schema        string                      `json:"schema"  validate:"required"`     // schema
	PrivacyPolicy string                      `json:"privacyPolicy"`                   // 隐私政策
	MethodNames   []string                    `json:"methodNames" validate:"required"` // API名称列表
}

func (req *MCPServerOpenAPIToolCreate) Check() error { return nil }
