package validation

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
)

const (
	EmailTag             = "email"
	PhoneTag             = "phone"
	e164PhoneRegexString = `^\+?[1-9]\d{1,14}$`
)

// GinValidator enables using go-playground's validator with Gin's struct validation.
type GinValidator struct {
	*validator.Validate
}

// Deprecated: Use GinValidator instead.
type DefaultValidator = GinValidator

var (
	enLocale            = en.New()
	universalTranslator = ut.New(enLocale, enLocale)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	Translator, _ = universalTranslator.GetTranslator("en")

	Validate *validator.Validate     = NewValidator()
	_        binding.StructValidator = &GinValidator{}
)

// ValidateStruct implements Gin's StructValidator interface.
func (v *GinValidator) ValidateStruct(obj interface{}) error {
	// If this GinValidator has no custom validator, use a default one.
	if v.Validate == nil {
		v.Validate = NewValidator()
	}

	if kindOfData(obj) == reflect.Struct {
		if err := v.Validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

// Creates a new validator.Validate with default behavior.
func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.SetTagName("binding")

	enTranslations.RegisterDefaultTranslations(validate, Translator)

	registerDefaultValidations(validate)
	registerDefaultTranslations(validate)

	return validate
}

func registerDefaultValidations(validate *validator.Validate) {
	validate.RegisterValidation(PhoneTag, IsE164Phone)
}

func registerDefaultTranslations(validate *validator.Validate) {
	validate.RegisterTranslation(EmailTag, Translator, func(ut ut.Translator) error {
		return nil
	}, translateFunc)

	validate.RegisterTranslation(PhoneTag, Translator, func(ut ut.Translator) error {
		return ut.Add(PhoneTag, "{0} number must be in E.164 format", true) // see universal-translator for details
	}, translateFunc)
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	field := fe.Field()

	if field == "" {
		field = "'" + fe.Tag() + "'"
	}

	t, err := ut.T(fe.Tag(), field)
	if err != nil {
		return fe.(error).Error()
	}

	return t
}
