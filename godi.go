package godi

import "reflect"

var application *container

// 	register constructor to the container with options.
// 	constructor must be function type and must haven:
// 		parameters is result of another constructor registered with the container.
// 		parameters can be map or slice for dependencies.
// 		result is one instance and error.
// example:
//	godi.Register(func (p map[string]*InjectType)(*Instance,error){
//		return &Instance{}, nil
//	})
// 	The above code will be injected map of *InjectType to p with keys is tags registered
// and the constructor of Instance will be registered.
// 	Can be check constructor is valid with godi.IsGodiConstructor()
func Register(constructor interface{}, opts ...Option) error {
	options := newRegisterOption()
	for _, opt := range opts {
		opt(&options)
	}
	return application.register(constructor, options.scope, options.tag, options.filters)
}

// remove registered constructors to the container
// the packetName parameter: is result godi.NameOf(instance)
// if 'tags' value is nil remove all
// and opposite remove tags in 'tags'
// example: godi.Unregister(godi.NameOf(SomeValue))
func Unregister(packetName string, tags ...string) {
	application.unregister(packetName, tags)
}

// retrieve an instance of the registered constructor
// with the first tag in filters matched with the registered tag
// example:
//  var instance *InstanceType
//	godi.InstanceOf(&instance)
func InstanceOf(value interface{}, filters ...string) error {
	return application.instanceOf(reflect.ValueOf(value), filters)
}

// retrieve a slice of an instance of the registered constructor
// with the tags in filters  matched with the registered tag
// if the 'tags' value is nil retrieve all
func SliceOf(value interface{}, filters ...string) error {
	return application.sliceOf(reflect.ValueOf(value), filters)
}

// retrieve a map of an instance of the registered constructor
// with the tags in filters matched with the registered tag
// if the 'tags' value is nil retrieve all
func MapOf(value interface{}, filters ...string) error {
	return application.mapOf(reflect.ValueOf(value), filters)
}

func init() {
	application = &container{}
}
