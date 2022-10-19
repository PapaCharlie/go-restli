// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Pdl
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// BasePdlListener is a complete listener for a parse tree produced by PdlParser.
type BasePdlListener struct{}

var _ PdlListener = &BasePdlListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasePdlListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasePdlListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasePdlListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasePdlListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterDocument is called when production document is entered.
func (s *BasePdlListener) EnterDocument(ctx *DocumentContext) {}

// ExitDocument is called when production document is exited.
func (s *BasePdlListener) ExitDocument(ctx *DocumentContext) {}

// EnterNamespaceDeclaration is called when production namespaceDeclaration is entered.
func (s *BasePdlListener) EnterNamespaceDeclaration(ctx *NamespaceDeclarationContext) {}

// ExitNamespaceDeclaration is called when production namespaceDeclaration is exited.
func (s *BasePdlListener) ExitNamespaceDeclaration(ctx *NamespaceDeclarationContext) {}

// EnterPackageDeclaration is called when production packageDeclaration is entered.
func (s *BasePdlListener) EnterPackageDeclaration(ctx *PackageDeclarationContext) {}

// ExitPackageDeclaration is called when production packageDeclaration is exited.
func (s *BasePdlListener) ExitPackageDeclaration(ctx *PackageDeclarationContext) {}

// EnterImportDeclarations is called when production importDeclarations is entered.
func (s *BasePdlListener) EnterImportDeclarations(ctx *ImportDeclarationsContext) {}

// ExitImportDeclarations is called when production importDeclarations is exited.
func (s *BasePdlListener) ExitImportDeclarations(ctx *ImportDeclarationsContext) {}

// EnterImportDeclaration is called when production importDeclaration is entered.
func (s *BasePdlListener) EnterImportDeclaration(ctx *ImportDeclarationContext) {}

// ExitImportDeclaration is called when production importDeclaration is exited.
func (s *BasePdlListener) ExitImportDeclaration(ctx *ImportDeclarationContext) {}

// EnterTypeReference is called when production typeReference is entered.
func (s *BasePdlListener) EnterTypeReference(ctx *TypeReferenceContext) {}

// ExitTypeReference is called when production typeReference is exited.
func (s *BasePdlListener) ExitTypeReference(ctx *TypeReferenceContext) {}

// EnterTypeDeclaration is called when production typeDeclaration is entered.
func (s *BasePdlListener) EnterTypeDeclaration(ctx *TypeDeclarationContext) {}

// ExitTypeDeclaration is called when production typeDeclaration is exited.
func (s *BasePdlListener) ExitTypeDeclaration(ctx *TypeDeclarationContext) {}

// EnterNamedTypeDeclaration is called when production namedTypeDeclaration is entered.
func (s *BasePdlListener) EnterNamedTypeDeclaration(ctx *NamedTypeDeclarationContext) {}

// ExitNamedTypeDeclaration is called when production namedTypeDeclaration is exited.
func (s *BasePdlListener) ExitNamedTypeDeclaration(ctx *NamedTypeDeclarationContext) {}

// EnterScopedNamedTypeDeclaration is called when production scopedNamedTypeDeclaration is entered.
func (s *BasePdlListener) EnterScopedNamedTypeDeclaration(ctx *ScopedNamedTypeDeclarationContext) {}

// ExitScopedNamedTypeDeclaration is called when production scopedNamedTypeDeclaration is exited.
func (s *BasePdlListener) ExitScopedNamedTypeDeclaration(ctx *ScopedNamedTypeDeclarationContext) {}

// EnterAnonymousTypeDeclaration is called when production anonymousTypeDeclaration is entered.
func (s *BasePdlListener) EnterAnonymousTypeDeclaration(ctx *AnonymousTypeDeclarationContext) {}

// ExitAnonymousTypeDeclaration is called when production anonymousTypeDeclaration is exited.
func (s *BasePdlListener) ExitAnonymousTypeDeclaration(ctx *AnonymousTypeDeclarationContext) {}

