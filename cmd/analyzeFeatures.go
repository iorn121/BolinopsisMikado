/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"bytes"
	"log"
	"net/http"
	"os/exec"
	"encoding/csv"
	chart "github.com/wcharczuk/go-chart/v2"
	"strconv"
	"path/filepath"
	"image"
	"image/color"
	"math"
	"github.com/spf13/cobra"
)

// analyzeEdgeCmd represents the analyzeEdge command
var analyzeFeaturesCmd = &cobra.Command{
	Use:   "analyzeFeatures",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dirpath:=args[0]
		showGraph, _ := cmd.Flags().GetBool("graph")
		for _, path := range getFiles(dirpath) {
			img := readImage(path)
			vertical, horizontal, sum := countEdges(img)
			// fmt.Printf("%s, %f, %f, %f\n", path, vertical, horizontal, sum)

			// csvファイルに書き込む
			file, err := os.OpenFile("ascii/features.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			fmt.Fprintf(file, "%s,%f,%f,%f\n", path, vertical, horizontal, sum)
		}
		if showGraph {
			showData()
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeFeaturesCmd)
	analyzeFeaturesCmd.Flags().BoolP("graph", "g", false, "graph image")
}

// 画像のパスを取得する関数
func getFiles(dirpath string) []string {
	var files []string
	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".jpg" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

// 画像を読み込む関数
func readImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return img
}

// 画像のエッジとみなすピクセルをカウントし、密度を求める関数
func countEdges(img image.Image) (float64, float64, float64) {
	// 画像の幅と高さを取得
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 縦成分と横成分、濃度の合計を格納する2次元配列を初期化
	vertical := 0.0
	horizontal := 0.0
	sum := 0.0

	// エッジ検出を行う
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y > 50 {
				sum++
			}
			if x==0 || y==0 || x==width-1 || y==height-1 {
				continue
			}
			// 縦成分のエッジ検出
			if math.Abs(float64(color.GrayModel.Convert(img.At(x, y+1)).(color.Gray).Y) - float64(color.GrayModel.Convert(img.At(x, y-1)).(color.Gray).Y)) > 10 {
				vertical++
			}
			// 横成分のエッジ検出
			if math.Abs(float64(color.GrayModel.Convert(img.At(x+1, y)).(color.Gray).Y) - float64(color.GrayModel.Convert(img.At(x-1, y)).(color.Gray).Y)) > 10 {
				horizontal++
			}
		}
	}

	return horizontal/float64(width*height), vertical/float64(width*height), sum/float64(width*height)
}




func showData() {
    // CSVファイルを開く
	isExist := pathExists("ascii/features.csv")
	if !isExist {
		log.Fatal("Error: ascii/features.csv does not exist")
	}
    f, err := os.Open("ascii/features.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

	fmt.Printf("Open file: %s\n", f.Name())

    // CSVファイルをパースする
    r := csv.NewReader(f)
    records, err := r.ReadAll()
    if err != nil {
        log.Fatal(err)
    }

    // データを抽出する
    var xData []float64
    var yData []float64
    for _, record := range records {
        x, err := strconv.ParseFloat(record[1], 64)
        if err != nil {
            log.Fatal(err)
        }
        y, err := strconv.ParseFloat(record[2], 64)
        if err != nil {
            log.Fatal(err)
        }
        xData = append(xData, x)
        yData = append(yData, y)
    }

    // データをグラフ化する
    graph := chart.Chart{
		Height: 500,
		Width:  500,
		Background: chart.Style{
			Padding: chart.Box{
				Top:  10,
				Bottom: 10,
				Left: 10,
				Right: 10,
			},
		},
        XAxis: chart.XAxis{
            Name: "x",
        },
        YAxis: chart.YAxis{
            Name: "y",
        },
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeWidth:      chart.Disabled,
					DotWidth:         5,
					DotColor:         chart.ColorRed,
				},
				XValues: xData,
				YValues: yData,
			},
		},
    }
	
    graph.Elements = []chart.Renderable{
        chart.Legend(&graph),
    }


    // PNG形式の画像を作成する
    buffer := bytes.Buffer{}
    err = graph.Render(chart.PNG, &buffer)
    if err != nil {
        log.Fatal(err)
    }

    // Webサーバーを起動して画像を表示する
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "image/png")
        w.Write(buffer.Bytes())
    })
    // Webブラウザを起動する
    err = exec.Command("open", "http://localhost:8080").Start()
    if err != nil {
		log.Fatal(err.Error())
    }
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}