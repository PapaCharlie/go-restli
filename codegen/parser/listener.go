package parser

import (
	"fmt"
	"log"
	"runtime"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

var _ PdlListener = &CustomListener{}

type CustomListener struct {
	*BasePdlListener
}

func caller() string {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if !ok || details == nil {
		log.Panicln("Could not get caller")
	}
	return details.Name()
}

func (cl *CustomListener) VisitErrorNode(c antlr.ErrorNode) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterDocument(c *DocumentContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterNamespaceDeclaration(c *NamespaceDeclarationContext) {
	fmt.Printf("%s %+v %q\n", caller(), c, c.TypeName().GetValue())
}

func (cl *CustomListener) EnterPackageDeclaration(c *PackageDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterImportDeclarations(c *ImportDeclarationsContext) {
	fmt.Printf("%s %+v\n", caller(), c)
	for _, imp := range c.AllImportDeclaration() {
		fmt.Printf("import %s\n", imp.GetType_().GetValue())
	}
}

func (cl *CustomListener) EnterTypeReference(c *TypeReferenceContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterTypeDeclaration(c *TypeDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterNamedTypeDeclaration(c *NamedTypeDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterScopedNamedTypeDeclaration(c *ScopedNamedTypeDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterAnonymousTypeDeclaration(c *AnonymousTypeDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterTypeAssignment(c *TypeAssignmentContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterPropDeclaration(c *PropDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterPropNameDeclaration(c *PropNameDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterPropJsonValue(c *PropJsonValueContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterRecordDeclaration(c *RecordDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterEnumDeclaration(c *EnumDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterEnumSymbolDeclarations(c *EnumSymbolDeclarationsContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterEnumSymbolDeclaration(c *EnumSymbolDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterEnumSymbol(c *EnumSymbolContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterTyperefDeclaration(c *TyperefDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterFixedDeclaration(c *FixedDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterUnionDeclaration(c *UnionDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterUnionTypeAssignments(c *UnionTypeAssignmentsContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterUnionMemberDeclaration(c *UnionMemberDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterUnionMemberAlias(c *UnionMemberAliasContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterArrayDeclaration(c *ArrayDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterArrayTypeAssignments(c *ArrayTypeAssignmentsContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterMapDeclaration(c *MapDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterMapTypeAssignments(c *MapTypeAssignmentsContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterFieldSelection(c *FieldSelectionContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterFieldIncludes(c *FieldIncludesContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterFieldDeclaration(c *FieldDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterFieldDefault(c *FieldDefaultContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterTypeName(c *TypeNameContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterIdentifier(c *IdentifierContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterPropName(c *PropNameContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterPropSegment(c *PropSegmentContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterSchemadoc(c *SchemadocContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterObject(c *ObjectContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterObjectEntry(c *ObjectEntryContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterArray(c *ArrayContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterJsonValue(c *JsonValueContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterString(c *StringContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterNumber(c *NumberContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterBool(c *BoolContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *CustomListener) EnterNullValue(c *NullValueContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}
