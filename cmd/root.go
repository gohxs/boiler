package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

const (
	bash_completion_func = `
__boiler_get_generator()
{
    if boiler_output=$(boiler generators); then
        out=$(echo "${boiler_output}")
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi

}
__custom_func() {
    case ${last_command} in
        boiler_add)
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
	RootCmd = &cobra.Command{Use: os.Args[0], BashCompletionFunction: bash_completion_func}
	// Stdin for cli app
	Stdin io.Reader = os.Stdin
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
