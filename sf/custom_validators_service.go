package sf

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"strings"
)

type CustomValidatorsService interface {
	RegisterCustomValidators()
}

type customValidatorsService struct{}

func NewCustomValidatorsService() CustomValidatorsService {
	return &customValidatorsService{}
}

func (*customValidatorsService) notBlank(field validator.FieldLevel) bool {
	value := field.Field().String()
	if len(value) == 0 {
		return true
	}
	trimmedValue := strings.TrimSpace(value)
	return len(trimmedValue) > 0
}

func (props *customValidatorsService) RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("notBlank", props.notBlank)
		if err != nil {
			panic(err)
		}
	} else {
		panic("Something went wrong while registering custom validator")
	}
}
