// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Pdl
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type PdlParser struct {
	*antlr.BaseParser
}

var pdlParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func pdlParserInit() {
	staticData := &pdlParserStaticData
	staticData.literalNames = []string{
		"", "'array'", "'enum'", "'fixed'", "'import'", "'optional'", "'package'",
		"'map'", "'namespace'", "'record'", "'typeref'", "'union'", "'includes'",
		"'('", "')'", "'{'", "'}'", "'['", "']'", "'@'", "':'", "'.'", "'='",
		"", "'null'",
	}
	staticData.symbolicNames = []string{
		"", "ARRAY", "ENUM", "FIXED", "IMPORT", "OPTIONAL", "PACKAGE", "MAP",
		"NAMESPACE", "RECORD", "TYPEREF", "UNION", "INCLUDES", "OPEN_PAREN",
		"CLOSE_PAREN", "OPEN_BRACE", "CLOSE_BRACE", "OPEN_BRACKET", "CLOSE_BRACKET",
		"AT", "COLON", "DOT", "EQ", "BOOLEAN_LITERAL", "NULL_LITERAL", "SCHEMADOC_COMMENT",
		"BLOCK_COMMENT", "LINE_COMMENT", "NUMBER_LITERAL", "STRING_LITERAL",
		"ID", "WS", "PROPERTY_ID", "ESCAPED_PROP_ID",
	}
	staticData.ruleNames = []string{
		"document", "namespaceDeclaration", "packageDeclaration", "importDeclarations",
		"importDeclaration", "typeReference", "typeDeclaration", "namedTypeDeclaration",
		"scopedNamedTypeDeclaration", "anonymousTypeDeclaration", "typeAssignment",
		"propDeclaration", "propNameDeclaration", "propJsonValue", "recordDeclaration",
		"enumDeclaration", "enumSymbolDeclarations", "enumSymbolDeclaration",
		"enumSymbol", "typerefDeclaration", "fixedDeclaration", "unionDeclaration",
		"unionTypeAssignments", "unionMemberDeclaration", "unionMemberAlias",
		"arrayDeclaration", "arrayTypeAssignments", "mapDeclaration", "mapTypeAssignments",
		"fieldSelection", "fieldIncludes", "fieldDeclaration", "fieldDefault",
		"typeName", "identifier", "propName", "propSegment", "schemadoc", "object",
		"objectEntry", "array", "jsonValue", "string", "number", "bool", "nullValue",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 33, 385, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 1, 0, 3, 0, 94, 8,
		0, 1, 0, 3, 0, 97, 8, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2,
		1, 2, 1, 3, 5, 3, 109, 8, 3, 10, 3, 12, 3, 112, 9, 3, 1, 4, 1, 4, 1, 4,
		1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 122, 8, 5, 1, 6, 1, 6, 1, 6, 3, 6,
		127, 8, 6, 1, 7, 3, 7, 130, 8, 7, 1, 7, 5, 7, 133, 8, 7, 10, 7, 12, 7,
		136, 9, 7, 1, 7, 1, 7, 1, 7, 1, 7, 3, 7, 142, 8, 7, 1, 8, 1, 8, 3, 8, 146,
		8, 8, 1, 8, 3, 8, 149, 8, 8, 1, 8, 1, 8, 1, 8, 1, 9, 5, 9, 155, 8, 9, 10,
		9, 12, 9, 158, 9, 9, 1, 9, 1, 9, 1, 9, 3, 9, 163, 8, 9, 1, 10, 1, 10, 3,
		10, 167, 8, 10, 1, 11, 1, 11, 3, 11, 171, 8, 11, 1, 11, 1, 11, 1, 12, 1,
		12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 3, 14, 185,
		8, 14, 1, 14, 1, 14, 3, 14, 189, 8, 14, 1, 14, 1, 14, 1, 15, 1, 15, 1,
		15, 1, 15, 1, 15, 1, 16, 1, 16, 5, 16, 200, 8, 16, 10, 16, 12, 16, 203,
		9, 16, 1, 16, 1, 16, 1, 17, 3, 17, 208, 8, 17, 1, 17, 5, 17, 211, 8, 17,
		10, 17, 12, 17, 214, 9, 17, 1, 17, 1, 17, 1, 18, 1, 18, 1, 18, 1, 19, 1,
		19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 21,
		1, 21, 1, 21, 1, 22, 1, 22, 5, 22, 237, 8, 22, 10, 22, 12, 22, 240, 9,
		22, 1, 22, 1, 22, 1, 23, 3, 23, 245, 8, 23, 1, 23, 1, 23, 1, 24, 3, 24,
		250, 8, 24, 1, 24, 5, 24, 253, 8, 24, 10, 24, 12, 24, 256, 9, 24, 1, 24,
		1, 24, 1, 24, 1, 25, 1, 25, 1, 25, 1, 26, 1, 26, 1, 26, 1, 26, 1, 27, 1,
		27, 1, 27, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 29, 1, 29, 5, 29, 278,
		8, 29, 10, 29, 12, 29, 281, 9, 29, 1, 29, 1, 29, 1, 30, 1, 30, 4, 30, 287,
		8, 30, 11, 30, 12, 30, 288, 1, 31, 3, 31, 292, 8, 31, 1, 31, 5, 31, 295,
		8, 31, 10, 31, 12, 31, 298, 9, 31, 1, 31, 1, 31, 1, 31, 3, 31, 303, 8,
		31, 1, 31, 1, 31, 3, 31, 307, 8, 31, 1, 31, 1, 31, 1, 32, 1, 32, 1, 32,
		1, 33, 1, 33, 1, 33, 5, 33, 317, 8, 33, 10, 33, 12, 33, 320, 9, 33, 1,
		33, 1, 33, 1, 34, 1, 34, 1, 34, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35,
		5, 35, 333, 8, 35, 10, 35, 12, 35, 336, 9, 35, 1, 36, 1, 36, 1, 36, 1,
		37, 1, 37, 1, 37, 1, 38, 1, 38, 5, 38, 346, 8, 38, 10, 38, 12, 38, 349,
		9, 38, 1, 38, 1, 38, 1, 39, 1, 39, 1, 39, 1, 39, 1, 40, 1, 40, 5, 40, 359,
		8, 40, 10, 40, 12, 40, 362, 9, 40, 1, 40, 1, 40, 1, 41, 1, 41, 1, 41, 1,
		41, 1, 41, 1, 41, 3, 41, 372, 8, 41, 1, 42, 1, 42, 1, 42, 1, 43, 1, 43,
		1, 43, 1, 44, 1, 44, 1, 44, 1, 45, 1, 45, 1, 45, 0, 0, 46, 0, 2, 4, 6,
		8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42,
		44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78,
		80, 82, 84, 86, 88, 90, 0, 1, 2, 0, 30, 30, 32, 33, 380, 0, 93, 1, 0, 0,
		0, 2, 101, 1, 0, 0, 0, 4, 104, 1, 0, 0, 0, 6, 110, 1, 0, 0, 0, 8, 113,
		1, 0, 0, 0, 10, 121, 1, 0, 0, 0, 12, 126, 1, 0, 0, 0, 14, 129, 1, 0, 0,
		0, 16, 143, 1, 0, 0, 0, 18, 156, 1, 0, 0, 0, 20, 166, 1, 0, 0, 0, 22, 168,
		1, 0, 0, 0, 24, 174, 1, 0, 0, 0, 26, 178, 1, 0, 0, 0, 28, 181, 1, 0, 0,
		0, 30, 192, 1, 0, 0, 0, 32, 197, 1, 0, 0, 0, 34, 207, 1, 0, 0, 0, 36, 217,
		1, 0, 0, 0, 38, 220, 1, 0, 0, 0, 40, 226, 1, 0, 0, 0, 42, 231, 1, 0, 0,
		0, 44, 234, 1, 0, 0, 0, 46, 244, 1, 0, 0, 0, 48, 249, 1, 0, 0, 0, 50, 260,
		1, 0, 0, 0, 52, 263, 1, 0, 0, 0, 54, 267, 1, 0, 0, 0, 56, 270, 1, 0, 0,
		0, 58, 275, 1, 0, 0, 0, 60, 284, 1, 0, 0, 0, 62, 291, 1, 0, 0, 0, 64, 310,
		1, 0, 0, 0, 66, 313, 1, 0, 0, 0, 68, 323, 1, 0, 0, 0, 70, 326, 1, 0, 0,
		0, 72, 337, 1, 0, 0, 0, 74, 340, 1, 0, 0, 0, 76, 343, 1, 0, 0, 0, 78, 352,
		1, 0, 0, 0, 80, 356, 1, 0, 0, 0, 82, 371, 1, 0, 0, 0, 84, 373, 1, 0, 0,
		0, 86, 376, 1, 0, 0, 0, 88, 379, 1, 0, 0, 0, 90, 382, 1, 0, 0, 0, 92, 94,
		3, 2, 1, 0, 93, 92, 1, 0, 0, 0, 93, 94, 1, 0, 0, 0, 94, 96, 1, 0, 0, 0,
		95, 97, 3, 4, 2, 0, 96, 95, 1, 0, 0, 0, 96, 97, 1, 0, 0, 0, 97, 98, 1,
		0, 0, 0, 98, 99, 3, 6, 3, 0, 99, 100, 3, 12, 6, 0, 100, 1, 1, 0, 0, 0,
		101, 102, 5, 8, 0, 0, 102, 103, 3, 66, 33, 0, 103, 3, 1, 0, 0, 0, 104,
		105, 5, 6, 0, 0, 105, 106, 3, 66, 33, 0, 106, 5, 1, 0, 0, 0, 107, 109,
		3, 8, 4, 0, 108, 107, 1, 0, 0, 0, 109, 112, 1, 0, 0, 0, 110, 108, 1, 0,
		0, 0, 110, 111, 1, 0, 0, 0, 111, 7, 1, 0, 0, 0, 112, 110, 1, 0, 0, 0, 113,
		114, 5, 4, 0, 0, 114, 115, 3, 66, 33, 0, 115, 9, 1, 0, 0, 0, 116, 117,
		5, 24, 0, 0, 117, 122, 6, 5, -1, 0, 118, 119, 3, 66, 33, 0, 119, 120, 6,
		5, -1, 0, 120, 122, 1, 0, 0, 0, 121, 116, 1, 0, 0, 0, 121, 118, 1, 0, 0,
		0, 122, 11, 1, 0, 0, 0, 123, 127, 3, 16, 8, 0, 124, 127, 3, 14, 7, 0, 125,
		127, 3, 18, 9, 0, 126, 123, 1, 0, 0, 0, 126, 124, 1, 0, 0, 0, 126, 125,
		1, 0, 0, 0, 127, 13, 1, 0, 0, 0, 128, 130, 3, 74, 37, 0, 129, 128, 1, 0,
		0, 0, 129, 130, 1, 0, 0, 0, 130, 134, 1, 0, 0, 0, 131, 133, 3, 22, 11,
		0, 132, 131, 1, 0, 0, 0, 133, 136, 1, 0, 0, 0, 134, 132, 1, 0, 0, 0, 134,
		135, 1, 0, 0, 0, 135, 141, 1, 0, 0, 0, 136, 134, 1, 0, 0, 0, 137, 142,
		3, 28, 14, 0, 138, 142, 3, 30, 15, 0, 139, 142, 3, 38, 19, 0, 140, 142,
		3, 40, 20, 0, 141, 137, 1, 0, 0, 0, 141, 138, 1, 0, 0, 0, 141, 139, 1,
		0, 0, 0, 141, 140, 1, 0, 0, 0, 142, 15, 1, 0, 0, 0, 143, 145, 5, 15, 0,
		0, 144, 146, 3, 2, 1, 0, 145, 144, 1, 0, 0, 0, 145, 146, 1, 0, 0, 0, 146,
		148, 1, 0, 0, 0, 147, 149, 3, 4, 2, 0, 148, 147, 1, 0, 0, 0, 148, 149,
		1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150, 151, 3, 14, 7, 0, 151, 152, 5, 16,
		0, 0, 152, 17, 1, 0, 0, 0, 153, 155, 3, 22, 11, 0, 154, 153, 1, 0, 0, 0,
		155, 158, 1, 0, 0, 0, 156, 154, 1, 0, 0, 0, 156, 157, 1, 0, 0, 0, 157,
		162, 1, 0, 0, 0, 158, 156, 1, 0, 0, 0, 159, 163, 3, 42, 21, 0, 160, 163,
		3, 50, 25, 0, 161, 163, 3, 54, 27, 0, 162, 159, 1, 0, 0, 0, 162, 160, 1,
		0, 0, 0, 162, 161, 1, 0, 0, 0, 163, 19, 1, 0, 0, 0, 164, 167, 3, 10, 5,
		0, 165, 167, 3, 12, 6, 0, 166, 164, 1, 0, 0, 0, 166, 165, 1, 0, 0, 0, 167,
		21, 1, 0, 0, 0, 168, 170, 3, 24, 12, 0, 169, 171, 3, 26, 13, 0, 170, 169,
		1, 0, 0, 0, 170, 171, 1, 0, 0, 0, 171, 172, 1, 0, 0, 0, 172, 173, 6, 11,
		-1, 0, 173, 23, 1, 0, 0, 0, 174, 175, 5, 19, 0, 0, 175, 176, 3, 70, 35,
		0, 176, 177, 6, 12, -1, 0, 177, 25, 1, 0, 0, 0, 178, 179, 5, 22, 0, 0,
		179, 180, 3, 82, 41, 0, 180, 27, 1, 0, 0, 0, 181, 182, 5, 9, 0, 0, 182,
		184, 3, 68, 34, 0, 183, 185, 3, 60, 30, 0, 184, 183, 1, 0, 0, 0, 184, 185,
		1, 0, 0, 0, 185, 186, 1, 0, 0, 0, 186, 188, 3, 58, 29, 0, 187, 189, 3,
		60, 30, 0, 188, 187, 1, 0, 0, 0, 188, 189, 1, 0, 0, 0, 189, 190, 1, 0,
		0, 0, 190, 191, 6, 14, -1, 0, 191, 29, 1, 0, 0, 0, 192, 193, 5, 2, 0, 0,
		193, 194, 3, 68, 34, 0, 194, 195, 3, 32, 16, 0, 195, 196, 6, 15, -1, 0,
		196, 31, 1, 0, 0, 0, 197, 201, 5, 15, 0, 0, 198, 200, 3, 34, 17, 0, 199,
		198, 1, 0, 0, 0, 200, 203, 1, 0, 0, 0, 201, 199, 1, 0, 0, 0, 201, 202,
		1, 0, 0, 0, 202, 204, 1, 0, 0, 0, 203, 201, 1, 0, 0, 0, 204, 205, 5, 16,
		0, 0, 205, 33, 1, 0, 0, 0, 206, 208, 3, 74, 37, 0, 207, 206, 1, 0, 0, 0,
		207, 208, 1, 0, 0, 0, 208, 212, 1, 0, 0, 0, 209, 211, 3, 22, 11, 0, 210,
		209, 1, 0, 0, 0, 211, 214, 1, 0, 0, 0, 212, 210, 1, 0, 0, 0, 212, 213,
		1, 0, 0, 0, 213, 215, 1, 0, 0, 0, 214, 212, 1, 0, 0, 0, 215, 216, 3, 36,
		18, 0, 216, 35, 1, 0, 0, 0, 217, 218, 3, 68, 34, 0, 218, 219, 6, 18, -1,
		0, 219, 37, 1, 0, 0, 0, 220, 221, 5, 10, 0, 0, 221, 222, 3, 68, 34, 0,
		222, 223, 5, 22, 0, 0, 223, 224, 3, 20, 10, 0, 224, 225, 6, 19, -1, 0,
		225, 39, 1, 0, 0, 0, 226, 227, 5, 3, 0, 0, 227, 228, 3, 68, 34, 0, 228,
		229, 5, 28, 0, 0, 229, 230, 6, 20, -1, 0, 230, 41, 1, 0, 0, 0, 231, 232,
		5, 11, 0, 0, 232, 233, 3, 44, 22, 0, 233, 43, 1, 0, 0, 0, 234, 238, 5,
		17, 0, 0, 235, 237, 3, 46, 23, 0, 236, 235, 1, 0, 0, 0, 237, 240, 1, 0,
		0, 0, 238, 236, 1, 0, 0, 0, 238, 239, 1, 0, 0, 0, 239, 241, 1, 0, 0, 0,
		240, 238, 1, 0, 0, 0, 241, 242, 5, 18, 0, 0, 242, 45, 1, 0, 0, 0, 243,
		245, 3, 48, 24, 0, 244, 243, 1, 0, 0, 0, 244, 245, 1, 0, 0, 0, 245, 246,
		1, 0, 0, 0, 246, 247, 3, 20, 10, 0, 247, 47, 1, 0, 0, 0, 248, 250, 3, 74,
		37, 0, 249, 248, 1, 0, 0, 0, 249, 250, 1, 0, 0, 0, 250, 254, 1, 0, 0, 0,
		251, 253, 3, 22, 11, 0, 252, 251, 1, 0, 0, 0, 253, 256, 1, 0, 0, 0, 254,
		252, 1, 0, 0, 0, 254, 255, 1, 0, 0, 0, 255, 257, 1, 0, 0, 0, 256, 254,
		1, 0, 0, 0, 257, 258, 3, 68, 34, 0, 258, 259, 5, 20, 0, 0, 259, 49, 1,
		0, 0, 0, 260, 261, 5, 1, 0, 0, 261, 262, 3, 52, 26, 0, 262, 51, 1, 0, 0,
		0, 263, 264, 5, 17, 0, 0, 264, 265, 3, 20, 10, 0, 265, 266, 5, 18, 0, 0,
		266, 53, 1, 0, 0, 0, 267, 268, 5, 7, 0, 0, 268, 269, 3, 56, 28, 0, 269,
		55, 1, 0, 0, 0, 270, 271, 5, 17, 0, 0, 271, 272, 3, 20, 10, 0, 272, 273,
		3, 20, 10, 0, 273, 274, 5, 18, 0, 0, 274, 57, 1, 0, 0, 0, 275, 279, 5,
		15, 0, 0, 276, 278, 3, 62, 31, 0, 277, 276, 1, 0, 0, 0, 278, 281, 1, 0,
		0, 0, 279, 277, 1, 0, 0, 0, 279, 280, 1, 0, 0, 0, 280, 282, 1, 0, 0, 0,
		281, 279, 1, 0, 0, 0, 282, 283, 5, 16, 0, 0, 283, 59, 1, 0, 0, 0, 284,
		286, 5, 12, 0, 0, 285, 287, 3, 20, 10, 0, 286, 285, 1, 0, 0, 0, 287, 288,
		1, 0, 0, 0, 288, 286, 1, 0, 0, 0, 288, 289, 1, 0, 0, 0, 289, 61, 1, 0,
		0, 0, 290, 292, 3, 74, 37, 0, 291, 290, 1, 0, 0, 0, 291, 292, 1, 0, 0,
		0, 292, 296, 1, 0, 0, 0, 293, 295, 3, 22, 11, 0, 294, 293, 1, 0, 0, 0,
		295, 298, 1, 0, 0, 0, 296, 294, 1, 0, 0, 0, 296, 297, 1, 0, 0, 0, 297,
		299, 1, 0, 0, 0, 298, 296, 1, 0, 0, 0, 299, 300, 3, 68, 34, 0, 300, 302,
		5, 20, 0, 0, 301, 303, 5, 5, 0, 0, 302, 301, 1, 0, 0, 0, 302, 303, 1, 0,
		0, 0, 303, 304, 1, 0, 0, 0, 304, 306, 3, 20, 10, 0, 305, 307, 3, 64, 32,
		0, 306, 305, 1, 0, 0, 0, 306, 307, 1, 0, 0, 0, 307, 308, 1, 0, 0, 0, 308,
		309, 6, 31, -1, 0, 309, 63, 1, 0, 0, 0, 310, 311, 5, 22, 0, 0, 311, 312,
		3, 82, 41, 0, 312, 65, 1, 0, 0, 0, 313, 318, 5, 30, 0, 0, 314, 315, 5,
		21, 0, 0, 315, 317, 5, 30, 0, 0, 316, 314, 1, 0, 0, 0, 317, 320, 1, 0,
		0, 0, 318, 316, 1, 0, 0, 0, 318, 319, 1, 0, 0, 0, 319, 321, 1, 0, 0, 0,
		320, 318, 1, 0, 0, 0, 321, 322, 6, 33, -1, 0, 322, 67, 1, 0, 0, 0, 323,
		324, 5, 30, 0, 0, 324, 325, 6, 34, -1, 0, 325, 69, 1, 0, 0, 0, 326, 327,
		3, 72, 36, 0, 327, 334, 6, 35, -1, 0, 328, 329, 5, 21, 0, 0, 329, 330,
		3, 72, 36, 0, 330, 331, 6, 35, -1, 0, 331, 333, 1, 0, 0, 0, 332, 328, 1,
		0, 0, 0, 333, 336, 1, 0, 0, 0, 334, 332, 1, 0, 0, 0, 334, 335, 1, 0, 0,
		0, 335, 71, 1, 0, 0, 0, 336, 334, 1, 0, 0, 0, 337, 338, 7, 0, 0, 0, 338,
		339, 6, 36, -1, 0, 339, 73, 1, 0, 0, 0, 340, 341, 5, 25, 0, 0, 341, 342,
		6, 37, -1, 0, 342, 75, 1, 0, 0, 0, 343, 347, 5, 15, 0, 0, 344, 346, 3,
		78, 39, 0, 345, 344, 1, 0, 0, 0, 346, 349, 1, 0, 0, 0, 347, 345, 1, 0,
		0, 0, 347, 348, 1, 0, 0, 0, 348, 350, 1, 0, 0, 0, 349, 347, 1, 0, 0, 0,
		350, 351, 5, 16, 0, 0, 351, 77, 1, 0, 0, 0, 352, 353, 3, 84, 42, 0, 353,
		354, 5, 20, 0, 0, 354, 355, 3, 82, 41, 0, 355, 79, 1, 0, 0, 0, 356, 360,
		5, 17, 0, 0, 357, 359, 3, 82, 41, 0, 358, 357, 1, 0, 0, 0, 359, 362, 1,
		0, 0, 0, 360, 358, 1, 0, 0, 0, 360, 361, 1, 0, 0, 0, 361, 363, 1, 0, 0,
		0, 362, 360, 1, 0, 0, 0, 363, 364, 5, 18, 0, 0, 364, 81, 1, 0, 0, 0, 365,
		372, 3, 84, 42, 0, 366, 372, 3, 86, 43, 0, 367, 372, 3, 76, 38, 0, 368,
		372, 3, 80, 40, 0, 369, 372, 3, 88, 44, 0, 370, 372, 3, 90, 45, 0, 371,
		365, 1, 0, 0, 0, 371, 366, 1, 0, 0, 0, 371, 367, 1, 0, 0, 0, 371, 368,
		1, 0, 0, 0, 371, 369, 1, 0, 0, 0, 371, 370, 1, 0, 0, 0, 372, 83, 1, 0,
		0, 0, 373, 374, 5, 29, 0, 0, 374, 375, 6, 42, -1, 0, 375, 85, 1, 0, 0,
		0, 376, 377, 5, 28, 0, 0, 377, 378, 6, 43, -1, 0, 378, 87, 1, 0, 0, 0,
		379, 380, 5, 23, 0, 0, 380, 381, 6, 44, -1, 0, 381, 89, 1, 0, 0, 0, 382,
		383, 5, 24, 0, 0, 383, 91, 1, 0, 0, 0, 34, 93, 96, 110, 121, 126, 129,
		134, 141, 145, 148, 156, 162, 166, 170, 184, 188, 201, 207, 212, 238, 244,
		249, 254, 279, 288, 291, 296, 302, 306, 318, 334, 347, 360, 371,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// PdlParserInit initializes any static state used to implement PdlParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewPdlParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func PdlParserInit() {
	staticData := &pdlParserStaticData
	staticData.once.Do(pdlParserInit)
}

// NewPdlParser produces a new parser instance for the optional input antlr.TokenStream.
func NewPdlParser(input antlr.TokenStream) *PdlParser {
	PdlParserInit()
	this := new(PdlParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &pdlParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "java-escape"

	return this
}

// PdlParser tokens.
const (
	PdlParserEOF               = antlr.TokenEOF
	PdlParserARRAY             = 1
	PdlParserENUM              = 2
	PdlParserFIXED             = 3
	PdlParserIMPORT            = 4
	PdlParserOPTIONAL          = 5
	PdlParserPACKAGE           = 6
	PdlParserMAP               = 7
	PdlParserNAMESPACE         = 8
	PdlParserRECORD            = 9
	PdlParserTYPEREF           = 10
	PdlParserUNION             = 11
	PdlParserINCLUDES          = 12
	PdlParserOPEN_PAREN        = 13
	PdlParserCLOSE_PAREN       = 14
	PdlParserOPEN_BRACE        = 15
	PdlParserCLOSE_BRACE       = 16
	PdlParserOPEN_BRACKET      = 17
	PdlParserCLOSE_BRACKET     = 18
	PdlParserAT                = 19
	PdlParserCOLON             = 20
	PdlParserDOT               = 21
	PdlParserEQ                = 22
	PdlParserBOOLEAN_LITERAL   = 23
	PdlParserNULL_LITERAL      = 24
	PdlParserSCHEMADOC_COMMENT = 25
	PdlParserBLOCK_COMMENT     = 26
	PdlParserLINE_COMMENT      = 27
	PdlParserNUMBER_LITERAL    = 28
	PdlParserSTRING_LITERAL    = 29
	PdlParserID                = 30
	PdlParserWS                = 31
	PdlParserPROPERTY_ID       = 32
	PdlParserESCAPED_PROP_ID   = 33
)

// PdlParser rules.
const (
	PdlParserRULE_document                   = 0
	PdlParserRULE_namespaceDeclaration       = 1
	PdlParserRULE_packageDeclaration         = 2
	PdlParserRULE_importDeclarations         = 3
	PdlParserRULE_importDeclaration          = 4
	PdlParserRULE_typeReference              = 5
	PdlParserRULE_typeDeclaration            = 6
	PdlParserRULE_namedTypeDeclaration       = 7
	PdlParserRULE_scopedNamedTypeDeclaration = 8
	PdlParserRULE_anonymousTypeDeclaration   = 9
	PdlParserRULE_typeAssignment             = 10
	PdlParserRULE_propDeclaration            = 11
	PdlParserRULE_propNameDeclaration        = 12
	PdlParserRULE_propJsonValue              = 13
	PdlParserRULE_recordDeclaration          = 14
	PdlParserRULE_enumDeclaration            = 15
	PdlParserRULE_enumSymbolDeclarations     = 16
	PdlParserRULE_enumSymbolDeclaration      = 17
	PdlParserRULE_enumSymbol                 = 18
	PdlParserRULE_typerefDeclaration         = 19
	PdlParserRULE_fixedDeclaration           = 20
	PdlParserRULE_unionDeclaration           = 21
	PdlParserRULE_unionTypeAssignments       = 22
	PdlParserRULE_unionMemberDeclaration     = 23
	PdlParserRULE_unionMemberAlias           = 24
	PdlParserRULE_arrayDeclaration           = 25
	PdlParserRULE_arrayTypeAssignments       = 26
	PdlParserRULE_mapDeclaration             = 27
	PdlParserRULE_mapTypeAssignments         = 28
	PdlParserRULE_fieldSelection             = 29
	PdlParserRULE_fieldIncludes              = 30
	PdlParserRULE_fieldDeclaration           = 31
	PdlParserRULE_fieldDefault               = 32
	PdlParserRULE_typeName                   = 33
	PdlParserRULE_identifier                 = 34
	PdlParserRULE_propName                   = 35
	PdlParserRULE_propSegment                = 36
	PdlParserRULE_schemadoc                  = 37
	PdlParserRULE_object                     = 38
	PdlParserRULE_objectEntry                = 39
	PdlParserRULE_array                      = 40
	PdlParserRULE_jsonValue                  = 41
	PdlParserRULE_string                     = 42
	PdlParserRULE_number                     = 43
	PdlParserRULE_bool                       = 44
	PdlParserRULE_nullValue                  = 45
)

// IDocumentContext is an interface to support dynamic dispatch.
type IDocumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDocumentContext differentiates from other interfaces.
	IsDocumentContext()
}

type DocumentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDocumentContext() *DocumentContext {
	var p = new(DocumentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_document
	return p
}

func (*DocumentContext) IsDocumentContext() {}

func NewDocumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DocumentContext {
	var p = new(DocumentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_document

	return p
}

func (s *DocumentContext) GetParser() antlr.Parser { return s.parser }

func (s *DocumentContext) ImportDeclarations() IImportDeclarationsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImportDeclarationsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImportDeclarationsContext)
}

func (s *DocumentContext) TypeDeclaration() ITypeDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeDeclarationContext)
}

