`cobra-cli` 是一个用于快速创建和管理基于 `cobra` 库的命令行工具的脚手架。下面详细介绍使用 `cobra-cli` 库封装代码生成器命令行工具的步骤：

### 1. 安装 `cobra-cli`
首先要确保你已经安装了 Go 环境，然后使用以下命令安装 `cobra-cli`：
```sh
go install github.com/spf13/cobra-cli@latest
```

### 2. 创建新的 `cobra` 项目
使用 `cobra-cli` 初始化一个新的命令行工具项目。在命令行中执行以下命令，这里将项目命名为 `code-generator`，你可以根据实际情况修改：
```sh
cobra-cli init 
```
此命令会创建一个基本的项目结构，包含必要的文件和目录，如 `cmd` 目录（用于存放命令相关代码）、`main.go` 文件（程序入口）等。

### 3. 项目结构
初始化完成后，项目结构大致如下：
```
code-generator/
├── cmd/
│   ├── root.go
│   └── ...
├── go.mod
├── go.sum
└── main.go
```

### 4. 定义命令和标志
#### 编辑 `cmd/root.go` 文件
`root.go` 是根命令的定义文件，以下是对其进行修改以适配代码生成器的示例：
```go
package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
    "github.com/yourusername/code-generator/generator"
)

// 定义命令行标志变量
var (
    pkg    string
    name   string
    fields string
    schema string
)

// rootCmd 代表基础命令，当没有指定子命令时调用
var rootCmd = &cobra.Command{
    Use:   "code-generator",
    Short: "Generate Go structs from JSON Schema or custom fields",
    Long:  `Generate Go structs from JSON Schema or custom fields.`,
    Run: func(cmd *cobra.Command, args []string) {
        var data []byte
        var err error
        // 根据是否提供 JSON Schema 文件选择不同的生成方式
        if schema != "" {
            data, err = generator.GenerateFromJSONSchema(pkg, name, schema)
        } else {
            data, err = generator.Generate(pkg, name, fields)
        }

        if err != nil {
            fmt.Printf("生成失败: %v\n", err)
            os.Exit(1)
        }

        // 创建 model 文件夹
        modelDir := "model"
        if _, err := os.Stat(modelDir); os.IsNotExist(err) {
            err := os.Mkdir(modelDir, 0755)
            if err != nil {
                fmt.Printf("创建 model 文件夹失败: %v\n", err)
                os.Exit(1)
            }
        }

        // 生成文件名
        fileName := strings.ToLower(name) + ".go"
        filePath := filepath.Join(modelDir, fileName)

        // 写入文件
        err = os.WriteFile(filePath, data, 0644)
        if err != nil {
            fmt.Printf("写入文件失败: %v\n", err)
            os.Exit(1)
        }

        fmt.Printf("代码已成功生成并写入 %s\n", filePath)
    },
}

// init 函数用于初始化命令行标志
func init() {
    rootCmd.Flags().StringVarP(&pkg, "pkg", "p", "main", "自定义包名")
    rootCmd.Flags().StringVarP(&name, "name", "n", "MyStruct", "自定义结构体名称")
    rootCmd.Flags().StringVarP(&fields, "fields", "f", "", "自定义字段列表，格式：name1:type1,name2:type2")
    rootCmd.Flags().StringVarP(&schema, "schema", "s", "", "JSON Schema 文件路径")
}

// Execute 函数用于执行根命令
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```
#### 代码解释：
- **标志定义**：在 `init` 函数中，使用 `rootCmd.Flags().StringVarP` 方法定义了四个命令行标志：`pkg`、`name`、`fields` 和 `schema`，分别用于指定包名、结构体名称、自定义字段列表和 JSON Schema 文件路径。
- **命令执行逻辑**：在 `rootCmd` 的 `Run` 函数中，根据是否提供了 `schema` 参数，调用不同的生成函数。生成代码后，将代码写入 `model` 文件夹下的文件。
- **错误处理**：在每个关键步骤都进行了错误处理，如果出现错误，会输出错误信息并退出程序。

### 5. 实现代码生成逻辑
在项目中创建 `generator` 包（如果还没有），并在其中实现代码生成的核心逻辑。以下是简单示例，假设 `generator/generator.go` 文件内容如下：
```go
package generator

import (
    "bytes"
    "fmt"
    "strings"
)

// Generate 从自定义字段生成代码
func Generate(pkg, name, fieldsStr string) ([]byte, error) {
    // 解析字段
    fields := parseFields(fieldsStr)
    buffer := &bytes.Buffer{}

    // 写入包名
    buffer.WriteString(fmt.Sprintf("package %s\n\n", pkg))
    // 写入结构体定义
    buffer.WriteString(fmt.Sprintf("type %s struct {\n", name))
    for _, field := range fields {
        buffer.WriteString(fmt.Sprintf("    %s %s\n", field.Name, field.Type))
    }
    buffer.WriteString("}\n")

    return buffer.Bytes(), nil
}

// GenerateFromJSONSchema 从 JSON Schema 生成代码
func GenerateFromJSONSchema(pkg, name, schemaPath string) ([]byte, error) {
    // 这里需要实现解析 JSON Schema 文件并生成代码的逻辑
    // 示例中简单返回一个占位信息
    buffer := &bytes.Buffer{}
    buffer.WriteString(fmt.Sprintf("package %s\n\n", pkg))
    buffer.WriteString(fmt.Sprintf("type %s struct {\n", name))
    buffer.WriteString("    // 待从 JSON Schema 生成字段\n")
    buffer.WriteString("}\n")

    return buffer.Bytes(), nil
}

// 解析字段字符串
func parseFields(fieldsStr string) []struct {
    Name string
    Type string
} {
    var fields []struct {
        Name string
        Type string
    }
    if fieldsStr == "" {
        return fields
    }
    parts := strings.Split(fieldsStr, ",")
    for _, part := range parts {
        fieldParts := strings.Split(part, ":")
        if len(fieldParts) == 2 {
            fields = append(fields, struct {
                Name string
                Type string
            }{
                Name: fieldParts[0],
                Type: fieldParts[1],
            })
        }
    }
    return fields
}
```

### 6. 修改 `main.go` 文件
在项目根目录下的 `main.go` 文件中，调用 `cmd.Execute()` 来启动命令行工具。示例代码如下：
```go
package main

import "github.com/yourusername/code-generator/cmd"

func main() {
    cmd.Execute()
}
```

### 7. 测试命令行工具
在终端中，进入项目根目录，运行以下命令来测试命令行工具：
- **使用自定义字段生成代码**：
```sh
go run main.go --pkg mypackage --name Person --fields "Name:string,Age:int"
```
- **使用 JSON Schema 文件生成代码**：
```sh
go run main.go --pkg mypackage --name Person --schema schema.json
```

### 8. 构建和安装命令行工具
如果你想将命令行工具安装到系统中，可以使用以下命令进行构建和安装：
```sh
go install
```
安装完成后，你可以在任何目录下直接使用 `code-generator` 命令来生成代码。

### 9. 扩展功能（可选）
你可以使用 `cobra-cli` 继续添加子命令、更多标志等。例如，添加一个 `version` 子命令来显示工具版本：
```sh
cobra-cli add version
```
然后编辑 `cmd/version.go` 文件，实现版本显示逻辑：
```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print the version number of code-generator",
    Long:  `All software has versions. This is code-generator's`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("code-generator v0.1.0")
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
```
之后就可以使用 `code-generator version` 命令查看版本信息。

通过以上步骤，你就可以使用 `cobra-cli` 封装一个功能较为完善的代码生成器命令行工具。 