package main

import (
	_ "embed"
	"log"

	"github.com/PapaCharlie/go-restli/v2/cmd"
)

// jar is embedded in the main package to ensure the 25MB jar isn't accidentally bundled in downstream builds via an
// unfortunate import
//
//go:embed go-restli-spec-parser.jar
var jar []byte

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	if err := cmd.CodeGenerator(jar).Execute(); err != nil {
		log.Fatalf("%+v", err)
	}
}
