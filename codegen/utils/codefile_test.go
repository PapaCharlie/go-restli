package utils

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanTargetDir(t *testing.T) {
	const targetDir = "a/"
	tests := []struct {
		Name          string
		InitialFiles  []string
		ExpectedFiles []string
	}{
		{
			Name: "Simple",
			InitialFiles: []string{
				targetDir,
				targetDir + "b/",
				targetDir + "b/c/",
				targetDir + "b/c/foo" + GeneratedFileSuffix,
			},
			ExpectedFiles: []string{},
		},
		{
			Name: "Multiple empty subdirs",
			InitialFiles: []string{
				targetDir,
				targetDir + "b/",
				targetDir + "b/c/",
				targetDir + "b/c/foo" + GeneratedFileSuffix,
				targetDir + "b/d/",
				targetDir + "b/d/bar" + GeneratedFileSuffix,
			},
			ExpectedFiles: []string{},
		},
		{
			Name: "Does not delete custom files",
			InitialFiles: []string{
				targetDir,
				targetDir + "foo.go",
				targetDir + "garbage.zip",
			},
			ExpectedFiles: []string{
				targetDir,
				targetDir + "foo.go",
				targetDir + "garbage.zip",
			},
		},
		{
			Name: "Only deletes dir if truly empty",
			InitialFiles: []string{
				targetDir,
				targetDir + "foo" + GeneratedFileSuffix,
				targetDir + "garbage.zip",
			},
			ExpectedFiles: []string{
				targetDir,
				targetDir + "garbage.zip",
			},
		},
		{
			Name:          "Initial target does not exist",
			InitialFiles:  []string{},
			ExpectedFiles: []string{},
		},
		{
			Name:          "Initial target empty",
			InitialFiles:  []string{targetDir},
			ExpectedFiles: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			tmp := t.TempDir() + "/"
			target := filepath.Join(tmp, targetDir)
			for _, f := range test.InitialFiles {
				out := filepath.Join(tmp, f)
				if strings.HasSuffix(f, "/") {
					require.NoError(t, os.Mkdir(out, os.ModePerm))
				} else {
					require.NoError(t, ioutil.WriteFile(out, []byte{}, ReadOnlyPermissions))
				}
			}

			t.Log(tree(t, tmp))

			require.NoError(t, CleanTargetDir(target))

			if len(test.ExpectedFiles) == 0 {
				_, err := os.Stat(target)
				require.True(t, os.IsNotExist(err))
			} else {
				require.Equal(t, test.ExpectedFiles, tree(t, tmp))
			}
		})
	}
}

func tree(t *testing.T, target string) (files []string) {
	require.NoError(t, filepath.WalkDir(target, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == target {
			return nil
		}

		path = strings.TrimPrefix(path, target)
		if d.IsDir() {
			path += "/"
		}
		files = append(files, path)
		return nil
	}))
	return files
}
