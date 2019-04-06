package cli

import (
	"bytes"
	"fmt"
	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/models"
	"github.com/PapaCharlie/go-restli/codegen/schema"
	"github.com/dave/jennifer/jen"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	packagePrefix string
	outputDir     string
	files         = make(map[string]func() io.Reader)
)

func CodeGenerator() *cobra.Command {
	cmd := &cobra.Command{
		Use: "go-restli [flags] RESTLI_SPECS",
		Args: func(_ *cobra.Command, args []string) error {
			for _, f := range args {
				data, err := ioutil.ReadFile(f)
				if err != nil {
					return err
				}
				files[f] = func() io.Reader { return bytes.NewBuffer(data) }
			}
			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) (err error) {
			codegen.SetPackagePrefix(packagePrefix)
			if outputDir == "" {
				return fmt.Errorf("must specify an output directory")
			} else {
				outputDir, err = filepath.Abs(outputDir)
				if err != nil {
					return fmt.Errorf("illegal path: %v", err)
				}
			}
			return
		},
		RunE: func(*cobra.Command, []string) error {
			return run()
		},
	}

	cmd.Flags().StringVarP(&packagePrefix, "package-prefix", "p", "", "The namespace to prefix all generated packages "+
		"with (e.g. github.com/PapaCharlie/go-restli)")
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "The directory in which to output the generated files")

	return cmd
}

func run() error {
	var codeFiles []*codegen.CodeFile

	for filename, buf := range files {
		log.Println(filename)
		loadedModels, err := models.LoadModels(buf())
		if err != nil {
			log.Fatalf("could not load %s: %+v", filename, err)
		}

		for _, m := range loadedModels {
			if code := m.GenerateModelCode(); code != nil {
				code.SourceFilename = filename
				codeFiles = append(codeFiles, code)
			}
		}

		loadedResources, err := schema.LoadResources(buf())
		if err != nil {
			log.Fatalf("%s: %+v", filename, err)
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
	return nil
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
	f.Func().Id("TestAllImports").Params(jen.Op("*").Qual("testing", "T")).Block()

	path := filepath.Join(outputDir, codegen.GetPackagePrefix(), "all_imports_test.go")
	out, err := os.Create(path)
	check(err)
	check(f.Render(out))
	check(out.Close())
	check(os.Chmod(path, os.FileMode(0555)))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
