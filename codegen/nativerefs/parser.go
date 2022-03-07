package nativerefs

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
)

func FindAllExternalImplementations(roots ...string) (err error) {
	for _, root := range roots {
		err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return err
			}

			return FindExternalImplementationsInDirectory(path)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func FindExternalImplementationsInDirectory(dir string) (err error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// ignore any "main" packages that can't be imported anyway
	delete(packages, "main")

	if len(packages) == 0 {
		return nil
	}

	var pkg *ast.Package
	for _, pkg = range packages {
		break
	}

	refs, err := locateAllNativeRefs(pkg)
	if err != nil {
		return err
	}

	for _, n := range refs {
		err = allNativeRefMethodsDefined(fset, n, dir, pkg)
		if err != nil {
			return err
		}
		types.RegisterNativeTyperef(n)
	}

	manifests, err := locateAllGeneratedTypes(pkg)
	if err != nil {
		return err
	}
	log.Println(len(manifests))

	for _, m := range manifests {
		utils.TypeRegistry.RegisterExternalImplementation(m.Identifier, utils.Identifier{
			Name:      m.Name,
			Namespace: m.Package,
		})
	}

	return nil
}

func locateAllNativeRefs(pkg *ast.Package) (refs []*types.NativeTypeRef, err error) {
	for _, f := range pkg.Files {
		for _, group := range f.Comments {
			_, data, found := strings.Cut(group.Text(), "go-restli:typeref")

			if !found {
				continue
			}

			ref := new(types.NativeTypeRef)
			err = json.Unmarshal([]byte(strings.TrimSpace(data)), ref)
			if err != nil {
				return nil, err
			}

			// ignore the error, we know it's valid JSON since we were able to deserialize it
			pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
			log.Printf("Found native typeref for %q: %s", ref.Ref, pretty)
			refs = append(refs, ref)
		}
	}
	return refs, nil
}

func locateAllGeneratedTypes(pkg *ast.Package) (manifests []*types.GeneratedTypeManifest, err error) {
	for _, f := range pkg.Files {
		for _, group := range f.Comments {
			_, data, found := strings.Cut(group.Text(), "go-restli:generated")

			if !found {
				continue
			}

			manifest := new(types.GeneratedTypeManifest)
			err = json.Unmarshal([]byte(strings.TrimSpace(data)), manifest)
			if err != nil {
				return nil, err
			}

			// ignore the error, we know it's valid JSON since we were able to deserialize it
			pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
			log.Printf("Found existing implementation for %q: %s", manifest.Identifier, pretty)
			manifests = append(manifests, manifest)
		}
	}
	return manifests, nil
}

func allNativeRefMethodsDefined(fset *token.FileSet, n *types.NativeTypeRef, pkgPath string, pkg *ast.Package) error {
	expectedType, err := locateUnmarshalRestLi(n, pkg, pkgPath, fset)
	if err != nil {
		return err
	}

	expectedFuncs := map[string]*ast.FuncDecl{}

	boolType := &ast.Ident{Name: "bool"}
	hashType := &ast.SelectorExpr{Sel: &ast.Ident{Name: "Hash"}}
	primitiveType := &ast.Ident{Name: n.Type.Type}
	errorType := &ast.Ident{Name: "error"}
	if n.IsCustomStruct() {
		expectedFuncs[utils.MarshalRestLi] = funcDecl(
			expectedType,
			utils.MarshalRestLi,
			[]ast.Expr{},
			[]ast.Expr{primitiveType, errorType},
		)
		expectedFuncs[utils.ComputeHash] = funcDecl(
			expectedType,
			utils.ComputeHash,
			[]ast.Expr{},
			[]ast.Expr{hashType},
		)
		expectedFuncs[utils.Equals] = funcDecl(
			expectedType,
			utils.Equals,
			[]ast.Expr{expectedType},
			[]ast.Expr{boolType},
		)
	} else {
		suffix := typeSuffix(n)
		expectedFuncs[utils.MarshalRestLi+suffix] = funcDecl(
			nil,
			utils.MarshalRestLi+suffix,
			[]ast.Expr{expectedType},
			[]ast.Expr{primitiveType, errorType},
		)
		expectedFuncs[utils.ComputeHash+suffix] = funcDecl(
			nil,
			utils.ComputeHash+suffix,
			[]ast.Expr{expectedType},
			[]ast.Expr{hashType},
		)
		expectedFuncs[utils.Equals+suffix] = funcDecl(
			nil,
			utils.Equals+suffix,
			[]ast.Expr{expectedType, expectedType},
			[]ast.Expr{boolType},
		)
	}

	typeDeclared := false

	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			switch decl := d.(type) {
			case *ast.FuncDecl:
				if expectedDecl, ok := expectedFuncs[decl.Name.Name]; ok && compareFuncs(expectedDecl, decl) {
					delete(expectedFuncs, decl.Name.Name)
				}
			case *ast.GenDecl:
				if !n.IsCustomStruct() || len(decl.Specs) == 0 {
					continue
				}
				if t, ok := decl.Specs[0].(*ast.TypeSpec); ok && t.Name.Name == n.Name {
					typeDeclared = true
				}
			}
		}
	}

	if n.IsCustomStruct() && !typeDeclared {
		return fmt.Errorf("go-restli: %q defines a typeref binding for %q but does not define a type even though no"+
			" non-receiver functions package is defined", pkgPath, n.Ref)
	}

	if len(expectedFuncs) != 0 {
		var funcNames []string
		for f := range expectedFuncs {
			funcNames = append(funcNames, f)
		}
		sort.Strings(funcNames)

		msg := fmt.Sprintf("go-restli: %q defines a typeref binding for %q but does not define the following funcs: %s",
			pkgPath, n.Ref, funcNames)
		if n.IsCustomStruct() {
			t := n.Name
			if n.ShouldReference.ShouldUsePointer() {
				t = "*" + t
			}
			msg += fmt.Sprintf(" (%s, %s and %s are expected to be defined on a receiver of type %s)",
				utils.MarshalRestLi, utils.ComputeHash, utils.Equals, t)
		}
		return fmt.Errorf(msg)
	}

	return nil
}

