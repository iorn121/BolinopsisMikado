/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/nfnt/resize"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	rootCmd.AddCommand(imgCmd)
	imgCmd.Flags().StringP("path", "p", "", "image path to convert")
	imgCmd.Flags().BoolP("colored", "c", true, "colored the ascii when output to the terminal")
	imgCmd.Flags().BoolP("default", "d", true, "print default image")
}

// imgCmd represents the img command
var imgCmd = &cobra.Command{
	Use:   "img",
	Short: "convert image into ascii art",
	Long: `
Convert image into ascii art. args [path] after "img" are the path to the image file.
-c bool : Colored the ascii when output to the terminal (default true)
-p string : Image path to be convert (default "../img/BolinopsisMikado.jpg")`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		def, _ := cmd.Flags().GetBool("default")
		colored, _ := cmd.Flags().GetBool("colored")
		def_path := defaultImage()
		if path == "" || def {
			width, height := decideSize(def_path)
			convertImageToAscii(def_path, width, height, colored)
		} else {
			width, height := decideSize(path)
			convertImageToAscii(path, width, height, colored)
		}
	},
}

func defaultImage() string {
	// 実行ファイル直下のパスを取得
	path, _ := os.Getwd()
	dataPath := filepath.Join(path, "img")
	defaultPath := filepath.Join(dataPath, "BolinopsisMikado.jpg")
	return defaultPath
}

func decideSize(path string) (int, int) {
	width_term, height_term := getTerminalSize()
	width_img, height_img := getImageSize(path)
	var width, height int
	if width_term/width_img < height_term/height_img {
		width = width_term
		height = int(float64(width) / float64(width_img) * float64(height_img))
	} else {
		height = height_term
		width = int(float64(height) / float64(height_img) * float64(width_img))
	}
	return width, height
}

type colour struct {
	r uint32
	g uint32
	b uint32
}
type ascii_dot struct {
	col   int
	row   int
	color colour
	char  string
}

type ascii_img struct {
	height int
	width  int
	dots   [][]string
}

func (ascii_img *ascii_img) addDot(dot *ascii_dot, img image.Image, colored bool, ch chan bool, wg *sync.WaitGroup) {
	dot.color = colorAverage(img)
	dot.char = decideChar(img)
	var ascii string
	if colored {
		ascii = fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", dot.color.r/256, dot.color.g/256, dot.color.b/256, dot.char)
	} else {
		ascii = fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", 0, 0, 0, dot.char)
	}
	ascii_img.dots[dot.row][dot.col] = ascii
	ch <- true
	wg.Done()
}

func colorAverage(img image.Image) colour {
	h := img.Bounds().Dy()
	w := img.Bounds().Dx()
	var r_sum uint64
	var g_sum uint64
	var b_sum uint64
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			r_sum += uint64(r)
			g_sum += uint64(g)
			b_sum += uint64(b)
		}
	}
	// fmt.Println(r_sum/uint64(w*h)/256, g_sum/uint64(w*h)/256, b_sum/uint64(w*h)/256)
	return colour{r: uint32(r_sum / uint64(w*h)), g: uint32(g_sum / uint64(w*h)), b: uint32(b_sum / uint64(w*h))}
}

func decideChar(img image.Image) string {
	h := img.Bounds().Dy()
	w := img.Bounds().Dx()
	var sum uint64
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			sum += uint64(r + g + b)
		}
	}
	return string([]rune("MNHQ$OC?7>!:-;. ")[sum/uint64(w*h*3)%16])
}

// convertImageToAscii convert image to ascii art
// path : path to the image file
// width : width of the ascii art
// height : height of the ascii art
// colored : colored the ascii when output to the terminal
func convertImageToAscii(path string, width int, height int, colored bool) {
	file, err := os.Open(path)
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
	img_h := img.Bounds().Dy()
	img_w := img.Bounds().Dx()
	dots := make([][]string, height)
	for line := range dots {
		dots[line] = make([]string, width)
	}
	output_ascii := ascii_img{height: height, width: width, dots: dots}

	ch := make(chan bool, height*width)
	fmt.Println(height, width, colored)
	var wg sync.WaitGroup

	wg.Add(height * width)
	h_len := img_h / height
	w_len := img_w / width
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			// crop the image
			trimmed_img := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(j*w_len, i*h_len, (j+1)*w_len, (i+1)*h_len))

			dot := ascii_dot{row: i, col: j}
			go output_ascii.addDot(&dot, trimmed_img, colored, ch, &wg)
		}
	}
	wg.Wait()

	// print the output to the terminal inserting newline
	for i := 0; i < height; i++ {
		fmt.Println(strings.Join(output_ascii.dots[i], ""))
	}
}

func convertLineToAscii(img image.Image, line int, colored bool, output []string, ch chan<- string) {
	var ascii string
	for i := 0; i < img.Bounds().Dx(); i++ {
		r, g, b, _ := img.At(i, line).RGBA()
		ascii += convertPixelToAscii(r, g, b, colored)
	}
	output[line] = ascii
	ch <- ascii
}

func convertPixelToAscii(r uint32, g uint32, b uint32, colored bool) string {
	var ascii string
	if colored {
		ascii = fmt.Sprintf("\033[38;2;%d;%d;%dmx\033[0m", r/256, g/256, b/256)
	} else {
		ascii = fmt.Sprintf("\033[38;2;%d;%d;%dmx\033[0m", 0, 0, 0)
	}
	return ascii
}

func resizeImage(img image.Image, width int, height int) image.Image {
	resized_img := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	return resized_img
}

func getTerminalSize() (int, int) {
	var width int
	var height int
	var err error
	width, height, err = terminal.GetSize(syscall.Stdin)

	if err != nil {
		fmt.Printf("Error : %+v", err)
		os.Exit(1)
	}

	return width, height
}

func getImageSize(path string) (int, int) {
	file, err := os.Open(path)
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

	return img.Bounds().Dx(), img.Bounds().Dy()
}
