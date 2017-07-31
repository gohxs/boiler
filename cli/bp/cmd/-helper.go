package cmd

import (
	"bufio"
	"fmt"

	"github.com/gohxs/boiler"
	"github.com/gohxs/boiler/config"
	"github.com/spf13/cobra"
)

func flagOrAsk(cmd *cobra.Command, userVars []config.UserVar, data map[string]interface{}) (err error) {
	in := bufio.NewReader(Stdin)

	for _, v := range userVars {
		fl := cmd.Flag(v.Name)
		if fl != nil && fl.Changed {
			data[v.Name] = fl.Value.String()
			return
		}
		question, err := boiler.ProcessString(v.Question, data)
		if err != nil {
			return err
		}
		value, ok := data[v.Name] // First we check on data, then on Processed default
		if !ok {
			value, err = boiler.ProcessString(v.Default, data)
			if err != nil {
				return err
			}
		}
		fmt.Printf("%s [%s] (%s)? ", question, v.Name, value)
		line, _, _ := in.ReadLine()
		if len(line) != 0 {
			value = string(line)
		}
		if value != "" {
			data[v.Name] = value
		}
	}

	return nil
}
