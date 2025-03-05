// 此文件由代码生成器自动生成，请勿手动修改！
package mypackage

type schema struct {
    Name string `json:"Name"`
    Age  int `json:"Age"`
}

func Newschema(Name string, Age int) *schema {
    return &schema{
        Name: Name,
        Age: Age,
    }
}

func (s *schema) Validate() error {
    // 这里可以添加具体的验证逻辑
    return nil
}

type schemaInterface interface {
    Validate() error
}