func (s *DocumentContext) NamespaceDeclaration() INamespaceDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamespaceDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamespaceDeclarationContext)
}

func (s *DocumentContext) PackageDeclaration() IPackageDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPackageDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPackageDeclarationContext)
}

func (s *DocumentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DocumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DocumentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterDocument(s)
	}
}

func (s *DocumentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitDocument(s)
	}
}

func (p *PdlParser) Document() (localctx IDocumentContext) {
	this := p
	_ = this

	localctx = NewDocumentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, PdlParserRULE_document)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(93)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserNAMESPACE {
		{
			p.SetState(92)
			p.NamespaceDeclaration()
		}

	}
	p.SetState(96)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserPACKAGE {
		{
			p.SetState(95)
			p.PackageDeclaration()
		}

	}
	{
		p.SetState(98)
		p.ImportDeclarations()
	}
	{
		p.SetState(99)
		p.TypeDeclaration()
	}

	return localctx
}

// INamespaceDeclarationContext is an interface to support dynamic dispatch.
type INamespaceDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNamespaceDeclarationContext differentiates from other interfaces.
	IsNamespaceDeclarationContext()
}

type NamespaceDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamespaceDeclarationContext() *NamespaceDeclarationContext {
	var p = new(NamespaceDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_namespaceDeclaration
	return p
}

func (*NamespaceDeclarationContext) IsNamespaceDeclarationContext() {}

func NewNamespaceDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamespaceDeclarationContext {
	var p = new(NamespaceDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_namespaceDeclaration

	return p
}

func (s *NamespaceDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *NamespaceDeclarationContext) NAMESPACE() antlr.TerminalNode {
	return s.GetToken(PdlParserNAMESPACE, 0)
}

func (s *NamespaceDeclarationContext) TypeName() ITypeNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeNameContext)
}

func (s *NamespaceDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamespaceDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamespaceDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterNamespaceDeclaration(s)
	}
}

func (s *NamespaceDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitNamespaceDeclaration(s)
	}
}

func (p *PdlParser) NamespaceDeclaration() (localctx INamespaceDeclarationContext) {
	this := p
	_ = this

	localctx = NewNamespaceDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, PdlParserRULE_namespaceDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(101)
		p.Match(PdlParserNAMESPACE)
	}
	{
		p.SetState(102)
		p.TypeName()
	}

	return localctx
}

// IPackageDeclarationContext is an interface to support dynamic dispatch.
type IPackageDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPackageDeclarationContext differentiates from other interfaces.
	IsPackageDeclarationContext()
}

type PackageDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPackageDeclarationContext() *PackageDeclarationContext {
	var p = new(PackageDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_packageDeclaration
	return p
}

func (*PackageDeclarationContext) IsPackageDeclarationContext() {}

func NewPackageDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PackageDeclarationContext {
	var p = new(PackageDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_packageDeclaration

	return p
}

func (s *PackageDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *PackageDeclarationContext) PACKAGE() antlr.TerminalNode {
	return s.GetToken(PdlParserPACKAGE, 0)
}

func (s *PackageDeclarationContext) TypeName() ITypeNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeNameContext)
}

func (s *PackageDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PackageDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PackageDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterPackageDeclaration(s)
	}
}

func (s *PackageDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitPackageDeclaration(s)
	}
}

func (p *PdlParser) PackageDeclaration() (localctx IPackageDeclarationContext) {
	this := p
	_ = this

	localctx = NewPackageDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, PdlParserRULE_packageDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(104)
		p.Match(PdlParserPACKAGE)
	}
	{
		p.SetState(105)
		p.TypeName()
	}

	return localctx
}

// IImportDeclarationsContext is an interface to support dynamic dispatch.
type IImportDeclarationsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImportDeclarationsContext differentiates from other interfaces.
	IsImportDeclarationsContext()
}

type ImportDeclarationsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportDeclarationsContext() *ImportDeclarationsContext {
	var p = new(ImportDeclarationsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_importDeclarations
	return p
}

func (*ImportDeclarationsContext) IsImportDeclarationsContext() {}

func NewImportDeclarationsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportDeclarationsContext {
	var p = new(ImportDeclarationsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_importDeclarations

	return p
}

func (s *ImportDeclarationsContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportDeclarationsContext) AllImportDeclaration() []IImportDeclarationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IImportDeclarationContext); ok {
			len++
		}
	}

	tst := make([]IImportDeclarationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IImportDeclarationContext); ok {
			tst[i] = t.(IImportDeclarationContext)
			i++
		}
	}

	return tst
}

func (s *ImportDeclarationsContext) ImportDeclaration(i int) IImportDeclarationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImportDeclarationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImportDeclarationContext)
}

func (s *ImportDeclarationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportDeclarationsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportDeclarationsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterImportDeclarations(s)
	}
}

func (s *ImportDeclarationsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitImportDeclarations(s)
	}
}

func (p *PdlParser) ImportDeclarations() (localctx IImportDeclarationsContext) {
	this := p
	_ = this

	localctx = NewImportDeclarationsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, PdlParserRULE_importDeclarations)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(110)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PdlParserIMPORT {
		{
			p.SetState(107)
			p.ImportDeclaration()
		}

		p.SetState(112)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IImportDeclarationContext is an interface to support dynamic dispatch.
type IImportDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetType_ returns the type_ rule contexts.
	GetType_() ITypeNameContext

	// SetType_ sets the type_ rule contexts.
	SetType_(ITypeNameContext)

	// IsImportDeclarationContext differentiates from other interfaces.
	IsImportDeclarationContext()
}

type ImportDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	type_  ITypeNameContext
}

func NewEmptyImportDeclarationContext() *ImportDeclarationContext {
	var p = new(ImportDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_importDeclaration
	return p
}

func (*ImportDeclarationContext) IsImportDeclarationContext() {}

func NewImportDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportDeclarationContext {
	var p = new(ImportDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_importDeclaration

	return p
}

func (s *ImportDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportDeclarationContext) GetType_() ITypeNameContext { return s.type_ }

func (s *ImportDeclarationContext) SetType_(v ITypeNameContext) { s.type_ = v }

func (s *ImportDeclarationContext) IMPORT() antlr.TerminalNode {
	return s.GetToken(PdlParserIMPORT, 0)
}

func (s *ImportDeclarationContext) TypeName() ITypeNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeNameContext)
}

func (s *ImportDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterImportDeclaration(s)
	}
}

func (s *ImportDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitImportDeclaration(s)
	}
}

func (p *PdlParser) ImportDeclaration() (localctx IImportDeclarationContext) {
	this := p
	_ = this

	localctx = NewImportDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, PdlParserRULE_importDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(113)
		p.Match(PdlParserIMPORT)
	}
	{
		p.SetState(114)

		var _x = p.TypeName()

		localctx.(*ImportDeclarationContext).type_ = _x
	}

	return localctx
}

// ITypeReferenceContext is an interface to support dynamic dispatch.
type ITypeReferenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_typeName returns the _typeName rule contexts.
	Get_typeName() ITypeNameContext

	// Set_typeName sets the _typeName rule contexts.
	Set_typeName(ITypeNameContext)

	// GetValue returns the value attribute.
	GetValue() string

	// SetValue sets the value attribute.
	SetValue(string)

	// IsTypeReferenceContext differentiates from other interfaces.
	IsTypeReferenceContext()
}

type TypeReferenceContext struct {
	*antlr.BaseParserRuleContext
	parser    antlr.Parser
	value     string
	_typeName ITypeNameContext
}

