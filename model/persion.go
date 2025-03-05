// 此文件由代码生成器自动生成，请勿手动修改！
package mypackage

type persion struct {
    Name string `json:"Name"`
    Age  int `json:"Age"`
}

func Newpersion(Name string, Age int) *persion {
    return &persion{
        Name: Name,
        Age: Age,
    }
}

func (s *persion) Validate() error {
    // 这里可以添加具体的验证逻辑
    return nil
}

type persionInterface interface {
    Validate() error
}

