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
go get github.com/iorn121/BolinopsisMikado
```

# Usage

mikado [arguments]

The commands are:
mikado (mikado -normal, mikado -n) : When you type this command, a jellyfish swims on your terminal window.

mikado -display (mikado -d) [filepath] : When you type this command, the image you specified will be displayed on your terminal window.

mikado -experiment (mikado -e) : When you type this command, a lissajous figure will be displayed on your terminal window.

Use "mikado -help (mikado -h)" for more infomation about a command`

## mikado -help (mikado)

## mikado (mikado -default) (mikado -d)

When you type this command, a Bolinopsis Mikado swims on your terminal window.

```bash
mikado
```

## mikado -display {filename}

This command allows you convert an image or video into ASCII art or animation and display it on your terminal window.

```bash
mikado -display sample.png
```

# Note

# Author

- iorn121(Io)
- github: @iorn121
- twitter: @121Tkn
- email: wmt.tkn.121@gmail.com
- portfolio: https://iorn121.github.io/

# License

This tool is under [MIT license](https://en.wikipedia.org/wiki/MIT_License).
