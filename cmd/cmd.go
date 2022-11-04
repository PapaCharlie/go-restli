package cmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PapaCharlie/go-restli/v2/codegen/types"
	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var Version string

type JarStdinParameters struct {
	PackageRoot  string   `json:"packageRoot,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	Inputs       []string `json:"inputs,omitempty"`
	RawRecords   []string `json:"rawRecords,omitempty"`
}

func CodeGenerator(jar []byte) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "go-restli INPUTS...",
		SilenceUsage: true,
		Version:      Version,
	}

	var outputDir string
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "The directory in which to output the generated files")

	var manifestDependencies []string
	cmd.Flags().StringArrayVarP(&manifestDependencies, "manifest-dependencies", "m", nil,
		`Files or directories that may contain other "`+utils.ManifestFile+`" manifest files that this manifest may `+
			`depend on. Note that this may simply be "$GOPATH" or the "vendor" directory after "go mod vendor" is run.`)

	var generateWithPackageRoot bool
	cmd.Flags().BoolVar(&generateWithPackageRoot, "generate-with-package-root", false,
		"If specified, the generated files will be generated with the package root directory structure.")

	var namespaceAllowList []string
	const namespaceAllowListFlag = "namespace-allow-list"
	cmd.Flags().StringArrayVar(&namespaceAllowList, namespaceAllowListFlag, nil,
		"HIDDEN FLAG, USE AT YOUR OWN RISK: if provided, any data type whose namespace is not in this list will not "+
			"be generated")
	cmd.Flag(namespaceAllowListFlag).Hidden = true

	var manifestBytes []byte

	var jarParams JarStdinParameters

	cmd.Short = strings.TrimSpace(`
This standalone executable will parse all the .pdsc/.pdl and .restspec.json
files in the given inputs and produce bindings for each model and resource.
Inputs can be directories, files or JARs`)
	cmd.Flags().StringVarP(&jarParams.PackageRoot, "package-root", "p", "",
		"All files will be generated as sub-packages of this package.")
	cmd.Flags().StringArrayVar(&jarParams.RawRecords, "raw-records", nil,
		"These records will be interpreted as `restli.RawRecord`s instead of their actual underlying type.")

	const dependenciesFlag = "dependencies"
	cmd.Flags().StringArrayVarP(&jarParams.Dependencies, dependenciesFlag, "d", nil,
		"The directories, files or JARs that contains all the PDSC/PDL schema definitions required to "+
			"generate the inputs.")
	var dependenciesFile string
	const dependenciesFileFlag = "dependencies-file"
	cmd.Flags().StringVar(&dependenciesFile, dependenciesFileFlag, "",
		"If specified, each line in the file will be treated as if it was provided via --"+dependenciesFlag)

	var inputsFile string
	const inputsFileFlag = "inputs-file"
	cmd.Flags().StringVar(&inputsFile, inputsFileFlag, "",
		"If specified, each line in the file will be treated as if it was provided as an INPUT")

	cmd.Args = func(_ *cobra.Command, args []string) error {
		readFileLines := func(file string) (files []string, err error) {
			data, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}

			for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
				line = strings.TrimSpace(line)
				if len(line) == 0 {
					continue
				}

				_, err = os.Stat(line)
				if err != nil {
					return nil, err
				}

				files = append(files, line)
			}

			return files, nil
		}

		dedupe := func(strings []string) (deduped []string) {
			seen := map[string]bool{}
			for _, s := range strings {
				if !seen[s] {
					deduped = append(deduped, s)
					seen[s] = true
				}
			}
			return deduped
		}

		if cmd.Flags().Changed(dependenciesFileFlag) {
			files, err := readFileLines(dependenciesFile)
			if err != nil {
				return err
			}
			jarParams.Dependencies = dedupe(append(jarParams.Dependencies, files...))
		}

		jarParams.Inputs = append(jarParams.Inputs, args...)
		if cmd.Flags().Changed(inputsFileFlag) {
			files, err := readFileLines(inputsFile)
			if err != nil {
				return err
			}
			jarParams.Inputs = dedupe(append(jarParams.Inputs, files...))
		}

		return nil
	}

	cmd.RunE = func(*cobra.Command, []string) (err error) {
		manifestBytes, err = ExecuteJar(jar, jarParams)
		if err != nil {
			return err
		}

		var manifests []*GoRestliManifest
		for _, d := range manifestDependencies {
			err = filepath.WalkDir(d, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if strings.HasSuffix(path, utils.ManifestFile) {
					var data []byte
					data, err = os.ReadFile(path)
					if err != nil {
						return err
					}

					// Reading the dependent manifests is enough to populate the TypeRegistry. There is no need to
					// capture the manifests directly.
					var m *GoRestliManifest
					m, err = ReadManifest(data)
					if err != nil {
						return errors.Wrapf(err, "go-restli: Could not deserialize manifest from %q", path)
					}
					manifests = append(manifests, m)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		inputManifest, err := ReadManifest(manifestBytes)
		if err != nil {
			return err
		}

		manifests = append(manifests, inputManifest)

		return GenerateCode(outputDir, manifests, generateWithPackageRoot, namespaceAllowList)
	}

	return cmd
}

func ExecuteJar(jar []byte, params JarStdinParameters) ([]byte, error) {
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrapf(err, "go-restli: Failed to serialize %+v", params)
	}

	f, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, bytes.NewReader(jar))
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

func GenerateCode(
	outputDir string,
	manifests []*GoRestliManifest,
	generateWithPackageRoot bool,
	namespaceAllowList []string,
) error {
	err := utils.CleanTargetDir(outputDir)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Could not clean up output dir %q", outputDir)
	}

	inputManifest := manifests[len(manifests)-1]
	if generateWithPackageRoot {
		outputDir = filepath.Join(outputDir, inputManifest.PackageRoot)
	}

	err = RegisterManifests(manifests)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Could not register all manifest types")
	}

	err = LocateCustomTyperefs(inputManifest, outputDir)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Failed to locate custom typerefs")
	}

	manifestBytes, err := json.MarshalIndent(inputManifest, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "go-restli: Could not marshal manifest")
	}

	_ = os.MkdirAll(outputDir, os.ModePerm)
	manifestFile := filepath.Join(outputDir, utils.ManifestFile)
	err = os.WriteFile(manifestFile, manifestBytes, utils.ReadOnlyPermissions)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Failed to write parsed manifest to %q", manifestFile)
	}
	log.Printf("Wrote manifest to: %q", manifestFile)

	allowList := map[string]bool{}
	for _, ns := range namespaceAllowList {
		allowList[ns] = true
	}

	var codeFiles []*utils.CodeFile
	for id := range utils.TypeRegistry.TypesInPackageRoot(inputManifest.PackageRoot) {
		if ok := allowList[id.Namespace]; len(allowList) > 0 && !ok {
			continue
		}
		t := id.Resolve()
		if code := t.GenerateCode(); code != nil {
			codeFiles = append(codeFiles, &utils.CodeFile{
				SourceFile:  t.GetSourceFile(),
				PackagePath: t.GetIdentifier().PackagePath(),
				PackageRoot: t.GetIdentifier().PackageRoot(),
				Filename:    t.GetIdentifier().TypeName(),
				Code:        t.GenerateCode(),
			})
		}
	}

	codeFiles = append(codeFiles, inputManifest.GenerateResourceCode()...)

	for _, code := range codeFiles {
		err = code.Write(outputDir)
		if err != nil {
			return errors.Wrapf(err, "go-restli: Could not generate code for %+v:\n%s", code, code.Code.GoString())
		}
	}

	err = GenerateCustomTyperefInit(outputDir, inputManifest)
	if err != nil {
		return err
	}

	err = GenerateAllImportsTest(outputDir, codeFiles)
	if err != nil {
		return err
	}

	return nil
}

func GenerateCustomTyperefInit(outputDir string, manifest *GoRestliManifest) (err error) {
	customTyperefs := map[string]utils.IdentifierSet{}

	for _, dt := range manifest.InputDataTypes {
		if dt.Typeref == nil || !dt.Typeref.IsCustom {
			continue
		}
		pkg := dt.Typeref.PackagePath()
		if _, ok := customTyperefs[pkg]; !ok {
			customTyperefs[pkg] = utils.IdentifierSet{}
		}
		customTyperefs[pkg].Add(dt.Typeref.Identifier)
	}

	for pkg, refs := range customTyperefs {
		f := jen.NewFilePathName(pkg, utils.PackageName(pkg))
		utils.AddImportNames(f)
		f.HeaderComment(utils.DoNotEditHeader)

		f.Func().Id("init").Call().BlockFunc(func(def *jen.Group) {
			refs.Range(func(id utils.Identifier) {
				def.Add(types.CallRegisterCustomTyperef(id))
			})
		})

		log.Println(pkg, outputDir, strings.TrimPrefix(pkg, manifest.PackageRoot))
		name := filepath.Join(outputDir, strings.TrimPrefix(pkg, manifest.PackageRoot), "init_custom_typerefs.gr.go")
		err = utils.WriteJenFile(name, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func GenerateAllImportsTest(outputDir string, codeFiles []*utils.CodeFile) (err error) {
	imports := make(map[string]bool)
	for _, code := range codeFiles {
		if code == nil {
			continue
		}
		imports[code.PackagePath] = true
	}
	f := jen.NewFile("main")
	f.HeaderComment(utils.DoNotEditHeader)

	for p := range imports {
		f.Anon(p)
	}
	f.Func().Id("TestAllImports").Params(jen.Op("*").Qual("testing", "T")).Block()

	out := filepath.Join(outputDir, "all_imports_test"+utils.GeneratedFileSuffix)
	err = utils.WriteJenFile(out, f)
	if err != nil {
		return errors.Wrapf(err, "Could not write all imports file: %+v", err)
	}
	return nil
}
