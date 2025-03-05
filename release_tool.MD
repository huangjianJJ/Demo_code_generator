将使用 `cobra-cli` 封装的代码生成器打包发布，让其他人能像使用 `gqlgen generate` 一样方便地使用，可按以下步骤操作：

### 1. 项目准备
确保你的代码生成器项目已经使用 `cobra-cli` 封装好，项目结构合理，核心生成逻辑和命令行交互部分分离。假设项目结构如下：
```
code-generator/
├── cmd/
│   ├── root.go
│   ├── version.go
│   └── ...
├── generator/
│   ├── generator.go
│   └── ...
├── go.mod
├── go.sum
└── main.go
```

### 2. 编写详细的文档
在项目根目录下创建 `README.md` 文件，详细说明代码生成器的功能、使用方法、命令行参数、安装方式等信息。示例如下：
```markdown
# Code Generator

这是一个用于生成 Go 结构体代码的命令行工具，支持自定义包名、结构体名称、字段列表以及从 JSON Schema 文件生成代码。

## 功能特性
- 自定义包名和结构体名称
- 支持自定义字段列表生成结构体
- 从 JSON Schema 文件生成结构体代码
- 自动生成构造函数、验证函数和接口方法

## 安装
### 使用 `go install`
```sh
go install github.com/yourusername/code-generator@latest
```

### 从发布版本下载
访问 [GitHub Releases](https://github.com/yourusername/code-generator/releases) 页面，下载适合你操作系统的预编译二进制文件，并将其添加到系统的 `PATH` 环境变量中。

## 使用方法
### 使用自定义字段生成
```sh
code-generator --pkg mypackage --name Person --fields "Name:string,Age:int"
```

### 使用 JSON Schema 文件生成
```sh
code-generator --pkg mypackage --name Person --schema schema.json
```

## 命令选项
- `--pkg`：指定生成代码的包名，默认为 `main`。
- `--name`：指定生成结构体的名称，默认为 `MyStruct`。
- `--fields`：自定义字段列表，格式为 `name1:type1,name2:type2`。
- `--schema`：指定 JSON Schema 文件的路径。
```

### 3. 版本控制与标签管理
#### 初始化 Git 仓库
如果你还没有使用 Git 管理项目，先初始化 Git 仓库：
```sh
git init
git add .
git commit -m "Initial commit"
```

#### 创建远程仓库
在 GitHub 或其他代码托管平台上创建一个新的仓库，并将本地仓库与远程仓库关联：
```sh
git remote add origin <GitHub 仓库地址>
git push -u origin main
```

#### 打标签
为项目打标签，以便用户可以指定版本进行安装。遵循语义化版本规则，例如：
```sh
git tag v1.0.0
git push origin v1.0.0
```

### 4. 打包构建
#### 跨平台构建
为了让不同操作系统的用户都能使用你的工具，需要进行跨平台构建。可以使用以下脚本进行构建：
```sh
#!/bin/bash

# 定义目标平台
PLATFORMS="linux/amd64 linux/arm64 windows/amd64 darwin/amd64 darwin/arm64"

# 创建 releases 目录
mkdir -p releases

# 遍历目标平台进行构建
for platform in $PLATFORMS; do
    IFS='/' read -r OS ARCH <<< "$platform"
    output_name="releases/code-generator-$OS-$ARCH"
    if [ $OS = "windows" ]; then
        output_name+=".exe"
    fi

    env GOOS=$OS GOARCH=$ARCH go build -o $output_name main.go
    if [ $? -ne 0 ]; then
        echo "构建 $platform 失败"
        exit 1
    fi
    echo "构建 $platform 成功"
done
```
将上述脚本保存为 `build.sh`，并赋予执行权限：
```sh
chmod +x build.sh
./build.sh
```
执行脚本后，会在 `releases` 目录下生成不同平台的可执行文件。

### 5. 发布到 GitHub Releases
#### 创建发布版本
在 GitHub 仓库的 `Releases` 页面，点击 `Draft a new release`。
- **选择标签**：选择之前打的标签（如 `v1.0.0`）。
- **添加标题和描述**：输入发布版本的标题和详细描述，介绍该版本的新特性、修复的问题等。
- **上传二进制文件**：将 `releases` 目录下生成的可执行文件上传到发布页面。

#### 发布
点击 `Publish release` 完成发布。

### 6. 让其他人使用
#### 使用 `go install` 安装
用户可以使用 `go install` 命令从 GitHub 仓库安装你的代码生成器：
```sh
go install github.com/yourusername/code-generator@latest
```
安装完成后，就可以在命令行中使用 `code-generator` 命令。

#### 使用预编译的二进制文件
用户也可以从 GitHub Releases 页面下载适合自己操作系统的预编译二进制文件，将其添加到系统的 `PATH` 环境变量中，然后就可以在命令行中使用。例如，在 Linux 系统中，可以将下载的二进制文件移动到 `/usr/local/bin` 目录：
```sh
sudo mv code-generator-linux-amd64 /usr/local/bin/code-generator
```

### 7. 持续维护与更新
随着项目的发展，你可能需要不断添加新功能、修复 bug。每次更新后，记得打新的标签并发布新版本，同时更新 `README.md` 文件中的安装和使用说明。

通过以上步骤，你就可以将使用 `cobra-cli` 封装的代码生成器打包发布，让其他人方便地使用。 