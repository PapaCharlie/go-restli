package utils

import (
	"log"
	"strings"

	"github.com/dave/jennifer/jen"
)

type ComplexType interface {
	GetIdentifier() Identifier
	GetSourceFile() string
	InnerTypes() IdentifierSet
	ShouldReference() ShouldUsePointer
	GenerateCode() *jen.Statement
}

var TypeRegistry = make(typeRegistry)

type registeredType struct {
	Type               ComplexType
	IsCyclic           bool
	ExternalIdentifier *Identifier
}

func (t *registeredType) innerTypes() IdentifierSet {
	if t.ExternalIdentifier != nil {
		return IdentifierSet{}
	} else {
		return t.Type.InnerTypes()
	}
}

type typeRegistry map[Identifier]*registeredType

func (reg typeRegistry) Register(t ComplexType) {
	id := t.GetIdentifier()
	if r, ok := reg[id]; ok {
		if r.Type != nil {
			Logger.Panicf("Cannot register type %s twice!", id)
		}
		r.Type = t
	} else {
		reg[id] = &registeredType{Type: t}
	}
	// Automatically register standard structs
	if strings.HasPrefix(t.GetIdentifier().Namespace, StdTypesPackage) {
		reg.RegisterExternalImplementation(id, id)
	}
}

func (reg typeRegistry) RegisterExternalImplementation(internal, external Identifier) {
	log.Printf("Registering external implementation for %q: %q", internal, external)
	if r, ok := reg[internal]; ok {
		if r.ExternalIdentifier != nil {
			Logger.Panicf("Cannot register external implementation twice! (%q is implemented by %+v and %+v)",
				internal, r.ExternalIdentifier, external)
		} else {
			r.ExternalIdentifier = &external
		}
	} else {
		reg[internal] = &registeredType{ExternalIdentifier: &external}
	}
}

func (reg typeRegistry) ExternalImplementation(id Identifier) *Identifier {
	return reg.get(id).ExternalIdentifier
}

func (reg typeRegistry) get(id Identifier) *registeredType {
	t, ok := reg[id]
	if !ok {
		Logger.Panicf("Unknown type: %s", id)
	}
	return t
}

func (reg typeRegistry) Resolve(id Identifier) ComplexType {
	return reg.get(id).Type
}

func (reg typeRegistry) IsCyclic(id Identifier) bool {
	return reg.get(id).IsCyclic
}

func (reg typeRegistry) GenerateTypeCode() (files []*CodeFile) {
	for id, t := range reg {
		if t.ExternalIdentifier != nil {
			continue
		}
		files = append(files, &CodeFile{
			SourceFile:  t.Type.GetSourceFile(),
			PackagePath: id.PackagePath(),
			Filename:    id.Name,
			Code:        t.Type.GenerateCode(),
		})
	}
	return files
}

func (reg typeRegistry) FindCycle(nextNode Identifier, path Path) []Identifier {
	if cycle := path.IntroducesCycle(nextNode); len(cycle) > 0 {
		return cycle
	}

	// We've already seen this node, but it didn't introduce a cycle. Don't descend into its children
	if path.SeenNode(nextNode) {
		return nil
	}

	newPath := path.Add(nextNode)
	for c := range reg.get(nextNode).innerTypes() {
		if !reg.IsCyclic(c) {
			if p := reg.FindCycle(c, newPath); len(p) > 0 {
				return p
			}
		}
	}

	return nil
}

func (reg typeRegistry) FlagCyclic(id Identifier) {
	node := reg.get(id)
	node.IsCyclic = true
	for c := range node.innerTypes() {
		if !reg.IsCyclic(c) {
			reg.FlagCyclic(c)
		}
	}
}

func (reg typeRegistry) FlagCyclicDependencies() {
	for id := range reg {
		for {
			cycle := reg.FindCycle(id, Path{})
			if len(cycle) > 0 {
				var identifiers []string
				for _, cyclicModel := range cycle {
					identifiers = append(identifiers, cyclicModel.String())
				}
				log.Println("Detected cyclic dependency:", strings.Join(identifiers, " -> "))

				for _, c := range cycle {
					reg.FlagCyclic(c)
				}
			} else {
				break
			}
		}
	}
}

func (reg typeRegistry) FindAllDependents(id Identifier) IdentifierSet {
	dependents := NewIdentifierSet(id)
	added := true
	for added {
		added = false
		for id, c := range reg {
			for inner := range c.innerTypes() {
				if dependents.Get(inner) {
					if !dependents.Get(id) {
						dependents.Add(id)
						added = true
					}
				}
			}
		}
	}
	return dependents
}

type Path []Identifier

func (p Path) Add(id Identifier) Path {
	return append(append(Path(nil), p...), id)
}

func (p Path) SeenNode(id Identifier) bool {
	// This can probably be made a little faster by having a lookup map on the side, but these dependency chains likely
	// won't grow past the dozen or so elements, which means this should be fast enough not to be noticeable during code
	// generation
	for _, node := range p {
		if node == id {
			return true
		}
	}
	return false
}

func (p Path) IntroducesCycle(nextNode Identifier) Path {
	inSameNamespace := true
	for i := len(p) - 1; i >= 0; i-- {
		if p[i].Namespace != nextNode.Namespace {
			inSameNamespace = false
			continue
		}
		if !inSameNamespace && p[i].Namespace == nextNode.Namespace {
			return append(append(Path(nil), p[i:]...), nextNode)
		}
	}
	return nil
}
