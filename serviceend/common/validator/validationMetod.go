package validator

import (
	"dataPanel/serviceend/model"
	"database/sql/driver"
	"reflect"
	"regexp"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
)

type validationMetod struct {
}

// AddValidationMethod 注册tag验证器
func AddValidationMethod(v **validator.Validate) {
	metod := NewValidationMetod()
	_ = (*v).RegisterValidation("checkPhone", metod.CheckPhoneMetod)
}

// ValidateValuer implements validator.CustomTypeFunc 自定义校验规则
func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}
	//spec.LocalTime 类型的自定义校验规则
	if field.Type() == reflect.TypeOf(model.LocalTime{}) {
		timeStr := field.Interface().(model.LocalTime).String()
		// 0001-01-01 00:00:00 是 go 中 time.Time 类型的空值
		// 这里返回 Nil 则会被 validator 判定为空值，而无法通过 `binding:"required"` 规则
		if timeStr == "0001-01-01 00:00:00" {
			return nil
		}
		return timeStr
	}
	return nil
}

// AddTranslation 自定义翻译器
func AddTranslation(v **validator.Validate, translator ut.Translator) {
	_ = (*v).RegisterTranslation("required", translator, registerTranslator("required", "{0}必填项[自定义翻译器]"), translate)
	_ = (*v).RegisterTranslation("checkPhone", translator, registerTranslator("checkPhone", "{0}校验不通过[自定义翻译器]"), translate)
	_ = (*v).RegisterTranslation("oneof", translator, registerTranslator("oneof", "{0}校验不通过[自定义翻译器]"), translate)
}

// registerTranslator 为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}

/*
	所有自定义tag 验证方法名
*/

func NewValidationMetod() *validationMetod {
	return &validationMetod{}
}

func (r validationMetod) CheckPhoneMetod(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(v)
}
