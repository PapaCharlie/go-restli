package main

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"go-restli/restli"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Panicln("Must specify at least one snapshot file")
	}
	task := restli.SnapshotParser{DestinationPackage: os.Args[1]}
	outputDir := os.Args[2]
	openFiles := map[string]*jen.File{}
	generatedTypes := map[string]bool{}

	for _, filename := range os.Args[3:] {
		task.GenerateTypes(filename)
		for _, m := range task.GeneratedTypes {
			fqn := restli.NsJoin(m.Namespace, m.Name)
			if generatedTypes[fqn] {
				continue
			} else {
				generatedTypes[fqn] = true
			}

			var f *jen.File
			var ok bool

			packageName := m.PackageName()
			fileName := packageName + "/" + string(m.Category) + ".go"

			if f, ok = openFiles[fileName]; !ok {
				f = jen.NewFilePath(packageName)
				openFiles[fileName] = f
			}
			f.Add(m.Definition...)
		}
	}

	for p, f := range openFiles {
		p, err := filepath.Abs(filepath.Join(outputDir, p))
		if err != nil {
			panic(err)
		}
		fmt.Println(p)

		if err := os.MkdirAll(filepath.Dir(p), os.ModePerm); err != nil {
			panic(err)
		}
		if _, err := os.Stat(p); err == nil {
			if err := os.Remove(p); err != nil {
				panic(err)
			}
		}

		if file, err := os.Create(p); err == nil {
			if err := f.Render(file); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}
