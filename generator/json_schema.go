package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// GenerateFromJSONSchema 从 JSON Schema 生成 Go 结构体代码
func GenerateFromJSONSchema(pkg, name, schemaPath string) ([]byte, error) {
	schemaData, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	var schema map[string]interface{}
	err = json.Unmarshal(schemaData, &schema)
	if err != nil {
		return nil, err
	}

	fields := parseSchemaFields(schema)
	buffer := &bytes.Buffer{}

	// 文件头
	buffer.WriteString("// 此文件由代码生成器自动生成，请勿手动修改！\n")
	buffer.WriteString(fmt.Sprintf("package %s\n\n", pkg))

	// 结构体定义
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", name))
	maxFieldLen := getMaxFieldLength(fields)
	for _, field := range fields {
		buffer.WriteString(fmt.Sprintf("    %-*s %s `json:\"%s\"`\n", maxFieldLen, field.Name, field.Type, field.Name))
	}
	buffer.WriteString("}\n\n")

	// 构造函数
	buffer.WriteString(fmt.Sprintf("func New%s(%s) *%s {\n", name, formatConstructorParams(fields), name))
	buffer.WriteString(fmt.Sprintf("    return &%s{\n", name))
	for _, field := range fields {
		buffer.WriteString(fmt.Sprintf("        %s: %s,\n", field.Name, field.Name))
	}
	buffer.WriteString("    }\n")
	buffer.WriteString("}\n\n")

	// 验证函数
	buffer.WriteString(fmt.Sprintf("func (s *%s) Validate() error {\n", name))
	buffer.WriteString("    // 这里可以添加具体的验证逻辑\n")
	buffer.WriteString("    return nil\n")
	buffer.WriteString("}\n\n")

	// 接口方法
	buffer.WriteString(fmt.Sprintf("type %sInterface interface {\n", name))
	buffer.WriteString("    Validate() error\n")
	buffer.WriteString("}\n\n")

	return buffer.Bytes(), nil
}

func parseSchemaFields(schema map[string]interface{}) []Field {
	var fields []Field
	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		return fields
	}
	for name, prop := range properties {
		propMap, ok := prop.(map[string]interface{})
		if !ok {
			continue
		}
		fieldType := getGoType(propMap)
		fields = append(fields, Field{
			Name: name,
			Type: fieldType,
		})
	}
	return fields
}

func getGoType(prop map[string]interface{}) string {
	schemaType, ok := prop["type"].(string)
	if !ok {
		return "interface{}"
	}
	switch schemaType {
	case "string":
		return "string"
	case "number":
		return "float64"
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	case "object":
		return "map[string]interface{}"
	case "array":
		items, ok := prop["items"].(map[string]interface{})
		if !ok {
			return "[]interface{}"
		}
		itemType := getGoType(items)
		return fmt.Sprintf("[]%s", itemType)
	default:
		return "interface{}"
	}
}
