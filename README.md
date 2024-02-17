# "mikado" (BolinopsisMikado) CLI by Go

画像や映像を ASCII アートに変換してターミナル上で表示できます。
Images and videos can be converted to ASCII art and displayed on the terminal.

# Demo

Todo: paste Demo Movie

# Features

ローカルファイルの画像や映像を ASCII アートに変換してターミナル上で表示できます。
デフォルトでカブトクラゲの ASCII アートをターミナルに表示できます。

# Requirement

Go 1.20

# Installation

```bash
go install github.com/iorn121/BolinopsisMikado@latest
```

# Usage

BolinopsisMikado [arguments]
binary path: "~/go/bin/"

The commands are:
BolinopsisMikado lissajous [path] : When you type this command, a random lissajous gif image is generated and saved to the specified path.

BolinopsisMikado img2ascii [filepath] : When you type this command, the image converted into ASCII you specified will be displayed on your terminal window. Please set option "-d" if you want to display same image.

Options:

```
-c Colored the ascii when output to the terminal (default true)
-p string
Image path to be convert (default "docs/images/lufei.jpg")

```

Use "mikado -help (mikado -h)" for more infomation about a command`

# Note

# Author

- iorn121(Io)
- github: @iorn121
- twitter: @121Tkn
- email: wmt.tkn.121@gmail.com
- portfolio: https://iorn121.github.io/

# License

This tool is under [MIT license](https://en.wikipedia.org/wiki/MIT_License).
