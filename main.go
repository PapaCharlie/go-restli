package main

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen/cli"
)

func main() {
	log.SetFlags(log.Lshortfile)
	cmd := cli.CodeGenerator()
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
