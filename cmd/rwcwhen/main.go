package main

import{

"github.com/spf13/cobra"
}

var (
	country       string
	group    	string
)

func init() {
	rootCmd.Flags().StringArrayVarP(&country, "country", "c", "", "country to search for")
	rootCmd.Flags().StringVarP(&group, "group", "g", "", "group to search for")
}