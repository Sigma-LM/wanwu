package openapi3_util

import (
	"context"
	"slices"

	eino_schema "github.com/cloudwego/eino/schema"
	"github.com/getkin/kin-openapi/openapi3"
)

func Schema2EinoTools(ctx context.Context, schema []byte) ([]*eino_schema.ToolInfo, error) {
	doc, err := LoadFromData(ctx, schema)
	if err != nil {
		return nil, err
	}
	return Doc2EinoTools(doc)
}

func Schema2EinoTool(ctx context.Context, schema []byte, operationID string) (*eino_schema.ToolInfo, error) {
	doc, err := LoadFromData(ctx, schema)
	if err != nil {
		return nil, err
	}
	return Doc2EinoTool(doc, operationID)
}

func Doc2EinoTools(doc *openapi3.T) ([]*eino_schema.ToolInfo, error) {
	var rets []*eino_schema.ToolInfo
	for _, pathItem := range doc.Paths {
		for _, operation := range pathItem.Operations() {
			rets = append(rets, Operation2EinoTool(operation))
		}
	}
	return rets, nil
}

func Doc2EinoTool(doc *openapi3.T, operationID string) (*eino_schema.ToolInfo, error) {
	var exist bool
	var ret *eino_schema.ToolInfo
	for _, pathItem := range doc.Paths {
		for _, operation := range pathItem.Operations() {
			if operation.OperationID != operationID {
				continue
			}
			exist = true
			ret = Operation2EinoTool(operation)
			break
		}
	}
	if !exist {
		return nil, nil
	}
	return ret, nil
}

func Operation2EinoTool(operation *openapi3.Operation) *eino_schema.ToolInfo {
	ret := &eino_schema.ToolInfo{
		Name: operation.OperationID,
		Desc: operation.Description,
	}
	// 处理description，保证非空
	if ret.Desc == "" {
		if operation.Summary != "" {
			ret.Desc = operation.Summary
		} else {
			ret.Desc = operation.OperationID
		}
	}
	params := make(map[string]*eino_schema.ParameterInfo)
	// 解析路径参数、查询参数、header 参数等
	for field, param := range Parameters2EinoParams(operation.Parameters) {
		params[field] = param
	}
	// 解析请求体
	if operation.RequestBody != nil && operation.RequestBody.Value != nil {
		for _, mediaType := range operation.RequestBody.Value.Content {
			if mediaType.Schema != nil && mediaType.Schema.Value != nil {
				for field, param := range Schemas2EinoParams(mediaType.Schema.Value.Properties, mediaType.Schema.Value.Required) {
					params[field] = param
				}
			}
		}
	}
	ret.ParamsOneOf = eino_schema.NewParamsOneOfByParams(params)
	return ret
}

func Parameters2EinoParams(parameters openapi3.Parameters) map[string]*eino_schema.ParameterInfo {
	if parameters == nil {
		return nil
	}

	rets := make(map[string]*eino_schema.ParameterInfo)
	for _, parameter := range parameters {
		if parameter.Value == nil {
			continue
		}
		rets[parameter.Value.Name] = Parameter2EinoParam(parameter.Value)
	}

	return rets
}

func Parameter2EinoParam(parameter *openapi3.Parameter) *eino_schema.ParameterInfo {
	if parameter == nil {
		return nil
	}

	dataType := ParameterType2EinoDataType(parameter)
	ret := &eino_schema.ParameterInfo{
		Type:     dataType,
		Desc:     parameter.Description,
		Required: parameter.Required,
		// todo enum
	}
	switch dataType {
	case eino_schema.Object:
		if parameter.Schema != nil && parameter.Schema.Value != nil {
			ret.SubParams = Schemas2EinoParams(parameter.Schema.Value.Properties, parameter.Schema.Value.Required)
		}
	case eino_schema.Array:
		if parameter.Schema != nil && parameter.Schema.Value != nil && parameter.Schema.Value.Items != nil {
			ret.ElemInfo = Schema2EinoParam(parameter.Schema.Value.Items.Value, false)
		}
	default:
	}

	return ret
}

func Schemas2EinoParams(schemas openapi3.Schemas, required []string) map[string]*eino_schema.ParameterInfo {
	if schemas == nil {
		return nil
	}

	rets := make(map[string]*eino_schema.ParameterInfo)
	for propName, propSchema := range schemas {
		if propSchema == nil {
			continue
		}
		rets[propName] = Schema2EinoParam(propSchema.Value, slices.Contains(required, propName))
	}

	return rets
}

func Schema2EinoParam(schema *openapi3.Schema, required bool) *eino_schema.ParameterInfo {
	if schema == nil {
		return nil
	}

	dataType := SchemaType2EinoDataType(schema)
	ret := &eino_schema.ParameterInfo{
		Type:     dataType,
		Desc:     schema.Description,
		Required: required,
		// todo enum
	}
	switch dataType {
	case eino_schema.Object:
		ret.SubParams = Schemas2EinoParams(schema.Properties, schema.Required)
	case eino_schema.Array:
		if schema.Items != nil {
			ret.ElemInfo = Schema2EinoParam(schema.Items.Value, false)
		}
	}

	return ret
}

func ParameterType2EinoDataType(parameter *openapi3.Parameter) eino_schema.DataType {
	if parameter.Schema == nil {
		return eino_schema.Null
	}
	return SchemaType2EinoDataType(parameter.Schema.Value)
}

func SchemaType2EinoDataType(schema *openapi3.Schema) eino_schema.DataType {
	if schema == nil {
		return eino_schema.Null
	}
	switch schema.Type {
	case openapi3.TypeObject:
		return eino_schema.Object
	case openapi3.TypeArray:
		return eino_schema.Array
	case openapi3.TypeString:
		return eino_schema.String
	case openapi3.TypeNumber:
		return eino_schema.Number
	case openapi3.TypeInteger:
		return eino_schema.Integer
	case openapi3.TypeBoolean:
		return eino_schema.Boolean
	default:
		return eino_schema.Null
	}
}
