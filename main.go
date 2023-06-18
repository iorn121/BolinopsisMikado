package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	"os"
	"BolinopsisMikado/lib"
	"BolinopsisMikado/experiment"
)


func flagUsage() {
usageText := `

Usage:
mikado [arguments]

The commands are:
mikado (mikado -default, mikado -d)	:	When you type this command, a jellyfish swims on your terminal window.

Use "mikado -help (mikado -h)" for more infomation about a command`

fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}



func main() {
	flag.Usage = flagUsage
	flag.Parse()
		switch flag.Arg(0) {
	case "" , "-normal" ,"-n":
		displayBolinopsisMikado()
		os.Exit(0)
	case "-display","-d":
		filepath:=flag.Args()[1:]
		display(filepath)
		os.Exit(0)
	case "-experiment","-e":
		experiment.makeLissajous()
	default:
		fmt.Println("error; can't find command")
	}
}

func display(filepath []string) {
	for _,path := range filepath {
		// fileの拡張子がjpeg, png, gifのいずれかであるかを確認する
		// 画像の場合はdisplayImage関数を呼び出す
		fmt.Println(path)
		lib.displayImage(path)
	}
}

// ターミナルにクラゲのアニメーションを表示する
func displayBolinopsisMikado() {
	for {
		for _, r := range "-\\|/" {
			fmt.Printf("\r%c", r)
		}
	}
}