// EnterTypeAssignment is called when production typeAssignment is entered.
func (s *BasePdlListener) EnterTypeAssignment(ctx *TypeAssignmentContext) {}

// ExitTypeAssignment is called when production typeAssignment is exited.
func (s *BasePdlListener) ExitTypeAssignment(ctx *TypeAssignmentContext) {}

// EnterPropDeclaration is called when production propDeclaration is entered.
func (s *BasePdlListener) EnterPropDeclaration(ctx *PropDeclarationContext) {}

// ExitPropDeclaration is called when production propDeclaration is exited.
func (s *BasePdlListener) ExitPropDeclaration(ctx *PropDeclarationContext) {}

// EnterPropNameDeclaration is called when production propNameDeclaration is entered.
func (s *BasePdlListener) EnterPropNameDeclaration(ctx *PropNameDeclarationContext) {}

// ExitPropNameDeclaration is called when production propNameDeclaration is exited.
func (s *BasePdlListener) ExitPropNameDeclaration(ctx *PropNameDeclarationContext) {}

// EnterPropJsonValue is called when production propJsonValue is entered.
func (s *BasePdlListener) EnterPropJsonValue(ctx *PropJsonValueContext) {}

// ExitPropJsonValue is called when production propJsonValue is exited.
func (s *BasePdlListener) ExitPropJsonValue(ctx *PropJsonValueContext) {}

// EnterRecordDeclaration is called when production recordDeclaration is entered.
func (s *BasePdlListener) EnterRecordDeclaration(ctx *RecordDeclarationContext) {}

// ExitRecordDeclaration is called when production recordDeclaration is exited.
func (s *BasePdlListener) ExitRecordDeclaration(ctx *RecordDeclarationContext) {}

// EnterEnumDeclaration is called when production enumDeclaration is entered.
func (s *BasePdlListener) EnterEnumDeclaration(ctx *EnumDeclarationContext) {}

// ExitEnumDeclaration is called when production enumDeclaration is exited.
func (s *BasePdlListener) ExitEnumDeclaration(ctx *EnumDeclarationContext) {}

// EnterEnumSymbolDeclarations is called when production enumSymbolDeclarations is entered.
func (s *BasePdlListener) EnterEnumSymbolDeclarations(ctx *EnumSymbolDeclarationsContext) {}

// ExitEnumSymbolDeclarations is called when production enumSymbolDeclarations is exited.
func (s *BasePdlListener) ExitEnumSymbolDeclarations(ctx *EnumSymbolDeclarationsContext) {}

// EnterEnumSymbolDeclaration is called when production enumSymbolDeclaration is entered.
func (s *BasePdlListener) EnterEnumSymbolDeclaration(ctx *EnumSymbolDeclarationContext) {}

// ExitEnumSymbolDeclaration is called when production enumSymbolDeclaration is exited.
func (s *BasePdlListener) ExitEnumSymbolDeclaration(ctx *EnumSymbolDeclarationContext) {}

// EnterEnumSymbol is called when production enumSymbol is entered.
func (s *BasePdlListener) EnterEnumSymbol(ctx *EnumSymbolContext) {}

// ExitEnumSymbol is called when production enumSymbol is exited.
func (s *BasePdlListener) ExitEnumSymbol(ctx *EnumSymbolContext) {}

// EnterTyperefDeclaration is called when production typerefDeclaration is entered.
func (s *BasePdlListener) EnterTyperefDeclaration(ctx *TyperefDeclarationContext) {}

// ExitTyperefDeclaration is called when production typerefDeclaration is exited.
func (s *BasePdlListener) ExitTyperefDeclaration(ctx *TyperefDeclarationContext) {}

// EnterFixedDeclaration is called when production fixedDeclaration is entered.
func (s *BasePdlListener) EnterFixedDeclaration(ctx *FixedDeclarationContext) {}

// ExitFixedDeclaration is called when production fixedDeclaration is exited.
func (s *BasePdlListener) ExitFixedDeclaration(ctx *FixedDeclarationContext) {}

