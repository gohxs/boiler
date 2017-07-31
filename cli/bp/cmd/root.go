package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/gohxs/boiler"
	"github.com/spf13/cobra"
)

const (
	bashCompleteFunc = `
__boiler_get_generator()
{
    if boiler_output=$(bp generators list); then
        out=$(echo "${boiler_output}")
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi

}
__custom_func() {
    case ${last_command} in
        bp_add)
            __boiler_get_generator 
            return
            ;;
        *)
            ;;
    esac
}
`
)

var (
	// RootCmd of application
	RootCmd = &cobra.Command{Use: os.Args[0], BashCompletionFunction: bashCompleteFunc}
	// Stdin for cli app
	Stdin   io.Reader = os.Stdin
	gboiler *boiler.Core
)

func init() {
	RootCmd.Flags().String("bashcompletion", "", "Generates bashcompletion helper")
	RootCmd.Run = func(cmd *cobra.Command, args []string) {
		fl := cmd.Flag("bashcompletion")
		if fl.Changed {
			// Specially reset commands generated from project path
			for _, c := range RootCmd.Commands() {
				if c.Name() == "add" {
					c.ResetCommands()
					break
				}
			}
			RootCmd.GenBashCompletionFile(fl.Value.String())
			return
		}
		RootCmd.Help()
	}
}

// Boiler Singleton?
func Boiler() *boiler.Core {
	var err error
	if gboiler != nil {
		return gboiler
	}
	gboiler, err = boiler.From(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return gboiler
}
