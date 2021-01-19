package godi

import (
	"github.com/godi-frame/godi/internal"
	"reflect"
)

// the result is the string of the package's name + type's name
// the value must be pointer or interface, it will panic if the value is another type
// example:
//		godi.NameOf(&gin.Context{}) => github.com/gin-gonic/gin/gin.Context
func NameOf(value interface{}) string {
	name, ok := internal.GetGodiPacketName(reflect.TypeOf(value))
	if !ok {
		panic("invalid packet name, value must be struct pointer or interface")
	}
	return name
}

// check constructor can register
// return true if constructor valid
// constructor valid:
//  - is function
// 	- parameters:
//		+ is struct pointer or interface,
//		+ is map[string] of struct pointer or interface,
//      + is slice[string] of struct pointer or interface,
//  - result is one instance and error
func IsGodiConstructor(constructor interface{}) bool {
	t := reflect.TypeOf(constructor)
	if t.Kind() != reflect.Func {
		return false
	}

	numIn := t.NumIn()
	numOut := t.NumOut()

	for i := 0; i < numIn; i++ {
		typeIn := t.In(i)
		typeInValid := false
		switch typeIn.Kind() {
		case reflect.Slice:
			typeInValid = internal.IsStructPtrOrInterface(typeIn.Elem())

		case reflect.Map:
			typeInValid = internal.IsStructPtrOrInterface(typeIn.Elem())

		default:
			typeInValid = internal.IsStructPtrOrInterface(typeIn)
		}
		if !typeInValid {
			return false
		}
	}

	if numOut != 2 {
		return false
	}

	tOut := t.Out(0)

	typeOutValid := false
	switch tOut.Kind() {
	default:
		typeOutValid = internal.IsStructPtrOrInterface(tOut)
	}
	if !typeOutValid {
		return false
	}
	return internal.IsError(t.Out(1))
}
