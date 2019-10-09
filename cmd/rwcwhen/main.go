package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"rwcwhen/rwcwhen"
)

var (
	country string
	group   string
)

func init() {
	rootCmd.Flags().StringVarP(&group, "group", "g", "", "group to search for")
	rootCmd.Flags().StringVarP(&country, "country", "c", "", "country to search for")
}

var rootCmd = &cobra.Command{
	Use:   "rwcwhen",
	Short: "A project to get the rugby world cup results",
	Long: `This is learning project on golang to assist in getting the results
			of the 2019 Rugby World Cup. It deals with both teams and groups 
			you can see the full documentation here: https://github.com/velcro-ie/RWCWhen`,
	Version: rwcwhen.GetVersion(),

	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	rwcwhen.Run(country, group)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
