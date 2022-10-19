// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Pdl
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// PdlListener is a complete listener for a parse tree produced by PdlParser.
type PdlListener interface {
	antlr.ParseTreeListener

	// EnterDocument is called when entering the document production.
	EnterDocument(c *DocumentContext)

	// EnterNamespaceDeclaration is called when entering the namespaceDeclaration production.
	EnterNamespaceDeclaration(c *NamespaceDeclarationContext)

	// EnterPackageDeclaration is called when entering the packageDeclaration production.
	EnterPackageDeclaration(c *PackageDeclarationContext)

	// EnterImportDeclarations is called when entering the importDeclarations production.
	EnterImportDeclarations(c *ImportDeclarationsContext)

	// EnterImportDeclaration is called when entering the importDeclaration production.
	EnterImportDeclaration(c *ImportDeclarationContext)

	// EnterTypeReference is called when entering the typeReference production.
	EnterTypeReference(c *TypeReferenceContext)

	// EnterTypeDeclaration is called when entering the typeDeclaration production.
	EnterTypeDeclaration(c *TypeDeclarationContext)

	// EnterNamedTypeDeclaration is called when entering the namedTypeDeclaration production.
	EnterNamedTypeDeclaration(c *NamedTypeDeclarationContext)

	// EnterScopedNamedTypeDeclaration is called when entering the scopedNamedTypeDeclaration production.
	EnterScopedNamedTypeDeclaration(c *ScopedNamedTypeDeclarationContext)

	// EnterAnonymousTypeDeclaration is called when entering the anonymousTypeDeclaration production.
	EnterAnonymousTypeDeclaration(c *AnonymousTypeDeclarationContext)

	// EnterTypeAssignment is called when entering the typeAssignment production.
	EnterTypeAssignment(c *TypeAssignmentContext)

	// EnterPropDeclaration is called when entering the propDeclaration production.
	EnterPropDeclaration(c *PropDeclarationContext)

	// EnterPropNameDeclaration is called when entering the propNameDeclaration production.
	EnterPropNameDeclaration(c *PropNameDeclarationContext)

	// EnterPropJsonValue is called when entering the propJsonValue production.
	EnterPropJsonValue(c *PropJsonValueContext)

	// EnterRecordDeclaration is called when entering the recordDeclaration production.
	EnterRecordDeclaration(c *RecordDeclarationContext)

	// EnterEnumDeclaration is called when entering the enumDeclaration production.
	EnterEnumDeclaration(c *EnumDeclarationContext)

	// EnterEnumSymbolDeclarations is called when entering the enumSymbolDeclarations production.
	EnterEnumSymbolDeclarations(c *EnumSymbolDeclarationsContext)

	// EnterEnumSymbolDeclaration is called when entering the enumSymbolDeclaration production.
	EnterEnumSymbolDeclaration(c *EnumSymbolDeclarationContext)

	// EnterEnumSymbol is called when entering the enumSymbol production.
	EnterEnumSymbol(c *EnumSymbolContext)

	// EnterTyperefDeclaration is called when entering the typerefDeclaration production.
	EnterTyperefDeclaration(c *TyperefDeclarationContext)

	// EnterFixedDeclaration is called when entering the fixedDeclaration production.
	EnterFixedDeclaration(c *FixedDeclarationContext)

	// EnterUnionDeclaration is called when entering the unionDeclaration production.
	EnterUnionDeclaration(c *UnionDeclarationContext)

	// EnterUnionTypeAssignments is called when entering the unionTypeAssignments production.
	EnterUnionTypeAssignments(c *UnionTypeAssignmentsContext)

	// EnterUnionMemberDeclaration is called when entering the unionMemberDeclaration production.
	EnterUnionMemberDeclaration(c *UnionMemberDeclarationContext)

	// EnterUnionMemberAlias is called when entering the unionMemberAlias production.
	EnterUnionMemberAlias(c *UnionMemberAliasContext)

	// EnterArrayDeclaration is called when entering the arrayDeclaration production.
	EnterArrayDeclaration(c *ArrayDeclarationContext)

	// EnterArrayTypeAssignments is called when entering the arrayTypeAssignments production.
	EnterArrayTypeAssignments(c *ArrayTypeAssignmentsContext)

	// EnterMapDeclaration is called when entering the mapDeclaration production.
	EnterMapDeclaration(c *MapDeclarationContext)

	// EnterMapTypeAssignments is called when entering the mapTypeAssignments production.
	EnterMapTypeAssignments(c *MapTypeAssignmentsContext)

	// EnterFieldSelection is called when entering the fieldSelection production.
	EnterFieldSelection(c *FieldSelectionContext)

	// EnterFieldIncludes is called when entering the fieldIncludes production.
	EnterFieldIncludes(c *FieldIncludesContext)

	// EnterFieldDeclaration is called when entering the fieldDeclaration production.
	EnterFieldDeclaration(c *FieldDeclarationContext)

	// EnterFieldDefault is called when entering the fieldDefault production.
	EnterFieldDefault(c *FieldDefaultContext)

	// EnterTypeName is called when entering the typeName production.
	EnterTypeName(c *TypeNameContext)

	// EnterIdentifier is called when entering the identifier production.
	EnterIdentifier(c *IdentifierContext)

	// EnterPropName is called when entering the propName production.
	EnterPropName(c *PropNameContext)

	// EnterPropSegment is called when entering the propSegment production.
	EnterPropSegment(c *PropSegmentContext)

	// EnterSchemadoc is called when entering the schemadoc production.
	EnterSchemadoc(c *SchemadocContext)

	// EnterObject is called when entering the object production.
	EnterObject(c *ObjectContext)

	// EnterObjectEntry is called when entering the objectEntry production.
	EnterObjectEntry(c *ObjectEntryContext)

	// EnterArray is called when entering the array production.
	EnterArray(c *ArrayContext)

	// EnterJsonValue is called when entering the jsonValue production.
	EnterJsonValue(c *JsonValueContext)

	// EnterString is called when entering the string production.
	EnterString(c *StringContext)

	// EnterNumber is called when entering the number production.
	EnterNumber(c *NumberContext)

	// EnterBool is called when entering the bool production.
	EnterBool(c *BoolContext)

	// EnterNullValue is called when entering the nullValue production.
	EnterNullValue(c *NullValueContext)

	// ExitDocument is called when exiting the document production.
	ExitDocument(c *DocumentContext)

	// ExitNamespaceDeclaration is called when exiting the namespaceDeclaration production.
	ExitNamespaceDeclaration(c *NamespaceDeclarationContext)

	// ExitPackageDeclaration is called when exiting the packageDeclaration production.
	ExitPackageDeclaration(c *PackageDeclarationContext)

	// ExitImportDeclarations is called when exiting the importDeclarations production.
	ExitImportDeclarations(c *ImportDeclarationsContext)

	// ExitImportDeclaration is called when exiting the importDeclaration production.
	ExitImportDeclaration(c *ImportDeclarationContext)

	// ExitTypeReference is called when exiting the typeReference production.
	ExitTypeReference(c *TypeReferenceContext)

	// ExitTypeDeclaration is called when exiting the typeDeclaration production.
	ExitTypeDeclaration(c *TypeDeclarationContext)

	// ExitNamedTypeDeclaration is called when exiting the namedTypeDeclaration production.
	ExitNamedTypeDeclaration(c *NamedTypeDeclarationContext)

	// ExitScopedNamedTypeDeclaration is called when exiting the scopedNamedTypeDeclaration production.
	ExitScopedNamedTypeDeclaration(c *ScopedNamedTypeDeclarationContext)

	// ExitAnonymousTypeDeclaration is called when exiting the anonymousTypeDeclaration production.
	ExitAnonymousTypeDeclaration(c *AnonymousTypeDeclarationContext)

	// ExitTypeAssignment is called when exiting the typeAssignment production.
	ExitTypeAssignment(c *TypeAssignmentContext)

	// ExitPropDeclaration is called when exiting the propDeclaration production.
	ExitPropDeclaration(c *PropDeclarationContext)

	// ExitPropNameDeclaration is called when exiting the propNameDeclaration production.
	ExitPropNameDeclaration(c *PropNameDeclarationContext)

	// ExitPropJsonValue is called when exiting the propJsonValue production.
	ExitPropJsonValue(c *PropJsonValueContext)

	// ExitRecordDeclaration is called when exiting the recordDeclaration production.
	ExitRecordDeclaration(c *RecordDeclarationContext)

	// ExitEnumDeclaration is called when exiting the enumDeclaration production.
	ExitEnumDeclaration(c *EnumDeclarationContext)

	// ExitEnumSymbolDeclarations is called when exiting the enumSymbolDeclarations production.
	ExitEnumSymbolDeclarations(c *EnumSymbolDeclarationsContext)

	// ExitEnumSymbolDeclaration is called when exiting the enumSymbolDeclaration production.
	ExitEnumSymbolDeclaration(c *EnumSymbolDeclarationContext)

	// ExitEnumSymbol is called when exiting the enumSymbol production.
	ExitEnumSymbol(c *EnumSymbolContext)

	// ExitTyperefDeclaration is called when exiting the typerefDeclaration production.
	ExitTyperefDeclaration(c *TyperefDeclarationContext)

	// ExitFixedDeclaration is called when exiting the fixedDeclaration production.
	ExitFixedDeclaration(c *FixedDeclarationContext)

	// ExitUnionDeclaration is called when exiting the unionDeclaration production.
	ExitUnionDeclaration(c *UnionDeclarationContext)

	// ExitUnionTypeAssignments is called when exiting the unionTypeAssignments production.
	ExitUnionTypeAssignments(c *UnionTypeAssignmentsContext)

	// ExitUnionMemberDeclaration is called when exiting the unionMemberDeclaration production.
	ExitUnionMemberDeclaration(c *UnionMemberDeclarationContext)

	// ExitUnionMemberAlias is called when exiting the unionMemberAlias production.
	ExitUnionMemberAlias(c *UnionMemberAliasContext)

	// ExitArrayDeclaration is called when exiting the arrayDeclaration production.
	ExitArrayDeclaration(c *ArrayDeclarationContext)

	// ExitArrayTypeAssignments is called when exiting the arrayTypeAssignments production.
	ExitArrayTypeAssignments(c *ArrayTypeAssignmentsContext)

	// ExitMapDeclaration is called when exiting the mapDeclaration production.
	ExitMapDeclaration(c *MapDeclarationContext)

	// ExitMapTypeAssignments is called when exiting the mapTypeAssignments production.
	ExitMapTypeAssignments(c *MapTypeAssignmentsContext)

	// ExitFieldSelection is called when exiting the fieldSelection production.
	ExitFieldSelection(c *FieldSelectionContext)

	// ExitFieldIncludes is called when exiting the fieldIncludes production.
	ExitFieldIncludes(c *FieldIncludesContext)

	// ExitFieldDeclaration is called when exiting the fieldDeclaration production.
	ExitFieldDeclaration(c *FieldDeclarationContext)

	// ExitFieldDefault is called when exiting the fieldDefault production.
	ExitFieldDefault(c *FieldDefaultContext)

	// ExitTypeName is called when exiting the typeName production.
	ExitTypeName(c *TypeNameContext)

	// ExitIdentifier is called when exiting the identifier production.
	ExitIdentifier(c *IdentifierContext)

	// ExitPropName is called when exiting the propName production.
	ExitPropName(c *PropNameContext)

	// ExitPropSegment is called when exiting the propSegment production.
	ExitPropSegment(c *PropSegmentContext)

	// ExitSchemadoc is called when exiting the schemadoc production.
	ExitSchemadoc(c *SchemadocContext)

	// ExitObject is called when exiting the object production.
	ExitObject(c *ObjectContext)

	// ExitObjectEntry is called when exiting the objectEntry production.
	ExitObjectEntry(c *ObjectEntryContext)

	// ExitArray is called when exiting the array production.
	ExitArray(c *ArrayContext)

	// ExitJsonValue is called when exiting the jsonValue production.
	ExitJsonValue(c *JsonValueContext)

	// ExitString is called when exiting the string production.
	ExitString(c *StringContext)

	// ExitNumber is called when exiting the number production.
	ExitNumber(c *NumberContext)

	// ExitBool is called when exiting the bool production.
	ExitBool(c *BoolContext)

	// ExitNullValue is called when exiting the nullValue production.
	ExitNullValue(c *NullValueContext)
}
