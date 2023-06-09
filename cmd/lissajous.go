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
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// lissajousCmd represents the lissajous command
var lissajousCmd = &cobra.Command{
	Use:   "lissajous [path1, path2,...]",
	Short: "make lissajous gif",
	Long:  `Create a GIF animation of a random Lissajous shape. args [path] are the path to the gif file created by lissajous.`,
	Run: func(cmd *cobra.Command, paths []string) {
		make_lissajous(paths)
	},
}

func init() {
	rootCmd.AddCommand(lissajousCmd)
}

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0
	blackIndex = 1
)

func make_lissajous(paths []string) {
	if len(paths) == 0 {
		path := "lissajous.gif"
		lissajous(path)
	}
	for _, path := range paths {
		ext := filepath.Ext(path)
		ext = strings.ToLower(ext)
		if ext == ".gif" {
			rand.Seed(time.Now().UnixNano())
			lissajous(path)
		} else {
			log.Fatal("Unsupported file extension")
		}
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

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Println(path)
	gif.EncodeAll(f, &anim)
}
