package cli

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/models"
	"github.com/PapaCharlie/go-restli/codegen/schema"
	"github.com/spf13/cobra"
)

func CodeGenerator() *cobra.Command {
	var outputDir string
	var packagePrefix string

	cmd := &cobra.Command{
		Use: "go-restli [flags] RESTLI_SPECS",
		PreRunE: func(*cobra.Command, []string) (err error) {
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
		RunE: func(_ *cobra.Command, args []string) error {
			return Run(args, outputDir, packagePrefix)
		},
	}

	cmd.Flags().StringVarP(&packagePrefix, "package-prefix", "p", "", "The namespace to prefix all generated packages "+
		"with (e.g. github.com/PapaCharlie/go-restli/generated)")
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "The directory in which to output the generated files")

	return cmd
}

func Run(files []string, outputDir string, packagePrefix string) error {
	codegen.SetPackagePrefix(packagePrefix)

	var allResources []*schema.Resource

	for _, filename := range files {
		log.Println(filename)
		snapshot, err := LoadSnapshotFromFile(filename)
		if err != nil {
			log.Fatalf("could not load %s: %+v", filename, err)
		}
		allResources = append(allResources, snapshot.Schema)
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

	codeFiles = codegen.DeduplicateFiles(codeFiles)

	for _, code := range codeFiles {
		code.SourceFilenames = files
		file, err := code.Write(outputDir)
		if err != nil {
			log.Fatalf("Could not generate code for %+v: %+v", code, err)
		} else {
			fmt.Println(file)
		}
	}

	codegen.GenerateAllImportsFile(outputDir, codeFiles)
	return nil
}
