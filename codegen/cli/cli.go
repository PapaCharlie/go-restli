package cli

import (
	"fmt"
	"path/filepath"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/schema"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func CodeGenerator() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:          "go-restli [flags] RESTLI_SPECS",
		SilenceUsage: true,
		PreRunE: func(*cobra.Command, []string) (err error) {
			if outputDir == "" {
				return fmt.Errorf("must specify an output directory")
			} else {
				outputDir, err = filepath.Abs(outputDir)
				if err != nil {
					return fmt.Errorf("illegal path: %v", err)
				}
			}
			if codegen.PdscDirectory == "" {
				return fmt.Errorf("must specify a pdsc directory")
			} else {
				codegen.PdscDirectory, err = filepath.Abs(codegen.PdscDirectory)
				if err != nil {
					return fmt.Errorf("illegal path: %v", err)
				}
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			return Run(args, outputDir)
		},
	}

	cmd.Flags().StringVarP(&codegen.PdscDirectory, "pdsc-dir", "d", "", "The directory containing all the pdsc files "+
		"required to generate the client bindings")
	cmd.Flags().StringVarP(&codegen.PackagePrefix, "package-prefix", "p", "", "The namespace to prefix all generated "+
		"packages with (e.g. github.com/PapaCharlie/go-restli/generated)")
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "The directory in which to output the generated files")

	return cmd
}

func Run(restSpecs []string, outputDir string) error {
	resources, types, err := schema.LoadRestSpecs(restSpecs)
	if err != nil {
		return err
	}

	var codeFiles []*codegen.CodeFile
	for _, m := range types {
		codeFiles = append(codeFiles, m.GenerateModelCode())
	}

	for _, r := range resources {
		for _, f := range r.GenerateCode() {
			if f != nil {
				codeFiles = append(codeFiles, f)
			}
		}
	}

	codeFiles = codegen.DeduplicateFiles(codeFiles)

	for _, code := range codeFiles {
		file, err := code.Write(outputDir)
		if err != nil {
			return errors.Wrapf(err, "Could not generate code for %+v", code)
		} else {
			fmt.Println(file)
		}
	}

	codegen.GenerateAllImportsFile(outputDir, codeFiles)
	return nil
}