func NewEmptyTypeReferenceContext() *TypeReferenceContext {
	var p = new(TypeReferenceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_typeReference
	return p
}

func (*TypeReferenceContext) IsTypeReferenceContext() {}

func NewTypeReferenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeReferenceContext {
	var p = new(TypeReferenceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_typeReference

	return p
}

func (s *TypeReferenceContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeReferenceContext) Get_typeName() ITypeNameContext { return s._typeName }

func (s *TypeReferenceContext) Set_typeName(v ITypeNameContext) { s._typeName = v }

func (s *TypeReferenceContext) GetValue() string { return s.value }

func (s *TypeReferenceContext) SetValue(v string) { s.value = v }

func (s *TypeReferenceContext) NULL_LITERAL() antlr.TerminalNode {
	return s.GetToken(PdlParserNULL_LITERAL, 0)
}

func (s *TypeReferenceContext) TypeName() ITypeNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeNameContext)
}

func (s *TypeReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeReferenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterTypeReference(s)
	}
}

func (s *TypeReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitTypeReference(s)
	}
}

func (p *PdlParser) TypeReference() (localctx ITypeReferenceContext) {
	this := p
	_ = this

	localctx = NewTypeReferenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, PdlParserRULE_typeReference)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(121)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PdlParserNULL_LITERAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(116)
			p.Match(PdlParserNULL_LITERAL)
		}
		localctx.(*TypeReferenceContext).SetValue("null")

	case PdlParserID:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(118)

			var _x = p.TypeName()

			localctx.(*TypeReferenceContext)._typeName = _x
		}

		localctx.(*TypeReferenceContext).SetValue(localctx.(*TypeReferenceContext).Get_typeName().GetValue())

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ITypeDeclarationContext is an interface to support dynamic dispatch.
type ITypeDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeDeclarationContext differentiates from other interfaces.
	IsTypeDeclarationContext()
}

type TypeDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeDeclarationContext() *TypeDeclarationContext {
	var p = new(TypeDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_typeDeclaration
	return p
}

func (*TypeDeclarationContext) IsTypeDeclarationContext() {}

func NewTypeDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeDeclarationContext {
	var p = new(TypeDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_typeDeclaration

	return p
}

func (s *TypeDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeDeclarationContext) ScopedNamedTypeDeclaration() IScopedNamedTypeDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IScopedNamedTypeDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IScopedNamedTypeDeclarationContext)
}

func (s *TypeDeclarationContext) NamedTypeDeclaration() INamedTypeDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamedTypeDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamedTypeDeclarationContext)
}

func (s *TypeDeclarationContext) AnonymousTypeDeclaration() IAnonymousTypeDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnonymousTypeDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAnonymousTypeDeclarationContext)
}

func (s *TypeDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterTypeDeclaration(s)
	}
}

func (s *TypeDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitTypeDeclaration(s)
	}
}

func (p *PdlParser) TypeDeclaration() (localctx ITypeDeclarationContext) {
	this := p
	_ = this

	localctx = NewTypeDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, PdlParserRULE_typeDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(123)
			p.ScopedNamedTypeDeclaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(124)
			p.NamedTypeDeclaration()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(125)
			p.AnonymousTypeDeclaration()
		}

	}

	return localctx
}

// INamedTypeDeclarationContext is an interface to support dynamic dispatch.
type INamedTypeDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetDoc returns the doc rule contexts.
	GetDoc() ISchemadocContext

	// Get_propDeclaration returns the _propDeclaration rule contexts.
	Get_propDeclaration() IPropDeclarationContext

	// SetDoc sets the doc rule contexts.
	SetDoc(ISchemadocContext)

	// Set_propDeclaration sets the _propDeclaration rule contexts.
	Set_propDeclaration(IPropDeclarationContext)

	// GetProps returns the props rule context list.
	GetProps() []IPropDeclarationContext

	// SetProps sets the props rule context list.
	SetProps([]IPropDeclarationContext)

	// IsNamedTypeDeclarationContext differentiates from other interfaces.
	IsNamedTypeDeclarationContext()
}

type NamedTypeDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser           antlr.Parser
	doc              ISchemadocContext
	_propDeclaration IPropDeclarationContext
	props            []IPropDeclarationContext
}

func NewEmptyNamedTypeDeclarationContext() *NamedTypeDeclarationContext {
	var p = new(NamedTypeDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_namedTypeDeclaration
	return p
}

func (*NamedTypeDeclarationContext) IsNamedTypeDeclarationContext() {}

func NewNamedTypeDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamedTypeDeclarationContext {
	var p = new(NamedTypeDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_namedTypeDeclaration

	return p
}

func (s *NamedTypeDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *NamedTypeDeclarationContext) GetDoc() ISchemadocContext { return s.doc }

func (s *NamedTypeDeclarationContext) Get_propDeclaration() IPropDeclarationContext {
	return s._propDeclaration
}

func (s *NamedTypeDeclarationContext) SetDoc(v ISchemadocContext) { s.doc = v }

func (s *NamedTypeDeclarationContext) Set_propDeclaration(v IPropDeclarationContext) {
	s._propDeclaration = v
}

func (s *NamedTypeDeclarationContext) GetProps() []IPropDeclarationContext { return s.props }

func (s *NamedTypeDeclarationContext) SetProps(v []IPropDeclarationContext) { s.props = v }

func (s *NamedTypeDeclarationContext) RecordDeclaration() IRecordDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRecordDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRecordDeclarationContext)
}

func (s *NamedTypeDeclarationContext) EnumDeclaration() IEnumDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnumDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnumDeclarationContext)
}

func (s *NamedTypeDeclarationContext) TyperefDeclaration() ITyperefDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITyperefDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITyperefDeclarationContext)
}

func (s *NamedTypeDeclarationContext) FixedDeclaration() IFixedDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFixedDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFixedDeclarationContext)
}

func (s *NamedTypeDeclarationContext) Schemadoc() ISchemadocContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISchemadocContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISchemadocContext)
}

func (s *NamedTypeDeclarationContext) AllPropDeclaration() []IPropDeclarationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPropDeclarationContext); ok {
			len++
		}
	}

	tst := make([]IPropDeclarationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPropDeclarationContext); ok {
			tst[i] = t.(IPropDeclarationContext)
			i++
		}
	}

	return tst
}

func (s *NamedTypeDeclarationContext) PropDeclaration(i int) IPropDeclarationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropDeclarationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropDeclarationContext)
}

func (s *NamedTypeDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamedTypeDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamedTypeDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterNamedTypeDeclaration(s)
	}
}

func (s *NamedTypeDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitNamedTypeDeclaration(s)
	}
}

func (p *PdlParser) NamedTypeDeclaration() (localctx INamedTypeDeclarationContext) {
	this := p
	_ = this

	localctx = NewNamedTypeDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, PdlParserRULE_namedTypeDeclaration)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(129)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserSCHEMADOC_COMMENT {
		{
			p.SetState(128)

			var _x = p.Schemadoc()

			localctx.(*NamedTypeDeclarationContext).doc = _x
		}

	}
	p.SetState(134)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PdlParserAT {
		{
			p.SetState(131)

			var _x = p.PropDeclaration()

			localctx.(*NamedTypeDeclarationContext)._propDeclaration = _x
		}
		localctx.(*NamedTypeDeclarationContext).props = append(localctx.(*NamedTypeDeclarationContext).props, localctx.(*NamedTypeDeclarationContext)._propDeclaration)

		p.SetState(136)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(141)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PdlParserRECORD:
		{
			p.SetState(137)
			p.RecordDeclaration()
		}

	case PdlParserENUM:
		{
			p.SetState(138)
			p.EnumDeclaration()
		}

	case PdlParserTYPEREF:
		{
			p.SetState(139)
			p.TyperefDeclaration()
		}

	case PdlParserFIXED:
		{
			p.SetState(140)
			p.FixedDeclaration()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IScopedNamedTypeDeclarationContext is an interface to support dynamic dispatch.
type IScopedNamedTypeDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsScopedNamedTypeDeclarationContext differentiates from other interfaces.
	IsScopedNamedTypeDeclarationContext()
}

type ScopedNamedTypeDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyScopedNamedTypeDeclarationContext() *ScopedNamedTypeDeclarationContext {
	var p = new(ScopedNamedTypeDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_scopedNamedTypeDeclaration
	return p
}

func (*ScopedNamedTypeDeclarationContext) IsScopedNamedTypeDeclarationContext() {}

func NewScopedNamedTypeDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ScopedNamedTypeDeclarationContext {
	var p = new(ScopedNamedTypeDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_scopedNamedTypeDeclaration

	return p
}

func (s *ScopedNamedTypeDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *ScopedNamedTypeDeclarationContext) OPEN_BRACE() antlr.TerminalNode {
	return s.GetToken(PdlParserOPEN_BRACE, 0)
}

func (s *ScopedNamedTypeDeclarationContext) NamedTypeDeclaration() INamedTypeDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamedTypeDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamedTypeDeclarationContext)
}

func (s *ScopedNamedTypeDeclarationContext) CLOSE_BRACE() antlr.TerminalNode {
	return s.GetToken(PdlParserCLOSE_BRACE, 0)
}

func (s *ScopedNamedTypeDeclarationContext) NamespaceDeclaration() INamespaceDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamespaceDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamespaceDeclarationContext)
}

func (s *ScopedNamedTypeDeclarationContext) PackageDeclaration() IPackageDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPackageDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPackageDeclarationContext)
}

func (s *ScopedNamedTypeDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScopedNamedTypeDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ScopedNamedTypeDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterScopedNamedTypeDeclaration(s)
	}
}

func (s *ScopedNamedTypeDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitScopedNamedTypeDeclaration(s)
	}
}

func (p *PdlParser) ScopedNamedTypeDeclaration() (localctx IScopedNamedTypeDeclarationContext) {
	this := p
	_ = this

	localctx = NewScopedNamedTypeDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, PdlParserRULE_scopedNamedTypeDeclaration)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(143)
		p.Match(PdlParserOPEN_BRACE)
	}
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserNAMESPACE {
		{
			p.SetState(144)
			p.NamespaceDeclaration()
		}

	}
	p.SetState(148)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserPACKAGE {
		{
			p.SetState(147)
			p.PackageDeclaration()
		}

	}
	{
		p.SetState(150)
		p.NamedTypeDeclaration()
	}
	{
		p.SetState(151)
		p.Match(PdlParserCLOSE_BRACE)
	}

	return localctx
}

// IAnonymousTypeDeclarationContext is an interface to support dynamic dispatch.
type IAnonymousTypeDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_propDeclaration returns the _propDeclaration rule contexts.
	Get_propDeclaration() IPropDeclarationContext

	// Set_propDeclaration sets the _propDeclaration rule contexts.
	Set_propDeclaration(IPropDeclarationContext)

	// GetProps returns the props rule context list.
	GetProps() []IPropDeclarationContext

	// SetProps sets the props rule context list.
	SetProps([]IPropDeclarationContext)

	// IsAnonymousTypeDeclarationContext differentiates from other interfaces.
	IsAnonymousTypeDeclarationContext()
}

type AnonymousTypeDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser           antlr.Parser
	_propDeclaration IPropDeclarationContext
	props            []IPropDeclarationContext
}

func NewEmptyAnonymousTypeDeclarationContext() *AnonymousTypeDeclarationContext {
	var p = new(AnonymousTypeDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_anonymousTypeDeclaration
	return p
}

func (*AnonymousTypeDeclarationContext) IsAnonymousTypeDeclarationContext() {}

func NewAnonymousTypeDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnonymousTypeDeclarationContext {
	var p = new(AnonymousTypeDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_anonymousTypeDeclaration

	return p
}

func (s *AnonymousTypeDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *AnonymousTypeDeclarationContext) Get_propDeclaration() IPropDeclarationContext {
	return s._propDeclaration
}

func (s *AnonymousTypeDeclarationContext) Set_propDeclaration(v IPropDeclarationContext) {
	s._propDeclaration = v
}

func (s *AnonymousTypeDeclarationContext) GetProps() []IPropDeclarationContext { return s.props }

func (s *AnonymousTypeDeclarationContext) SetProps(v []IPropDeclarationContext) { s.props = v }

func (s *AnonymousTypeDeclarationContext) UnionDeclaration() IUnionDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionDeclarationContext)
}

func (s *AnonymousTypeDeclarationContext) ArrayDeclaration() IArrayDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrayDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrayDeclarationContext)
}

func (s *AnonymousTypeDeclarationContext) MapDeclaration() IMapDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMapDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMapDeclarationContext)
}

func (s *AnonymousTypeDeclarationContext) AllPropDeclaration() []IPropDeclarationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPropDeclarationContext); ok {
			len++
		}
	}

	tst := make([]IPropDeclarationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPropDeclarationContext); ok {
			tst[i] = t.(IPropDeclarationContext)
			i++
		}
	}

	return tst
}

func (s *AnonymousTypeDeclarationContext) PropDeclaration(i int) IPropDeclarationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropDeclarationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropDeclarationContext)
}

func (s *AnonymousTypeDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnonymousTypeDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AnonymousTypeDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterAnonymousTypeDeclaration(s)
	}
}

func (s *AnonymousTypeDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitAnonymousTypeDeclaration(s)
	}
}

func (p *PdlParser) AnonymousTypeDeclaration() (localctx IAnonymousTypeDeclarationContext) {
	this := p
	_ = this

	localctx = NewAnonymousTypeDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, PdlParserRULE_anonymousTypeDeclaration)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(156)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PdlParserAT {
		{
			p.SetState(153)

			var _x = p.PropDeclaration()

			localctx.(*AnonymousTypeDeclarationContext)._propDeclaration = _x
		}
		localctx.(*AnonymousTypeDeclarationContext).props = append(localctx.(*AnonymousTypeDeclarationContext).props, localctx.(*AnonymousTypeDeclarationContext)._propDeclaration)

		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(162)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PdlParserUNION:
		{
			p.SetState(159)
			p.UnionDeclaration()
		}

	case PdlParserARRAY:
		{
			p.SetState(160)
			p.ArrayDeclaration()
		}

	case PdlParserMAP:
		{
			p.SetState(161)
			p.MapDeclaration()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ITypeAssignmentContext is an interface to support dynamic dispatch.
type ITypeAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeAssignmentContext differentiates from other interfaces.
	IsTypeAssignmentContext()
}

type TypeAssignmentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeAssignmentContext() *TypeAssignmentContext {
	var p = new(TypeAssignmentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_typeAssignment
	return p
}

func (*TypeAssignmentContext) IsTypeAssignmentContext() {}

func NewTypeAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeAssignmentContext {
	var p = new(TypeAssignmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_typeAssignment

	return p
}

func (s *TypeAssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeAssignmentContext) TypeReference() ITypeReferenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeReferenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeReferenceContext)
}

func (s *TypeAssignmentContext) TypeDeclaration() ITypeDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeDeclarationContext)
}

func (s *TypeAssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeAssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeAssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterTypeAssignment(s)
	}
}

func (s *TypeAssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitTypeAssignment(s)
	}
}

func (p *PdlParser) TypeAssignment() (localctx ITypeAssignmentContext) {
	this := p
	_ = this

	localctx = NewTypeAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, PdlParserRULE_typeAssignment)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(166)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PdlParserNULL_LITERAL, PdlParserID:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(164)
			p.TypeReference()
		}

	case PdlParserARRAY, PdlParserENUM, PdlParserFIXED, PdlParserMAP, PdlParserRECORD, PdlParserTYPEREF, PdlParserUNION, PdlParserOPEN_BRACE, PdlParserAT, PdlParserSCHEMADOC_COMMENT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(165)
			p.TypeDeclaration()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IPropDeclarationContext is an interface to support dynamic dispatch.
type IPropDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_propNameDeclaration returns the _propNameDeclaration rule contexts.
	Get_propNameDeclaration() IPropNameDeclarationContext

	// Set_propNameDeclaration sets the _propNameDeclaration rule contexts.
	Set_propNameDeclaration(IPropNameDeclarationContext)

	// GetPath returns the path attribute.
	GetPath() []string

	// SetPath sets the path attribute.
	SetPath([]string)

	// IsPropDeclarationContext differentiates from other interfaces.
	IsPropDeclarationContext()
}

type PropDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser               antlr.Parser
	path                 []string
	_propNameDeclaration IPropNameDeclarationContext
}

func NewEmptyPropDeclarationContext() *PropDeclarationContext {
	var p = new(PropDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_propDeclaration
	return p
}

func (*PropDeclarationContext) IsPropDeclarationContext() {}

func NewPropDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropDeclarationContext {
	var p = new(PropDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_propDeclaration

	return p
}

func (s *PropDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *PropDeclarationContext) Get_propNameDeclaration() IPropNameDeclarationContext {
	return s._propNameDeclaration
}

func (s *PropDeclarationContext) Set_propNameDeclaration(v IPropNameDeclarationContext) {
	s._propNameDeclaration = v
}

func (s *PropDeclarationContext) GetPath() []string { return s.path }

func (s *PropDeclarationContext) SetPath(v []string) { s.path = v }

func (s *PropDeclarationContext) PropNameDeclaration() IPropNameDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropNameDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropNameDeclarationContext)
}

func (s *PropDeclarationContext) PropJsonValue() IPropJsonValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropJsonValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropJsonValueContext)
}

func (s *PropDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterPropDeclaration(s)
	}
}

func (s *PropDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitPropDeclaration(s)
	}
}

func (p *PdlParser) PropDeclaration() (localctx IPropDeclarationContext) {
	this := p
	_ = this

	localctx = NewPropDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, PdlParserRULE_propDeclaration)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(168)

		var _x = p.PropNameDeclaration()

		localctx.(*PropDeclarationContext)._propNameDeclaration = _x
	}
	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserEQ {
		{
			p.SetState(169)
			p.PropJsonValue()
		}

	}

	localctx.(*PropDeclarationContext).SetPath(localctx.(*PropDeclarationContext).Get_propNameDeclaration().GetPath())

	return localctx
}

