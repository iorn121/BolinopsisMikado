package main

import (
	"fmt"
	"image"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	width, height := getTerminalSize()
	fmt.Printf("width : %d, height : %d\n", width, height)
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

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Printf("Error : %+v", err)
		os.Exit(1)
	}

	return img.Width, img.Height
}
