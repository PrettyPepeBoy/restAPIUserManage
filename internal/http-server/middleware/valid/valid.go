package valid

import (
	"github.com/go-playground/validator/v10"
	"log"
	"time"
)

func CreateValidator() *validator.Validate {
	validate := validator.New()
	validText(validate)
	validData(validate)
	return validate
}

func validText(validate *validator.Validate) {
	vErr := validate.RegisterValidation("correct_text", func(fl validator.FieldLevel) bool {
		text := fl.Field().String()
		if len(text) > 20 {
			return false
		}
		for _, a := range text {
			if (a < 'A' || a > 'Z') && (a < 'a' || a > 'z') {
				return false
			}
		}
		return true
	})
	if vErr != nil {
		log.Fatal("register ValidName", vErr)
	}
}

func validData(validate *validator.Validate) {
	vErr := validate.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		text := fl.Field().String()
		layout := "20060102"
		data, err := time.Parse(layout, text)
		if err != nil {
			return false
		}
		if data.Year() < 1900 {
			return false
		}
		return true
	})
	if vErr != nil {
		log.Fatal("register data", vErr)
	}
}
