package main

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"go-restli/restli"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Panicln("Must specify at least one snapshot file")
	}
	task := restli.NewGenerator(os.Args[1], "generated", "typerefs")
	outputDir := os.Args[2]
	openFiles := map[string]*jen.File{}
	generatedTypes := map[string]bool{}

	for _, filename := range os.Args[3:] {
		task.DecodeSnapshotModels(filename)
		for _, m := range task.GeneratedTypes {
			fqn := restli.NsJoin(m.Namespace, m.Name)
			if generatedTypes[fqn] {
				continue
			} else {
				generatedTypes[fqn] = true
			}

			var f *jen.File
			var ok bool

			packagePath := strings.Replace(restli.NsJoin(task.GeneratedTypesNamespacePrefix, m.Namespace), restli.NamespaceSep, "/", -1)

			if f, ok = openFiles[packagePath]; !ok {
				f = jen.NewFilePath(packagePath)
				openFiles[packagePath] = f
			}
			f.Add(m.Definition...)
		}
	}

	for p, f := range openFiles {
		p := filepath.Join(outputDir, p, "types.go")
		fmt.Println(filepath.Abs(p))
		if err := os.MkdirAll(filepath.Dir(p), os.ModePerm); err != nil {
			panic(err)
		}
		if _, err := os.Stat(p); err == nil {
			if err := os.Remove(p); err != nil {
				panic(err)
			}
		}

		if file, err := os.Create(p); err == nil {
			if _, err := fmt.Fprintf(file, "%#v\n", f); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}
