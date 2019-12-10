package codegen

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

type ComplexType interface {
	GetIdentifier() Identifier
	GetSourceFile() string
	InnerTypes() IdentifierSet
	GenerateCode() *jen.Statement
}

var TypeRegistry = make(typeRegistry)

type registeredType struct {
	Type     ComplexType
	IsCyclic bool
}

type typeRegistry map[Identifier]*registeredType

func (reg *typeRegistry) Register(t ComplexType) {
	if _, ok := (*reg)[t.GetIdentifier()]; ok {
		Logger.Panicf("Cannot register type %s twice!", t.GetIdentifier())
	}
	(*reg)[t.GetIdentifier()] = &registeredType{Type: t}
}

func (reg *typeRegistry) get(id Identifier) *registeredType {
	t, ok := (*reg)[id]
	if !ok {
		Logger.Panicf("Unknown type: %s", id)
	}
	return t
}

func (reg *typeRegistry) Resolve(id Identifier) ComplexType {
	return reg.get(id).Type
}

func (reg *typeRegistry) IsCyclic(id Identifier) bool {
	t := (*reg)[id]
	if t == nil {
		return false
	} else {
		return t.IsCyclic
	}
}

func (reg *typeRegistry) GenerateTypeCode() (files []*CodeFile) {
	for _, t := range *reg {
		files = append(files, &CodeFile{
			SourceFile:  t.Type.GetSourceFile(),
			PackagePath: t.Type.GetIdentifier().PackagePath(),
			Filename:    t.Type.GetIdentifier().Name,
			Code:        t.Type.GenerateCode(),
		})
	}
	return files
}

func (reg *typeRegistry) FindCycle(nextNode Identifier, depth int, path Path) []Identifier {
	if cycle := path.IntroducesCycle(nextNode, depth); len(cycle) > 0 {
		return cycle
	}

	// We've already seen this node, but it didn't introduce a cycle. Don't descend into its children
	if _, ok := path.VisitedNodes[nextNode]; ok {
		return nil
	}

	newPath := path.CopyWith(nextNode, depth)
	for c := range reg.get(nextNode).Type.InnerTypes() {
		if !reg.IsCyclic(c) {
			if p := reg.FindCycle(c, depth+1, newPath); len(p) > 0 {
				return p
			}
		}
	}
	return nil
}

func (reg *typeRegistry) FlagCyclic(id Identifier) {
	node := reg.get(id)
	node.IsCyclic = true
	for c := range node.Type.InnerTypes() {
		if !reg.IsCyclic(c) {
			reg.FlagCyclic(c)
		}
	}
}

func (reg *typeRegistry) FlagCyclicDependencies() {
	for id := range *reg {
		for {
			cycle := reg.FindCycle(id, 0, Path{})
			if len(cycle) > 0 {
				var identifiers []string
				for _, cyclicModel := range cycle {
					identifiers = append(identifiers, cyclicModel.String())
				}
				Logger.Println("Detected cyclic dependency:", strings.Join(identifiers, " -> "))

				for _, c := range cycle {
					reg.FlagCyclic(c)
				}
			} else {
				break
			}
		}
	}
}

type Path struct {
	VisitedNodes      map[Identifier]int
	VisitedNamespaces map[string]int
}

func (p *Path) CopyWith(id Identifier, depth int) Path {
	newPath := Path{
		VisitedNodes:      make(map[Identifier]int, len(p.VisitedNodes)+1),
		VisitedNamespaces: make(map[string]int),
	}

	for i, d := range p.VisitedNodes {
		newPath.VisitedNodes[i] = d
	}
	newPath.VisitedNodes[id] = depth

	for n, d := range p.VisitedNamespaces {
		newPath.VisitedNamespaces[n] = d
	}
	if _, ok := newPath.VisitedNamespaces[id.Namespace]; !ok {
		newPath.VisitedNamespaces[id.Namespace] = depth
	}

	return newPath
}

func (p *Path) IntroducesCycle(nextNode Identifier, depth int) []Identifier {
	if cycleStart, ok := p.VisitedNamespaces[nextNode.Namespace]; ok {
		for _, d := range p.VisitedNamespaces {
			if d > cycleStart {
				// This means we visited this namespace strictly after the original visit to nextNode's namespace. In
				// other words: nextNode is introducing a cycle
				seq := make([]Identifier, len(p.VisitedNodes)-cycleStart+1)
				for id, index := range p.VisitedNodes {
					if index >= cycleStart {
						seq[index-cycleStart] = id
					}
				}
				seq[len(seq)-1] = nextNode
				return seq
			}
		}
	}
	return nil
}
