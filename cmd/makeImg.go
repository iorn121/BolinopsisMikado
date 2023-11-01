/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
 	// "bufio"
 	"image"
 	"image/color"
 	"image/draw"
 	"image/jpeg"
	// "image/png"
 	"io/ioutil"
	"os/user"
    "path/filepath"
	// "strings"
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
		isascii, _ := cmd.Flags().GetBool("ascii")
		message, _ := cmd.Flags().GetString("message")
		path, _ := cmd.Flags().GetString("path")

		if path == "" {
				path = getDefaultDownloadPath()+"/output.jpg"
				fmt.Printf("path is not specified, use default path: %s\n", path)
		} else if !pathExists(path) {
			panic(fmt.Errorf("path does not exist: %s", path))
		} else {
			path+="/output.jpg"
		}
		if isascii {
			writeAscii()
		} else {
			img:=writeImage(message)
			saveImage(img, path)
		}
	},
}

func init() {
	rootCmd.AddCommand(makeImgCmd)
	makeImgCmd.Flags().BoolP("ascii", "a", false, "save ascii image")
	makeImgCmd.Flags().StringP("message", "s", "hello, world!", "string to convert")
	makeImgCmd.Flags().StringP("path", "p", "", "image path to save")
	makeImgCmd.Flags().MarkHidden("path")
}


func writeAscii() {
	face := readFont()

	imageHeight, imageWidth := 100, 100

	for i := 33; i < 127; i++ {
		img := createBaseImage(imageHeight, imageWidth)
		ascii := string([]byte{byte(i)})
		img_text:=writeText(ascii , img , face)
		saveImage(img_text, fmt.Sprintf("ascii/%d.jpg", i))
		fmt.Printf("%c image generation started \n", i)
	}
}

func writeImage(message string) draw.Image {
	face := readFont()
	message_length := len(message)
	imageHeight, imageWidth := 100, 70 * message_length

	img := createBaseImage(imageHeight, imageWidth)
	img_text:=writeText(message , img , face)
	fmt.Printf("%s image generation started \n", message)
	return img_text
}

func writeText(message string, img draw.Image, face font.Face) draw.Image {
	col := color.RGBA{33, 33, 33, 1}

	dot := fixed.Point26_6{X: fixed.Int26_6(1200), Y: fixed.Int26_6(5000)}
	d := font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  dot,
	}
	d.DrawString(message)
	return img
}

func readFont() font.Face {
	ftBin, err := ioutil.ReadFile("./font/MPLUS1-VariableFont_wght.ttf")
	if err != nil {
		log.Fatalf("failed to load font: %s", err.Error())
	}
	ft, err := truetype.Parse(ftBin)
	if err != nil {
		log.Fatalf("failed to parse font: %s", err.Error())
	}
	opt := truetype.Options{
		Size: 100,
	}
	face := truetype.NewFace(ft, &opt)
	return face
}

func createBaseImage(imageHeight int, imageWidth int) draw.Image {
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	basicColor := color.RGBA{255, 255, 255, 1}
	base:= image.NewUniform(basicColor)
	draw.Draw(img, img.Bounds(), base, image.Pt(0, 0), draw.Src)
	return img
}

func saveImage(img image.Image, filepath string) {
	fso, err := os.Create(filepath)
	if err != nil {
		fmt.Println("create:", err)
		return
	}
	defer fso.Close()
	jpeg.Encode(fso, img, &jpeg.Options{Quality: 100})
	// img_type := strings.Split(filepath, ".")[1]
	// if img_type == "jpg" {
	// 	jpeg.Encode(fso, img, &jpeg.Options{Quality: 100})
	// } else if img_type == "png" {
	// 	png.Encode(fso, img)
	// } else {
	// 	fmt.Println("Error: invalid file type")
	// }
}

func getDefaultDownloadPath() string {
    usr, err := user.Current()
    if err != nil {
        panic(err)
    }
    return filepath.Join(usr.HomeDir, "Downloads")
}

func pathExists(path string) bool {
    _, err := os.Stat(path)
    if err == nil {
        return true
    }
    if os.IsNotExist(err) {
        return false
    }
    panic(err)
}