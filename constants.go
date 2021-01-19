package godi

const (
	// the default value for Register without Tag()
	DefaultTag = "default"
)

type provideScope uint8

const (
	// only one instance initialized for one tag registered
	// but the same type can have many instances as long as it tags difference
	Singleton provideScope = 0

	// one instance initialized for one retrieve call
	Prototype provideScope = 1
)

const (
	_ALPHA_             = "a"
	_BETA_              = "b"
	_RELEASE_CANDIDATE_ = "rc"
	_RELEASE_           = ""
)
