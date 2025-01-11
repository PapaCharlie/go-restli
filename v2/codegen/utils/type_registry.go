package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/dave/jennifer/jen"
)

type ComplexType interface {
	GetIdentifier() Identifier
	GetSourceFile() string
	ReferencedTypes() IdentifierSet
	ShouldReference() ShouldUsePointer
	GenerateCode() *jen.Statement
}

var TypeRegistry = &typeRegistry{
	types:        map[Identifier]*registeredType{},
	packageRoots: map[string]IdentifierSet{},
}

type registeredType struct {
	Type             ComplexType
	PackageRoot      string
	IsCyclic         bool
	TypeNameOverride string
	IsCustomTyperef  bool
}

type typeRegistry struct {
	types        map[Identifier]*registeredType
	packageRoots map[string]IdentifierSet
}

func (reg *typeRegistry) Register(t ComplexType, packageRoot string) error {
	id := t.GetIdentifier()
	if ct, ok := reg.types[id]; !ok {
		reg.types[id] = &registeredType{Type: t, PackageRoot: packageRoot}
		if _, ok = reg.packageRoots[packageRoot]; !ok {
			reg.packageRoots[packageRoot] = IdentifierSet{}
		}
		reg.packageRoots[packageRoot].Add(id)
		addImportName(id.PackagePath())
		return nil
	} else {
		return fmt.Errorf("go-restli: %q has already been registered in package %q", id, ct.PackageRoot)
	}
}

func (reg *typeRegistry) get(id Identifier) *registeredType {
	t, ok := reg.types[id]
	if !ok {
		log.Panicf("Unknown type: %q", id.FullName())
	}
	return t
}

func (reg *typeRegistry) Resolve(id Identifier) ComplexType {
	return reg.get(id).Type
}

func (reg *typeRegistry) PackageRoot(id Identifier) string {
	return reg.get(id).PackageRoot
}

func (reg *typeRegistry) IsCyclic(id Identifier) bool {
	return reg.get(id).IsCyclic
}

func (reg *typeRegistry) TypeNameOverride(id Identifier) string {
	if t, ok := reg.types[id]; ok {
		return t.TypeNameOverride
	} else {
		return ""
	}
}

func (reg *typeRegistry) SetCustomTyperef(id Identifier) {
	reg.get(id).IsCustomTyperef = true
}

func (reg *typeRegistry) IsCustomTyperef(id Identifier) bool {
	return reg.get(id).IsCustomTyperef
}

func (reg *typeRegistry) findCycle(nextNode Identifier, path Path) []Identifier {
	if cycle := path.IntroducesCycle(nextNode); len(cycle) > 0 {
		return cycle
	}

	// We've already seen this node, but it didn't introduce a cycle. Don't descend into its children
	if path.SeenNode(nextNode) {
		return nil
	}

	newPath := path.Add(nextNode)
	for c := range reg.get(nextNode).Type.ReferencedTypes() {
		if !reg.IsCyclic(c) {
			if p := reg.findCycle(c, newPath); len(p) > 0 {
				return p
			}
		}
	}

	return nil
}

func (reg *typeRegistry) flagCyclic(id Identifier) {
	node := reg.get(id)
	if !node.IsCyclic {
		node.IsCyclic = true
		log.Printf("Flagging %q as cyclic", id.FullName())
	}
	for c := range node.Type.ReferencedTypes() {
		child := reg.get(c)
		if !child.IsCyclic && node.PackageRoot == child.PackageRoot {
			reg.flagCyclic(c)
		}
	}
}

func (reg *typeRegistry) Finalize() (err error) {
	err = reg.validateAllTypesSatisfied()
	if err != nil {
		return err
	}

	err = reg.flagCyclicDependencies()
	if err != nil {
		return err
	}

	err = reg.remediateConflictingNames()
	if err != nil {
		return err
	}

	return nil
}

func (reg *typeRegistry) validateAllTypesSatisfied() error {
	for id, rt := range reg.types {
		for dep := range rt.Type.ReferencedTypes() {
			if _, ok := reg.types[dep]; !ok {
				return fmt.Errorf("go-restli: %q depends on unknown type %q", id.FullName(), dep.FullName())
			}
		}
	}
	return nil
}

