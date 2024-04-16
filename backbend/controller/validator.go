package controller

import (
	"fmt"
	"reflect"
	"strings"
	"backbend/models"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/gin-gonic/gin/binding"
 	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

/**
 * @Author wkwar
 * @Description //TODO 自定义参数校验
 * @Date 14:00 2022/3/17
 **/

// 定义一个全局翻译器
var trans ut.Translator

// locale 通常取决于 http 请求头的 'Accept-Language'
func InitValidator(locale string) (err error) {
	//注册自定义验证
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义 json tag 函数
		//reflect.StructField 可以获取结构体类型信息
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			//字符串按照，分割，最多返回两个
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		//注册自定义结构体验证
		v.RegisterStructValidation(SignUpParamStructLevelValidate, models.RegisterForm{})

		//注册翻译器
		err = InitTrans(v, locale)	
	}

	return
}

//初始化翻译器
func InitTrans(v *validator.Validate, locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//自定义验证翻译器 --- 将英文翻译为中文，用于测试的时候使用
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// // 添加额外翻译
		// _ = v.RegisterTranslation("required_with", Trans, func(ut ut.Translator) error {
		// 	return ut.Add("required_with", "{0} 为必填字段!", true)
		// }, func(ut ut.Translator, fe validator.FieldError) string {
		// 	t, _ := ut.T("required_with", fe.Field())
		// 	return t
		// })
		// _ = v.RegisterTranslation("required_without", Trans, func(ut ut.Translator) error {
		// 	return ut.Add("required_without", "{0} 为必填字段!", true)
		// }, func(ut ut.Translator, fe validator.FieldError) string {
		// 	t, _ := ut.T("required_without", fe.Field())
		// 	return t
		// })
		// _ = v.RegisterTranslation("required_without_all", Trans, func(ut ut.Translator) error {
		// 	return ut.Add("required_without_all", "{0} 为必填字段!", true)
		// }, func(ut ut.Translator, fe validator.FieldError) string {
		// 	t, _ := ut.T("required_without_all", fe.Field())
		// 	return t
		// })

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

//定义一个去掉结构体名称前缀的自定义方法：
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

//自定义注册 密码 与 再次输入密码 验证方法
func SignUpParamStructLevelValidate(sl validator.StructLevel) {
	su := sl.Current().Interface().(models.RegisterForm)
	//判断注册用户密码是输入否正确
	if su.Password != su.ConfirmPassword {
		//字段数据，字段名，结构体字段名，tag，参数
		sl.ReportError(su.ConfirmPassword, "confirm_password", "ConfirmPassword", "eqfield", "password")
	}
}