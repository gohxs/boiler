package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	// RootCmd of application
	RootCmd = &cobra.Command{Use: os.Args[0]}
	// Stdin for cli app
	Stdin io.Reader = os.Stdin
)
