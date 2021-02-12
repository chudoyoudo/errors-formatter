package errors_formatter

import (
    "unicode"

    enLocale "github.com/go-playground/locales/en"
    translator "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
)

func FormatErrors(err error) map[string][]string {
    result := map[string][]string{}
    en := enLocale.New()

    trans, found := translator.New(en, en).GetTranslator("en")
    if !found {
        log.Error(errors.Wrap(errors.New("Translator not found"), "Get translator for errors formatter"))
        result["system"] = append(result["system"], "Error translator not found")
        return result
    }

    switch err.(type) {
    case validator.ValidationErrors:
        for _, e := range err.(validator.ValidationErrors) {
            field := e.StructField()
            fieldName := correctFieldName(field)
            if _, keyExist := result[fieldName]; !keyExist {
                result[fieldName] = []string{}
            }
            result[fieldName] = append(result[fieldName], e.Translate(trans))
        }

    default:
        log.Error(errors.Wrap(err, "Unknown error"), err)
        result["system"] = append(result["system"], "Unknown error type")
    }

    return result
}

func correctFieldName(field string) string {
    a := []rune(field)
    a[0] = unicode.ToLower(a[0])
    return string(a)
}