// EnterUnionDeclaration is called when production unionDeclaration is entered.
func (s *BasePdlListener) EnterUnionDeclaration(ctx *UnionDeclarationContext) {}

// ExitUnionDeclaration is called when production unionDeclaration is exited.
func (s *BasePdlListener) ExitUnionDeclaration(ctx *UnionDeclarationContext) {}

// EnterUnionTypeAssignments is called when production unionTypeAssignments is entered.
func (s *BasePdlListener) EnterUnionTypeAssignments(ctx *UnionTypeAssignmentsContext) {}

// ExitUnionTypeAssignments is called when production unionTypeAssignments is exited.
func (s *BasePdlListener) ExitUnionTypeAssignments(ctx *UnionTypeAssignmentsContext) {}

// EnterUnionMemberDeclaration is called when production unionMemberDeclaration is entered.
func (s *BasePdlListener) EnterUnionMemberDeclaration(ctx *UnionMemberDeclarationContext) {}

// ExitUnionMemberDeclaration is called when production unionMemberDeclaration is exited.
func (s *BasePdlListener) ExitUnionMemberDeclaration(ctx *UnionMemberDeclarationContext) {}

// EnterUnionMemberAlias is called when production unionMemberAlias is entered.
func (s *BasePdlListener) EnterUnionMemberAlias(ctx *UnionMemberAliasContext) {}

// ExitUnionMemberAlias is called when production unionMemberAlias is exited.
func (s *BasePdlListener) ExitUnionMemberAlias(ctx *UnionMemberAliasContext) {}

// EnterArrayDeclaration is called when production arrayDeclaration is entered.
func (s *BasePdlListener) EnterArrayDeclaration(ctx *ArrayDeclarationContext) {}

// ExitArrayDeclaration is called when production arrayDeclaration is exited.
func (s *BasePdlListener) ExitArrayDeclaration(ctx *ArrayDeclarationContext) {}

// EnterArrayTypeAssignments is called when production arrayTypeAssignments is entered.
func (s *BasePdlListener) EnterArrayTypeAssignments(ctx *ArrayTypeAssignmentsContext) {}

// ExitArrayTypeAssignments is called when production arrayTypeAssignments is exited.
func (s *BasePdlListener) ExitArrayTypeAssignments(ctx *ArrayTypeAssignmentsContext) {}

// EnterMapDeclaration is called when production mapDeclaration is entered.
func (s *BasePdlListener) EnterMapDeclaration(ctx *MapDeclarationContext) {}

// ExitMapDeclaration is called when production mapDeclaration is exited.
func (s *BasePdlListener) ExitMapDeclaration(ctx *MapDeclarationContext) {}

// EnterMapTypeAssignments is called when production mapTypeAssignments is entered.
func (s *BasePdlListener) EnterMapTypeAssignments(ctx *MapTypeAssignmentsContext) {}

// ExitMapTypeAssignments is called when production mapTypeAssignments is exited.
func (s *BasePdlListener) ExitMapTypeAssignments(ctx *MapTypeAssignmentsContext) {}

// EnterFieldSelection is called when production fieldSelection is entered.
func (s *BasePdlListener) EnterFieldSelection(ctx *FieldSelectionContext) {}

// ExitFieldSelection is called when production fieldSelection is exited.
func (s *BasePdlListener) ExitFieldSelection(ctx *FieldSelectionContext) {}

// EnterFieldIncludes is called when production fieldIncludes is entered.
func (s *BasePdlListener) EnterFieldIncludes(ctx *FieldIncludesContext) {}

// ExitFieldIncludes is called when production fieldIncludes is exited.
func (s *BasePdlListener) ExitFieldIncludes(ctx *FieldIncludesContext) {}

// EnterFieldDeclaration is called when production fieldDeclaration is entered.
func (s *BasePdlListener) EnterFieldDeclaration(ctx *FieldDeclarationContext) {}

