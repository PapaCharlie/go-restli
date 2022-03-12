package main

import (
	"flag"
	"testing"

	"github.com/PapaCharlie/go-restli/cmd"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	flag.Parse()
	command := cmd.CodeGenerator()
	command.SetArgs(flag.Args())
	require.NoError(t, command.Execute())
}
