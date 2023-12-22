package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var (
	uni           *ut.UniversalTranslator
	Trans         ut.Translator
	CustomMsgFlag = "[@#CSMTM#]"
)

func Init() error {
	// 注册翻译器
	zhT := zh.New()
	uni = ut.New(zhT, zhT)

	Trans, _ = uni.GetTranslator("zh")

	// 获取gin的校验器
	validate := binding.Validator.Engine().(*validator.Validate)
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		customMsg := field.Tag.Get("errMsg")
		if customMsg != "" {
			return CustomMsgFlag + customMsg + CustomMsgFlag
		}

		customName := field.Tag.Get("name")
		if customName != "" {
			return customName
		}
		jsonName := field.Tag.Get("json")
		if jsonName != "" {
			return jsonName
		}
		return field.Name
	})

	// 注册翻译器
	return zhtranslations.RegisterDefaultTranslations(validate, Trans)
}

// GetCustomMsg 获取自定义的错误信息
func GetCustomMsg(msg string) string {
	start := strings.Index(msg, CustomMsgFlag)
	if start < 0 {
		return ""
	}

	end := strings.LastIndex(msg, CustomMsgFlag)
	return msg[start+len(CustomMsgFlag) : end]
}
