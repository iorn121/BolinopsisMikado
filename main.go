package main

import (
	"flag"
	"fmt"
	"os"
)

func flagUsage() {
usageText := `

Usage:
example command [arguments]
The commands are:
convhex    convert Number to Hex
convbinary convert Number to binary
Use "Example [command] --help" for more infomation about a command`

fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

var i int
var s string
var b bool

func init() {
	// Var() 関数の引数は以下の通り
    // 第一引数は束縛する変数のポインタ
    // 第二引数はフラグ名
    // 第三引数はデフォルト値
    // 第四引数はフラグの説明
	flag.IntVar(&i, "i", 0, "数値" )
	flag.StringVar(&s, "s", "default", "文字列")
	flag.BoolVar(&b, "b", false, "真偽値")
}

func main() {
	flag.Parse()
	fmt.Println(i, s, b)
}