package main

import (
	"fmt"
	"go-restli/codegen"
	"go-restli/codegen/models"
	"go-restli/codegen/schema"
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if len(os.Args) == 1 {
		log.Fatalf("Usage: %s PACKAGE_PREFIX OUTPUT_DIR SNAPSHOT_FILES...", os.Args[0])
	}

	if len(os.Args) < 2 {
		log.Fatal("Must specify the package prefix")
	}
	packagePrefix := os.Args[1]

	if len(os.Args) < 3 {
		log.Fatal("Must specify the output dir")
	}
	outputDir, err := filepath.Abs(os.Args[2])
	if err != nil {
		log.Fatal("Illegal path", err)
	}

	if len(os.Args) < 3 {
		log.Fatalf("Must specify at least one snapshot file")
	}
	snapshotFiles := os.Args[3:]

	var codeFiles []*codegen.CodeFile

	for _, filename := range snapshotFiles {
		loadedModels, err := models.LoadModels(readFile(filename))
		if err != nil {
			log.Fatalf("could not load %s: %+v", filename, err)
		}

		for _, m := range loadedModels {
			if code := m.GenerateModelCode(packagePrefix); code != nil {
				codeFiles = append(codeFiles, code)
			}
		}

		loadedSchema, err := schema.LoadSchema(readFile(filename))
		if err != nil {
			log.Panicf("%+v", err)
		}

		if loadedSchema != nil {
			if code := loadedSchema.GenerateCode(packagePrefix); code != nil {
				codeFiles = append(codeFiles, code)
			}
		}
	}

	for _, code := range codeFiles {
		file, err := code.Write(outputDir)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(file)

		}
	}
}

func readFile(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
