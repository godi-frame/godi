package godi

import (
	"github.com/godi-frame/godi/internal"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

// save info of constructor registered and singleton instance
type provider struct {
	// constructor registered
	constructor interface{}

	// key is packName string is filter registered with godi.Filter
	filters map[string][]string

	// types of dependencies
	// initialize when the first call of getInstance
	inTypes []reflect.Type

	// scope of provider
	scope provideScope

	// instance of singleton
	// initialize when the first call of getInstance
	instance reflect.Value

	// make sure have only one instance
	once sync.Once
}

// provider constructor
func newProvider(constructor interface{}, scope provideScope, filters map[string][]string) *provider {
	return &provider{
		constructor: constructor,
		filters:     filters,
		scope:       scope,
	}
}

// return result of constructor
// if scope is singleton
//		initialize if is the first time call
// 		return instance
func (p *provider) getInstance(ctn *container) (reflect.Value, error) {
	if p.scope == Singleton {
		var err error
		p.once.Do(func() {
			p.instance, err = p.call(ctn)
		})
		return p.instance, err
	} else {
		return p.call(ctn)
	}
}

// find dependencies in the container and initialize they
func (p *provider) inject(ctn *container) ([]reflect.Value, error) {
	if p.inTypes == nil {
		p.inTypes = internal.InputOf(p.constructor)
	}
	typesLen := len(p.inTypes)
	inParam := make([]reflect.Value, typesLen)
	var err error
	for i, inType := range p.inTypes {
		v := reflect.New(inType)
		switch inType.Kind() {
		case reflect.Slice:
			err = ctn.sliceOf(v, p.filters[inType.Elem().String()])
		case reflect.Map:
			err = ctn.mapOf(v, p.filters[inType.Elem().String()])
		case reflect.Interface:
			err = ctn.instanceOf(v, p.filters[inType.String()])
		case reflect.Ptr:
			err = ctn.instanceOf(v, p.filters[inType.Elem().String()])
		default:
			err = InvalidDependencies
		}
		if err != nil {
			return inParam, errors.Wrapf(err, "inject %s failed", v.Type().String())
		}
		inParam[i] = v.Elem()
	}
	return inParam, nil
}

// call the constructor and inject dependencies
func (p *provider) call(ctn *container) (reflect.Value, error) {
	injectParam, err := p.inject(ctn)
	if err != nil {
		return reflect.Value{}, err
	}
	constructorResults := reflect.ValueOf(p.constructor).Call(injectParam)
	result := constructorResults[0]
	err, ok := constructorResults[1].Interface().(error)
	if !ok {
		err = nil
	}
	if err != nil {
		comp, _ := internal.GetOutputFieldType(p.constructor, 0)
		return result, errors.Wrapf(err, "initialize %s failed", comp.String())
	}
	return result, nil
}
