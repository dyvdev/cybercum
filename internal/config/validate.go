package config

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

const (
	mapstructTag    = "mapstructure"
	validateTagName = "validate"
	optValue        = "optional"
)

// Validator реализуется внутри структур конфигураций сервисов.
type Validator interface {
	// Validate проверяет структуру конфигурации.
	Validate() error
}

// Validate рекурсивно вызывает методы Validate у структуры
// конфига и её составных частей.
func Validate(src interface{}) error {
	if src == nil {
		return nil
	}

	return traverse(reflect.ValueOf(src), true)
}

func traverse(v reflect.Value, parent bool) (err error) {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if !v.IsNil() && v.CanInterface() {
			if err := tryValidateInterface(v.Interface()); err != nil {
				return err
			}
			if err := traverse(v.Elem(), true); err != nil {
				return err
			}
		}
	case reflect.Struct:
		if !parent && v.CanInterface() {
			if err := tryValidateInterface(v.Interface()); err != nil {
				return err
			}
		}
		for j := 0; j < v.NumField(); j++ {
			optTag := v.Type().Field(j).Tag.Get(validateTagName)
			if optTag == optValue && v.Field(j).IsNil() {
				continue
			}

			err := traverse(v.Field(j), false)
			if err != nil {
				tagValue := v.Type().Field(j).Tag.Get(mapstructTag)
				return errors.New(fmt.Sprintf("invalid section '%s': %v", tagValue, err))
			}

			// вызываем Validate() у детей
			child := v.Field(j)
			if child.CanAddr() {
				if child.Addr().MethodByName("Validate").Kind() != reflect.Invalid {
					child.Addr().MethodByName("Validate").Call([]reflect.Value{})
				}
			}
		}
	default:
		if v.CanInterface() {
			return tryValidateInterface(v.Interface())
		}
	}
	return nil
}

func tryValidateInterface(v interface{}) (err error) {
	pr, ok := v.(Validator)
	if ok {
		err = pr.Validate()
	}
	return
}
