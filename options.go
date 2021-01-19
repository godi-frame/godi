package godi

// define Option type
type Option func(option *registerOption)

type registerOption struct {
	scope provideScope
	//types   provideType
	tag     string
	filters map[string][]string
}

func newRegisterOption() registerOption {
	return registerOption{
		scope: Singleton,
		//types:   Normal,
		tag:     DefaultTag,
		filters: nil,
	}
}

// option setup scope
func Scope(scope provideScope) Option {
	return func(option *registerOption) {
		option.scope = scope
	}
}

// option setup tag
func Tag(tag string) Option {
	return func(option *registerOption) {
		option.tag = tag
	}
}

// option setup filter for dependency
// packName is result of godi.NameOf
func Filter(packName string, tags ...string) Option {
	return func(option *registerOption) {
		if option.filters == nil {
			option.filters = make(map[string][]string)
		}
		option.filters[packName] = tags
	}
}
