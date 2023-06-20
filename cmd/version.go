/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of BolinopsisMikado",
	Long:  "All software has versions. This is BolinopsisMikado's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version 1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