// ExitFieldDeclaration is called when production fieldDeclaration is exited.
func (s *BasePdlListener) ExitFieldDeclaration(ctx *FieldDeclarationContext) {}

// EnterFieldDefault is called when production fieldDefault is entered.
func (s *BasePdlListener) EnterFieldDefault(ctx *FieldDefaultContext) {}

// ExitFieldDefault is called when production fieldDefault is exited.
func (s *BasePdlListener) ExitFieldDefault(ctx *FieldDefaultContext) {}

// EnterTypeName is called when production typeName is entered.
func (s *BasePdlListener) EnterTypeName(ctx *TypeNameContext) {}

// ExitTypeName is called when production typeName is exited.
func (s *BasePdlListener) ExitTypeName(ctx *TypeNameContext) {}

// EnterIdentifier is called when production identifier is entered.
func (s *BasePdlListener) EnterIdentifier(ctx *IdentifierContext) {}

// ExitIdentifier is called when production identifier is exited.
func (s *BasePdlListener) ExitIdentifier(ctx *IdentifierContext) {}

// EnterPropName is called when production propName is entered.
func (s *BasePdlListener) EnterPropName(ctx *PropNameContext) {}

// ExitPropName is called when production propName is exited.
func (s *BasePdlListener) ExitPropName(ctx *PropNameContext) {}

// EnterPropSegment is called when production propSegment is entered.
func (s *BasePdlListener) EnterPropSegment(ctx *PropSegmentContext) {}

// ExitPropSegment is called when production propSegment is exited.
func (s *BasePdlListener) ExitPropSegment(ctx *PropSegmentContext) {}

// EnterSchemadoc is called when production schemadoc is entered.
func (s *BasePdlListener) EnterSchemadoc(ctx *SchemadocContext) {}

// ExitSchemadoc is called when production schemadoc is exited.
func (s *BasePdlListener) ExitSchemadoc(ctx *SchemadocContext) {}

// EnterObject is called when production object is entered.
func (s *BasePdlListener) EnterObject(ctx *ObjectContext) {}

// ExitObject is called when production object is exited.
func (s *BasePdlListener) ExitObject(ctx *ObjectContext) {}

// EnterObjectEntry is called when production objectEntry is entered.
func (s *BasePdlListener) EnterObjectEntry(ctx *ObjectEntryContext) {}

// ExitObjectEntry is called when production objectEntry is exited.
func (s *BasePdlListener) ExitObjectEntry(ctx *ObjectEntryContext) {}

// EnterArray is called when production array is entered.
func (s *BasePdlListener) EnterArray(ctx *ArrayContext) {}

// ExitArray is called when production array is exited.
func (s *BasePdlListener) ExitArray(ctx *ArrayContext) {}

// EnterJsonValue is called when production jsonValue is entered.
func (s *BasePdlListener) EnterJsonValue(ctx *JsonValueContext) {}

// ExitJsonValue is called when production jsonValue is exited.
func (s *BasePdlListener) ExitJsonValue(ctx *JsonValueContext) {}

// EnterString is called when production string is entered.
func (s *BasePdlListener) EnterString(ctx *StringContext) {}

// ExitString is called when production string is exited.
func (s *BasePdlListener) ExitString(ctx *StringContext) {}

// EnterNumber is called when production number is entered.
func (s *BasePdlListener) EnterNumber(ctx *NumberContext) {}

// ExitNumber is called when production number is exited.
func (s *BasePdlListener) ExitNumber(ctx *NumberContext) {}

// EnterBool is called when production bool is entered.
func (s *BasePdlListener) EnterBool(ctx *BoolContext) {}

// ExitBool is called when production bool is exited.
func (s *BasePdlListener) ExitBool(ctx *BoolContext) {}

// EnterNullValue is called when production nullValue is entered.
func (s *BasePdlListener) EnterNullValue(ctx *NullValueContext) {}

// ExitNullValue is called when production nullValue is exited.
func (s *BasePdlListener) ExitNullValue(ctx *NullValueContext) {}
