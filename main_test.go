package main_test

import (
	"fmt"
	"os"
)


func convert2ascii() {
	s,sep := "",""
	for _,arg := range os.Args[1:] {
		s+=sep+arg
		sep=" "
	}
	fmt.Println(s)
}