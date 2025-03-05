// 此文件由代码生成器自动生成，请勿手动修改！
package mypackage

type Person struct {
    Name string `json:"Name"`
    Age  int `json:"Age"`
}

func NewPerson(Name string, Age int) *Person {
    return &Person{
        Name: Name,
        Age: Age,
    }
}

func (s *Person) Validate() error {
    // 这里可以添加具体的验证逻辑
    return nil
}

type PersonInterface interface {
    Validate() error
}