// IPropNameDeclarationContext is an interface to support dynamic dispatch.
type IPropNameDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_propName returns the _propName rule contexts.
	Get_propName() IPropNameContext

	// Set_propName sets the _propName rule contexts.
	Set_propName(IPropNameContext)

	// GetPath returns the path attribute.
	GetPath() []string

	// SetPath sets the path attribute.
	SetPath([]string)

	// IsPropNameDeclarationContext differentiates from other interfaces.
	IsPropNameDeclarationContext()
}

type PropNameDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser    antlr.Parser
	path      []string
	_propName IPropNameContext
}

func NewEmptyPropNameDeclarationContext() *PropNameDeclarationContext {
	var p = new(PropNameDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_propNameDeclaration
	return p
}

func (*PropNameDeclarationContext) IsPropNameDeclarationContext() {}

func NewPropNameDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropNameDeclarationContext {
	var p = new(PropNameDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_propNameDeclaration

	return p
}

func (s *PropNameDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *PropNameDeclarationContext) Get_propName() IPropNameContext { return s._propName }

func (s *PropNameDeclarationContext) Set_propName(v IPropNameContext) { s._propName = v }

func (s *PropNameDeclarationContext) GetPath() []string { return s.path }

func (s *PropNameDeclarationContext) SetPath(v []string) { s.path = v }

func (s *PropNameDeclarationContext) AT() antlr.TerminalNode {
	return s.GetToken(PdlParserAT, 0)
}

func (s *PropNameDeclarationContext) PropName() IPropNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropNameContext)
}

func (s *PropNameDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropNameDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropNameDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterPropNameDeclaration(s)
	}
}

func (s *PropNameDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitPropNameDeclaration(s)
	}
}

func (p *PdlParser) PropNameDeclaration() (localctx IPropNameDeclarationContext) {
	this := p
	_ = this

	localctx = NewPropNameDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, PdlParserRULE_propNameDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(174)
		p.Match(PdlParserAT)
	}
	{
		p.SetState(175)

		var _x = p.PropName()

		localctx.(*PropNameDeclarationContext)._propName = _x
	}

	localctx.(*PropNameDeclarationContext).SetPath(localctx.(*PropNameDeclarationContext).Get_propName().GetPath())

	return localctx
}

// IPropJsonValueContext is an interface to support dynamic dispatch.
type IPropJsonValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPropJsonValueContext differentiates from other interfaces.
	IsPropJsonValueContext()
}

type PropJsonValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPropJsonValueContext() *PropJsonValueContext {
	var p = new(PropJsonValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_propJsonValue
	return p
}

func (*PropJsonValueContext) IsPropJsonValueContext() {}

func NewPropJsonValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropJsonValueContext {
	var p = new(PropJsonValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_propJsonValue

	return p
}

func (s *PropJsonValueContext) GetParser() antlr.Parser { return s.parser }

func (s *PropJsonValueContext) EQ() antlr.TerminalNode {
	return s.GetToken(PdlParserEQ, 0)
}

func (s *PropJsonValueContext) JsonValue() IJsonValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJsonValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJsonValueContext)
}

func (s *PropJsonValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropJsonValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropJsonValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterPropJsonValue(s)
	}
}

func (s *PropJsonValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitPropJsonValue(s)
	}
}

func (p *PdlParser) PropJsonValue() (localctx IPropJsonValueContext) {
	this := p
	_ = this

	localctx = NewPropJsonValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, PdlParserRULE_propJsonValue)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(178)
		p.Match(PdlParserEQ)
	}
	{
		p.SetState(179)
		p.JsonValue()
	}

	return localctx
}

// IRecordDeclarationContext is an interface to support dynamic dispatch.
type IRecordDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_identifier returns the _identifier rule contexts.
	Get_identifier() IIdentifierContext

	// GetBeforeIncludes returns the beforeIncludes rule contexts.
	GetBeforeIncludes() IFieldIncludesContext

	// GetRecordDecl returns the recordDecl rule contexts.
	GetRecordDecl() IFieldSelectionContext

	// GetAfterIncludes returns the afterIncludes rule contexts.
	GetAfterIncludes() IFieldIncludesContext

	// Set_identifier sets the _identifier rule contexts.
	Set_identifier(IIdentifierContext)

	// SetBeforeIncludes sets the beforeIncludes rule contexts.
	SetBeforeIncludes(IFieldIncludesContext)

	// SetRecordDecl sets the recordDecl rule contexts.
	SetRecordDecl(IFieldSelectionContext)

	// SetAfterIncludes sets the afterIncludes rule contexts.
	SetAfterIncludes(IFieldIncludesContext)

	// GetName returns the name attribute.
	GetName() string

	// SetName sets the name attribute.
	SetName(string)

	// IsRecordDeclarationContext differentiates from other interfaces.
	IsRecordDeclarationContext()
}

type RecordDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser         antlr.Parser
	name           string
	_identifier    IIdentifierContext
	beforeIncludes IFieldIncludesContext
	recordDecl     IFieldSelectionContext
	afterIncludes  IFieldIncludesContext
}

func NewEmptyRecordDeclarationContext() *RecordDeclarationContext {
	var p = new(RecordDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_recordDeclaration
	return p
}

func (*RecordDeclarationContext) IsRecordDeclarationContext() {}

func NewRecordDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RecordDeclarationContext {
	var p = new(RecordDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_recordDeclaration

	return p
}

func (s *RecordDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *RecordDeclarationContext) Get_identifier() IIdentifierContext { return s._identifier }

func (s *RecordDeclarationContext) GetBeforeIncludes() IFieldIncludesContext { return s.beforeIncludes }

func (s *RecordDeclarationContext) GetRecordDecl() IFieldSelectionContext { return s.recordDecl }

func (s *RecordDeclarationContext) GetAfterIncludes() IFieldIncludesContext { return s.afterIncludes }

func (s *RecordDeclarationContext) Set_identifier(v IIdentifierContext) { s._identifier = v }

func (s *RecordDeclarationContext) SetBeforeIncludes(v IFieldIncludesContext) { s.beforeIncludes = v }

func (s *RecordDeclarationContext) SetRecordDecl(v IFieldSelectionContext) { s.recordDecl = v }

func (s *RecordDeclarationContext) SetAfterIncludes(v IFieldIncludesContext) { s.afterIncludes = v }

func (s *RecordDeclarationContext) GetName() string { return s.name }

func (s *RecordDeclarationContext) SetName(v string) { s.name = v }

func (s *RecordDeclarationContext) RECORD() antlr.TerminalNode {
	return s.GetToken(PdlParserRECORD, 0)
}

func (s *RecordDeclarationContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *RecordDeclarationContext) FieldSelection() IFieldSelectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldSelectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldSelectionContext)
}

func (s *RecordDeclarationContext) AllFieldIncludes() []IFieldIncludesContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldIncludesContext); ok {
			len++
		}
	}

	tst := make([]IFieldIncludesContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldIncludesContext); ok {
			tst[i] = t.(IFieldIncludesContext)
			i++
		}
	}

	return tst
}

func (s *RecordDeclarationContext) FieldIncludes(i int) IFieldIncludesContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldIncludesContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldIncludesContext)
}

func (s *RecordDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecordDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RecordDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterRecordDeclaration(s)
	}
}

func (s *RecordDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitRecordDeclaration(s)
	}
}

func (p *PdlParser) RecordDeclaration() (localctx IRecordDeclarationContext) {
	this := p
	_ = this

	localctx = NewRecordDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, PdlParserRULE_recordDeclaration)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(181)
		p.Match(PdlParserRECORD)
	}
	{
		p.SetState(182)

		var _x = p.Identifier()

		localctx.(*RecordDeclarationContext)._identifier = _x
	}
	p.SetState(184)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserINCLUDES {
		{
			p.SetState(183)

			var _x = p.FieldIncludes()

			localctx.(*RecordDeclarationContext).beforeIncludes = _x
		}

	}
	{
		p.SetState(186)

		var _x = p.FieldSelection()

		localctx.(*RecordDeclarationContext).recordDecl = _x
	}
	p.SetState(188)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserINCLUDES {
		{
			p.SetState(187)

			var _x = p.FieldIncludes()

			localctx.(*RecordDeclarationContext).afterIncludes = _x
		}

	}

	localctx.(*RecordDeclarationContext).SetName(localctx.(*RecordDeclarationContext).Get_identifier().GetValue())

	return localctx
}

// IEnumDeclarationContext is an interface to support dynamic dispatch.
type IEnumDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_identifier returns the _identifier rule contexts.
	Get_identifier() IIdentifierContext

	// GetEnumDecl returns the enumDecl rule contexts.
	GetEnumDecl() IEnumSymbolDeclarationsContext

	// Set_identifier sets the _identifier rule contexts.
	Set_identifier(IIdentifierContext)

	// SetEnumDecl sets the enumDecl rule contexts.
	SetEnumDecl(IEnumSymbolDeclarationsContext)

	// GetName returns the name attribute.
	GetName() string

	// SetName sets the name attribute.
	SetName(string)

	// IsEnumDeclarationContext differentiates from other interfaces.
	IsEnumDeclarationContext()
}

type EnumDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	name        string
	_identifier IIdentifierContext
	enumDecl    IEnumSymbolDeclarationsContext
}

func NewEmptyEnumDeclarationContext() *EnumDeclarationContext {
	var p = new(EnumDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_enumDeclaration
	return p
}

func (*EnumDeclarationContext) IsEnumDeclarationContext() {}

func NewEnumDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumDeclarationContext {
	var p = new(EnumDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_enumDeclaration

	return p
}

func (s *EnumDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *EnumDeclarationContext) Get_identifier() IIdentifierContext { return s._identifier }

func (s *EnumDeclarationContext) GetEnumDecl() IEnumSymbolDeclarationsContext { return s.enumDecl }

func (s *EnumDeclarationContext) Set_identifier(v IIdentifierContext) { s._identifier = v }

func (s *EnumDeclarationContext) SetEnumDecl(v IEnumSymbolDeclarationsContext) { s.enumDecl = v }

func (s *EnumDeclarationContext) GetName() string { return s.name }

func (s *EnumDeclarationContext) SetName(v string) { s.name = v }

func (s *EnumDeclarationContext) ENUM() antlr.TerminalNode {
	return s.GetToken(PdlParserENUM, 0)
}

func (s *EnumDeclarationContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *EnumDeclarationContext) EnumSymbolDeclarations() IEnumSymbolDeclarationsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnumSymbolDeclarationsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnumSymbolDeclarationsContext)
}

func (s *EnumDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EnumDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EnumDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterEnumDeclaration(s)
	}
}

func (s *EnumDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitEnumDeclaration(s)
	}
}

func (p *PdlParser) EnumDeclaration() (localctx IEnumDeclarationContext) {
	this := p
	_ = this

	localctx = NewEnumDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, PdlParserRULE_enumDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(192)
		p.Match(PdlParserENUM)
	}
	{
		p.SetState(193)

		var _x = p.Identifier()

		localctx.(*EnumDeclarationContext)._identifier = _x
	}
	{
		p.SetState(194)

		var _x = p.EnumSymbolDeclarations()

		localctx.(*EnumDeclarationContext).enumDecl = _x
	}

	localctx.(*EnumDeclarationContext).SetName(localctx.(*EnumDeclarationContext).Get_identifier().GetValue())

	return localctx
}

// IEnumSymbolDeclarationsContext is an interface to support dynamic dispatch.
type IEnumSymbolDeclarationsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_enumSymbolDeclaration returns the _enumSymbolDeclaration rule contexts.
	Get_enumSymbolDeclaration() IEnumSymbolDeclarationContext

	// Set_enumSymbolDeclaration sets the _enumSymbolDeclaration rule contexts.
	Set_enumSymbolDeclaration(IEnumSymbolDeclarationContext)

	// GetSymbolDecls returns the symbolDecls rule context list.
	GetSymbolDecls() []IEnumSymbolDeclarationContext

	// SetSymbolDecls sets the symbolDecls rule context list.
	SetSymbolDecls([]IEnumSymbolDeclarationContext)

	// IsEnumSymbolDeclarationsContext differentiates from other interfaces.
	IsEnumSymbolDeclarationsContext()
}

type EnumSymbolDeclarationsContext struct {
	*antlr.BaseParserRuleContext
	parser                 antlr.Parser
	_enumSymbolDeclaration IEnumSymbolDeclarationContext
	symbolDecls            []IEnumSymbolDeclarationContext
}

func NewEmptyEnumSymbolDeclarationsContext() *EnumSymbolDeclarationsContext {
	var p = new(EnumSymbolDeclarationsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_enumSymbolDeclarations
	return p
}

func (*EnumSymbolDeclarationsContext) IsEnumSymbolDeclarationsContext() {}

func NewEnumSymbolDeclarationsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumSymbolDeclarationsContext {
	var p = new(EnumSymbolDeclarationsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_enumSymbolDeclarations

	return p
}

func (s *EnumSymbolDeclarationsContext) GetParser() antlr.Parser { return s.parser }

func (s *EnumSymbolDeclarationsContext) Get_enumSymbolDeclaration() IEnumSymbolDeclarationContext {
	return s._enumSymbolDeclaration
}

func (s *EnumSymbolDeclarationsContext) Set_enumSymbolDeclaration(v IEnumSymbolDeclarationContext) {
	s._enumSymbolDeclaration = v
}

func (s *EnumSymbolDeclarationsContext) GetSymbolDecls() []IEnumSymbolDeclarationContext {
	return s.symbolDecls
}

func (s *EnumSymbolDeclarationsContext) SetSymbolDecls(v []IEnumSymbolDeclarationContext) {
	s.symbolDecls = v
}

func (s *EnumSymbolDeclarationsContext) OPEN_BRACE() antlr.TerminalNode {
	return s.GetToken(PdlParserOPEN_BRACE, 0)
}

func (s *EnumSymbolDeclarationsContext) CLOSE_BRACE() antlr.TerminalNode {
	return s.GetToken(PdlParserCLOSE_BRACE, 0)
}

func (s *EnumSymbolDeclarationsContext) AllEnumSymbolDeclaration() []IEnumSymbolDeclarationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IEnumSymbolDeclarationContext); ok {
			len++
		}
	}

	tst := make([]IEnumSymbolDeclarationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IEnumSymbolDeclarationContext); ok {
			tst[i] = t.(IEnumSymbolDeclarationContext)
			i++
		}
	}

	return tst
}

func (s *EnumSymbolDeclarationsContext) EnumSymbolDeclaration(i int) IEnumSymbolDeclarationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnumSymbolDeclarationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnumSymbolDeclarationContext)
}

func (s *EnumSymbolDeclarationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EnumSymbolDeclarationsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EnumSymbolDeclarationsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterEnumSymbolDeclarations(s)
	}
}

func (s *EnumSymbolDeclarationsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitEnumSymbolDeclarations(s)
	}
}

func (p *PdlParser) EnumSymbolDeclarations() (localctx IEnumSymbolDeclarationsContext) {
	this := p
	_ = this

	localctx = NewEnumSymbolDeclarationsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, PdlParserRULE_enumSymbolDeclarations)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(197)
		p.Match(PdlParserOPEN_BRACE)
	}
	p.SetState(201)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1107820544) != 0 {
		{
			p.SetState(198)

			var _x = p.EnumSymbolDeclaration()

			localctx.(*EnumSymbolDeclarationsContext)._enumSymbolDeclaration = _x
		}
		localctx.(*EnumSymbolDeclarationsContext).symbolDecls = append(localctx.(*EnumSymbolDeclarationsContext).symbolDecls, localctx.(*EnumSymbolDeclarationsContext)._enumSymbolDeclaration)

		p.SetState(203)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(204)
		p.Match(PdlParserCLOSE_BRACE)
	}

	return localctx
}

// IEnumSymbolDeclarationContext is an interface to support dynamic dispatch.
type IEnumSymbolDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetDoc returns the doc rule contexts.
	GetDoc() ISchemadocContext

	// Get_propDeclaration returns the _propDeclaration rule contexts.
	Get_propDeclaration() IPropDeclarationContext

	// GetSymbol returns the symbol rule contexts.
	GetSymbol() IEnumSymbolContext

	// SetDoc sets the doc rule contexts.
	SetDoc(ISchemadocContext)

	// Set_propDeclaration sets the _propDeclaration rule contexts.
	Set_propDeclaration(IPropDeclarationContext)

	// SetSymbol sets the symbol rule contexts.
	SetSymbol(IEnumSymbolContext)

	// GetProps returns the props rule context list.
	GetProps() []IPropDeclarationContext

	// SetProps sets the props rule context list.
	SetProps([]IPropDeclarationContext)

	// IsEnumSymbolDeclarationContext differentiates from other interfaces.
	IsEnumSymbolDeclarationContext()
}

type EnumSymbolDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser           antlr.Parser
	doc              ISchemadocContext
	_propDeclaration IPropDeclarationContext
	props            []IPropDeclarationContext
	symbol           IEnumSymbolContext
}

func NewEmptyEnumSymbolDeclarationContext() *EnumSymbolDeclarationContext {
	var p = new(EnumSymbolDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_enumSymbolDeclaration
	return p
}

func (*EnumSymbolDeclarationContext) IsEnumSymbolDeclarationContext() {}

func NewEnumSymbolDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumSymbolDeclarationContext {
	var p = new(EnumSymbolDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_enumSymbolDeclaration

	return p
}

func (s *EnumSymbolDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *EnumSymbolDeclarationContext) GetDoc() ISchemadocContext { return s.doc }

func (s *EnumSymbolDeclarationContext) Get_propDeclaration() IPropDeclarationContext {
	return s._propDeclaration
}

func (s *EnumSymbolDeclarationContext) GetSymbol() IEnumSymbolContext { return s.symbol }

func (s *EnumSymbolDeclarationContext) SetDoc(v ISchemadocContext) { s.doc = v }

func (s *EnumSymbolDeclarationContext) Set_propDeclaration(v IPropDeclarationContext) {
	s._propDeclaration = v
}

func (s *EnumSymbolDeclarationContext) SetSymbol(v IEnumSymbolContext) { s.symbol = v }

func (s *EnumSymbolDeclarationContext) GetProps() []IPropDeclarationContext { return s.props }

func (s *EnumSymbolDeclarationContext) SetProps(v []IPropDeclarationContext) { s.props = v }

func (s *EnumSymbolDeclarationContext) EnumSymbol() IEnumSymbolContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnumSymbolContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnumSymbolContext)
}

