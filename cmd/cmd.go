package cmd

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var Version string

type JarStdinParameters struct {
	ResolverPath               string                         `json:"resolverPath"`
	RestSpecPaths              []string                       `json:"restSpecPaths"`
	NamedDataSchemasToGenerate []string                       `json:"namedDataSchemasToGenerate"`
	RawRecords                 []string                       `json:"rawRecords"`
	NativeTyperefs             map[string]types.NativeTyperef `json:"nativeTyperefs"`
}

func CodeGenerator() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "go-restli",
		SilenceUsage: true,
		Version:      Version,
	}

	var specBytes []byte
	var outputDir string

	cmd.Flags().StringVarP(&utils.PackagePrefix, "package-prefix", "p", "",
		"All files will be generated as sub-packages of this package.")
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "The directory in which to output the generated files")
	var typerefBindingPluginPath string
	cmd.Flags().StringVar(&typerefBindingPluginPath, "typeref-binding-plugin", "",
		"If specified, the plugin file provided by this flag will be loaded and called for each typeref so custom "+
			"bindings can be generated for typerefs.")

	if len(Jar) > 0 {
		var params JarStdinParameters

		cmd.Use += " REST_SPEC [REST_SPEC...]"
		cmd.Short = strings.TrimSpace(`
This standalone executable will parse all the .pdsc/.pdl files in -r/--resolver-path and produce
bindings for the given rest specs, and all necessary associated schemas. By default, bindings are
generated only for the schemas required to interact with the given resources, but this behavior can
be overridden with -n/--named-schemas-to-generate which can be used to specify some extra schemas
that should also have bindings. If named schemas are specified, it's not necessary to specify rest
specs.`)
		cmd.Flags().StringVarP(&params.ResolverPath, "resolver-path", "r", "",
			"The directory that contains all the .pdsc/.pdl files that may be needed")
		cmd.Flags().StringArrayVarP(&params.NamedDataSchemasToGenerate, "named-schemas-to-generate", "n", nil,
			"Bindings for these schemas will be generated (can be used without .restspec.json files)")
		cmd.Flags().StringArrayVar(&params.RawRecords, "raw-records", nil,
			"These records will be interpreted as `protocol.RawRecord`s instead of their actual underlying type.")
		var nativeTyperefs string
		cmd.Flags().StringVar(&nativeTyperefs, "native-typerefs", "",
			"This parameter expects a file containing a JSON map of fully classified name to a native typeref.")

		cmd.Args = func(_ *cobra.Command, args []string) (err error) {
			params.RestSpecPaths = args
			if len(params.RestSpecPaths) == 0 && len(params.NamedDataSchemasToGenerate) == 0 {
				return errors.New("go-restli: Must specify at least one restspec file or named data schema")
			}

			if params.ResolverPath == "" {
				return errors.New("go-restli: Must specify a schema dir")
			} else if _, err = os.Stat(params.ResolverPath); err != nil {
				return errors.Wrap(err, "go-restli: Must specify a valid schema dir: %w")
			}

			if nativeTyperefs != "" {
				var f *os.File
				f, err = os.Open(nativeTyperefs)
				if err != nil {
					return err
				}
				err = json.NewDecoder(f).Decode(&params.NativeTyperefs)
				if err != nil {
					return err
				}
				for _, v := range params.NativeTyperefs {
					// Ignore whatever was passed in for these two parameters as they will be provided by the spec
					// parser
					v.OriginalTypeName = nil
					v.Primitive = nil
				}
			}

			return nil
		}

		cmd.PreRunE = func(_ *cobra.Command, args []string) (err error) {
			specBytes, err = ExecuteJar(params)
			return err
		}
	} else {
		cmd.Use += " [SPEC_FILE]"
		cmd.Short = "Generate rest.li bindings for the given parsed specs"

		cmd.Args = func(_ *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				stat, err := os.Stdin.Stat()
				if err != nil {
					return errors.Wrap(err, "go-restli: Could not stat stdin")
				}
				if (stat.Mode() & os.ModeCharDevice) != 0 {
					return errors.New("go-restli: No stdin and no spec file given")
				}
				return nil
			case 1:
				if _, err := os.Stat(args[0]); err != nil {
					return errors.Wrap(err, "go-restli: Must specify a valid spec file")
				}
				return nil
			default:
				return errors.New("go-restli: Too many arguments")
			}
		}

		cmd.PreRunE = func(_ *cobra.Command, args []string) (err error) {
			specBytes, err = ReadSpec(args)
			return err
		}
	}

	cmd.RunE = func(*cobra.Command, []string) (err error) {
		if typerefBindingPluginPath != "" {
			log.Printf("Loading plugin from %q", typerefBindingPluginPath)
			types.TyperefBindingPlugin, err = plugin.Open(typerefBindingPluginPath)
			if err != nil {
				return err
			}
		}

		return GenerateCode(specBytes, outputDir)
	}

	return cmd
}

