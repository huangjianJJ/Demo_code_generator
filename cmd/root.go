package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code-generator/generator"

	"github.com/spf13/cobra"
)

// 定义命令行标志变量
var (
	pkg           string
	name          string
	fields        string
	jsonSchemaDir string
	jsonSchema    string
)

// rootCmd 代表基础命令，当没有指定子命令时调用
var rootCmd = &cobra.Command{
	Use:   "code-generator",
	Short: "Generate Go structs from JSON Schema or custom fields",
	Long:  `Generate Go structs from JSON Schema or custom fields.`,
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		var err error

		if jsonSchemaDir != "" {
			fmt.Printf("正在处理 JSON Schema 文件夹: %s\n", jsonSchemaDir)
			// 处理 JSON Schema 文件夹路径的逻辑
			err := filepath.Walk(jsonSchemaDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
					fmt.Printf("找到 JSON 文件: %s\n", path)
					structName := strings.TrimSuffix(info.Name(), ".json")
					data, err = generator.GenerateFromJSONSchema(pkg, structName, path)
					if err != nil {
						fmt.Printf("生成 %s 失败: %v\n", path, err)
						return nil
					}

					// 创建 model 文件夹
					modelDir := "model"
					if _, err := os.Stat(modelDir); os.IsNotExist(err) {
						err := os.Mkdir(modelDir, 0755)
						if err != nil {
							fmt.Printf("创建 model 文件夹失败: %v\n", err)
							return err
						}
					}

					// 生成文件名
					fileName := strings.ToLower(structName) + ".go"
					filePath := filepath.Join(modelDir, fileName)

					// 写入文件
					err = os.WriteFile(filePath, data, 0644)
					if err != nil {
						fmt.Printf("写入 %s 失败: %v\n", filePath, err)
						return nil
					}

					fmt.Printf("代码已成功生成并写入 %s\n", filePath)
				}
				return nil
			})
			if err != nil {
				fmt.Printf("遍历文件夹失败: %v\n", err)
				os.Exit(1)
			}
			return
		}

		if jsonSchema != "" {
			data, err = generator.GenerateFromJSONSchema(pkg, name, jsonSchema)
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
	rootCmd.Flags().StringVarP(&jsonSchemaDir, "schema-dir", "d", "", "JSON Schema 文件所在文件夹路径")
	rootCmd.Flags().StringVarP(&jsonSchema, "schema", "s", "", "单个 JSON Schema 文件路径")
	// 打印调试信息
	fmt.Printf("pkg 标志: %v\n", rootCmd.Flags().Lookup("pkg"))
	fmt.Printf("schema-dir 标志: %v\n", rootCmd.Flags().Lookup("schema-dir"))
}

// Execute 函数用于执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