func (s *EnumSymbolDeclarationContext) Schemadoc() ISchemadocContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISchemadocContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISchemadocContext)
}

func (s *EnumSymbolDeclarationContext) AllPropDeclaration() []IPropDeclarationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPropDeclarationContext); ok {
			len++
		}
	}

	tst := make([]IPropDeclarationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPropDeclarationContext); ok {
			tst[i] = t.(IPropDeclarationContext)
			i++
		}
	}

	return tst
}

func (s *EnumSymbolDeclarationContext) PropDeclaration(i int) IPropDeclarationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropDeclarationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropDeclarationContext)
}

func (s *EnumSymbolDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EnumSymbolDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EnumSymbolDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterEnumSymbolDeclaration(s)
	}
}

func (s *EnumSymbolDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitEnumSymbolDeclaration(s)
	}
}

func (p *PdlParser) EnumSymbolDeclaration() (localctx IEnumSymbolDeclarationContext) {
	this := p
	_ = this

	localctx = NewEnumSymbolDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, PdlParserRULE_enumSymbolDeclaration)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(207)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserSCHEMADOC_COMMENT {
		{
			p.SetState(206)

			var _x = p.Schemadoc()

			localctx.(*EnumSymbolDeclarationContext).doc = _x
		}

	}
	p.SetState(212)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PdlParserAT {
		{
			p.SetState(209)

			var _x = p.PropDeclaration()

			localctx.(*EnumSymbolDeclarationContext)._propDeclaration = _x
		}
		localctx.(*EnumSymbolDeclarationContext).props = append(localctx.(*EnumSymbolDeclarationContext).props, localctx.(*EnumSymbolDeclarationContext)._propDeclaration)

		p.SetState(214)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(215)

		var _x = p.EnumSymbol()

		localctx.(*EnumSymbolDeclarationContext).symbol = _x
	}

	return localctx
}

// IEnumSymbolContext is an interface to support dynamic dispatch.
type IEnumSymbolContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_identifier returns the _identifier rule contexts.
	Get_identifier() IIdentifierContext

	// Set_identifier sets the _identifier rule contexts.
	Set_identifier(IIdentifierContext)

	// GetValue returns the value attribute.
	GetValue() string

	// SetValue sets the value attribute.
	SetValue(string)

	// IsEnumSymbolContext differentiates from other interfaces.
	IsEnumSymbolContext()
}

type EnumSymbolContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	value       string
	_identifier IIdentifierContext
}

func NewEmptyEnumSymbolContext() *EnumSymbolContext {
	var p = new(EnumSymbolContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_enumSymbol
	return p
}

func (*EnumSymbolContext) IsEnumSymbolContext() {}

func NewEnumSymbolContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumSymbolContext {
	var p = new(EnumSymbolContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_enumSymbol

	return p
}

func (s *EnumSymbolContext) GetParser() antlr.Parser { return s.parser }

func (s *EnumSymbolContext) Get_identifier() IIdentifierContext { return s._identifier }

func (s *EnumSymbolContext) Set_identifier(v IIdentifierContext) { s._identifier = v }

func (s *EnumSymbolContext) GetValue() string { return s.value }

func (s *EnumSymbolContext) SetValue(v string) { s.value = v }

func (s *EnumSymbolContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *EnumSymbolContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EnumSymbolContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EnumSymbolContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterEnumSymbol(s)
	}
}

func (s *EnumSymbolContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitEnumSymbol(s)
	}
}

func (p *PdlParser) EnumSymbol() (localctx IEnumSymbolContext) {
	this := p
	_ = this

	localctx = NewEnumSymbolContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, PdlParserRULE_enumSymbol)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(217)

		var _x = p.Identifier()

		localctx.(*EnumSymbolContext)._identifier = _x
	}

	localctx.(*EnumSymbolContext).SetValue(localctx.(*EnumSymbolContext).Get_identifier().GetValue())

	return localctx
}

// ITyperefDeclarationContext is an interface to support dynamic dispatch.
type ITyperefDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_identifier returns the _identifier rule contexts.
	Get_identifier() IIdentifierContext

	// GetRef returns the ref rule contexts.
	GetRef() ITypeAssignmentContext

	// Set_identifier sets the _identifier rule contexts.
	Set_identifier(IIdentifierContext)

	// SetRef sets the ref rule contexts.
	SetRef(ITypeAssignmentContext)

	// GetName returns the name attribute.
	GetName() string

	// SetName sets the name attribute.
	SetName(string)

	// IsTyperefDeclarationContext differentiates from other interfaces.
	IsTyperefDeclarationContext()
}

type TyperefDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	name        string
	_identifier IIdentifierContext
	ref         ITypeAssignmentContext
}

func NewEmptyTyperefDeclarationContext() *TyperefDeclarationContext {
	var p = new(TyperefDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_typerefDeclaration
	return p
}

func (*TyperefDeclarationContext) IsTyperefDeclarationContext() {}

func NewTyperefDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TyperefDeclarationContext {
	var p = new(TyperefDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_typerefDeclaration

	return p
}

func (s *TyperefDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *TyperefDeclarationContext) Get_identifier() IIdentifierContext { return s._identifier }

func (s *TyperefDeclarationContext) GetRef() ITypeAssignmentContext { return s.ref }

func (s *TyperefDeclarationContext) Set_identifier(v IIdentifierContext) { s._identifier = v }

func (s *TyperefDeclarationContext) SetRef(v ITypeAssignmentContext) { s.ref = v }

func (s *TyperefDeclarationContext) GetName() string { return s.name }

func (s *TyperefDeclarationContext) SetName(v string) { s.name = v }

func (s *TyperefDeclarationContext) TYPEREF() antlr.TerminalNode {
	return s.GetToken(PdlParserTYPEREF, 0)
}

func (s *TyperefDeclarationContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *TyperefDeclarationContext) EQ() antlr.TerminalNode {
	return s.GetToken(PdlParserEQ, 0)
}

func (s *TyperefDeclarationContext) TypeAssignment() ITypeAssignmentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeAssignmentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeAssignmentContext)
}

func (s *TyperefDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TyperefDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TyperefDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterTyperefDeclaration(s)
	}
}

func (s *TyperefDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitTyperefDeclaration(s)
	}
}

func (p *PdlParser) TyperefDeclaration() (localctx ITyperefDeclarationContext) {
	this := p
	_ = this

	localctx = NewTyperefDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, PdlParserRULE_typerefDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(220)
		p.Match(PdlParserTYPEREF)
	}
	{
		p.SetState(221)

		var _x = p.Identifier()

		localctx.(*TyperefDeclarationContext)._identifier = _x
	}
	{
		p.SetState(222)
		p.Match(PdlParserEQ)
	}
	{
		p.SetState(223)

		var _x = p.TypeAssignment()

		localctx.(*TyperefDeclarationContext).ref = _x
	}

	localctx.(*TyperefDeclarationContext).SetName(localctx.(*TyperefDeclarationContext).Get_identifier().GetValue())

	return localctx
}

// IFixedDeclarationContext is an interface to support dynamic dispatch.
type IFixedDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetSizeStr returns the sizeStr token.
	GetSizeStr() antlr.Token

	// SetSizeStr sets the sizeStr token.
	SetSizeStr(antlr.Token)

	// Get_identifier returns the _identifier rule contexts.
	Get_identifier() IIdentifierContext

	// Set_identifier sets the _identifier rule contexts.
	Set_identifier(IIdentifierContext)

	// GetName returns the name attribute.
	GetName() string

	// GetSize returns the size attribute.
	GetSize() int

	// SetName sets the name attribute.
	SetName(string)

	// SetSize sets the size attribute.
	SetSize(int)

	// IsFixedDeclarationContext differentiates from other interfaces.
	IsFixedDeclarationContext()
}

type FixedDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	name        string
	size        int
	_identifier IIdentifierContext
	sizeStr     antlr.Token
}

func NewEmptyFixedDeclarationContext() *FixedDeclarationContext {
	var p = new(FixedDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_fixedDeclaration
	return p
}

func (*FixedDeclarationContext) IsFixedDeclarationContext() {}

func NewFixedDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FixedDeclarationContext {
	var p = new(FixedDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_fixedDeclaration

	return p
}

func (s *FixedDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *FixedDeclarationContext) GetSizeStr() antlr.Token { return s.sizeStr }

func (s *FixedDeclarationContext) SetSizeStr(v antlr.Token) { s.sizeStr = v }

func (s *FixedDeclarationContext) Get_identifier() IIdentifierContext { return s._identifier }

func (s *FixedDeclarationContext) Set_identifier(v IIdentifierContext) { s._identifier = v }

func (s *FixedDeclarationContext) GetName() string { return s.name }

func (s *FixedDeclarationContext) GetSize() int { return s.size }

func (s *FixedDeclarationContext) SetName(v string) { s.name = v }

func (s *FixedDeclarationContext) SetSize(v int) { s.size = v }

func (s *FixedDeclarationContext) FIXED() antlr.TerminalNode {
	return s.GetToken(PdlParserFIXED, 0)
}

func (s *FixedDeclarationContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *FixedDeclarationContext) NUMBER_LITERAL() antlr.TerminalNode {
	return s.GetToken(PdlParserNUMBER_LITERAL, 0)
}

func (s *FixedDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FixedDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FixedDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterFixedDeclaration(s)
	}
}

func (s *FixedDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitFixedDeclaration(s)
	}
}

func (p *PdlParser) FixedDeclaration() (localctx IFixedDeclarationContext) {
	this := p
	_ = this

	localctx = NewFixedDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, PdlParserRULE_fixedDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(226)
		p.Match(PdlParserFIXED)
	}
	{
		p.SetState(227)

		var _x = p.Identifier()

		localctx.(*FixedDeclarationContext)._identifier = _x
	}
	{
		p.SetState(228)

		var _m = p.Match(PdlParserNUMBER_LITERAL)

		localctx.(*FixedDeclarationContext).sizeStr = _m
	}

	localctx.(*FixedDeclarationContext).SetName(localctx.(*FixedDeclarationContext).Get_identifier().GetValue())
	localctx.(*FixedDeclarationContext).SetSize((func() int {
		if localctx.(*FixedDeclarationContext).GetSizeStr() == nil {
			return 0
		} else {
			i, _ := strconv.Atoi(localctx.(*FixedDeclarationContext).GetSizeStr().GetText())
			return i
		}
	}()))

	return localctx
}

// IUnionDeclarationContext is an interface to support dynamic dispatch.
type IUnionDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetTypeParams returns the typeParams rule contexts.
	GetTypeParams() IUnionTypeAssignmentsContext

	// SetTypeParams sets the typeParams rule contexts.
	SetTypeParams(IUnionTypeAssignmentsContext)

	// IsUnionDeclarationContext differentiates from other interfaces.
	IsUnionDeclarationContext()
}

type UnionDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser     antlr.Parser
	typeParams IUnionTypeAssignmentsContext
}

func NewEmptyUnionDeclarationContext() *UnionDeclarationContext {
	var p = new(UnionDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_unionDeclaration
	return p
}

func (*UnionDeclarationContext) IsUnionDeclarationContext() {}

func NewUnionDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionDeclarationContext {
	var p = new(UnionDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_unionDeclaration

	return p
}

func (s *UnionDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *UnionDeclarationContext) GetTypeParams() IUnionTypeAssignmentsContext { return s.typeParams }

func (s *UnionDeclarationContext) SetTypeParams(v IUnionTypeAssignmentsContext) { s.typeParams = v }

func (s *UnionDeclarationContext) UNION() antlr.TerminalNode {
	return s.GetToken(PdlParserUNION, 0)
}

func (s *UnionDeclarationContext) UnionTypeAssignments() IUnionTypeAssignmentsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionTypeAssignmentsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionTypeAssignmentsContext)
}

func (s *UnionDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnionDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterUnionDeclaration(s)
	}
}

func (s *UnionDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitUnionDeclaration(s)
	}
}

func (p *PdlParser) UnionDeclaration() (localctx IUnionDeclarationContext) {
	this := p
	_ = this

	localctx = NewUnionDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, PdlParserRULE_unionDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(231)
		p.Match(PdlParserUNION)
	}
	{
		p.SetState(232)

		var _x = p.UnionTypeAssignments()

		localctx.(*UnionDeclarationContext).typeParams = _x
	}

	return localctx
}

// IUnionTypeAssignmentsContext is an interface to support dynamic dispatch.
type IUnionTypeAssignmentsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_unionMemberDeclaration returns the _unionMemberDeclaration rule contexts.
	Get_unionMemberDeclaration() IUnionMemberDeclarationContext

	// Set_unionMemberDeclaration sets the _unionMemberDeclaration rule contexts.
	Set_unionMemberDeclaration(IUnionMemberDeclarationContext)

	// GetMembers returns the members rule context list.
	GetMembers() []IUnionMemberDeclarationContext

	// SetMembers sets the members rule context list.
	SetMembers([]IUnionMemberDeclarationContext)

	// IsUnionTypeAssignmentsContext differentiates from other interfaces.
	IsUnionTypeAssignmentsContext()
}

type UnionTypeAssignmentsContext struct {
	*antlr.BaseParserRuleContext
	parser                  antlr.Parser
	_unionMemberDeclaration IUnionMemberDeclarationContext
	members                 []IUnionMemberDeclarationContext
}

func NewEmptyUnionTypeAssignmentsContext() *UnionTypeAssignmentsContext {
	var p = new(UnionTypeAssignmentsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_unionTypeAssignments
	return p
}

func (*UnionTypeAssignmentsContext) IsUnionTypeAssignmentsContext() {}

func NewUnionTypeAssignmentsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionTypeAssignmentsContext {
	var p = new(UnionTypeAssignmentsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_unionTypeAssignments

	return p
}

func (s *UnionTypeAssignmentsContext) GetParser() antlr.Parser { return s.parser }

func (s *UnionTypeAssignmentsContext) Get_unionMemberDeclaration() IUnionMemberDeclarationContext {
	return s._unionMemberDeclaration
}

func (s *UnionTypeAssignmentsContext) Set_unionMemberDeclaration(v IUnionMemberDeclarationContext) {
	s._unionMemberDeclaration = v
}

func (s *UnionTypeAssignmentsContext) GetMembers() []IUnionMemberDeclarationContext { return s.members }

func (s *UnionTypeAssignmentsContext) SetMembers(v []IUnionMemberDeclarationContext) { s.members = v }

func (s *UnionTypeAssignmentsContext) OPEN_BRACKET() antlr.TerminalNode {
	return s.GetToken(PdlParserOPEN_BRACKET, 0)
}

func (s *UnionTypeAssignmentsContext) CLOSE_BRACKET() antlr.TerminalNode {
	return s.GetToken(PdlParserCLOSE_BRACKET, 0)
}

func (s *UnionTypeAssignmentsContext) AllUnionMemberDeclaration() []IUnionMemberDeclarationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IUnionMemberDeclarationContext); ok {
			len++
		}
	}

	tst := make([]IUnionMemberDeclarationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IUnionMemberDeclarationContext); ok {
			tst[i] = t.(IUnionMemberDeclarationContext)
			i++
		}
	}

	return tst
}

func (s *UnionTypeAssignmentsContext) UnionMemberDeclaration(i int) IUnionMemberDeclarationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionMemberDeclarationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionMemberDeclarationContext)
}

func (s *UnionTypeAssignmentsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionTypeAssignmentsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnionTypeAssignmentsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterUnionTypeAssignments(s)
	}
}

func (s *UnionTypeAssignmentsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitUnionTypeAssignments(s)
	}
}

func (p *PdlParser) UnionTypeAssignments() (localctx IUnionTypeAssignmentsContext) {
	this := p
	_ = this

	localctx = NewUnionTypeAssignmentsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, PdlParserRULE_unionTypeAssignments)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(234)
		p.Match(PdlParserOPEN_BRACKET)
	}
	p.SetState(238)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1124634254) != 0 {
		{
			p.SetState(235)

			var _x = p.UnionMemberDeclaration()

			localctx.(*UnionTypeAssignmentsContext)._unionMemberDeclaration = _x
		}
		localctx.(*UnionTypeAssignmentsContext).members = append(localctx.(*UnionTypeAssignmentsContext).members, localctx.(*UnionTypeAssignmentsContext)._unionMemberDeclaration)

		p.SetState(240)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(241)
		p.Match(PdlParserCLOSE_BRACKET)
	}

	return localctx
}

// IUnionMemberDeclarationContext is an interface to support dynamic dispatch.
type IUnionMemberDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetAlias returns the alias rule contexts.
	GetAlias() IUnionMemberAliasContext

	// GetMember returns the member rule contexts.
	GetMember() ITypeAssignmentContext

	// SetAlias sets the alias rule contexts.
	SetAlias(IUnionMemberAliasContext)

	// SetMember sets the member rule contexts.
	SetMember(ITypeAssignmentContext)

	// IsUnionMemberDeclarationContext differentiates from other interfaces.
	IsUnionMemberDeclarationContext()
}

type UnionMemberDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	alias  IUnionMemberAliasContext
	member ITypeAssignmentContext
}

func NewEmptyUnionMemberDeclarationContext() *UnionMemberDeclarationContext {
	var p = new(UnionMemberDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_unionMemberDeclaration
	return p
}

func (*UnionMemberDeclarationContext) IsUnionMemberDeclarationContext() {}

func NewUnionMemberDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionMemberDeclarationContext {
	var p = new(UnionMemberDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_unionMemberDeclaration

	return p
}

func (s *UnionMemberDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *UnionMemberDeclarationContext) GetAlias() IUnionMemberAliasContext { return s.alias }

func (s *UnionMemberDeclarationContext) GetMember() ITypeAssignmentContext { return s.member }

func (s *UnionMemberDeclarationContext) SetAlias(v IUnionMemberAliasContext) { s.alias = v }

func (s *UnionMemberDeclarationContext) SetMember(v ITypeAssignmentContext) { s.member = v }

func (s *UnionMemberDeclarationContext) TypeAssignment() ITypeAssignmentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeAssignmentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeAssignmentContext)
}

func (s *UnionMemberDeclarationContext) UnionMemberAlias() IUnionMemberAliasContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionMemberAliasContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionMemberAliasContext)
}

