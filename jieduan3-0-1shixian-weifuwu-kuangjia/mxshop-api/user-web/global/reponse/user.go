package reponse

import (
	"fmt"
	"time"
)

type JsonTime time.Time

// 只要一个类型 实现了 MarshalJSON() 方法Go 的 json.Marshal 就会自动调用它，不用你手动调
// Go 的 json 序列化工具 看到这个字段是 JsonTime 类型，发现它实现了 MarshalJSON()，就自动调用
// 该项目中的c.JSON(data) 函数就会把 data 序列化成 JSON 字符串(并调用这个实现MarshalJSON函数)，并返回给前端
func (j JsonTime) MarshalJSON() ([]byte, error) {
	// 1. 把自定义类型 JsonTime 转回 time.Time
	// 2. 调用格式化方法成 2006-01-02 字符串
	// // - 只有time.Time类型才有格式化方法，自定义类型没有这个方法，虽然是基于time.Time类型也没有
	// 3. 给字符串前后套上双引号 " （JSON 字符串必须带引号）

	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	// 把字符串转成 []byte 返回给 JSON 序列化工具
	// 将字符串转成 []byte要求的字节切片类型
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32  `json:"id"`
	NickName string `json:"name"`
	//Birthday string `json:"birthday"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}
