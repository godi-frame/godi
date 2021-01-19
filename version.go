package godi

import "fmt"

const (
	_MAJOR_         = 1
	_MINOR_         = 0
	_PATCH_         = 0
	_STATE_         = _ALPHA_
	_STATE_VERSION_ = 1
)

// return version of godi
func Version() (version string) {
	threePart := fmt.Sprintf("%d.%d.%d", _MAJOR_, _MINOR_, _PATCH_)
	statePart := fmt.Sprintf("%s.%d", _STATE_, _STATE_VERSION_)
	if _STATE_ != _RELEASE_ {
		version = threePart + "-" + statePart
	} else {
		version = threePart
	}
	return "godi/" + version
}
