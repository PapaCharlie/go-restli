package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/PapaCharlie/go-restli/codegen/cli"
)

const (
	restLiClientTestSuite  = "rest.li-test-suite/client-testsuite"
	generatedPackageSuffix = "generated"
	packagePrefix          = "github.com/PapaCharlie/go-restli/tests/" + generatedPackageSuffix
)

//go:generate go run .
func main() {
	snapshotsDir := filepath.Join(restLiClientTestSuite, "snapshots")
	if _, err := os.Stat(snapshotsDir); err != nil {
		log.Panicln(snapshotsDir, "does not exist! Did you use `git clone --recurse-submodules`?", err)
	}

	snapshots, err := filepath.Glob(filepath.Join(snapshotsDir, "/*.snapshot.json"))
	if err != nil {
		log.Panicln("Could list read", snapshotsDir, err)
	}
	if len(snapshots) == 0 {
		log.Panicln("No snapshots found in", snapshotsDir)
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Panicln("Failed to create temp directory", err)
	}

	err = cli.Run(snapshots, tmpDir, packagePrefix)
	if err != nil {
		log.Panicln("Could not generate code from snapshots", err)
	}

	_ = os.RemoveAll(generatedPackageSuffix)
	err = os.Rename(filepath.Join(tmpDir, packagePrefix), generatedPackageSuffix)
	if err != nil {
		log.Panicln("Failed to move the generated code", err)
	}

	_ = os.RemoveAll(tmpDir)
}
