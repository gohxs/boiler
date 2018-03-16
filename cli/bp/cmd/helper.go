package cmd

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey"
	"github.com/gohxs/boiler"
	"github.com/gohxs/boiler/config"
	"github.com/spf13/cobra"
)

func flagOrAsk(cmd *cobra.Command, userVars []config.UserVar, data map[string]interface{}) (err error) {
	//in := bufio.NewReader(Stdin)

	for _, v := range userVars {
		fl := cmd.Flag(v.Name)
		if fl != nil && fl.Changed {
			data[v.Name] = fl.Value.String()
			continue
		}
		/*if v.Question == "" && v.Default == "" {
			return fmt.Errorf("No flag provided for '%s' and no default/question defined", v.Name)
		}*/
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
			question = v.Name
		}

		var answer string
		if len(v.Choices) > 0 {
			choices := make([]string, len(v.Choices))
			for i, c := range v.Choices {
				choices[i], err = boiler.ProcessString(c, data)
				if err != nil {
					return err
				}
			}
			survey.AskOne(&survey.Select{
				Message: question,
				Options: choices,
			},
				&answer,
				nil,
			)
			log.Println("Choosen:", answer)
			data[v.Name] = answer
			continue
		}
		// Perform question
		survey.AskOne(&survey.Input{
			Message: question + ":",
			Default: fmt.Sprintf("%s", value),
		},
			&answer,
			nil,
		)
		data[v.Name] = answer
	}

	return nil
}
