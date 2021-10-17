package main

import (
	"log"

	"github.com/PapaCharlie/go-restli/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := cmd.CodeGenerator().Execute(); err != nil {
		log.Fatalf("%+v", err)
	}
}
