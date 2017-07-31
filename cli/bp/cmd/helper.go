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
			continue
		}
		if v.Question == "" && v.Default == "" {
			return fmt.Errorf("No flag provided for '%s' and no default/question defined", v.Name)
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
		if question == "" { // will not ask and try to set default
			data[v.Name] = value
			return nil
		}

		if value != "" {
			fmt.Printf("%s (%s): ", question, value)
		} else {
			fmt.Printf("%s: ", question)

		}
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
