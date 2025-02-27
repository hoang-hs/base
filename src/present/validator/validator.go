package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/hoang-hs/base/src/common/log"
	"regexp"
	"time"
)

func NewValidator() *validator.Validate {
	return validator.New()
}

func registerValidation(validator *validator.Validate, tag string, fn validator.Func) {
	if err := validator.RegisterValidation(tag, fn); err != nil {
		log.GetLogger().GetZap().Fatalf("Register custom validation %s failed with error: %s", tag, err.Error())
	}
	return
}
func registerStructValidation(validator *validator.Validate, fn validator.StructLevelFunc, in ...interface{}) {
	validator.RegisterStructValidation(fn, in...)
}
func RegisterValidations(validator *validator.Validate) {
	registerValidation(validator, "valid-by-regex-pattern", validByRegexPattern)
	registerValidation(validator, "valid-today", validToday)
}

var validByRegexPattern validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	regexPattern := fl.Param()
	match, err := regexp.MatchString(regexPattern, value)
	if err != nil {
		return false
	}
	return match
}

var validToday validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	date := value.Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	return date == today
}
