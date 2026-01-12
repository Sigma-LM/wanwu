package constant

// openapi type
const (
	OpenAPITypeChatflow  = "chatflow"  // 对话问答
	OpenAPITypeWorkflow  = "workflow"  // 工作流
	OpenAPITypeAgent     = "agent"     // 智能体
	OpenAPITypeRag       = "rag"       // 文本问答
	OpenAPITypeKnowledge = "knowledge" // 知识库
)

// app type
const (
	AppTypeAgent     = "agent"     // 智能体
	AppTypeRag       = "rag"       // 文本问答
	AppTypeWorkflow  = "workflow"  // 工作流
	AppTypeChatflow  = "chatflow"  // 对话流
	AppTypeMCPServer = "mcpserver" // mcp server
)

// app publish type
const (
	AppPublishPublic       = "public"       // 系统公开发布
	AppPublishOrganization = "organization" // 组织公开发布
	AppPublishPrivate      = "private"      // 私密发布
)

// tool type
const (
	ToolTypeBuiltIn = "builtin" // 内置工具
	ToolTypeCustom  = "custom"  // 自定义工具
)

// mcp type
const (
	MCPTypeMCP       = "mcp"       // mcp
	MCPTypeMCPServer = "mcpserver" // mcp server
)

// mcp server tool type
const (
	MCPServerToolTypeCustomTool  = "custom"  // 自定义工具
	MCPServerToolTypeBuiltInTool = "builtin" // 内置工具
	MCPServerToolTypeOpenAPI     = "openapi" // 用户导入的openapi
)

// model experience template
const (
	FileItemTemplate = "[file name]: {{.FileName}}\n" +
		"[file content begin]\n" +
		"{{.FileContent}}\n" +
		"[file content end]\n"
	ModelExperienceTemplate = "## 任务说明\n你将根据提供的一个或多个文档内容，准确、详细地回答用户提出的问题。\n\n## 文档内容\n以下是多个文档片段的合并内容，可能来源于多个不同文件或章节：\n\n```\n{{.Context}}\n```\n\n## 用户提问\n{{.Question}}\n\n## 回答要求\n请严格按照以下要求作答：\n1. **仅依据提供的文档内容**进行回答，**禁止使用常识、经验或外部知识补充内容**。\n2. **答案必须完整、详细**，清晰说明推理依据。"
)
