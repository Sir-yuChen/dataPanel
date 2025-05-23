package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	vk := make(map[string]interface{})
	var v = "{\"name\":\"主题颜色\",\"value\":\"white\"}"
	err := json.Unmarshal([]byte(fmt.Sprintf("%v", v)), &vk)
	if err != nil {
		fmt.Sprintf("转化异常" + err.Error())
	}
	fmt.Sprintf("数据：%s", vk)
}
