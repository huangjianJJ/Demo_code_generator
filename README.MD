以下是一个满足你需求的 Go 代码生成器项目的详细实现，包含完整的目录结构和代码示例：

### 目录结构
```
code-generator/
├── cmd/
│   └── main.go
├── generator/
│   ├── generator.go
│   └── json_schema.go
└── go.mod
```

### 使用方法

#### 1. 初始化项目
在项目根目录下执行以下命令初始化 Go 模块：
```sh
go mod init code-generator
```
然后下载依赖：
```sh
go mod tidy
```

#### 2. 生成代码
#### 1. 自定义参数生成
```sh
go run cmd/main.go -pkg mypackage -name Person -fields "Name:string,Age:int"
```
这将在 `model` 文件夹下生成 `person.go` 文件。

#### 2. 单个 JSON Schema 文件生成
```sh
go run cmd/main.go -pkg mypackage -name Person -schema schema.json
```
同样会在 `model` 文件夹下生成 `person.go` 文件。

#### 3. 从 JSON Schema 文件夹生成
假设你有一个 `schemas` 文件夹，里面包含多个 JSON Schema 文件，执行以下命令：
```sh
go run cmd/main.go -pkg mypackage -schema-dir schemas
```
这将遍历 `schemas` 文件夹下的所有 JSON 文件，为每个文件生成对应的 Go 结构体代码，并将其写入 `model` 文件夹下的文件。

### 代码解释
- **`cmd/main.go`**：作为项目的入口，解析命令行参数并调用 `generator` 包中的函数生成代码。
- **`generator/generator.go`**：包含核心的代码生成逻辑，支持自定义包名、结构体名称和字段列表，同时生成构造函数、验证函数和接口方法。
- **`generator/json_schema.go`**：实现了从 JSON Schema 生成 Go 结构体代码的功能，通过解析 JSON Schema 文件，将其转换为 Go 结构体。

通过以上代码和步骤，你可以轻松实现一个支持多种功能的 Go 代码生成器。