package godi

import "github.com/pkg/errors"

var (

	// retrieve not registered pack
	PacketNotFound = errors.New("packet not found")

	// SliceOf parameter is not slice pointer
	VariableIsNotSlice = errors.New("retrieve parameter must be slice")

	// InstanceOf parameter is map pointer
	VariableIsNotMap = errors.New("retrieve parameter must be map")

	// InstanceOf parameter is not struct pointer
	VariableIsPtr = errors.New("retrieve parameter must be pointer")

	// retrieve value is not a reference type
	ValueIsNotPtrOrInterface = errors.New("value must be pointer or interface")

	// IsGodiConstructor return false
	InvalidConstructor = errors.New("invalid godi constructor")

	// can't initialize instance dependencies
	InvalidDependencies = errors.New("invalid dependencies")

	// retrieve tag unregistered
	ProviderNotFound = errors.New("provider not found")

	// register duplicate tag
	DuplicatedTag = errors.New("duplicated tag")
)
