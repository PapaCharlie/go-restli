package cli

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/models"
	"github.com/PapaCharlie/go-restli/codegen/schema"
	"github.com/dave/jennifer/jen"
	"github.com/spf13/cobra"
)

var (
	packagePrefix string
	outputDir     string
	snapshotMode  bool
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
			return nil
		},
		RunE: func(*cobra.Command, []string) error {
			if snapshotMode {
				return run(models.LoadSnapshotModels, schema.LoadSnapshotResource)
			} else {
				return run(models.LoadModels, schema.LoadResources)
			}
		},
	}

	cmd.Flags().StringVarP(&packagePrefix, "package-prefix", "p", "", "The namespace to prefix all generated packages "+
		"with (e.g. github.com/PapaCharlie/go-restli/generated)")
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "The directory in which to output the generated files")
	cmd.Flags().BoolVarP(&snapshotMode, "snapshot-mode", "s", false, "Assumes input is in the shape of .snapshot.json "+
		"files (statically generated during normal java build)")

	return cmd
}

func run(modelLoader func(io.Reader) error, resourceLoader func(io.Reader) ([]*schema.Resource, error)) error {
	var filenames []string
	for f := range files {
		filenames = append(filenames, f)
	}

	var allResources []*schema.Resource

	for filename, buf := range files {
		log.Println(filename)
		err := modelLoader(buf())
		if err != nil {
			log.Fatalf("could not load %s: %+v", filename, err)
		}

		loadedResources, err := resourceLoader(buf())
		if err != nil {
			log.Fatalf("%s: %+v", filename, err)
		}
		allResources = append(allResources, loadedResources...)
	}

	var codeFiles []*codegen.CodeFile
	for _, m := range models.GetRegisteredModels() {
		if f := models.GenerateModelCode(m); f != nil {
			codeFiles = append(codeFiles, f)
		}
	}
	for _, r := range allResources {
		for _, f := range r.GenerateCode() {
			if f != nil {
				codeFiles = append(codeFiles, f)
			}
		}
	}

	codeFiles = deduplicateFiles(codeFiles)

	for _, code := range codeFiles {
		code.SourceFilenames = filenames
		file, err := code.Write(outputDir)
		if err != nil {
			log.Fatalf("Could not generate code for %+v: %+v", code, err)
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

	err := codegen.Write(filepath.Join(outputDir, codegen.GetPackagePrefix(), "all_imports_test.go"), f)
	if err != nil {
		log.Panicf("Could not write all imports file: %+v", err)
	}
}

func renderCode(s *jen.Statement) []byte {
	b := bytes.NewBuffer(nil)
	if err := s.Render(b); err != nil {
		log.Panicln(err)
	}
	return b.Bytes()
}

func deduplicateFiles(files []*codegen.CodeFile) []*codegen.CodeFile {
	idToFile := make(map[string]*codegen.CodeFile)

	for _, file := range files {
		id := file.Identifier()
		if existingFile, ok := idToFile[id]; ok {
			existingCode := renderCode(existingFile.Code)
			code := renderCode(file.Code)
			if !bytes.Equal(existingCode, code) {
				log.Fatalf("Conflicting defitions of %s: %s\n\n-----------\n\n%s",
					id, string(existingCode), string(code))
			}
		} else {
			idToFile[id] = file
		}
	}

	identifiers := make([]string, 0, len(idToFile))
	for id := range idToFile {
		identifiers = append(identifiers, id)
	}
	sort.Strings(identifiers)

	uniqueCodeFiles := make([]*codegen.CodeFile, 0, len(idToFile))
	for _, id := range identifiers {
		uniqueCodeFiles = append(uniqueCodeFiles, idToFile[id])
	}

	return uniqueCodeFiles
}
