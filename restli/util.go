package restli

import (
	"github.com/dave/jennifer/jen"
	"log"
	"strings"
)

func getStringField(m map[string]interface{}, field string) string {
	if i, hasField := m[field]; hasField {
		if s, isString := i.(string); isString {
			return s
		} else {
			log.Panicf("%s was not a string in %v", field, m)
		}
	}
	return ""
}

func getFieldName(f map[string]interface{}) string {
	return capitalizeFirstLetter(f[Name].(string))
}

func jsonTag(fieldName string) map[string]string {
	return map[string]string{"json": fieldName}
}

func capitalizeFirstLetter(name string) string {
	return strings.ToUpper(name[:1]) + name[1:]
}

func getQual(namespace, name string) *jen.Statement {
	return jen.Qual(strings.Replace(namespace, NamespaceSep, "/", -1), name)
}

func NsJoin(elements ...string) string {
	return strings.Join(elements, NamespaceSep)
}
