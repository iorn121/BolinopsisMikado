/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
 	"bufio"
 	"image"
 	"image/color"
 	"image/draw"
 	"image/png"
 	"io/ioutil"
 	"log"
 	"os"
 	"github.com/golang/freetype/truetype"
 	"golang.org/x/image/font"
 	"golang.org/x/image/math/fixed"
	"github.com/spf13/cobra"
)

// makeImgCmd represents the makeImg command
var makeImgCmd = &cobra.Command{
	Use:   "makeImg",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("makeImg called")
	},
}

func init() {
	rootCmd.AddCommand(makeImgCmd)


 	ftBin, err := ioutil.ReadFile("./font/MPLUS1-VariableFont_wght.ttf")
 	if err != nil {
 		log.Fatalf("failed to load font: %s", err.Error())
 	}
 	ft, err := truetype.Parse(ftBin)
 	if err != nil {
 		log.Fatalf("failed to parse font: %s", err.Error())
 	}
	imageHeight, imageWidth := 100, 100
 	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
 	dst := image.NewRGBA(img.Bounds())
 	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)
	
 	text := "a"
	
 	col := color.RGBA{33, 33, 33, 1}
	
 	opt := truetype.Options{
 		Size: 20,
 	}
 	face := truetype.NewFace(ft, &opt)
	
 	x, y := 100, 100
 	dot := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 26)}
	
 	d := &font.Drawer{
 		Dst:  dst,
 		Src:  image.NewUniform(col),
 		Face: face,
 		Dot:  dot,
 	}
 	d.DrawString(text)
	
 	newFile, err := os.Create("output/sample.png")
 	if err != nil {
 		log.Fatalf("failed to create file: %s", err.Error())
 	}
 	defer newFile.Close()
	
 	b := bufio.NewWriter(newFile)
 	if err := png.Encode(b, dst); err != nil {
 		log.Fatalf("failed to encode image: %s", err.Error())
 	}
}