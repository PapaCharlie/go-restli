package utils

import (
	"log"
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

func (i Identifier) FullName() string {
	return i.Namespace + "." + i.Name
}

func (i Identifier) GetIdentifier() Identifier {
	return i
}

func (i Identifier) PackagePath() string {
	if i.Namespace == "" {
		log.Panicf("%+v has no namespace!", i)
	}
	var p string
	if TypeRegistry.IsCyclic(i) {
		p = "conflictResolution"
	} else {
		p = i.Namespace
	}
	return FqcpToPackagePath(i.PackageRoot(), p)
}

func (i Identifier) TypeName() string {
	override := TypeRegistry.TypeNameOverride(i)
	if override != "" {
		return override
	} else {
		return i.Name
	}
}

func (i Identifier) IsCustomTyperef() bool {
	return TypeRegistry.IsCustomTyperef(i)
}

func (i Identifier) IsEmptyRecord() bool {
	return i == EmptyRecordIdentifier
}

func (i Identifier) Qual() *jen.Statement {
	return jen.Qual(i.PackagePath(), i.TypeName())
}

func (i Identifier) Receiver() string {
	return ReceiverName(i.TypeName())
}

func (i Identifier) Resolve() ComplexType {
	return TypeRegistry.Resolve(i)
}

func (i Identifier) PackageRoot() string {
	return TypeRegistry.PackageRoot(i)
}

func FqcpToPackagePath(packageRoot string, fqcp string) string {
	fqcp = strings.Replace(namespaceEscape.ReplaceAllString(fqcp, "${1}_internal${2}"), ".", "/", -1)

	if packageRoot != "" {
		fqcp = filepath.Join(packageRoot, fqcp)
	}

	return fqcp
}

type IdentifierSet map[Identifier]bool

func NewIdentifierSet(ids ...Identifier) (set IdentifierSet) {
	set = IdentifierSet{}
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
		set.Add(id)
	}
}

func (set IdentifierSet) Remove(id Identifier) {
	delete(set, id)
}

func (set IdentifierSet) String() string {
	var classes []string
	set.Range(func(id Identifier) {
		classes = append(classes, id.FullName())
	})
	return "{" + strings.Join(classes, ", ") + "}"
}

func (set IdentifierSet) Range(f func(id Identifier)) {
	var ids []Identifier
	for id := range set {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].FullName() < ids[j].FullName()
	})
	for _, id := range ids {
		f(id)
	}
}
