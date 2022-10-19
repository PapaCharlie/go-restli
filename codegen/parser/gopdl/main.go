package main

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PapaCharlie/go-restli/codegen/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

func main() {
	for _, f := range loadPDLs(os.Args[1:]...) {
		lexer := parser.NewPdlLexer(f)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.NewPdlParser(stream)
		p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
		p.BuildParseTrees = true
		tree := p.Document()
		antlr.ParseTreeWalkerDefault.Walk(new(parser.CustomListener), tree)
	}
	start := time.Now()
	pdls := loadPDLs(os.Args[1:]...)
	log.Printf("Loaded %d PDLs in %s", len(pdls), time.Since(start))
}

func loadPDLs(inputs ...string) (streams []*antlr.InputStream) {
	for _, input := range inputs {
		err := filepath.WalkDir(input, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			switch filepath.Ext(path) {
			case ".jar":
				r, err := zip.OpenReader(path)
				if err != nil {
					log.Panicln(err)
				}
				defer r.Close()
				files := map[string]*zip.File{}
				for _, zipFile := range r.File {
					files[zipFile.Name] = zipFile
				}
				for name, zipFile := range files {
					switch filepath.Ext(name) {
					case ".pdl":
						r, err := zipFile.Open()
						if err != nil {
							log.Panicln(err)
						}
						defer r.Close()
						data, err := ioutil.ReadAll(r)
						if err != nil {
							log.Panicln(err)
						}
						streams = append(streams, antlr.NewInputStream(string(data)))
					case ".pdsc":
						segments := strings.Split(name, "/")
						if segments[0] == "legacyPegasusSchemas" {
							newPath := filepath.Join("pegasus", filepath.Join(segments[1:]...))
							newPath = strings.TrimSuffix(newPath, ".pdsc") + ".pdl"
							if _, ok := files[newPath]; ok {
								continue
							}
						}
						fmt.Printf("Ignored pdsc: \"%s!%s\"\n", path, zipFile.Name)
					}
				}
			case ".pdl":
				fmt.Println(path)
				data, err := os.ReadFile(path)
				if err != nil {
					log.Panicln(err)
				}
				streams = append(streams, antlr.NewInputStream(string(data)))
			case ".pdsc":
				fmt.Printf("Ignored pdsc: %q\n", path)
			}
			return nil
		})
		if err != nil {
			log.Panicln(err)
		}
	}

	return streams
}
