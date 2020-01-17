package cmd

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/PapaCharlie/go-restli/internal/codegen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var Version string

func CodeGenerator() *cobra.Command {
	var specBytes []byte
	var outputDir string
	var schemaDir string

	cmd := &cobra.Command{
		Use:          "go-restli",
		SilenceUsage: true,
		Version:      Version,
		Args: func(_ *cobra.Command, args []string) error {
			if len(Jar) > 0 {
				if len(args) == 0 {
					return errors.New("go-restli: Must specify at least one restspec file")
				}
			} else {
				switch len(args) {
				case 0:
					stat, err := os.Stdin.Stat()
					if err != nil {
						return errors.Wrap(err, "go-restli: Could not stat stdin")
					}
					if (stat.Mode() & os.ModeCharDevice) != 0 {
						return errors.New("go-restli: No stdin and no spec file given")
					}
				case 1:
					if _, err := os.Stat(args[0]); err != nil {
						return errors.Wrap(err, "go-restli: Must specify a valid spec file")
					}
				default:
					return errors.New("go-restli: Too many arguments")
				}
			}

			if schemaDir == "" {
				return errors.New("go-restli: Must specify a schema dir")
			} else if _, err := os.Stat(schemaDir); err != nil {
				return errors.Wrap(err, "go-restli: Must specify a valid schema dir: %w")
			}

			return nil
		},
		PreRunE: func(_ *cobra.Command, args []string) (err error) {
			if len(Jar) > 0 {
				specBytes, err = ExecuteJar(schemaDir, args)
			} else {
				specBytes, err = ReadSpec(args)
			}
			return err
		},
		RunE: func(*cobra.Command, []string) error {
			return codegen.GenerateCode(specBytes, outputDir)
		},
	}

	if len(Jar) > 0 {
		cmd.Use += " REST_SPEC [REST_SPEC...]"
	} else {
		cmd.Use += " [SPEC_FILE]"
	}

	cmd.Flags().StringVarP(&codegen.PackagePrefix, "package-prefix", "p", "", "The namespace to prefix all generated "+
		"packages with (e.g. github.com/PapaCharlie/go-restli/generated)")
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "The directory in which to output the generated files")
	cmd.Flags().StringVarP(&schemaDir, "schema-dir", "s", "", "The directory that contains all the .pdsc/.pdl files "+
		"that may be needed")

	return cmd
}

func ExecuteJar(schemaDir string, restSpecs []string) ([]byte, error) {
	if len(Jar) == 0 {
		log.Panicln("No jar!")
	}

	r, err := gzip.NewReader(bytes.NewBuffer(Jar))
	if err != nil {
		return nil, err
	}

	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	c := exec.Command("java", append([]string{"-jar", f.Name(), schemaDir}, restSpecs...)...)
	stdout, err := c.Output()
	if err != nil {
		return nil, err
	}

	err = os.Remove(f.Name())
	if err != nil {
		return nil, err
	}

	return stdout, nil
}

func ReadSpec(args []string) ([]byte, error) {
	if len(args) == 0 {
		specBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, errors.Wrap(err, "go-restli: Could not read spec from stdin")
		}
		return specBytes, nil
	} else {
		specByes, err := ioutil.ReadFile(args[0])
		if err != nil {
			return nil, errors.Wrap(err, "go-restli: Could not open spec file")
		}
		return specByes, nil
	}
}
