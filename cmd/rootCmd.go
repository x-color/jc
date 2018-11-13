package cmd

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/x-color/jc/formatter"
	"github.com/x-color/jc/parser"
)

var coloring bool

func filter(list []string, checker func(string) bool) []string {
	filteredList := []string{}
	for _, str := range list {
		if checker(str) {
			filteredList = append(filteredList, str)
		}
	}
	return filteredList
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "jc FILTER",
		Long: "jc is a tool to format, color and print JSON input.",
		Args: cobra.ExactArgs(1),
		Example: `  cat data.json | jc .
  cat data.json | jc .foo
  cat data.json | jc .foo[0]
  cat data.json | jc .[1].bar`,
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			data, err := parser.ParseJSON(b)
			if err != nil {
				return err
			}
			keys, err := parser.ParseKeys(filter(strings.Split(args[0], "."),
				func(str string) bool {
					return !(str == "")
				}))
			if err != nil {
				return err
			}
			data, err = parser.ChoiceFromJSON(data, keys)
			if err != nil {
				return err
			}
			s := formatter.JSONToString(data, coloring)
			cmd.Println(s)
			return nil
		},
	}

	cmd.Flags().BoolVar(&coloring, "color", true, "color JSON object")

	return cmd
}

func Execute() {
	cmd := rootCmd()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}
}