func locateUnmarshalRestLi(n *types.NativeTypeRef, pkg *ast.Package, pkgPath string, fset *token.FileSet) (expectedType ast.Expr, err error) {
	var nonPointerType ast.Expr = &ast.Ident{Name: n.Name}
	if !n.IsCustomStruct() {
		nonPointerType = &ast.SelectorExpr{Sel: nonPointerType.(*ast.Ident)}
	}
	pointerType := &ast.StarExpr{X: nonPointerType}

	expectedName := utils.UnmarshalRestLi + typeSuffix(n)
	unmarshalFunc := func(t ast.Expr) *ast.FuncDecl {
		return funcDecl(
			nil,
			expectedName,
			[]ast.Expr{&ast.Ident{Name: n.Type.Type}},
			[]ast.Expr{t, &ast.Ident{Name: "error"}},
		)
	}

	nonPointerFunc := unmarshalFunc(nonPointerType)
	pointerFunc := unmarshalFunc(pointerType)

	expectedTypes := fmt.Sprintf("expected `func %s(%s) (%s, error)` or `func %s(%s) (*%s, error)`",
		expectedName, n.Type.Type, n.Name,
		expectedName, n.Type.Type, n.Name,
	)

	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			decl, ok := d.(*ast.FuncDecl)
			if !ok || decl.Name.Name != expectedName {
				continue
			}

			if compareFuncs(nonPointerFunc, decl) {
				n.ShouldReference = utils.No
				return nonPointerType, nil
			}
			if compareFuncs(pointerFunc, decl) {
				n.ShouldReference = utils.Yes
				return pointerType, nil
			}
			declCopy := *decl
			declCopy.Body = nil
			return nil, fmt.Errorf("go-restli: Package %q defines an invalid %s (got `%s`, %s)",
				pkgPath, expectedName, prettyPrintAST(fset, &declCopy), expectedTypes,
			)
		}
	}

	return nil, fmt.Errorf("go-restli: Package %q does not define %s (%s)", pkgPath, expectedName, expectedTypes)
}

