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
	return genericMatches(p, path, func(s string) string { return s })
}

// func (p PathSpec) PrefixAll(prefix string) PathSpec {
// 	spec := NewPathSpec(prefix)
// 	child := spec
// 	for {
// 		for k, v := range child {
// 			if len(v) == 0 {
// 				child[k] = p
// 				return spec
// 			} else {
// 				child = v
// 				break
// 			}
// 		}
// 	}
// }

func genericMatches[T any](p PathSpec, path []T, extract func(t T) string) bool {
	if len(p) == 0 {
		return false
	}
	p0 := extract(path[0])
	if p0 == "$set" || p0 == "$delete" {
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
				return genericMatches(spec, path[1:], extract)
			}
		}
	}
	return matches(WildCard) || matches(p0)
}
