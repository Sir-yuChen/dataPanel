package utils

import (
	"dataPanel/serviceend/common/ApiReturn"
	"dataPanel/serviceend/common/response"
	"dataPanel/serviceend/global"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func addValueToMap(fields map[string]string) map[string]interface{} {
	res := make(map[string]interface{})
	for field, err := range fields {
		fieldArr := strings.SplitN(field, ".", 2)
		if len(fieldArr) > 1 {
			NewFields := map[string]string{fieldArr[1]: err}
			returnMap := addValueToMap(NewFields)
			if res[fieldArr[0]] != nil {
				for k, v := range returnMap {
					res[fieldArr[0]].(map[string]interface{})[k] = v
				}
			} else {
				res[fieldArr[0]] = returnMap
			}
			continue
		} else {
			res[field] = err
			continue
		}
	}
	return res
}

// 去掉结构体名称前缀
func removeTopStruct(fields map[string]string) map[string]interface{} {
	lowerMap := map[string]string{}
	for field, err := range fields {
		fieldArr := strings.SplitN(field, ".", 2)
		lowerMap[fieldArr[1]] = err
	}
	res := addValueToMap(lowerMap)
	return res
}

// ErrValidatorResp ErrResp handler中调用的错误翻译方法
func ErrValidatorResp(err error, method string, req interface{}, c *gin.Context) {
	global.GvaLog.Error(method+"参数校验不通过", zap.Any("Request", req), zap.Any("error", err))
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		global.GvaLog.Error(method+"翻译ERROR异常", zap.Any("error", err))
		response.WithApiReturn(ApiReturn.ErrSystem, c)
		return
	}
	returnCode := ApiReturn.ErrCheckParameterFailed
	translate := errs.Translate(global.GvaTrans)
	m := removeTopStruct(translate)
	msg := ""
	for _, value := range m {
		if msg == "" {
			msg = fmt.Sprintf("%s", value)
		} else {
			msg += "," + fmt.Sprintf("%s", value)
		}
	}
	returnCode.Msg = msg
	response.WithApiReturn(returnCode, c)
}
