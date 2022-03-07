package utils

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
)

var namespaceEscape = regexp.MustCompile("([/.])_?internal([/.]?)")

type Identifier struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (i Identifier) String() string {
	return i.Namespace + "." + i.Name
}

func (i *Identifier) GetIdentifier() Identifier {
	return *i
}

func (i Identifier) pkg() (string, *Identifier) {
	if i.Namespace == "" {
		Logger.Panicf("%+v has no namespace!", i)
	}
	if i.Name == "" {
		Logger.Panicf("%+v has no name!", i)
	}
	if ex := TypeRegistry.ExternalImplementation(i); ex != nil {
		return ex.Namespace, ex
	} else {
		var p string
		if TypeRegistry.IsCyclic(i) {
			p = "conflictResolution"
		} else {
			p = i.Namespace
		}
		return FqcpToPackagePath(p), nil
	}
}

func (i Identifier) PackagePath() string {
	pkg, _ := i.pkg()
	return pkg
}

func (i Identifier) Qual() *jen.Statement {
	pkg, ex := i.pkg()
	var name string
	if ex != nil {
		name = ex.Name
	} else {
		name = i.Name
	}
	return jen.Qual(pkg, name)
}

func (i Identifier) UnmarshalerFuncName() string {
	return UnmarshalRestLi + i.Name
}

func (i Identifier) UnmarshalerFunc() *jen.Statement {
	return jen.Qual(i.PackagePath(), i.UnmarshalerFuncName())
}

func (i *Identifier) Receiver() string {
	return ReceiverName(i.Name)
}

func (i *Identifier) Resolve() ComplexType {
	return TypeRegistry.Resolve(*i)
}

func FqcpToPackagePath(fqcp string) string {
	fqcp = strings.Replace(namespaceEscape.ReplaceAllString(fqcp, "${1}_internal${2}"), ".", "/", -1)

	if PackagePrefix != "" {
		fqcp = filepath.Join(PackagePrefix, fqcp)
	}

	return fqcp
}

type IdentifierSet map[Identifier]bool

func NewIdentifierSet(ids ...Identifier) IdentifierSet {
	set := make(IdentifierSet)
	for _, id := range ids {
		set.Add(id)
	}
	return set
}

func (set IdentifierSet) Add(id Identifier) {
	set[id] = true
}

func (set IdentifierSet) AddAll(other IdentifierSet) {
	for id := range other {
		set[id] = true
	}
}

func (set IdentifierSet) Remove(id Identifier) {
	delete(set, id)
}

func (set IdentifierSet) Get(id Identifier) bool {
	return set[id]
}

func (set IdentifierSet) String() string {
	var classes []string
	for s := range set {
		classes = append(classes, s.String())
	}
	sort.Strings(classes)
	return "{" + strings.Join(classes, ", ") + "}"
}
