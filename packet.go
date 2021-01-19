package godi

import (
	"github.com/godi-frame/godi/internal"
	"github.com/pkg/errors"
	"reflect"
	"sync"
)

type packet struct {
	ctn *container

	// type of packet
	prototype reflect.Type

	// map[string]*provider
	// key is tag
	providers sync.Map
}

// packet constructor
func newPacket(ctn *container, iType reflect.Type) *packet {
	return &packet{
		ctn:       ctn,
		prototype: iType,
	}
}

// find the first tag matched and get an instance
func (p *packet) retrieveInstance(value reflect.Value, tags []string) error {
	prv, err := p.firstProvider(tags)
	if err != nil {
		return err
	}
	instance, err := prv.getInstance(p.ctn)
	if err != nil {
		return err
	}
	value.Elem().Set(instance)
	return nil
}

// find tags matched and get a slice of instances
func (p *packet) retrieveSlice(value reflect.Value, tags []string) error {
	matchedInstances := reflect.MakeSlice(value.Type().Elem(), 0, 0)
	for _, prv := range p.getProviders(tags) {
		instance, err := prv.getInstance(p.ctn)
		if err != nil {
			return err
		}
		matchedInstances = reflect.Append(matchedInstances, instance)
	}
	value.Elem().Set(matchedInstances)
	return nil
}

//find tags matched and get a map of instances with keys is tags
func (p *packet) retrieveMap(value reflect.Value, tags []string) error {
	matchedInstances := reflect.MakeMap(value.Type().Elem())
	for tag, prv := range p.getProviders(tags) {
		instance, err := prv.getInstance(p.ctn)
		if err != nil {
			return err
		}
		matchedInstances.SetMapIndex(reflect.ValueOf(tag), instance)
	}
	value.Elem().Set(matchedInstances)
	return nil
}

// store provider with key is tag
func (p *packet) provide(tag string, src *provider) error {
	if _, ok := p.providers.Load(tag); ok {
		return errors.Wrapf(DuplicatedTag, "'%s'", tag)
	}
	p.providers.Store(tag, src)
	return nil
}

// remove key in tags
func (p *packet) remove(tags []string) {
	for _, tag := range tags {
		p.providers.Delete(tag)
	}
}

// get providers match with tags
func (p *packet) getProviders(tags []string) map[string]*provider {
	result := make(map[string]*provider)
	getAll := tags == nil
	p.providers.Range(func(key, value interface{}) bool {
		keyStr := key.(string)
		if getAll || internal.KeyInStrings(keyStr, tags) {
			prv := value.(*provider)
			result[keyStr] = prv
		}
		return true
	})
	return result
}

// get provider matched with tag
func (p *packet) getProvider(tag string) (*provider, error) {
	prv, ok := p.providers.Load(tag)
	if !ok {
		return nil, ProviderNotFound
	}
	return prv.(*provider), nil
}

// get the first provider matched
func (p *packet) firstProvider(tags []string) (*provider, error) {
	for _, tag := range tags {
		prv, err := p.getProvider(tag)
		if err == nil {
			return prv, nil
		}
	}
	return nil, ProviderNotFound
}