func (s *UnionMemberDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionMemberDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnionMemberDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterUnionMemberDeclaration(s)
	}
}

func (s *UnionMemberDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitUnionMemberDeclaration(s)
	}
}

func (p *PdlParser) UnionMemberDeclaration() (localctx IUnionMemberDeclarationContext) {
	this := p
	_ = this

	localctx = NewUnionMemberDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, PdlParserRULE_unionMemberDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(244)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(243)

			var _x = p.UnionMemberAlias()

			localctx.(*UnionMemberDeclarationContext).alias = _x
		}

	}
	{
		p.SetState(246)

		var _x = p.TypeAssignment()

		localctx.(*UnionMemberDeclarationContext).member = _x
	}

	return localctx
}

// IUnionMemberAliasContext is an interface to support dynamic dispatch.
type IUnionMemberAliasContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetDoc returns the doc rule contexts.
	GetDoc() ISchemadocContext

	// Get_propDeclaration returns the _propDeclaration rule contexts.
	Get_propDeclaration() IPropDeclarationContext

	// GetName returns the name rule contexts.
	GetName() IIdentifierContext

	// SetDoc sets the doc rule contexts.
	SetDoc(ISchemadocContext)

	// Set_propDeclaration sets the _propDeclaration rule contexts.
	Set_propDeclaration(IPropDeclarationContext)

	// SetName sets the name rule contexts.
	SetName(IIdentifierContext)

	// GetProps returns the props rule context list.
	GetProps() []IPropDeclarationContext

	// SetProps sets the props rule context list.
	SetProps([]IPropDeclarationContext)

	// IsUnionMemberAliasContext differentiates from other interfaces.
	IsUnionMemberAliasContext()
}

type UnionMemberAliasContext struct {
	*antlr.BaseParserRuleContext
	parser           antlr.Parser
	doc              ISchemadocContext
	_propDeclaration IPropDeclarationContext
	props            []IPropDeclarationContext
	name             IIdentifierContext
}

func NewEmptyUnionMemberAliasContext() *UnionMemberAliasContext {
	var p = new(UnionMemberAliasContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_unionMemberAlias
	return p
}

func (*UnionMemberAliasContext) IsUnionMemberAliasContext() {}

func NewUnionMemberAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionMemberAliasContext {
	var p = new(UnionMemberAliasContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_unionMemberAlias

	return p
}

func (s *UnionMemberAliasContext) GetParser() antlr.Parser { return s.parser }

func (s *UnionMemberAliasContext) GetDoc() ISchemadocContext { return s.doc }

func (s *UnionMemberAliasContext) Get_propDeclaration() IPropDeclarationContext {
	return s._propDeclaration
}

func (s *UnionMemberAliasContext) GetName() IIdentifierContext { return s.name }

func (s *UnionMemberAliasContext) SetDoc(v ISchemadocContext) { s.doc = v }

func (s *UnionMemberAliasContext) Set_propDeclaration(v IPropDeclarationContext) {
	s._propDeclaration = v
}

func (s *UnionMemberAliasContext) SetName(v IIdentifierContext) { s.name = v }

func (s *UnionMemberAliasContext) GetProps() []IPropDeclarationContext { return s.props }

func (s *UnionMemberAliasContext) SetProps(v []IPropDeclarationContext) { s.props = v }

func (s *UnionMemberAliasContext) COLON() antlr.TerminalNode {
	return s.GetToken(PdlParserCOLON, 0)
}

func (s *UnionMemberAliasContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *UnionMemberAliasContext) Schemadoc() ISchemadocContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISchemadocContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISchemadocContext)
}

func (s *UnionMemberAliasContext) AllPropDeclaration() []IPropDeclarationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPropDeclarationContext); ok {
			len++
		}
	}

	tst := make([]IPropDeclarationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPropDeclarationContext); ok {
			tst[i] = t.(IPropDeclarationContext)
			i++
		}
	}

	return tst
}

func (s *UnionMemberAliasContext) PropDeclaration(i int) IPropDeclarationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropDeclarationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropDeclarationContext)
}

func (s *UnionMemberAliasContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionMemberAliasContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnionMemberAliasContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterUnionMemberAlias(s)
	}
}

func (s *UnionMemberAliasContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitUnionMemberAlias(s)
	}
}

func (p *PdlParser) UnionMemberAlias() (localctx IUnionMemberAliasContext) {
	this := p
	_ = this

	localctx = NewUnionMemberAliasContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, PdlParserRULE_unionMemberAlias)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(249)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserSCHEMADOC_COMMENT {
		{
			p.SetState(248)

			var _x = p.Schemadoc()

			localctx.(*UnionMemberAliasContext).doc = _x
		}

	}
	p.SetState(254)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PdlParserAT {
		{
			p.SetState(251)

			var _x = p.PropDeclaration()

			localctx.(*UnionMemberAliasContext)._propDeclaration = _x
		}
		localctx.(*UnionMemberAliasContext).props = append(localctx.(*UnionMemberAliasContext).props, localctx.(*UnionMemberAliasContext)._propDeclaration)

		p.SetState(256)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(257)

		var _x = p.Identifier()

		localctx.(*UnionMemberAliasContext).name = _x
	}
	{
		p.SetState(258)
		p.Match(PdlParserCOLON)
	}

	return localctx
}

// IArrayDeclarationContext is an interface to support dynamic dispatch.
type IArrayDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetTypeParams returns the typeParams rule contexts.
	GetTypeParams() IArrayTypeAssignmentsContext

	// SetTypeParams sets the typeParams rule contexts.
	SetTypeParams(IArrayTypeAssignmentsContext)

	// IsArrayDeclarationContext differentiates from other interfaces.
	IsArrayDeclarationContext()
}

type ArrayDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser     antlr.Parser
	typeParams IArrayTypeAssignmentsContext
}

func NewEmptyArrayDeclarationContext() *ArrayDeclarationContext {
	var p = new(ArrayDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_arrayDeclaration
	return p
}

func (*ArrayDeclarationContext) IsArrayDeclarationContext() {}

func NewArrayDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayDeclarationContext {
	var p = new(ArrayDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_arrayDeclaration

	return p
}

func (s *ArrayDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayDeclarationContext) GetTypeParams() IArrayTypeAssignmentsContext { return s.typeParams }

func (s *ArrayDeclarationContext) SetTypeParams(v IArrayTypeAssignmentsContext) { s.typeParams = v }

func (s *ArrayDeclarationContext) ARRAY() antlr.TerminalNode {
	return s.GetToken(PdlParserARRAY, 0)
}

func (s *ArrayDeclarationContext) ArrayTypeAssignments() IArrayTypeAssignmentsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrayTypeAssignmentsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrayTypeAssignmentsContext)
}

func (s *ArrayDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterArrayDeclaration(s)
	}
}

func (s *ArrayDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitArrayDeclaration(s)
	}
}

func (p *PdlParser) ArrayDeclaration() (localctx IArrayDeclarationContext) {
	this := p
	_ = this

	localctx = NewArrayDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, PdlParserRULE_arrayDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(260)
		p.Match(PdlParserARRAY)
	}
	{
		p.SetState(261)

		var _x = p.ArrayTypeAssignments()

		localctx.(*ArrayDeclarationContext).typeParams = _x
	}

	return localctx
}

// IArrayTypeAssignmentsContext is an interface to support dynamic dispatch.
type IArrayTypeAssignmentsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItems returns the items rule contexts.
	GetItems() ITypeAssignmentContext

	// SetItems sets the items rule contexts.
	SetItems(ITypeAssignmentContext)

	// IsArrayTypeAssignmentsContext differentiates from other interfaces.
	IsArrayTypeAssignmentsContext()
}

type ArrayTypeAssignmentsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	items  ITypeAssignmentContext
}

func NewEmptyArrayTypeAssignmentsContext() *ArrayTypeAssignmentsContext {
	var p = new(ArrayTypeAssignmentsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_arrayTypeAssignments
	return p
}

func (*ArrayTypeAssignmentsContext) IsArrayTypeAssignmentsContext() {}

func NewArrayTypeAssignmentsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayTypeAssignmentsContext {
	var p = new(ArrayTypeAssignmentsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_arrayTypeAssignments

	return p
}

func (s *ArrayTypeAssignmentsContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayTypeAssignmentsContext) GetItems() ITypeAssignmentContext { return s.items }

func (s *ArrayTypeAssignmentsContext) SetItems(v ITypeAssignmentContext) { s.items = v }

func (s *ArrayTypeAssignmentsContext) OPEN_BRACKET() antlr.TerminalNode {
	return s.GetToken(PdlParserOPEN_BRACKET, 0)
}

func (s *ArrayTypeAssignmentsContext) CLOSE_BRACKET() antlr.TerminalNode {
	return s.GetToken(PdlParserCLOSE_BRACKET, 0)
}

func (s *ArrayTypeAssignmentsContext) TypeAssignment() ITypeAssignmentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeAssignmentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeAssignmentContext)
}

func (s *ArrayTypeAssignmentsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayTypeAssignmentsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayTypeAssignmentsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterArrayTypeAssignments(s)
	}
}

func (s *ArrayTypeAssignmentsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitArrayTypeAssignments(s)
	}
}

func (p *PdlParser) ArrayTypeAssignments() (localctx IArrayTypeAssignmentsContext) {
	this := p
	_ = this

	localctx = NewArrayTypeAssignmentsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, PdlParserRULE_arrayTypeAssignments)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(263)
		p.Match(PdlParserOPEN_BRACKET)
	}
	{
		p.SetState(264)

		var _x = p.TypeAssignment()

		localctx.(*ArrayTypeAssignmentsContext).items = _x
	}
	{
		p.SetState(265)
		p.Match(PdlParserCLOSE_BRACKET)
	}

	return localctx
}

// IMapDeclarationContext is an interface to support dynamic dispatch.
type IMapDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetTypeParams returns the typeParams rule contexts.
	GetTypeParams() IMapTypeAssignmentsContext

	// SetTypeParams sets the typeParams rule contexts.
	SetTypeParams(IMapTypeAssignmentsContext)

	// IsMapDeclarationContext differentiates from other interfaces.
	IsMapDeclarationContext()
}

type MapDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser     antlr.Parser
	typeParams IMapTypeAssignmentsContext
}

func NewEmptyMapDeclarationContext() *MapDeclarationContext {
	var p = new(MapDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_mapDeclaration
	return p
}

func (*MapDeclarationContext) IsMapDeclarationContext() {}

func NewMapDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapDeclarationContext {
	var p = new(MapDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_mapDeclaration

	return p
}

func (s *MapDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *MapDeclarationContext) GetTypeParams() IMapTypeAssignmentsContext { return s.typeParams }

func (s *MapDeclarationContext) SetTypeParams(v IMapTypeAssignmentsContext) { s.typeParams = v }

func (s *MapDeclarationContext) MAP() antlr.TerminalNode {
	return s.GetToken(PdlParserMAP, 0)
}

func (s *MapDeclarationContext) MapTypeAssignments() IMapTypeAssignmentsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMapTypeAssignmentsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMapTypeAssignmentsContext)
}

func (s *MapDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MapDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MapDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterMapDeclaration(s)
	}
}

func (s *MapDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitMapDeclaration(s)
	}
}

func (p *PdlParser) MapDeclaration() (localctx IMapDeclarationContext) {
	this := p
	_ = this

	localctx = NewMapDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, PdlParserRULE_mapDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(267)
		p.Match(PdlParserMAP)
	}
	{
		p.SetState(268)

		var _x = p.MapTypeAssignments()

		localctx.(*MapDeclarationContext).typeParams = _x
	}

	return localctx
}

// IMapTypeAssignmentsContext is an interface to support dynamic dispatch.
type IMapTypeAssignmentsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetKey returns the key rule contexts.
	GetKey() ITypeAssignmentContext

	// GetValue returns the value rule contexts.
	GetValue() ITypeAssignmentContext

	// SetKey sets the key rule contexts.
	SetKey(ITypeAssignmentContext)

	// SetValue sets the value rule contexts.
	SetValue(ITypeAssignmentContext)

	// IsMapTypeAssignmentsContext differentiates from other interfaces.
	IsMapTypeAssignmentsContext()
}

type MapTypeAssignmentsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	key    ITypeAssignmentContext
	value  ITypeAssignmentContext
}

func NewEmptyMapTypeAssignmentsContext() *MapTypeAssignmentsContext {
	var p = new(MapTypeAssignmentsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_mapTypeAssignments
	return p
}

func (*MapTypeAssignmentsContext) IsMapTypeAssignmentsContext() {}

func NewMapTypeAssignmentsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapTypeAssignmentsContext {
	var p = new(MapTypeAssignmentsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_mapTypeAssignments

	return p
}

func (s *MapTypeAssignmentsContext) GetParser() antlr.Parser { return s.parser }

func (s *MapTypeAssignmentsContext) GetKey() ITypeAssignmentContext { return s.key }

func (s *MapTypeAssignmentsContext) GetValue() ITypeAssignmentContext { return s.value }

func (s *MapTypeAssignmentsContext) SetKey(v ITypeAssignmentContext) { s.key = v }

func (s *MapTypeAssignmentsContext) SetValue(v ITypeAssignmentContext) { s.value = v }

func (s *MapTypeAssignmentsContext) OPEN_BRACKET() antlr.TerminalNode {
	return s.GetToken(PdlParserOPEN_BRACKET, 0)
}

func (s *MapTypeAssignmentsContext) CLOSE_BRACKET() antlr.TerminalNode {
	return s.GetToken(PdlParserCLOSE_BRACKET, 0)
}

func (s *MapTypeAssignmentsContext) AllTypeAssignment() []ITypeAssignmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITypeAssignmentContext); ok {
			len++
		}
	}

	tst := make([]ITypeAssignmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITypeAssignmentContext); ok {
			tst[i] = t.(ITypeAssignmentContext)
			i++
		}
	}

	return tst
}

func (s *MapTypeAssignmentsContext) TypeAssignment(i int) ITypeAssignmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeAssignmentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeAssignmentContext)
}

func (s *MapTypeAssignmentsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MapTypeAssignmentsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MapTypeAssignmentsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterMapTypeAssignments(s)
	}
}

func (s *MapTypeAssignmentsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitMapTypeAssignments(s)
	}
}

func (p *PdlParser) MapTypeAssignments() (localctx IMapTypeAssignmentsContext) {
	this := p
	_ = this

	localctx = NewMapTypeAssignmentsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, PdlParserRULE_mapTypeAssignments)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(270)
		p.Match(PdlParserOPEN_BRACKET)
	}
	{
		p.SetState(271)

		var _x = p.TypeAssignment()

		localctx.(*MapTypeAssignmentsContext).key = _x
	}
	{
		p.SetState(272)

		var _x = p.TypeAssignment()

		localctx.(*MapTypeAssignmentsContext).value = _x
	}
	{
		p.SetState(273)
		p.Match(PdlParserCLOSE_BRACKET)
	}

	return localctx
}

// IFieldSelectionContext is an interface to support dynamic dispatch.
type IFieldSelectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_fieldDeclaration returns the _fieldDeclaration rule contexts.
	Get_fieldDeclaration() IFieldDeclarationContext

	// Set_fieldDeclaration sets the _fieldDeclaration rule contexts.
	Set_fieldDeclaration(IFieldDeclarationContext)

	// GetFields returns the fields rule context list.
	GetFields() []IFieldDeclarationContext

	// SetFields sets the fields rule context list.
	SetFields([]IFieldDeclarationContext)

	// IsFieldSelectionContext differentiates from other interfaces.
	IsFieldSelectionContext()
}

type FieldSelectionContext struct {
	*antlr.BaseParserRuleContext
	parser            antlr.Parser
	_fieldDeclaration IFieldDeclarationContext
	fields            []IFieldDeclarationContext
}

func NewEmptyFieldSelectionContext() *FieldSelectionContext {
	var p = new(FieldSelectionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_fieldSelection
	return p
}

func (*FieldSelectionContext) IsFieldSelectionContext() {}

func NewFieldSelectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldSelectionContext {
	var p = new(FieldSelectionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_fieldSelection

	return p
}

func (s *FieldSelectionContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldSelectionContext) Get_fieldDeclaration() IFieldDeclarationContext {
	return s._fieldDeclaration
}

func (s *FieldSelectionContext) Set_fieldDeclaration(v IFieldDeclarationContext) {
	s._fieldDeclaration = v
}

func (s *FieldSelectionContext) GetFields() []IFieldDeclarationContext { return s.fields }

func (s *FieldSelectionContext) SetFields(v []IFieldDeclarationContext) { s.fields = v }

func (s *FieldSelectionContext) OPEN_BRACE() antlr.TerminalNode {
	return s.GetToken(PdlParserOPEN_BRACE, 0)
}

func (s *FieldSelectionContext) CLOSE_BRACE() antlr.TerminalNode {
	return s.GetToken(PdlParserCLOSE_BRACE, 0)
}

func (s *FieldSelectionContext) AllFieldDeclaration() []IFieldDeclarationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldDeclarationContext); ok {
			len++
		}
	}

	tst := make([]IFieldDeclarationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldDeclarationContext); ok {
			tst[i] = t.(IFieldDeclarationContext)
			i++
		}
	}

	return tst
}

func (s *FieldSelectionContext) FieldDeclaration(i int) IFieldDeclarationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldDeclarationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldDeclarationContext)
}

func (s *FieldSelectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldSelectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldSelectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterFieldSelection(s)
	}
}

func (s *FieldSelectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitFieldSelection(s)
	}
}

func (p *PdlParser) FieldSelection() (localctx IFieldSelectionContext) {
	this := p
	_ = this

	localctx = NewFieldSelectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, PdlParserRULE_fieldSelection)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(275)
		p.Match(PdlParserOPEN_BRACE)
	}
	p.SetState(279)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1107820544) != 0 {
		{
			p.SetState(276)

			var _x = p.FieldDeclaration()

			localctx.(*FieldSelectionContext)._fieldDeclaration = _x
		}
		localctx.(*FieldSelectionContext).fields = append(localctx.(*FieldSelectionContext).fields, localctx.(*FieldSelectionContext)._fieldDeclaration)

		p.SetState(281)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(282)
		p.Match(PdlParserCLOSE_BRACE)
	}

	return localctx
}

