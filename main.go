package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/models"
	"github.com/PapaCharlie/go-restli/codegen/schema"
	"github.com/dave/jennifer/jen"
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
	codegen.SetPackagePrefix(os.Args[1])

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
		log.Println(filename)
		loadedModels, err := models.LoadModels(readFile(filename))
		if err != nil {
			log.Fatalf("could not load %s: %+v", filename, err)
		}

		for _, m := range loadedModels {
			if code := m.GenerateModelCode(); code != nil {
				code.SourceFilename = filename
				codeFiles = append(codeFiles, code)
			}
		}

		loadedResources, err := schema.LoadResources(readFile(filename))
		if err != nil {
			log.Panicf("%s: %+v", filename, err)
		}

		if len(loadedResources) > 0 {
			for _, r := range loadedResources {
				k := r.GenerateCode()
				for _, code := range k {
					code.SourceFilename = filename
					codeFiles = append(codeFiles, code)
				}
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

	generateAllImportsFile(outputDir, codeFiles)
	unzipProtocol(outputDir)
}

func readFile(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func generateAllImportsFile(outputDir string, codeFiles []*codegen.CodeFile) {
	imports := make(map[string]bool)
	for _, code := range codeFiles {
		imports[code.PackagePath] = true
	}
	f := jen.NewFile("main")
	for p := range imports {
		f.Anon(p)
	}
	f.Func().Id("main").Params().Block(jen.Qual("fmt", "Println").Call(jen.Lit("success!")))

	path := filepath.Join(outputDir, codegen.GetPackagePrefix(), "all_imports.go")
	out, err := os.Create(path)
	check(err)
	check(f.Render(out))
	check(out.Close())
	check(os.Chmod(path, os.FileMode(0555)))
}

func unzipProtocol(outputDir string) {
	reader, err := zip.NewReader(bytes.NewReader(codegen.ProtocolZip), int64(len(codegen.ProtocolZip)))
	check(err)

	for _, zipFile := range reader.File {
		name := filepath.Join(outputDir, codegen.GetPackagePrefix(), zipFile.Name)
		check(os.MkdirAll(filepath.Dir(name), os.ModePerm))

		f, err := os.Create(name)
		check(err)

		zipFileReader, err := zipFile.Open()
		check(err)

		_, err = io.Copy(f, zipFileReader)
		check(err)
		check(zipFileReader.Close())
		check(f.Close())

		check(os.Chmod(name, os.FileMode(0555)))
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