func ExecuteJar(params JarStdinParameters) ([]byte, error) {
	if len(Jar) == 0 {
		log.Panicln("No jar!")
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrapf(err, "go-restli: Failed to serialize %+v", params)
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

	c := exec.Command("java", "-jar", f.Name())
	c.Stdin = bytes.NewReader(paramsBytes)
	c.Stderr = os.Stderr
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

func GenerateCode(specBytes []byte, outputDir string) error {
	var schemas GoRestliSpec

	_ = os.MkdirAll(outputDir, os.ModePerm)
	parsedSpecs := filepath.Join(outputDir, "parsed-specs.json")
	_ = os.Remove(parsedSpecs)
	err := ioutil.WriteFile(parsedSpecs, specBytes, utils.ReadOnlyPermissions)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Failed to write parsed specs to %q", parsedSpecs)
	}

	// Use a Reader regardless since it'll handle leading/trailing whitespace and other niceties
	err = json.NewDecoder(bytes.NewBuffer(specBytes)).Decode(&schemas)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Could not deserialize GoRestliSpec")
	}

	tmpOutputDir, err := ioutil.TempDir("", "go-restli_*")
	if err != nil {
		return errors.Wrapf(err, "go-restli: Failed to create temporary directory")
	}
	defer os.RemoveAll(tmpOutputDir)

	codeFiles := append(utils.TypeRegistry.GenerateTypeCode(), schemas.GenerateClientCode()...)

	for _, code := range codeFiles {
		err = code.Write(tmpOutputDir)
		if err != nil {
			return errors.Wrapf(err, "go-restli: Could not generate code for %+v:\n%s", code, jen.Add(code.Code).GoString())
		}
	}

	err = GenerateAllImportsTest(tmpOutputDir, codeFiles)
	if err != nil {
		return err
	}

	children, err := ioutil.ReadDir(tmpOutputDir)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Could not list %q", tmpOutputDir)
	}

	for _, c := range children {
		source := filepath.Join(tmpOutputDir, c.Name())
		destination := filepath.Join(outputDir, c.Name())

		err = os.RemoveAll(destination)
		if err != nil {
			return errors.Wrapf(err, "go-restli: Failed to delete %q", destination)
		}

		err = os.Rename(source, destination)
		if err != nil {
			return errors.Wrapf(err, "go-restli: Failed to move %q to %q", source, destination)
		}
	}

	return nil
}

func GenerateAllImportsTest(outputDir string, codeFiles []*utils.CodeFile) error {
	imports := make(map[string]bool)
	for _, code := range codeFiles {
		if code == nil {
			continue
		}
		imports[code.PackagePath] = true
	}
	f := jen.NewFile("main")
	for p := range imports {
		f.Anon(p)
	}
	f.Func().Id("TestAllImports").Params(jen.Op("*").Qual("testing", "T")).Block()

	out := filepath.Join(outputDir, "all_imports_test.go")
	_ = os.Remove(out)
	err := utils.WriteJenFile(out, f)
	if err != nil {
		return errors.Wrapf(err, "Could not write all imports file: %+v", err)
	}
	return nil
}