// IFieldIncludesContext is an interface to support dynamic dispatch.
type IFieldIncludesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldIncludesContext differentiates from other interfaces.
	IsFieldIncludesContext()
}

type FieldIncludesContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldIncludesContext() *FieldIncludesContext {
	var p = new(FieldIncludesContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_fieldIncludes
	return p
}

func (*FieldIncludesContext) IsFieldIncludesContext() {}

func NewFieldIncludesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldIncludesContext {
	var p = new(FieldIncludesContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_fieldIncludes

	return p
}

func (s *FieldIncludesContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldIncludesContext) INCLUDES() antlr.TerminalNode {
	return s.GetToken(PdlParserINCLUDES, 0)
}

func (s *FieldIncludesContext) AllTypeAssignment() []ITypeAssignmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITypeAssignmentContext); ok {
			len++
		}
	}

	tst := make([]ITypeAssignmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITypeAssignmentContext); ok {
			tst[i] = t.(ITypeAssignmentContext)
			i++
		}
	}

	return tst
}

func (s *FieldIncludesContext) TypeAssignment(i int) ITypeAssignmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeAssignmentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeAssignmentContext)
}

func (s *FieldIncludesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldIncludesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldIncludesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterFieldIncludes(s)
	}
}

func (s *FieldIncludesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitFieldIncludes(s)
	}
}

func (p *PdlParser) FieldIncludes() (localctx IFieldIncludesContext) {
	this := p
	_ = this

	localctx = NewFieldIncludesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, PdlParserRULE_fieldIncludes)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(284)
		p.Match(PdlParserINCLUDES)
	}
	p.SetState(286)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(285)
				p.TypeAssignment()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(288)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext())
	}

	return localctx
}

// IFieldDeclarationContext is an interface to support dynamic dispatch.
type IFieldDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_OPTIONAL returns the _OPTIONAL token.
	Get_OPTIONAL() antlr.Token

	// Set_OPTIONAL sets the _OPTIONAL token.
	Set_OPTIONAL(antlr.Token)

	// GetDoc returns the doc rule contexts.
	GetDoc() ISchemadocContext

	// Get_propDeclaration returns the _propDeclaration rule contexts.
	Get_propDeclaration() IPropDeclarationContext

	// GetFieldName returns the fieldName rule contexts.
	GetFieldName() IIdentifierContext

	// Get_identifier returns the _identifier rule contexts.
	Get_identifier() IIdentifierContext

	// GetType_ returns the type_ rule contexts.
	GetType_() ITypeAssignmentContext

	// SetDoc sets the doc rule contexts.
	SetDoc(ISchemadocContext)

	// Set_propDeclaration sets the _propDeclaration rule contexts.
	Set_propDeclaration(IPropDeclarationContext)

	// SetFieldName sets the fieldName rule contexts.
	SetFieldName(IIdentifierContext)

	// Set_identifier sets the _identifier rule contexts.
	Set_identifier(IIdentifierContext)

	// SetType_ sets the type_ rule contexts.
	SetType_(ITypeAssignmentContext)

	// GetProps returns the props rule context list.
	GetProps() []IPropDeclarationContext

	// SetProps sets the props rule context list.
	SetProps([]IPropDeclarationContext)

	// GetName returns the name attribute.
	GetName() string

	// GetIsOptional returns the isOptional attribute.
	GetIsOptional() bool

	// SetName sets the name attribute.
	SetName(string)

	// SetIsOptional sets the isOptional attribute.
	SetIsOptional(bool)

	// IsFieldDeclarationContext differentiates from other interfaces.
	IsFieldDeclarationContext()
}

type FieldDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser           antlr.Parser
	name             string
	isOptional       bool
	doc              ISchemadocContext
	_propDeclaration IPropDeclarationContext
	props            []IPropDeclarationContext
	fieldName        IIdentifierContext
	_identifier      IIdentifierContext
	_OPTIONAL        antlr.Token
	type_            ITypeAssignmentContext
}

func NewEmptyFieldDeclarationContext() *FieldDeclarationContext {
	var p = new(FieldDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_fieldDeclaration
	return p
}

func (*FieldDeclarationContext) IsFieldDeclarationContext() {}

func NewFieldDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldDeclarationContext {
	var p = new(FieldDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_fieldDeclaration

	return p
}

func (s *FieldDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldDeclarationContext) Get_OPTIONAL() antlr.Token { return s._OPTIONAL }

func (s *FieldDeclarationContext) Set_OPTIONAL(v antlr.Token) { s._OPTIONAL = v }

func (s *FieldDeclarationContext) GetDoc() ISchemadocContext { return s.doc }

func (s *FieldDeclarationContext) Get_propDeclaration() IPropDeclarationContext {
	return s._propDeclaration
}

func (s *FieldDeclarationContext) GetFieldName() IIdentifierContext { return s.fieldName }

func (s *FieldDeclarationContext) Get_identifier() IIdentifierContext { return s._identifier }

func (s *FieldDeclarationContext) GetType_() ITypeAssignmentContext { return s.type_ }

func (s *FieldDeclarationContext) SetDoc(v ISchemadocContext) { s.doc = v }

func (s *FieldDeclarationContext) Set_propDeclaration(v IPropDeclarationContext) {
	s._propDeclaration = v
}

func (s *FieldDeclarationContext) SetFieldName(v IIdentifierContext) { s.fieldName = v }

func (s *FieldDeclarationContext) Set_identifier(v IIdentifierContext) { s._identifier = v }

func (s *FieldDeclarationContext) SetType_(v ITypeAssignmentContext) { s.type_ = v }

func (s *FieldDeclarationContext) GetProps() []IPropDeclarationContext { return s.props }

func (s *FieldDeclarationContext) SetProps(v []IPropDeclarationContext) { s.props = v }

func (s *FieldDeclarationContext) GetName() string { return s.name }

func (s *FieldDeclarationContext) GetIsOptional() bool { return s.isOptional }

func (s *FieldDeclarationContext) SetName(v string) { s.name = v }

func (s *FieldDeclarationContext) SetIsOptional(v bool) { s.isOptional = v }

func (s *FieldDeclarationContext) COLON() antlr.TerminalNode {
	return s.GetToken(PdlParserCOLON, 0)
}

func (s *FieldDeclarationContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *FieldDeclarationContext) TypeAssignment() ITypeAssignmentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeAssignmentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeAssignmentContext)
}

func (s *FieldDeclarationContext) OPTIONAL() antlr.TerminalNode {
	return s.GetToken(PdlParserOPTIONAL, 0)
}

func (s *FieldDeclarationContext) FieldDefault() IFieldDefaultContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldDefaultContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldDefaultContext)
}

func (s *FieldDeclarationContext) Schemadoc() ISchemadocContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISchemadocContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISchemadocContext)
}

func (s *FieldDeclarationContext) AllPropDeclaration() []IPropDeclarationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPropDeclarationContext); ok {
			len++
		}
	}

	tst := make([]IPropDeclarationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPropDeclarationContext); ok {
			tst[i] = t.(IPropDeclarationContext)
			i++
		}
	}

	return tst
}

func (s *FieldDeclarationContext) PropDeclaration(i int) IPropDeclarationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropDeclarationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropDeclarationContext)
}

func (s *FieldDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterFieldDeclaration(s)
	}
}

func (s *FieldDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitFieldDeclaration(s)
	}
}

func (p *PdlParser) FieldDeclaration() (localctx IFieldDeclarationContext) {
	this := p
	_ = this

	localctx = NewFieldDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, PdlParserRULE_fieldDeclaration)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(291)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserSCHEMADOC_COMMENT {
		{
			p.SetState(290)

			var _x = p.Schemadoc()

			localctx.(*FieldDeclarationContext).doc = _x
		}

	}
	p.SetState(296)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PdlParserAT {
		{
			p.SetState(293)

			var _x = p.PropDeclaration()

			localctx.(*FieldDeclarationContext)._propDeclaration = _x
		}
		localctx.(*FieldDeclarationContext).props = append(localctx.(*FieldDeclarationContext).props, localctx.(*FieldDeclarationContext)._propDeclaration)

		p.SetState(298)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(299)

		var _x = p.Identifier()

		localctx.(*FieldDeclarationContext).fieldName = _x
		localctx.(*FieldDeclarationContext)._identifier = _x
	}
	{
		p.SetState(300)
		p.Match(PdlParserCOLON)
	}
	p.SetState(302)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserOPTIONAL {
		{
			p.SetState(301)

			var _m = p.Match(PdlParserOPTIONAL)

			localctx.(*FieldDeclarationContext)._OPTIONAL = _m
		}

	}
	{
		p.SetState(304)

		var _x = p.TypeAssignment()

		localctx.(*FieldDeclarationContext).type_ = _x
	}
	p.SetState(306)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PdlParserEQ {
		{
			p.SetState(305)
			p.FieldDefault()
		}

	}

	localctx.(*FieldDeclarationContext).SetName(localctx.(*FieldDeclarationContext).Get_identifier().GetValue())
	localctx.(*FieldDeclarationContext).SetIsOptional(localctx.(*FieldDeclarationContext).Get_OPTIONAL() != nil)

	return localctx
}

// IFieldDefaultContext is an interface to support dynamic dispatch.
type IFieldDefaultContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldDefaultContext differentiates from other interfaces.
	IsFieldDefaultContext()
}

type FieldDefaultContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldDefaultContext() *FieldDefaultContext {
	var p = new(FieldDefaultContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_fieldDefault
	return p
}

func (*FieldDefaultContext) IsFieldDefaultContext() {}

func NewFieldDefaultContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldDefaultContext {
	var p = new(FieldDefaultContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_fieldDefault

	return p
}

func (s *FieldDefaultContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldDefaultContext) EQ() antlr.TerminalNode {
	return s.GetToken(PdlParserEQ, 0)
}

func (s *FieldDefaultContext) JsonValue() IJsonValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJsonValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJsonValueContext)
}

func (s *FieldDefaultContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldDefaultContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldDefaultContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterFieldDefault(s)
	}
}

func (s *FieldDefaultContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitFieldDefault(s)
	}
}

func (p *PdlParser) FieldDefault() (localctx IFieldDefaultContext) {
	this := p
	_ = this

	localctx = NewFieldDefaultContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, PdlParserRULE_fieldDefault)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(310)
		p.Match(PdlParserEQ)
	}
	{
		p.SetState(311)
		p.JsonValue()
	}

	return localctx
}

// ITypeNameContext is an interface to support dynamic dispatch.
type ITypeNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetValue returns the value attribute.
	GetValue() string

	// SetValue sets the value attribute.
	SetValue(string)

	// IsTypeNameContext differentiates from other interfaces.
	IsTypeNameContext()
}

type TypeNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	value  string
}

func NewEmptyTypeNameContext() *TypeNameContext {
	var p = new(TypeNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_typeName
	return p
}

func (*TypeNameContext) IsTypeNameContext() {}

func NewTypeNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeNameContext {
	var p = new(TypeNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_typeName

	return p
}

func (s *TypeNameContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeNameContext) GetValue() string { return s.value }

func (s *TypeNameContext) SetValue(v string) { s.value = v }

func (s *TypeNameContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(PdlParserID)
}

func (s *TypeNameContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(PdlParserID, i)
}

func (s *TypeNameContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(PdlParserDOT)
}

func (s *TypeNameContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(PdlParserDOT, i)
}

func (s *TypeNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterTypeName(s)
	}
}

func (s *TypeNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitTypeName(s)
	}
}

func (p *PdlParser) TypeName() (localctx ITypeNameContext) {
	this := p
	_ = this

	localctx = NewTypeNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, PdlParserRULE_typeName)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(313)
		p.Match(PdlParserID)
	}
	p.SetState(318)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PdlParserDOT {
		{
			p.SetState(314)
			p.Match(PdlParserDOT)
		}
		{
			p.SetState(315)
			p.Match(PdlParserID)
		}

		p.SetState(320)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	localctx.(*TypeNameContext).SetValue(validatePegasusId(unescapeIdentifier(p.GetTokenStream().GetTextFromTokens(localctx.GetStart(), p.GetTokenStream().LT(-1)))))

	return localctx
}

// IIdentifierContext is an interface to support dynamic dispatch.
type IIdentifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetValue returns the value attribute.
	GetValue() string

	// SetValue sets the value attribute.
	SetValue(string)

	// IsIdentifierContext differentiates from other interfaces.
	IsIdentifierContext()
}

type IdentifierContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	value  string
}

func NewEmptyIdentifierContext() *IdentifierContext {
	var p = new(IdentifierContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_identifier
	return p
}

func (*IdentifierContext) IsIdentifierContext() {}

func NewIdentifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifierContext {
	var p = new(IdentifierContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_identifier

	return p
}

func (s *IdentifierContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentifierContext) GetValue() string { return s.value }

func (s *IdentifierContext) SetValue(v string) { s.value = v }

func (s *IdentifierContext) ID() antlr.TerminalNode {
	return s.GetToken(PdlParserID, 0)
}

func (s *IdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterIdentifier(s)
	}
}

func (s *IdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitIdentifier(s)
	}
}

func (p *PdlParser) Identifier() (localctx IIdentifierContext) {
	this := p
	_ = this

	localctx = NewIdentifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, PdlParserRULE_identifier)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(323)
		p.Match(PdlParserID)
	}

	localctx.(*IdentifierContext).SetValue(validatePegasusId(unescapeIdentifier(p.GetTokenStream().GetTextFromTokens(localctx.GetStart(), p.GetTokenStream().LT(-1)))))

	return localctx
}

// IPropNameContext is an interface to support dynamic dispatch.
type IPropNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_propSegment returns the _propSegment rule contexts.
	Get_propSegment() IPropSegmentContext

	// Set_propSegment sets the _propSegment rule contexts.
	Set_propSegment(IPropSegmentContext)

	// GetPath returns the path attribute.
	GetPath() []string

	// SetPath sets the path attribute.
	SetPath([]string)

	// IsPropNameContext differentiates from other interfaces.
	IsPropNameContext()
}

type PropNameContext struct {
	*antlr.BaseParserRuleContext
	parser       antlr.Parser
	path         []string
	_propSegment IPropSegmentContext
}

func NewEmptyPropNameContext() *PropNameContext {
	var p = new(PropNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_propName
	return p
}

func (*PropNameContext) IsPropNameContext() {}

func NewPropNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropNameContext {
	var p = new(PropNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_propName

	return p
}

func (s *PropNameContext) GetParser() antlr.Parser { return s.parser }

func (s *PropNameContext) Get_propSegment() IPropSegmentContext { return s._propSegment }

func (s *PropNameContext) Set_propSegment(v IPropSegmentContext) { s._propSegment = v }

func (s *PropNameContext) GetPath() []string { return s.path }

func (s *PropNameContext) SetPath(v []string) { s.path = v }

func (s *PropNameContext) AllPropSegment() []IPropSegmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPropSegmentContext); ok {
			len++
		}
	}

	tst := make([]IPropSegmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPropSegmentContext); ok {
			tst[i] = t.(IPropSegmentContext)
			i++
		}
	}

	return tst
}

func (s *PropNameContext) PropSegment(i int) IPropSegmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropSegmentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropSegmentContext)
}

func (s *PropNameContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(PdlParserDOT)
}

func (s *PropNameContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(PdlParserDOT, i)
}

func (s *PropNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterPropName(s)
	}
}

func (s *PropNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitPropName(s)
	}
}

func (p *PdlParser) PropName() (localctx IPropNameContext) {
	this := p
	_ = this

	localctx = NewPropNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, PdlParserRULE_propName)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(326)

		var _x = p.PropSegment()

		localctx.(*PropNameContext)._propSegment = _x
	}
	localctx.(*PropNameContext).SetPath(append(localctx.(*PropNameContext).path, localctx.(*PropNameContext).Get_propSegment().GetValue()))
	p.SetState(334)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PdlParserDOT {
		{
			p.SetState(328)
			p.Match(PdlParserDOT)
		}
		{
			p.SetState(329)

			var _x = p.PropSegment()

			localctx.(*PropNameContext)._propSegment = _x
		}
		localctx.(*PropNameContext).SetPath(append(localctx.(*PropNameContext).path, localctx.(*PropNameContext).Get_propSegment().GetValue()))

		p.SetState(336)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IPropSegmentContext is an interface to support dynamic dispatch.
type IPropSegmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetValue returns the value attribute.
	GetValue() string

	// SetValue sets the value attribute.
	SetValue(string)

	// IsPropSegmentContext differentiates from other interfaces.
	IsPropSegmentContext()
}

type PropSegmentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	value  string
}

func NewEmptyPropSegmentContext() *PropSegmentContext {
	var p = new(PropSegmentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_propSegment
	return p
}

func (*PropSegmentContext) IsPropSegmentContext() {}

func NewPropSegmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropSegmentContext {
	var p = new(PropSegmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_propSegment

	return p
}

func (s *PropSegmentContext) GetParser() antlr.Parser { return s.parser }

func (s *PropSegmentContext) GetValue() string { return s.value }

func (s *PropSegmentContext) SetValue(v string) { s.value = v }

func (s *PropSegmentContext) ID() antlr.TerminalNode {
	return s.GetToken(PdlParserID, 0)
}

func (s *PropSegmentContext) PROPERTY_ID() antlr.TerminalNode {
	return s.GetToken(PdlParserPROPERTY_ID, 0)
}

func (s *PropSegmentContext) ESCAPED_PROP_ID() antlr.TerminalNode {
	return s.GetToken(PdlParserESCAPED_PROP_ID, 0)
}

func (s *PropSegmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropSegmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropSegmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterPropSegment(s)
	}
}

func (s *PropSegmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitPropSegment(s)
	}
}

func (p *PdlParser) PropSegment() (localctx IPropSegmentContext) {
	this := p
	_ = this

	localctx = NewPropSegmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, PdlParserRULE_propSegment)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(337)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&13958643712) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	localctx.(*PropSegmentContext).SetValue(unescapeIdentifier(p.GetTokenStream().GetTextFromTokens(localctx.GetStart(), p.GetTokenStream().LT(-1))))

	return localctx
}

