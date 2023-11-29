/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"sort"
	"encoding/csv"

	"github.com/nfnt/resize"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	rootCmd.AddCommand(img2asciiCmd)
	img2asciiCmd.Flags().StringP("path", "p", "", "image path to convert")
	img2asciiCmd.Flags().BoolP("colored", "c", true, "colored the ascii when output to the terminal")
	img2asciiCmd.Flags().BoolP("default", "d", false, "print default image")
}

// img2asciiCmd represents the img2ascii command
var img2asciiCmd = &cobra.Command{
	Use:   "img2ascii",
	Short: "convert image into ascii art",
	Long: `
Convert image into ascii art. args [path] after "img2ascii" are the path to the image file.
-c bool : Colored the ascii when output to the terminal (default true)
-p string : Image path to be convert (default "../img2ascii/BolinopsisMikado.jpg")`,
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
	width_term, height_term := getTerminalSize()
	width_img, height_img := getImageSize(path)
	var width, height int
	// 横幅と縦幅の比率を保ち、ターミナルで全体が見えるようにする
	if width_img/width_term > height_img/height_term {
		width = width_term
		height = int(float64(height_img) / float64(width_img) * float64(width_term))
	} else {
		height = height_term
		width = int(float64(width_img) / float64(height_img) * float64(height_term))
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

func (ascii_img *ascii_img) addDot(dot *ascii_dot, img image.Image,ref []ascii, ch chan bool, wg *sync.WaitGroup) {
	// fmt.Println(dot)
	dot.char = decideChar(img,ref)
	var ascii string
	ascii = fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", dot.color.r/256, dot.color.g/256, dot.color.b/256, dot.char)

	ascii_img.dots[dot.row][dot.col] = ascii
	ch <- true
	wg.Done()
}



type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type ascii struct {
	char string
	ratio float64
}

// ascii/ratio.csvを読み込んで、asciiの配列を返す
func readRatio() []ascii {
	// CSVファイルを開く
	filepath := filepath.Join("ascii", "ratio_data.csv")
	isExist := pathExists(filepath)
	if !isExist {
		fmt.Println("Error: ratio.csv does not exist")
		os.Exit(1)
	}
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error: failed to open ratio.csv")
		os.Exit(1)
	}
	defer f.Close()

	// CSVファイルをパースする
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println("Error: failed to parse ratio.csv")
		os.Exit(1)
	}
	asciis := make([]ascii, len(records))
	for i, record := range records {
		// record[0]からascii文字に変換
		c,_:=strconv.Atoi(record[0])
		asciis[i].char=fmt.Sprintf("%c",c)
		asciis[i].ratio, _ = strconv.ParseFloat(record[1], 64)
	}

	// asciisをratioの昇順にソート
	sort.Slice(asciis, func(i, j int) bool {
		return asciis[i].ratio < asciis[j].ratio
	})
	return asciis
}




func decideChar(img image.Image, ref []ascii) string {
	// analyzeFeatures.goのedgeDensity関数を使う
	horizontal, vertical := edgeDensity(img)
	var ratio float64
	if horizontal == 0 {
		ratio = 0
	} else {
		ratio = vertical / horizontal
	}
	// 昇順のref[i].ratioを参照して選択
	// 二分探索の結果、最も近いref[i].ratioを持つasciiを選択
	l:=0
	r:=len(ref)
	for l+1<r {
		mid:=(l+r)/2
		if ref[mid].ratio>ratio {
			r=mid
		} else {
			l=mid
		}
	}
	// char:=ref[l].char
	// fmt.Printf("ratio:%f, char:%s\n",ratio,char)

	return "x"
}



// convertImageToAscii convert image to ascii art
// path : path to the image file
// width : width of the ascii art
// height : height of the ascii art
// colored : colored the ascii when output to the terminal
func convertImageToAscii(path string, width int, height int, colored bool) {
	// 時間の計測開始
	start := time.Now()
	// 画像を開く
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
	ref := readRatio()
	ch := make(chan bool, height*width)
	// fmt.Println(height, width)
	var wg sync.WaitGroup

	wg.Add(height * width)
	h_len := img_h / height
	w_len := img_w / width
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			trimmed_img := img.(SubImager).SubImage(image.Rect(j*w_len, i*h_len, (j+1)*w_len, (i+1)*h_len))
			// save_image(trimmed_img, i*width+j)
			r,g,b,_ := img_resized.At(j,i).RGBA()
			dot := ascii_dot{row: i, col: j, color: colour{r:r,g:g,b:b}, char: "."}
			go output_ascii.addDot(&dot, trimmed_img,ref, ch, &wg)
		}
	}
	wg.Wait()

	// print the output to the terminal inserting newline
	for i := 0; i < height; i++ {
		fmt.Println(strings.Join(output_ascii.dots[i], ""))
	}
	// 時間の計測終了
	end := time.Now()
	fmt.Printf("height:%d, width:%d\ntime: %f sec\n",height, width, (end.Sub(start)).Seconds())
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
