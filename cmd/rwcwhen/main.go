package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	country  string
	group    string
	upcoming bool
)

func init() {
	rootCmd.Flags().StringVarP(&group, "group", "g", "", "group to search for")
	rootCmd.Flags().StringVarP(&country, "country", "c", "", "country to search for")
	rootCmd.Flags().BoolVarP(&upcoming, "upcoming", "u", false, "return upcoming matches in a group")
}

var rootCmd = &cobra.Command{
	Use:   "rwcwhen",
	Short: "A project to get the rugby world cup results",
	Long: `This is learning project on golang to assist in getting the results
			of the 2019 Rugby World Cup. It deals with both teams and groups 
			you can see the full documentation here: https://github.com/velcro-ie/RWCWhen`,
	Version: GetVersion(),

	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	err := RunAll(country, group, upcoming)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
