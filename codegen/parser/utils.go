package parser

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

//go:generate java -Xmx500M -cp antlr4-4.11.1-complete.jar org.antlr.v4.Tool -Dlanguage=Go -no-visitor -package parser Pdl.g4

type Number = json.Number

func parseNumber(literal string) Number {
	n := Number(literal)
	_, err := n.Float64()
	if err != nil {
		log.Panicf("Invalid number literal: %q", literal)
	}
	return n
}

func parseStringLiteral(literal string) (parsed string) {
	err := json.Unmarshal([]byte(literal), &parsed)
	if err != nil {
		log.Panicf("Invalid string literal: %q", literal)
	}
	return parsed
}

func parseBool(literal string) bool {
	b, err := strconv.ParseBool(literal)
	if err != nil {
		log.Panicf("Invalid bool literal: %q", literal)
	}
	return b
}

func unescapeIdentifier(identifier string) string {
	return strings.ReplaceAll(identifier, "`", "")
}

func validatePegasusId(identifier string) string {
	if strings.Contains(identifier, "-") {
		log.Panicf("Invalid identifier contains \"-\": %q", identifier)
	}
	return identifier
}

func extractMarkdown(mkdwn string) string {
	return mkdwn
}