func prettyPrintAST(fset *token.FileSet, node any) string {
	out := new(strings.Builder)
	err := format.Node(out, fset, node)
	if err != nil {
		log.Panicf("Failed to print %+v: %s", node, err)
	}
	return out.String()
}

func typeSuffix(n *types.NativeTypeRef) string {
	return n.Ref[strings.LastIndex(n.Ref, ".")+1:]
}

func compareFuncs(left, right *ast.FuncDecl) bool {
	if !compareFieldList(left.Recv, right.Recv) {
		return false
	}

	if left.Name.Name != right.Name.Name {
		return false
	}

	if !compareFieldList(left.Type.TypeParams, right.Type.TypeParams) {
		return false
	}

	if !compareFieldList(left.Type.Params, right.Type.Params) {
		return false
	}

	if !compareFieldList(left.Type.Results, right.Type.Results) {
		return false
	}

	return true
}

func extractFieldTypes(l *ast.FieldList) (types []ast.Expr) {
	for _, f := range l.List {
		fields := len(f.Names)
		// Unnamed fields (like output fields) declare no names but should still be counted
		if fields == 0 {
			fields = 1
		}
		for i := 0; i < fields; i++ {
			types = append(types, f.Type)
		}
	}
	return types
}

func compareFieldList(left, right *ast.FieldList) bool {
	if left == nil || right == nil {
		return left == right
	}

	leftFields, rightFields := extractFieldTypes(left), extractFieldTypes(right)
	if len(leftFields) != len(rightFields) {
		return false
	}

	for i := range leftFields {
		if !compareTypes(leftFields[i], rightFields[i]) {
			return false
		}
	}

	return true
}

func compareTypes(left, right ast.Expr) bool {
	leftStar, leftOk := left.(*ast.StarExpr)
	rightStar, rightOk := right.(*ast.StarExpr)

	if leftOk != rightOk {
		return false
	}

	if leftOk {
		left, right = leftStar.X, rightStar.X
	}

	leftSel, leftOk := left.(*ast.SelectorExpr)
	rightSel, rightOk := right.(*ast.SelectorExpr)

	if leftOk != rightOk {
		return false
	}

	var leftIdent, rightIdent *ast.Ident
	if leftOk {
		// TODO: Find a way to compare the package names as well. Technically because package imports can be renamed,
		//  it's not possible to check leftSel.X == rightSet.X directly. For now, checking the literal type names will
		//  have to do since they are not allowed to be renamed
		leftIdent, rightIdent = leftSel.Sel, rightSel.Sel
	} else {
		leftIdent, leftOk = left.(*ast.Ident)
		if !leftOk {
			log.Panicf("Unknown type expression: %+v", left)
		}
		rightIdent, rightOk = right.(*ast.Ident)
		if !rightOk {
			log.Panicf("Unknown type expression: %+v", right)
		}
	}

	return leftIdent.Name == rightIdent.Name
}

func funcDecl(recv ast.Expr, name string, params []ast.Expr, returns []ast.Expr) *ast.FuncDecl {
	newFieldList := func(types ...ast.Expr) *ast.FieldList {
		l := &ast.FieldList{}
		for _, t := range types {
			l.List = append(l.List, &ast.Field{
				Names: []*ast.Ident{nil},
				Type:  t,
			})
		}
		return l
	}

	decl := &ast.FuncDecl{
		Name: &ast.Ident{Name: name},
		Type: &ast.FuncType{
			Params:  newFieldList(params...),
			Results: newFieldList(returns...),
		},
	}

	if recv != nil {
		decl.Recv = newFieldList(recv)
	}

	return decl
}