func (reg *typeRegistry) flagCyclicDependencies() error {
	for id := range reg.types {
		for {
			cycle := reg.findCycle(id, Path{})
			if len(cycle) > 0 {
				packageRoots := map[string]bool{}
				var identifiers []string
				for _, cyclicModel := range cycle {
					packageRoots[cyclicModel.PackageRoot()] = true
					identifiers = append(identifiers, cyclicModel.FullName())
				}

				path := strings.Join(identifiers, " -> ")
				if len(packageRoots) > 1 {
					return fmt.Errorf("go-restli: The following cyclic dependency between packages was detected but "+
						"cannot be remediated due to type definitions being in different manifests: %s", path)
				}
				log.Printf("Detected cyclic dependency: %s", path)

				for _, c := range cycle {
					reg.flagCyclic(c)
				}
			} else {
				break
			}
		}
	}

	return nil
}

func (reg *typeRegistry) remediateConflictingNames() error {
	for _, types := range reg.packageRoots {
		conflictingTypes := map[string]IdentifierSet{}
		for id := range types {
			if reg.IsCyclic(id) {
				name := strings.ToLower(id.Name)
				if _, ok := conflictingTypes[name]; !ok {
					conflictingTypes[name] = IdentifierSet{}
				}
				conflictingTypes[name].Add(id)
			}
		}

		for _, v := range conflictingTypes {
			groups := make(map[string]IdentifierSet)
			for id := range v {
				name := strings.ToLower(id.Name)
				group, ok := groups[name]
				if !ok {
					group = IdentifierSet{}
					groups[name] = group
				}
				group.Add(id)
			}
			for _, group := range groups {
				err := reg.resolveConflicts(group)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (reg *typeRegistry) resolveConflicts(types IdentifierSet) error {
	if len(types) == 1 {
		return nil
	}

	log.Printf("WARNING: The following types have conflicting names: %s", types)

	var maxAttempts int
	namespaces := make(map[Identifier][]string, len(types))
	for id := range types {
		ns := strings.Split(id.Namespace, ".")
		for i, s := range ns {
			ns[i] = ExportedIdentifier(s)
		}
		namespaces[id] = ns
		if l := len(ns); l > maxAttempts {
			maxAttempts = l
		}
	}
	getOverriddenName := func(id Identifier, attempt int) string {
		ns := namespaces[id]
		log.Println("Before", id, attempt, ns)
		if len(ns) > attempt {
			ns = ns[len(ns)-attempt:]
		}
		log.Println("After", id, attempt, ns)
		return strings.Join(ns, "") + ExportedIdentifier(id.Name)
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		overrides := map[string]Identifier{}
		success := true
		for id := range types {
			override := getOverriddenName(id, attempt)
			if conflict, ok := overrides[override]; ok {
				log.Printf("Could not rename conflicting type %q to \"%s.%s\" as it would "+
					"conflict with other renamed type %q", id.FullName(), id.Namespace, override,
					conflict.FullName())
				success = false
				break
			} else {
				overrides[override] = id
			}
		}

		if success {
			for override, id := range overrides {
				reg.get(id).TypeNameOverride = override
				log.Printf("Conflicting type %q successfully renamed to \"%s.%s\"",
					id.FullName(), id.Namespace, override)
			}
			return nil
		}
	}

	return fmt.Errorf("go-restli: Failed to rename types in import cycle with conflicting type names %s", types)
}

func (reg *typeRegistry) TypesInPackageRoot(packageRoot string) IdentifierSet {
	return reg.packageRoots[packageRoot]
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
	nextPkg := nextNode.PackagePath()
	for i := len(p) - 1; i >= 0; i-- {
		pIPkg := p[i].PackagePath()
		if pIPkg != nextPkg {
			inSameNamespace = false
			continue
		}
		if !inSameNamespace && pIPkg == nextPkg {
			return append(append(Path(nil), p[i:]...), nextNode)
		}
	}
	return nil
}

var PagingContextIdentifier = Identifier{
	Name:      "PagingContext",
	Namespace: "restlidata",
}

var RawRecordIdentifier = Identifier{
	Name:      "RawRecord",
	Namespace: "restlidata",
}

var EmptyRecordIdentifier = Identifier{
	Name:      "EmptyRecord",
	Namespace: "com.linkedin.restli.common",
}
