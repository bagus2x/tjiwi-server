package app

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

func ValidateAndTranslate(validate *validator.Validate, err error) error {
	if err == nil {
		return nil
	}

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	enTranslations.RegisterDefaultTranslations(validate, trans)
	msg := TranslateValidationErrors(err, trans)

	err = &Error{
		Code:     EBadRequest,
		Messages: msg,
	}

	return err
}

func TranslateValidationErrors(err error, trans ut.Translator) []string {
	validatorErrs := err.(validator.ValidationErrors)
	msg := make([]string, 0)

	for _, e := range validatorErrs {
		msg = append(msg, e.Translate(trans))
	}

	return msg
}
