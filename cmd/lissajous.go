/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// lissajousCmd represents the lissajous command
var lissajousCmd = &cobra.Command{
	Use:   "lissajous",
	Short: "make lissajous gif",
	Long:  `Create a GIF animation of a random Lissajous shape. args [path] are the path to the gif file created by lissajous.`,
	Run: func(cmd *cobra.Command, paths []string) {
		make_lissajous(paths)
	},
}

func init() {
	rootCmd.AddCommand(lissajousCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lissajousCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lissajousCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0
	blackIndex = 1
)

func make_lissajous(paths []string) {
	for _, path := range paths {
		rand.Seed(time.Now().UnixNano())
		lissajous(path)
	}
}

func lissajous(path string) {
	rand.Seed(time.Now().UnixNano())
	const (
		cycles  = 5
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 8
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	if path != "" {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		fmt.Println(path)
		gif.EncodeAll(f, &anim)
	} else {
		gif.EncodeAll(os.Stdout, &anim)
	}
}
