package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	zh2 "github.com/go-playground/validator/v10/translations/zh"
	"github.com/go-playground/validator/v10/translations/zh_tw"
)

var uni *ut.UniversalTranslator

func init() {
	uni = ut.New(en.New(), zh.New(), zh_Hant_TW.New())
	supportLocale := []string{"zh", "en", "zh_tw", "zh_hk"}
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	for _, locale := range supportLocale {
		trans, _ := uni.GetTranslator(locale)
		if ok {
			switch locale {
			case "zh":
				_ = zh2.RegisterDefaultTranslations(v, trans)
			case "en":
				_ = en2.RegisterDefaultTranslations(v, trans)
			case "zh_tw":
				_ = zh_tw.RegisterDefaultTranslations(v, trans)
			case "zh_hk":
				_ = zh_tw.RegisterDefaultTranslations(v, trans)
			default:
				_ = en2.RegisterDefaultTranslations(v, trans)
			}
		}
	}
}

func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.GetHeader("locale")
		trans, ok := uni.GetTranslator(locale)
		if ok {
			c.Set("trans", trans)
		}
		c.Next()
	}
}
