package code

import (
	"dataPanel/serviceend/common/validator"
	"dataPanel/serviceend/global"
	"database/sql"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validatorv10 "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	trans ut.Translator
)

// InitTrans 初始化翻译器
func InitTrans(locale string) {
	//修改gin框架中的Validator属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validatorv10.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			label := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]
			if label == "-" || label == "" {
				label = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			}
			return label
		})

		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		uni := ut.New(zhT, zhT, enT) // 万能翻译器，保存所有的语言环境和翻译数据

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			global.GvaLog.Error("未查找到指定校验器！！")
			os.Exit(-1)
		}

		// register all sql.Null* types to use the ValidateValuer CustomTypeFunc
		v.RegisterCustomTypeFunc(validator.ValidateValuer, sql.NullString{}, sql.NullInt64{}, sql.NullInt32{}, sql.NullBool{}, sql.NullFloat64{})
		// 注意！因为这里会使用到trans实例
		// 所以这一步注册要放到trans初始化的后面
		// 添加额外翻译
		validator.AddTranslation(&v, trans)
		//自定义验证方法
		validator.AddValidationMethod(&v)
		// 注册翻译器
		switch locale {
		case "en":
			enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			zhTranslations.RegisterDefaultTranslations(v, trans)
		}
		global.GvaTrans = &trans
		return
	}
	return
}
