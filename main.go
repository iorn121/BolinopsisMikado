package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
)

func flagUsage() {
usageText := `

Usage:
bmikado [arguments]

The commands are:
swimjelly	:	When you type this command, a jellyfish swims on your terminal window.

Use "bmikado [command] --help" for more infomation about a command`

fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

var i int
var s string
var b bool

// サブコマンドを受け取る
func init() {
	flag.Usage = flagUsage
	flag.Parse()
	switch flag.Arg(0) {
	case "swimjelly":
		displayJellyfish()
		os.Exit(0)
	}
}


func main() {
	flag.Usage = flagUsage
	flag.Parse()
	fmt.Println(i, s, b)
}

// jpgファイルを読み込んで、その画像をアスキーアートに変換してターミナルに表示する
func displayImage(image_path string) {
	// imgで画像を読み込む
	data, err := os.Open(image_path)
	if err != nil {
		fmt.Println("画像の読み込みに失敗しました")
		os.Exit(1)
	}
	defer data.Close()
	img,_,err := image.Decode(data)
	// 画像をリサイズする
	img = resizeImage(img,100,100)

	// 画像をグレースケールに変換する

	// 画像をアスキーアートに変換する

	// 画像をターミナルに表示する

}

// 画像を指定したサイズにリサイズする
func resizeImage(img image.Image, height int,width int) image.Image {
	var ratio float64

	// 画像のサイズを取得する
	imgHeight := img.Bounds().Max.Y
	imgWidth := img.Bounds().Max.X

	// 縦横の比率を計算する
	heightRatio := float64(height) / float64(imgHeight)
	widthRatio := float64(width) / float64(imgWidth)

	// 縦横の比率のうち、小さい方を採用する
	if heightRatio < widthRatio {
		ratio = heightRatio
	} else {
		ratio = widthRatio
	}
	resizedImage := &image.RGBA{}
	// 縦横の比率を元にリサイズする
	resizedImage=image.NewRGBA(image.Rect(0, 0, int(float64(imgWidth)*ratio), int(float64(imgHeight)*ratio)))

	return resizedImage
}

// ターミナルにクラゲのアニメーションを表示する
func displayJellyfish() {
	for {
		for _, r := range "-\\|/" {
			fmt.Printf("\r%c", r)
		}
	}
}