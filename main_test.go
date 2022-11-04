package main

import (
	"flag"
	"testing"

	"github.com/PapaCharlie/go-restli/v2/cmd"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	flag.Parse()
	command := cmd.CodeGenerator(jar)
	if len(flag.Args()) == 0 {
		return
	}
	command.SetArgs(flag.Args())
	require.NoError(t, command.Execute())
}
