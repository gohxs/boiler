package cliapp

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gohxs/boiler/internal/config"
	"github.com/gohxs/boiler/internal/core"
	"github.com/urfave/cli"
)

func flagOrAsk(cl *cli.Context, userVars []config.UserVar, data map[string]interface{}) (err error) {

	in := bufio.NewReader(os.Stdin)
	for _, v := range userVars {
		if cl.IsSet(v.Name) { // We have a flag, continue
			data[v.Name] = cl.String(v.Name)
			continue
		}
		question, err := core.ProcessString(v.Question, data)
		if err != nil {
			return err
		}
		value, err := core.ProcessString(v.Default, data)
		if err != nil {
			return err
		}
		if question == "" {
			question = "Type value for:"
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
