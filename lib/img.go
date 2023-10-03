package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"syscall"

	"github.com/nfnt/resize"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	// 画像を読み込む
	file, err := os.Open("./img/BolinopsisMikado.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	// 画像をリサイズする
	width := 100
	height := 0
	if img.Bounds().Dx() > width {
		height = img.Bounds().Dy() * width / img.Bounds().Dx()
	} else {
		width = img.Bounds().Dx()
		height = img.Bounds().Dy()
	}
	resizedImg := resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor)

	// ターミナルのサイズを取得する
	terminalWidth, terminalHeight := decideSize("./img/BolinopsisMikado.jpg")
	if err != nil {
		panic(err)
	}

	// 画像をエリアごとに分割する
	areaWidth := resizedImg.Bounds().Dx() / terminalWidth
	areaHeight := resizedImg.Bounds().Dy() / terminalHeight

	// 画像をピクセル単位で処理する
	asciiArt := ""
	for y := 0; y < terminalHeight; y++ {
		for x := 0; x < terminalWidth; x++ {
			// エリアの範囲を計算する
			areaX := x * areaWidth
			areaY := y * areaHeight
			areaWidthEnd := areaX + areaWidth
			areaHeightEnd := areaY + areaHeight
			if x == terminalWidth-1 {
				areaWidthEnd = resizedImg.Bounds().Dx()
			}
			if y == terminalHeight-1 {
				areaHeightEnd = resizedImg.Bounds().Dy()
			}

			// エリアの平均色を計算する
			avgColor := getAverageColor(resizedImg, areaX, areaY, areaWidthEnd, areaHeightEnd)

			// ASCII文字を選択する
			asciiChar := getAsciiChar(avgColor)

			// ASCIIアートに追加する
			asciiArt += asciiChar
		}
		asciiArt += "\n"
	}

	// ASCIIアートを出力する
	fmt.Println(asciiArt)
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

func getAverageColor(img image.Image, x1, y1, x2, y2 int) color.RGBA {
	var r, g, b, a uint32
	var count uint32
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			pixel := img.At(x, y)
			pr, pg, pb, pa := pixel.RGBA()
			r += pr
			g += pg
			b += pb
			a += pa
			count++
		}
	}
	if count == 0 {
		return color.RGBA{}
	}
	r /= count
	g /= count
	b /= count
	a /= count
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func getAsciiChar(c color.RGBA) string {
	asciiChars := []string{" ", ".", ":", "-", "=", "+", "*", "#", "%", "@"}
	gray := uint8(0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B))
	charIndex := int(gray) * (len(asciiChars) - 1) / 255
	return asciiChars[charIndex]
}
