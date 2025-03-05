package generator

import (
	"bytes"
	"fmt"
	"strings"
)

// Generate 生成 Go 结构体代码
func Generate(pkg, name, fieldsStr string) ([]byte, error) {
	fields := parseFields(fieldsStr)
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

type Field struct {
	Name string
	Type string
}

func parseFields(fieldsStr string) []Field {
	var fields []Field
	if fieldsStr == "" {
		return fields
	}
	for _, fieldStr := range strings.Split(fieldsStr, ",") {
		parts := strings.Split(fieldStr, ":")
		if len(parts) == 2 {
			fields = append(fields, Field{
				Name: parts[0],
				Type: parts[1],
			})
		}
	}
	return fields
}

func getMaxFieldLength(fields []Field) int {
	maxLen := 0
	for _, field := range fields {
		if len(field.Name) > maxLen {
			maxLen = len(field.Name)
		}
	}
	return maxLen
}

func formatConstructorParams(fields []Field) string {
	var params []string
	for _, field := range fields {
		params = append(params, fmt.Sprintf("%s %s", field.Name, field.Type))
	}
	return strings.Join(params, ", ")
}
