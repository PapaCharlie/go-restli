package main

import (
	"fmt"
	"go-restli/restli/models"
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

	for _, filename := range snapshotFiles {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}

		loadedModels, err := models.LoadModels(file)
		if err != nil {
			log.Panicf("%+v", err)
		}

		for _, m := range loadedModels {
			file, err := m.GenerateModelCode(outputDir, packagePrefix)
			if err != nil {
				log.Fatal(err)
			}
			if file != "" {
				fmt.Println(file)
			}
		}
	}
}
