/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"image"
	"math"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
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
	imgCmd.Flags().BoolP("default", "d", false, "print default image")
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
			fmt.Println(path)
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
	width_term, _ := getTerminalSize()
	width_img, height_img := getImageSize(path)
	var width, height int
	width = width_term
	height = width * height_img / width_img
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

func (ascii_img *ascii_img) addDot(dot *ascii_dot, img image.Image, ch chan bool, wg *sync.WaitGroup) {
	// fmt.Println(dot)
	dot.char = decideChar(img)
	var ascii string
	ascii = fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", dot.color.r/256, dot.color.g/256, dot.color.b/256, dot.char)

	ascii_img.dots[dot.row][dot.col] = ascii
	ch <- true
	wg.Done()
}

func (ascii_img *ascii_img) addDots(dot *ascii_dot, img image.Image, colored bool) {
	dot.color = colorAverage(img)
	dot.char = decideChar(img)
	var ascii string
	if colored {
		ascii = fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", dot.color.r/256, dot.color.g/256, dot.color.b/256, dot.char)
		fmt.Println("color", ascii)
	} else {
		ascii = fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", 0, 0, 0, dot.char)
	}
	// fmt.Println(ascii)
	ascii_img.dots[dot.row][dot.col] = ascii
}

func P(t interface{}) {
	fmt.Println(reflect.TypeOf(t))
}

// 画像の平均色を計算する
// 画像の平均色は、画像の各ピクセルのRGB値の平均値で求める
func colorAverage(img image.Image) colour {
	h := img.Bounds().Dy()
	w := img.Bounds().Dx()
	var r_average uint32
	var g_average uint32
	var b_average uint32
	for i := 0; i < h; i++ {
		// convertLineToAscii(img, i, true)
		for j := 0; j < w; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			r_average += uint32(r)
			g_average += uint32(g)
			b_average += uint32(b)
		}
	}
	fmt.Println(r_average/uint32(h*w), g_average/uint32(h*w), b_average/uint32(h*w))
	return colour{r: r_average / uint32(h*w), g: g_average / uint32(h*w), b: b_average / uint32(h*w)}
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// 画像をASCII文字に変換する
// 画像の濃度は、画像の各ピクセルのRGB値の合計値で求める
// 画像の濃度に応じて、文字を選択する
// 画像の濃度が高いほど、濃い文字を選択する
// 画像の濃度が低いほど、薄い文字を選択する
// 画像の濃度が0の場合、空白文字を選択する

func decideChar(img image.Image) string {
	h := img.Bounds().Dy()
	w := img.Bounds().Dx()
	var sum uint64
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			sum += uint64(r) + uint64(g) + uint64(b)
		}
	}


	// ascii_chars := "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	// 各文字列を画像処理して、エッジを計算して縦成分と横成分の合計を計算する

	// vertical, horizontal := detectEdges(img)

	// fmt.Println(vertical, horizontal)

	return "x"
}

// 画像のエッジ検出を行う関数
func detectEdges(img image.Image) (float64, float64) {
	// 画像の幅と高さを取得
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 縦成分と横成分の合計を格納する2次元配列を初期化
	vertical := 0.0
	horizontal := 0.0

	// エッジ検出を行う
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// 縦成分のエッジ検出
			vertical += math.Abs(float64(color.GrayModel.Convert(img.At(x, y+1)).(color.Gray).Y) - float64(color.GrayModel.Convert(img.At(x, y-1)).(color.Gray).Y))
			// 横成分のエッジ検出
			horizontal += math.Abs(float64(color.GrayModel.Convert(img.At(x+1, y)).(color.Gray).Y) - float64(color.GrayModel.Convert(img.At(x-1, y)).(color.Gray).Y))
		}
	}

	return vertical, horizontal
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
	img_resized := resizeImage(img, width, height)
	// save_image(img_resized, 0)
	dots := make([][]string, height)
	for line := range dots {
		dots[line] = make([]string, width)
	}
	output_ascii := ascii_img{height: height, width: width, dots: dots}

	ch := make(chan bool, height*width)
	// fmt.Println(height, width)
	var wg sync.WaitGroup

	wg.Add(height * width)
	h_len := img_h / height
	w_len := img_w / width
	fmt.Println(h_len, w_len)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			trimmed_img := img.(SubImager).SubImage(image.Rect(j*w_len, i*h_len, (j+1)*w_len, (i+1)*h_len))
			save_image(trimmed_img, i*width+j)
			r,g,b,_ := img_resized.At(j,i).RGBA()
			dot := ascii_dot{row: i, col: j, color: colour{r:r,g:g,b:b}, char: "."}
			go output_ascii.addDot(&dot, trimmed_img, ch, &wg)
		}
	}
	wg.Wait()

	// print the output to the terminal inserting newline
	for i := 0; i < height; i++ {
		fmt.Println(strings.Join(output_ascii.dots[i], ""))
	}
}

// save_image save image to file
func save_image(img image.Image, index int) {
	fso, err := os.Create(fmt.Sprintf("output/out_%s.jpg", strconv.Itoa(index)))
	if err != nil {
		fmt.Println("create:", err)
		return
	}
	defer fso.Close()
	jpeg.Encode(fso, img, &jpeg.Options{Quality: 100})
}

func convertLineToAscii(img image.Image, line int, colored bool) {
	var ascii string
	for i := 0; i < img.Bounds().Dx(); i++ {
		r, g, b, _ := img.At(i, line).RGBA()
		ascii += convertPixelToAscii(r, g, b, colored)
	}
	// fmt.Println(ascii)
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

// 画像を縦height, 横widthに分割して、それぞれの画像をgoroutineで処理する
// 分割した各画像は、それぞれのgoroutineで処理され、結果はchannelに送られる
// goroutineでは、画像の平均色を計算し、濃度に応じて文字を選択して、ASCII文字に変換する