// ISchemadocContext is an interface to support dynamic dispatch.
type ISchemadocContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_SCHEMADOC_COMMENT returns the _SCHEMADOC_COMMENT token.
	Get_SCHEMADOC_COMMENT() antlr.Token

	// Set_SCHEMADOC_COMMENT sets the _SCHEMADOC_COMMENT token.
	Set_SCHEMADOC_COMMENT(antlr.Token)

	// GetValue returns the value attribute.
	GetValue() string

	// SetValue sets the value attribute.
	SetValue(string)

	// IsSchemadocContext differentiates from other interfaces.
	IsSchemadocContext()
}

type SchemadocContext struct {
	*antlr.BaseParserRuleContext
	parser             antlr.Parser
	value              string
	_SCHEMADOC_COMMENT antlr.Token
}

func NewEmptySchemadocContext() *SchemadocContext {
	var p = new(SchemadocContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_schemadoc
	return p
}

func (*SchemadocContext) IsSchemadocContext() {}

func NewSchemadocContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SchemadocContext {
	var p = new(SchemadocContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_schemadoc

	return p
}

func (s *SchemadocContext) GetParser() antlr.Parser { return s.parser }

func (s *SchemadocContext) Get_SCHEMADOC_COMMENT() antlr.Token { return s._SCHEMADOC_COMMENT }

func (s *SchemadocContext) Set_SCHEMADOC_COMMENT(v antlr.Token) { s._SCHEMADOC_COMMENT = v }

func (s *SchemadocContext) GetValue() string { return s.value }

func (s *SchemadocContext) SetValue(v string) { s.value = v }

func (s *SchemadocContext) SCHEMADOC_COMMENT() antlr.TerminalNode {
	return s.GetToken(PdlParserSCHEMADOC_COMMENT, 0)
}

func (s *SchemadocContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SchemadocContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SchemadocContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterSchemadoc(s)
	}
}

func (s *SchemadocContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitSchemadoc(s)
	}
}

func (p *PdlParser) Schemadoc() (localctx ISchemadocContext) {
	this := p
	_ = this

	localctx = NewSchemadocContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, PdlParserRULE_schemadoc)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(340)

		var _m = p.Match(PdlParserSCHEMADOC_COMMENT)

		localctx.(*SchemadocContext)._SCHEMADOC_COMMENT = _m
	}

	localctx.(*SchemadocContext).SetValue(extractMarkdown((func() string {
		if localctx.(*SchemadocContext).Get_SCHEMADOC_COMMENT() == nil {
			return ""
		} else {
			return localctx.(*SchemadocContext).Get_SCHEMADOC_COMMENT().GetText()
		}
	}())))

	return localctx
}

// IObjectContext is an interface to support dynamic dispatch.
type IObjectContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsObjectContext differentiates from other interfaces.
	IsObjectContext()
}

type ObjectContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObjectContext() *ObjectContext {
	var p = new(ObjectContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_object
	return p
}

func (*ObjectContext) IsObjectContext() {}

func NewObjectContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectContext {
	var p = new(ObjectContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_object

	return p
}

func (s *ObjectContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectContext) OPEN_BRACE() antlr.TerminalNode {
	return s.GetToken(PdlParserOPEN_BRACE, 0)
}

func (s *ObjectContext) CLOSE_BRACE() antlr.TerminalNode {
	return s.GetToken(PdlParserCLOSE_BRACE, 0)
}

func (s *ObjectContext) AllObjectEntry() []IObjectEntryContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IObjectEntryContext); ok {
			len++
		}
	}

	tst := make([]IObjectEntryContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IObjectEntryContext); ok {
			tst[i] = t.(IObjectEntryContext)
			i++
		}
	}

	return tst
}

func (s *ObjectContext) ObjectEntry(i int) IObjectEntryContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectEntryContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectEntryContext)
}

func (s *ObjectContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterObject(s)
	}
}

func (s *ObjectContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitObject(s)
	}
}

func (p *PdlParser) Object() (localctx IObjectContext) {
	this := p
	_ = this

	localctx = NewObjectContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, PdlParserRULE_object)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(343)
		p.Match(PdlParserOPEN_BRACE)
	}
	p.SetState(347)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PdlParserSTRING_LITERAL {
		{
			p.SetState(344)
			p.ObjectEntry()
		}

		p.SetState(349)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(350)
		p.Match(PdlParserCLOSE_BRACE)
	}

	return localctx
}

// IObjectEntryContext is an interface to support dynamic dispatch.
type IObjectEntryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetKey returns the key rule contexts.
	GetKey() IStringContext

	// GetValue returns the value rule contexts.
	GetValue() IJsonValueContext

	// SetKey sets the key rule contexts.
	SetKey(IStringContext)

	// SetValue sets the value rule contexts.
	SetValue(IJsonValueContext)

	// IsObjectEntryContext differentiates from other interfaces.
	IsObjectEntryContext()
}

type ObjectEntryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	key    IStringContext
	value  IJsonValueContext
}

func NewEmptyObjectEntryContext() *ObjectEntryContext {
	var p = new(ObjectEntryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_objectEntry
	return p
}

func (*ObjectEntryContext) IsObjectEntryContext() {}

func NewObjectEntryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectEntryContext {
	var p = new(ObjectEntryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_objectEntry

	return p
}

func (s *ObjectEntryContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectEntryContext) GetKey() IStringContext { return s.key }

func (s *ObjectEntryContext) GetValue() IJsonValueContext { return s.value }

func (s *ObjectEntryContext) SetKey(v IStringContext) { s.key = v }

func (s *ObjectEntryContext) SetValue(v IJsonValueContext) { s.value = v }

func (s *ObjectEntryContext) COLON() antlr.TerminalNode {
	return s.GetToken(PdlParserCOLON, 0)
}

func (s *ObjectEntryContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *ObjectEntryContext) JsonValue() IJsonValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJsonValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJsonValueContext)
}

func (s *ObjectEntryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectEntryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectEntryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterObjectEntry(s)
	}
}

func (s *ObjectEntryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitObjectEntry(s)
	}
}

func (p *PdlParser) ObjectEntry() (localctx IObjectEntryContext) {
	this := p
	_ = this

	localctx = NewObjectEntryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, PdlParserRULE_objectEntry)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(352)

		var _x = p.String_()

		localctx.(*ObjectEntryContext).key = _x
	}
	{
		p.SetState(353)
		p.Match(PdlParserCOLON)
	}
	{
		p.SetState(354)

		var _x = p.JsonValue()

		localctx.(*ObjectEntryContext).value = _x
	}

	return localctx
}

// IArrayContext is an interface to support dynamic dispatch.
type IArrayContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItems returns the items rule contexts.
	GetItems() IJsonValueContext

	// SetItems sets the items rule contexts.
	SetItems(IJsonValueContext)

	// IsArrayContext differentiates from other interfaces.
	IsArrayContext()
}

type ArrayContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	items  IJsonValueContext
}

func NewEmptyArrayContext() *ArrayContext {
	var p = new(ArrayContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_array
	return p
}

func (*ArrayContext) IsArrayContext() {}

func NewArrayContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayContext {
	var p = new(ArrayContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_array

	return p
}

func (s *ArrayContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayContext) GetItems() IJsonValueContext { return s.items }

func (s *ArrayContext) SetItems(v IJsonValueContext) { s.items = v }

func (s *ArrayContext) OPEN_BRACKET() antlr.TerminalNode {
	return s.GetToken(PdlParserOPEN_BRACKET, 0)
}

func (s *ArrayContext) CLOSE_BRACKET() antlr.TerminalNode {
	return s.GetToken(PdlParserCLOSE_BRACKET, 0)
}

func (s *ArrayContext) AllJsonValue() []IJsonValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IJsonValueContext); ok {
			len++
		}
	}

	tst := make([]IJsonValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IJsonValueContext); ok {
			tst[i] = t.(IJsonValueContext)
			i++
		}
	}

	return tst
}

func (s *ArrayContext) JsonValue(i int) IJsonValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJsonValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJsonValueContext)
}

func (s *ArrayContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterArray(s)
	}
}

func (s *ArrayContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitArray(s)
	}
}

func (p *PdlParser) Array() (localctx IArrayContext) {
	this := p
	_ = this

	localctx = NewArrayContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, PdlParserRULE_array)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(356)
		p.Match(PdlParserOPEN_BRACKET)
	}
	p.SetState(360)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&830636032) != 0 {
		{
			p.SetState(357)

			var _x = p.JsonValue()

			localctx.(*ArrayContext).items = _x
		}

		p.SetState(362)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(363)
		p.Match(PdlParserCLOSE_BRACKET)
	}

	return localctx
}

// IJsonValueContext is an interface to support dynamic dispatch.
type IJsonValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsJsonValueContext differentiates from other interfaces.
	IsJsonValueContext()
}

type JsonValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyJsonValueContext() *JsonValueContext {
	var p = new(JsonValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_jsonValue
	return p
}

func (*JsonValueContext) IsJsonValueContext() {}

func NewJsonValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *JsonValueContext {
	var p = new(JsonValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_jsonValue

	return p
}

func (s *JsonValueContext) GetParser() antlr.Parser { return s.parser }

func (s *JsonValueContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *JsonValueContext) Number() INumberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INumberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INumberContext)
}

func (s *JsonValueContext) Object() IObjectContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *JsonValueContext) Array() IArrayContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrayContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrayContext)
}

func (s *JsonValueContext) Bool_() IBoolContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBoolContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBoolContext)
}

func (s *JsonValueContext) NullValue() INullValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INullValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INullValueContext)
}

func (s *JsonValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JsonValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *JsonValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterJsonValue(s)
	}
}

func (s *JsonValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitJsonValue(s)
	}
}

func (p *PdlParser) JsonValue() (localctx IJsonValueContext) {
	this := p
	_ = this

	localctx = NewJsonValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, PdlParserRULE_jsonValue)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(371)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PdlParserSTRING_LITERAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(365)
			p.String_()
		}

	case PdlParserNUMBER_LITERAL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(366)
			p.Number()
		}

	case PdlParserOPEN_BRACE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(367)
			p.Object()
		}

	case PdlParserOPEN_BRACKET:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(368)
			p.Array()
		}

	case PdlParserBOOLEAN_LITERAL:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(369)
			p.Bool_()
		}

	case PdlParserNULL_LITERAL:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(370)
			p.NullValue()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IStringContext is an interface to support dynamic dispatch.
type IStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_STRING_LITERAL returns the _STRING_LITERAL token.
	Get_STRING_LITERAL() antlr.Token

	// Set_STRING_LITERAL sets the _STRING_LITERAL token.
	Set_STRING_LITERAL(antlr.Token)

	// GetValue returns the value attribute.
	GetValue() string

	// SetValue sets the value attribute.
	SetValue(string)

	// IsStringContext differentiates from other interfaces.
	IsStringContext()
}

type StringContext struct {
	*antlr.BaseParserRuleContext
	parser          antlr.Parser
	value           string
	_STRING_LITERAL antlr.Token
}

func NewEmptyStringContext() *StringContext {
	var p = new(StringContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_string
	return p
}

func (*StringContext) IsStringContext() {}

func NewStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringContext {
	var p = new(StringContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_string

	return p
}

func (s *StringContext) GetParser() antlr.Parser { return s.parser }

func (s *StringContext) Get_STRING_LITERAL() antlr.Token { return s._STRING_LITERAL }

func (s *StringContext) Set_STRING_LITERAL(v antlr.Token) { s._STRING_LITERAL = v }

func (s *StringContext) GetValue() string { return s.value }

func (s *StringContext) SetValue(v string) { s.value = v }

func (s *StringContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(PdlParserSTRING_LITERAL, 0)
}

func (s *StringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterString(s)
	}
}

func (s *StringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitString(s)
	}
}

func (p *PdlParser) String_() (localctx IStringContext) {
	this := p
	_ = this

	localctx = NewStringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, PdlParserRULE_string)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(373)

		var _m = p.Match(PdlParserSTRING_LITERAL)

		localctx.(*StringContext)._STRING_LITERAL = _m
	}

	localctx.(*StringContext).value = parseStringLiteral((func() string {
		if localctx.(*StringContext).Get_STRING_LITERAL() == nil {
			return ""
		} else {
			return localctx.(*StringContext).Get_STRING_LITERAL().GetText()
		}
	}()))

	return localctx
}

// INumberContext is an interface to support dynamic dispatch.
type INumberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_NUMBER_LITERAL returns the _NUMBER_LITERAL token.
	Get_NUMBER_LITERAL() antlr.Token

	// Set_NUMBER_LITERAL sets the _NUMBER_LITERAL token.
	Set_NUMBER_LITERAL(antlr.Token)

	// GetValue returns the value attribute.
	GetValue() Number

	// SetValue sets the value attribute.
	SetValue(Number)

	// IsNumberContext differentiates from other interfaces.
	IsNumberContext()
}

type NumberContext struct {
	*antlr.BaseParserRuleContext
	parser          antlr.Parser
	value           Number
	_NUMBER_LITERAL antlr.Token
}

func NewEmptyNumberContext() *NumberContext {
	var p = new(NumberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_number
	return p
}

func (*NumberContext) IsNumberContext() {}

func NewNumberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumberContext {
	var p = new(NumberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_number

	return p
}

func (s *NumberContext) GetParser() antlr.Parser { return s.parser }

func (s *NumberContext) Get_NUMBER_LITERAL() antlr.Token { return s._NUMBER_LITERAL }

func (s *NumberContext) Set_NUMBER_LITERAL(v antlr.Token) { s._NUMBER_LITERAL = v }

func (s *NumberContext) GetValue() Number { return s.value }

func (s *NumberContext) SetValue(v Number) { s.value = v }

func (s *NumberContext) NUMBER_LITERAL() antlr.TerminalNode {
	return s.GetToken(PdlParserNUMBER_LITERAL, 0)
}

func (s *NumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NumberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterNumber(s)
	}
}

func (s *NumberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitNumber(s)
	}
}

func (p *PdlParser) Number() (localctx INumberContext) {
	this := p
	_ = this

	localctx = NewNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, PdlParserRULE_number)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(376)

		var _m = p.Match(PdlParserNUMBER_LITERAL)

		localctx.(*NumberContext)._NUMBER_LITERAL = _m
	}

	localctx.(*NumberContext).SetValue(parseNumber((func() string {
		if localctx.(*NumberContext).Get_NUMBER_LITERAL() == nil {
			return ""
		} else {
			return localctx.(*NumberContext).Get_NUMBER_LITERAL().GetText()
		}
	}())))

	return localctx
}

// IBoolContext is an interface to support dynamic dispatch.
type IBoolContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_BOOLEAN_LITERAL returns the _BOOLEAN_LITERAL token.
	Get_BOOLEAN_LITERAL() antlr.Token

	// Set_BOOLEAN_LITERAL sets the _BOOLEAN_LITERAL token.
	Set_BOOLEAN_LITERAL(antlr.Token)

	// GetValue returns the value attribute.
	GetValue() bool

	// SetValue sets the value attribute.
	SetValue(bool)

	// IsBoolContext differentiates from other interfaces.
	IsBoolContext()
}

type BoolContext struct {
	*antlr.BaseParserRuleContext
	parser           antlr.Parser
	value            bool
	_BOOLEAN_LITERAL antlr.Token
}

func NewEmptyBoolContext() *BoolContext {
	var p = new(BoolContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_bool
	return p
}

func (*BoolContext) IsBoolContext() {}

func NewBoolContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoolContext {
	var p = new(BoolContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_bool

	return p
}

func (s *BoolContext) GetParser() antlr.Parser { return s.parser }

func (s *BoolContext) Get_BOOLEAN_LITERAL() antlr.Token { return s._BOOLEAN_LITERAL }

func (s *BoolContext) Set_BOOLEAN_LITERAL(v antlr.Token) { s._BOOLEAN_LITERAL = v }

func (s *BoolContext) GetValue() bool { return s.value }

func (s *BoolContext) SetValue(v bool) { s.value = v }

func (s *BoolContext) BOOLEAN_LITERAL() antlr.TerminalNode {
	return s.GetToken(PdlParserBOOLEAN_LITERAL, 0)
}

func (s *BoolContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BoolContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterBool(s)
	}
}

func (s *BoolContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitBool(s)
	}
}

func (p *PdlParser) Bool_() (localctx IBoolContext) {
	this := p
	_ = this

	localctx = NewBoolContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, PdlParserRULE_bool)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(379)

		var _m = p.Match(PdlParserBOOLEAN_LITERAL)

		localctx.(*BoolContext)._BOOLEAN_LITERAL = _m
	}

	localctx.(*BoolContext).SetValue(parseBool((func() string {
		if localctx.(*BoolContext).Get_BOOLEAN_LITERAL() == nil {
			return ""
		} else {
			return localctx.(*BoolContext).Get_BOOLEAN_LITERAL().GetText()
		}
	}())))

	return localctx
}

// INullValueContext is an interface to support dynamic dispatch.
type INullValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNullValueContext differentiates from other interfaces.
	IsNullValueContext()
}

type NullValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNullValueContext() *NullValueContext {
	var p = new(NullValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PdlParserRULE_nullValue
	return p
}

func (*NullValueContext) IsNullValueContext() {}

func NewNullValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NullValueContext {
	var p = new(NullValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PdlParserRULE_nullValue

	return p
}

func (s *NullValueContext) GetParser() antlr.Parser { return s.parser }

func (s *NullValueContext) NULL_LITERAL() antlr.TerminalNode {
	return s.GetToken(PdlParserNULL_LITERAL, 0)
}

func (s *NullValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NullValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NullValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.EnterNullValue(s)
	}
}

func (s *NullValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PdlListener); ok {
		listenerT.ExitNullValue(s)
	}
}

func (p *PdlParser) NullValue() (localctx INullValueContext) {
	this := p
	_ = this

	localctx = NewNullValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, PdlParserRULE_nullValue)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(382)
		p.Match(PdlParserNULL_LITERAL)
	}

	return localctx
}
