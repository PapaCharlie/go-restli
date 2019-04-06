package main

import (
	"github.com/PapaCharlie/go-restli/codegen/cli"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	cmd := cli.CodeGenerator()
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
