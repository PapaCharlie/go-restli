package restlicodec

import (
	"strings"
)

const WildCard = "*"

type PathSpec map[string]PathSpec

var NoExcludedFields PathSpec

func NewPathSpec(directives ...string) (p PathSpec) {
	p = make(PathSpec)
	for _, s := range directives {
		s = strings.TrimPrefix(s, "/")
		inner := p
		for _, segment := range strings.Split(s, "/") {
			if _, ok := inner[segment]; !ok {
				inner[segment] = make(PathSpec)
			}
			inner = inner[segment]
		}
	}
	return p
}

func (p PathSpec) Matches(path []string) bool {
	if len(p) == 0 {
		return false
	}
	if path[0] == "$set" {
		if len(path) == 1 {
			return false
		} else {
			path = path[1:]
		}
	}
	matches := func(s string) bool {
		spec, ok := p[s]
		switch {
		case !ok:
			return false
		case len(spec) == 0:
			return true
		default:
			if len(path) == 1 {
				return false
			} else {
				return spec.Matches(path[1:])
			}
		}
	}
	return matches(WildCard) || matches(path[0])
}
