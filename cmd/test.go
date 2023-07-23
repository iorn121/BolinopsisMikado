/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")
		check()
	},
}

func check() {
	path, _ := os.Getwd()
	dataPath := filepath.Join(path, "img")
	defaultPath := filepath.Join(dataPath, "BolinopsisMikado.jpg")
	file, err := os.Open(defaultPath)
	if err != nil {
		fmt.Printf("Error : %+v", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error : %+v", err)
		os.Exit(1)
	}

	h := img.Bounds().Dy()
	w := img.Bounds().Dx()
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			fmt.Println(img.At(i, j).RGBA())
		}
	}
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
