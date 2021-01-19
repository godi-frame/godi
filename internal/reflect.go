package internal

import (
	"reflect"

	"github.com/pkg/errors"
)

var (
	IsNotFunction = errors.New("is not function")
	ValueIsNotPtr = errors.New("value is not pointer")
)

func InputOf(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	numIn := t.NumIn()
	result := make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		result[i] = t.In(i)
	}
	return result
}

func TypeOfPtr(v reflect.Value) (reflect.Type, error) {
	t := v.Type()
	if t.Kind() != reflect.Ptr {
		return nil, errors.Wrapf(ValueIsNotPtr, "'%s'", v.Type().Name())
	}
	return t, nil
}

func GetOutputFieldType(f interface{}, idx int) (reflect.Type, error) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return nil, IsNotFunction
	}
	numOut := t.NumOut()
	if idx < 0 || idx >= numOut {
		return nil, errors.New("out of range")
	}
	return t.Out(idx), nil
}

func IsStructPtrOrInterface(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Interface:
		return true
	case reflect.Ptr:
		return t.Elem().Kind() == reflect.Struct
	default:
		return false
	}
}

func IsError(t reflect.Type) bool {
	var err error
	return t.Implements(reflect.TypeOf(&err).Elem())
}

func GetGodiPacketName(t reflect.Type) (string, bool) {
	k := t.Kind()
	switch k {
	case reflect.Ptr:
		return t.Elem().PkgPath() + "/" + t.Elem().String(), true
	case reflect.Interface:
		return t.PkgPath() + "/" + t.String(), true
	default:
		return "", false
	}
}
