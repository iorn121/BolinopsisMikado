/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io"
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
		dirpath, _ := cmd.Flags().GetString("path")
		outputpath, _ := cmd.Flags().GetString("output")
		data,_:=NormalizeCSV(dirpath)
		save_csv(data,outputpath)
	// 	showGraph, _ := cmd.Flags().GetBool("graph")
	// 	if dirpath == "" {
	// 		fmt.Printf("dir path is not specified")
	// 		os.Exit(1)
	// 	}
	// 	filepath:= fmt.Sprintf("%s/%s", outputpath, "features.csv")
	// 	for _, path := range getFiles(dirpath) {
	// 		img := readImage(path)
	// 		vertical, horizontal:= edgeDensity(img)

	// 		file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		defer file.Close()
	// 		fmt.Fprintf(file, "%s,%f,%f\n", path, vertical, horizontal)
	// 	}
	// 	if showGraph {
	// 		showData(filepath)
	// 	}
	},
}

func init() {
	rootCmd.AddCommand(analyzeFeaturesCmd)
	analyzeFeaturesCmd.Flags().StringP("path", "p", "", "dir path to convert")
	// デフォルトはDownloadディレクトリ
	analyzeFeaturesCmd.Flags().StringP("output", "o", getDefaultDownloadPath() , "output dir path")
	analyzeFeaturesCmd.Flags().BoolP("graph", "g", false, "graph image")
}

func save_csv(records []Feature,path string) {
	// CSVファイルを開く
	filepath:= fmt.Sprintf("%s/%s", path, "ratio_data.csv")
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// CSVファイルを書き込み可能な形で開く
	w := csv.NewWriter(f)
	defer w.Flush()

	// データを書き込む
	for _, record := range records {
		// record.Y/record.X をratioとして書き込む
		// record.Xが0の場合は0を書き込む
		var ratio float64
		if record.X == 0 {
			ratio = 0
		} else {
			ratio = record.Y/record.X
		}

		content:=[]string{strconv.Itoa(record.Label),strconv.FormatFloat(ratio, 'f', -1, 64)}
		if err := w.Write(content); err != nil {
			log.Fatal(err)
		}
	}
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

// 画像の縦横それぞれの明暗の差を計算して密度として返す関数
func edgeDensity(img image.Image) (float64, float64) {
	// 画像の幅と高さを取得
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 縦成分と横成分、濃度の合計を格納する2次元配列を初期化
	vertical := 0.0
	horizontal := 0.0

	// エッジ検出を行う
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x==0 || y==0 || x==width-1 || y==height-1 {
				continue
			}
			// 縦成分のエッジ検出
			vertical+= math.Abs(float64(color.GrayModel.Convert(img.At(x, y+1)).(color.Gray).Y) - float64(color.GrayModel.Convert(img.At(x, y-1)).(color.Gray).Y))
			// 横成分のエッジ検出
			horizontal+= math.Abs(float64(color.GrayModel.Convert(img.At(x+1, y)).(color.Gray).Y) - float64(color.GrayModel.Convert(img.At(x-1, y)).(color.Gray).Y))
		}
	}

	return horizontal/float64(width*height), vertical/float64(width*height)
}


func showData(filepath string) {
    // CSVファイルを開く
	isExist := pathExists(filepath)
	if !isExist {
		log.Fatal(fmt.Sprintf("Error: %s does not exist", filepath))
	}
    f, err := os.Open(filepath)
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

type Feature struct {
	Label int
	X float64
	Y float64
}

// CSVファイルを読み込んで、x値とy値を0〜1の範囲に正規化する関数
func NormalizeCSV(filename string) ([]Feature, error) {
    // CSVファイルを開く
    file, err := os.Open(filename)
    if err != nil {
		return nil, err
    }
    defer file.Close()

    // CSVリーダーを作成
    reader := csv.NewReader(file)

    // x値とy値を格納するスライス
    var xData []float64
    var yData []float64
	var xMin float64
	var xMax float64
	var yMin float64
	var yMax float64

    // CSVファイルを1行ずつ読み込む
    for {
        // 行を読み込む
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
			return nil, err
        }
        // x値とy値を取得
        x, err := strconv.ParseFloat(record[1], 64)
        if err != nil {
			return nil,  err
        }
		if x < xMin {
			xMin = x
		}
		if x > xMax {
			xMax = x
		}
        y, err := strconv.ParseFloat(record[2], 64)
        if err != nil {
			return nil, err
        }
		if y < yMin {
			yMin = y
		}
		if y > yMax {
			yMax = y
		}

        // 正規化された値をスライスに追加
        xData = append(xData, x)
        yData = append(yData, y)
    }
	var normalizeData []Feature

	for i, x := range xData {
		if i<len(yData) {
			y := yData[i]
			// x値を0〜1の範囲に正規化
			xNormalized := (x - xMin) / (xMax - xMin)

			// y値を0〜1の範囲に正規化
			yNormalized := (y - yMin) / (yMax - yMin)

			// 正規化された値をスライスに追加
			normalizeData = append(normalizeData, Feature{Label: i+33, X: xNormalized, Y: yNormalized})
		}
	}
	fmt.Printf("xMin: %f, xMax: %f, yMin: %f, yMax: %f\n", xMin, xMax, yMin, yMax)
	// fmt.Printf("Normalize data: %+v\n", normalizeData)
	return normalizeData, nil
}
