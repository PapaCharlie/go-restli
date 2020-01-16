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

	cmd := &cobra.Command{
		Use:          "go-restli",
		SilenceUsage: true,
		Version:      Version,
		PreRunE: func(_ *cobra.Command, args []string) (err error) {
			if len(Jar) > 0 {
				specBytes, err = ExecuteJar(args)
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
		cmd.Use += " PEGASUS_DIR REST_SPEC,[REST_SPEC...]"
	} else {
		cmd.Use += " [SPEC_FILE]"
	}

	cmd.Flags().StringVarP(&codegen.PackagePrefix, "package-prefix", "p", "", "The namespace to prefix all generated "+
		"packages with (e.g. github.com/PapaCharlie/go-restli/generated)")
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "The directory in which to output the generated files")

	return cmd
}

func ExecuteJar(args []string) ([]byte, error) {
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

	c := exec.Command("java", append([]string{"-jar", f.Name()}, args...)...)
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
		stat, err := os.Stdin.Stat()
		if err != nil {
			return nil, errors.Wrap(err, "go-restli: Could not stat stdin")
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return nil, errors.New("go-restli: No stdin and no spec file given")
		}

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
