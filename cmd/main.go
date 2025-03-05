package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code-generator/generator"
)

func main() {
	pkg := flag.String("pkg", "main", "自定义包名")
	name := flag.String("name", "MyStruct", "自定义结构体名称")
	fields := flag.String("fields", "", "自定义字段列表，格式：name1:type1,name2:type2")
	jsonSchemaDir := flag.String("schema-dir", "", "JSON Schema 文件所在文件夹路径")
	jsonSchema := flag.String("schema", "", "单个 JSON Schema 文件路径")

	flag.Parse()

	if *jsonSchemaDir != "" {
		// 遍历文件夹下的所有 JSON 文件
		err := filepath.Walk(*jsonSchemaDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
				structName := strings.TrimSuffix(info.Name(), ".json")
				data, err := generator.GenerateFromJSONSchema(*pkg, structName, path)
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

	var data []byte
	var err error
	if *jsonSchema != "" {
		data, err = generator.GenerateFromJSONSchema(*pkg, *name, *jsonSchema)
	} else {
		data, err = generator.Generate(*pkg, *name, *fields)
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
	fileName := strings.ToLower(*name) + ".go"
	filePath := filepath.Join(modelDir, fileName)

	// 写入文件
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("代码已成功生成并写入 %s\n", filePath)
}
