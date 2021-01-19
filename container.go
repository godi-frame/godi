package godi

import (
	"github.com/godi-frame/godi/internal"
	"github.com/pkg/errors"
	"reflect"
	"sync"
)

// save package
type container struct {
	// map[string]*packet
	// key is packName (result of godi.NameOf)
	packets sync.Map
}

// check valid constructor
// get or create a package if it does not exist and provide a new provider for packet
func (c *container) register(constructor interface{}, scope provideScope, tag string, filters map[string][]string) error {
	if !IsGodiConstructor(constructor) {
		return errors.Wrapf(InvalidConstructor, "'%s'", reflect.TypeOf(constructor).String())
	}
	prototype, _ := internal.GetOutputFieldType(constructor, 0)
	prototypeName, _ := internal.GetGodiPacketName(prototype)
	pack, ok := c.packets.Load(prototypeName)
	if !ok {
		c.packets.Store(prototypeName, newPacket(c, prototype))
	}
	err := pack.(*packet).provide(tag, newProvider(constructor, scope, filters))
	if err != nil {
		return err
	}
	return nil
}

// remove the matched package of tags is nil
func (c *container) unregister(packetName string, tags []string) {
	if tags == nil {
		c.packets.Delete(packetName)
	} else {
		pack, ok := c.packets.Load(packetName)
		if ok {
			pack.(*packet).remove(tags)
		}
	}
}

// check valid retrieve argument
// get package matched and retrieve an instance of the packet
func (c *container) instanceOf(value reflect.Value, filters []string) error {
	t, err := internal.TypeOfPtr(value)
	if err != nil {
		return err
	}
	packetType := t.Elem()
	name, ok := internal.GetGodiPacketName(packetType)
	if !ok {
		return errors.Wrapf(VariableIsPtr, "instance of '%s'", value.String())
	}
	packets, ok := c.packets.Load(name)
	if !ok {
		return errors.Wrapf(PacketNotFound, "'%s'", name)
	}

	if len(filters) == 0 {
		filters = []string{DefaultTag}
	}
	return packets.(*packet).retrieveInstance(value, filters)
}

// check valid retrieve argument
// get package matched and list instance of the packet
func (c *container) sliceOf(value reflect.Value, filters []string) error {
	t, err := internal.TypeOfPtr(value)
	if err != nil {
		return err
	}
	sliceType := t.Elem()
	if sliceType.Kind() != reflect.Slice {
		return errors.Wrapf(VariableIsNotSlice, "slice of %s", sliceType.String())
	}
	packetType := sliceType.Elem()
	name, ok := internal.GetGodiPacketName(packetType)
	if !ok {
		return errors.Wrapf(ValueIsNotPtrOrInterface, "slice elements of '%s'", sliceType.String())
	}
	cpn, ok := c.packets.Load(name)
	if !ok {
		return errors.Wrapf(PacketNotFound, "'%s'", name)
	}
	return cpn.(*packet).retrieveSlice(value, filters)
}

// check valid retrieve argument
// get package matched and map instance of the packet
func (c *container) mapOf(value reflect.Value, filters []string) error {
	t, err := internal.TypeOfPtr(value)
	if err != nil {
		return err
	}
	mapType := t.Elem()
	if mapType.Kind() != reflect.Map {
		return errors.Wrapf(VariableIsNotMap, "map of '%s'", mapType.String())
	}
	packetType := mapType.Elem()
	name, ok := internal.GetGodiPacketName(packetType)
	if !ok {
		return errors.Wrapf(ValueIsNotPtrOrInterface, "map elements of '%s'", mapType.String())
	}
	cpn, ok := c.packets.Load(name)
	if !ok {
		return errors.Wrapf(PacketNotFound, "'%s'", name)
	}
	return cpn.(*packet).retrieveMap(value, filters)
}
