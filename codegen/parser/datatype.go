package parser

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type DataType struct {
	Enum            *types.Enum            `json:"enum"`
	Fixed           *types.Fixed           `json:"fixed"`
	Record          *types.Record          `json:"record"`
	ComplexKey      *types.ComplexKey      `json:"complexKey"`
	StandaloneUnion *types.StandaloneUnion `json:"standaloneUnion"`
	Typeref         *types.Typeref         `json:"typeref"`
}

func Parse(input antlr.CharStream) (*DocumentContext, error) {
	lexer := NewPdlLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewPdlParser(stream)
	new(parser.ErrorListener)
	p.AddErrorListener()
	p.BuildParseTrees = true
	tree := p.Document()
}

type ErrorListener struct {
	*antlr.DefaultErrorListener
	Err error
}

func (el *ErrorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	el.Err = errors.New("line " + strconv.Itoa(line) + ":" + strconv.Itoa(column) + " " + msg)
	fmt.Fprintln(os.Stderr, el.Err.Error())
}

var _ PdlListener = &FileListener{}

type FileListener struct {
	*BasePdlListener
	ParsedTypes []DataType
}

func (cl *FileListener) VisitErrorNode(c antlr.ErrorNode) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterDocument(c *DocumentContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterNamespaceDeclaration(c *NamespaceDeclarationContext) {
	fmt.Printf("%s %+v %q\n", caller(), c, c.TypeName().GetValue())
}

func (cl *FileListener) EnterPackageDeclaration(c *PackageDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterImportDeclarations(c *ImportDeclarationsContext) {
	fmt.Printf("%s %+v\n", caller(), c)
	for _, imp := range c.AllImportDeclaration() {
		fmt.Printf("import %s\n", imp.GetType_().GetValue())
	}
}

func (cl *FileListener) EnterTypeReference(c *TypeReferenceContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterTypeDeclaration(c *TypeDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterNamedTypeDeclaration(c *NamedTypeDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterScopedNamedTypeDeclaration(c *ScopedNamedTypeDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterAnonymousTypeDeclaration(c *AnonymousTypeDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterTypeAssignment(c *TypeAssignmentContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterPropDeclaration(c *PropDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterPropNameDeclaration(c *PropNameDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterPropJsonValue(c *PropJsonValueContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterRecordDeclaration(c *RecordDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterEnumDeclaration(c *EnumDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterEnumSymbolDeclarations(c *EnumSymbolDeclarationsContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterEnumSymbolDeclaration(c *EnumSymbolDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterEnumSymbol(c *EnumSymbolContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterTyperefDeclaration(c *TyperefDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterFixedDeclaration(c *FixedDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterUnionDeclaration(c *UnionDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterUnionTypeAssignments(c *UnionTypeAssignmentsContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterUnionMemberDeclaration(c *UnionMemberDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterUnionMemberAlias(c *UnionMemberAliasContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterArrayDeclaration(c *ArrayDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterArrayTypeAssignments(c *ArrayTypeAssignmentsContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterMapDeclaration(c *MapDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterMapTypeAssignments(c *MapTypeAssignmentsContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterFieldSelection(c *FieldSelectionContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterFieldIncludes(c *FieldIncludesContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterFieldDeclaration(c *FieldDeclarationContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterFieldDefault(c *FieldDefaultContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterTypeName(c *TypeNameContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterIdentifier(c *IdentifierContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterPropName(c *PropNameContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterPropSegment(c *PropSegmentContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterSchemadoc(c *SchemadocContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterObject(c *ObjectContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterObjectEntry(c *ObjectEntryContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterArray(c *ArrayContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterJsonValue(c *JsonValueContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterString(c *StringContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterNumber(c *NumberContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterBool(c *BoolContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}

func (cl *FileListener) EnterNullValue(c *NullValueContext) {
	fmt.Printf("%s %+v\n", caller(), c)
}
