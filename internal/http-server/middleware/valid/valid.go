package valid

import (
	"github.com/go-playground/validator/v10"
	"log"
	"time"
)

func CreateValidator() *validator.Validate {
	validate := validator.New()
	validName(validate)
	validSurname(validate)
	validData(validate)
	return validate
}

func validName(validate *validator.Validate) {
	vErr := validate.RegisterValidation("name", func(fl validator.FieldLevel) bool {
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

func validSurname(validate *validator.Validate) {
	vErr := validate.RegisterValidation("surname", func(fl validator.FieldLevel) bool {
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
		log.Fatal("register ValidSurname", vErr)
	}
}

func validData(validate *validator.Validate) {
	vErr := validate.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		text := fl.Field().String()
		layout := "20060102"
		_, err := time.Parse(layout, text)
		if err != nil {
			return false
		}
		return true
	})
	if vErr != nil {
		log.Fatal("register data", vErr)
	}
}

//func validCash(validate *validator.Validate){
//	vErr := validate.RegisterValidation("cash", func(fl validator.FieldLevel) bool {
//		cash := fl.Field().Float()
//	})
//}
