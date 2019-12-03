package internal

import (
	"log"
	"sort"
	"strings"
)

type IdentifierSet map[Identifier]bool

func (set *IdentifierSet) Add(id Identifier) {
	(*set)[id] = true
}

func (set *IdentifierSet) AddAll(other IdentifierSet) {
	for id := range other {
		(*set)[id] = true
	}
}

func (set *IdentifierSet) Remove(id Identifier) {
	delete(*set, id)
}

func (set *IdentifierSet) Get(id Identifier) bool {
	return (*set)[id]
}

func (set IdentifierSet) String() string {
	var classes []string
	for s := range set {
		classes = append(classes, s.GetQualifiedClasspath())
	}
	sort.Strings(classes)
	return "{" + strings.Join(classes, ", ") + "}"
}

type GraphNode struct {
	Identifier Identifier
	Parents    IdentifierSet
	Children   IdentifierSet
	IsCyclic   bool
}

type Graph map[Identifier]*GraphNode

var DependencyGraph = make(Graph)

func (g *Graph) getOrCreate(id Identifier) *GraphNode {
	if _, ok := (*g)[id]; !ok {
		(*g)[id] = &GraphNode{
			Identifier: id,
			Parents:    make(IdentifierSet),
		}
	}
	return (*g)[id]
}

func (g *Graph) AddParent(id Identifier, parent Identifier) {
	node := g.getOrCreate(id)
	node.Parents.Add(parent)
}

func (g *Graph) SetChildren(id Identifier, set IdentifierSet) {
	node := g.getOrCreate(id)
	node.Children = set
}

func (g *Graph) AllDependencies(id Identifier, set IdentifierSet) IdentifierSet {
	if set == nil {
		set = make(IdentifierSet)
	}
	set[id] = true
	for child := range g.getOrCreate(id).Children {
		if !set[child] {
			g.AllDependencies(child, set)
		}
	}
	return set
}

func (g *Graph) FindCycle(nextNode Identifier, depth int, path Path) []Identifier {
	if cycle := path.IntroducesCycle(nextNode, depth); len(cycle) > 0 {
		return cycle
	}

	// We've already seen this node, but it didn't introduce a cycle. Don't descend into its children
	if _, ok := path.VisitedNodes[nextNode]; ok {
		return nil
	}

	newPath := path.CopyWith(nextNode, depth)
	for c := range (*g)[nextNode].Children {
		if !g.IsCyclic(c) {
			if p := g.FindCycle(c, depth+1, newPath); len(p) > 0 {
				return p
			}
		}
	}
	return nil
}

func (g *Graph) IsCyclic(id Identifier) bool {
	gn, ok := (*g)[id]
	return ok && gn.IsCyclic
}

func (g *Graph) FlagCyclic(id Identifier) {
	node := g.getOrCreate(id)
	node.IsCyclic = true
	for c := range node.Children {
		if !g.IsCyclic(c) {
			g.FlagCyclic(c)
		}
	}
}

func buildDependencyGraph() {
	DependencyGraph = make(Graph)
	for id, t := range ModelRegistry.resolvedTypes {
		children := (&Model{ComplexType: t.Type}).flattenInnerModels()
		DependencyGraph.SetChildren(id, children)
		for child := range children {
			DependencyGraph.AddParent(child, id)
		}
	}
}

func flagCyclicDependencies() {
	for id := range DependencyGraph {
		for {
			cycle := DependencyGraph.FindCycle(id, 0, Path{})
			if len(cycle) > 0 {
				var identifiers []string
				for _, cyclicModel := range cycle {
					identifiers = append(identifiers, cyclicModel.String())
				}
				log.Println("Detected cyclic dependency:", strings.Join(identifiers, " -> "))

				for _, c := range cycle {
					DependencyGraph.FlagCyclic(c)
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
