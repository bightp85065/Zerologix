package http

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var _ binding.StructValidator = &CustomValidator{}

type CustomValidator struct {
	Once     sync.Once
	Validate *validator.Validate
}

type sliceValidateError []error

func (err sliceValidateError) Error() string {
	var errMsgs []string
	for i, e := range err {
		if e == nil {
			continue
		}
		errMsgs = append(errMsgs, fmt.Sprintf("[%d]: %s", i, e.Error()))
	}
	return strings.Join(errMsgs, "\n")
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{}
}

func (v *CustomValidator) ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		v.lazyinit()
		return v.Validate.Struct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(sliceValidateError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

func (v *CustomValidator) Engine() interface{} {
	v.lazyinit()
	return v.Validate
}

func (v *CustomValidator) lazyinit() {
	v.Once.Do(func() {
		v.Validate = validator.New()
		v.Validate.SetTagName("binding")
	})
}
