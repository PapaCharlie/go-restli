package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/dave/jennifer/jen"
	"go-restli/codegen"
	"go-restli/codegen/models"
	"go-restli/codegen/schema"
	"io"
	"log"
	"os"
	"path/filepath"
)

//go:generate go run codegen/protocol/protocol_zipper.go protocol/ codegen/zipped_protocol.go
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

	if len(os.Args) < 4 {
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
			if code := m.GenerateModelCode(packagePrefix, filename); code != nil {
				codeFiles = append(codeFiles, code)
			}
		}

		loadedSchema, err := schema.LoadSchema(readFile(filename))
		if err != nil {
			log.Panicf("%s: %+v", filename, err)
		}

		if loadedSchema != nil {
			if code := loadedSchema.GenerateCode(packagePrefix, filename); code != nil {
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

	generateAllImportsFile(outputDir, packagePrefix, codeFiles)
	unzipProtocol(outputDir, packagePrefix)
}

func readFile(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func generateAllImportsFile(outputDir, packagePrefix string, codeFiles []*codegen.CodeFile) {
	imports := make(map[string]bool)
	for _, code := range codeFiles {
		imports[code.PackagePath] = true
	}
	f := jen.NewFile("main")
	for p := range imports {
		f.Anon(p)
	}
	f.Func().Id("main").Params().Block(jen.Qual("fmt", "Println").Call(jen.Lit("success!")))

	out, err := os.Create(filepath.Join(outputDir, packagePrefix, "all_imports.go"))
	check(err)
	check(f.Render(out))
}

func unzipProtocol(outputDir, packagePrefix string) {
	reader, err := zip.NewReader(bytes.NewReader(codegen.ProtocolZip), int64(len(codegen.ProtocolZip)))
	check(err)

	for _, zipFile := range reader.File {
		name := filepath.Join(outputDir, packagePrefix, zipFile.Name)
		check(os.MkdirAll(filepath.Dir(name), os.ModePerm))

		f, err := os.Create(name)
		check(err)

		zipFileReader, err := zipFile.Open()
		check(err)

		_, err = io.Copy(f, zipFileReader)
		check(err)
		check(zipFileReader.Close())
		check(f.Close())
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
