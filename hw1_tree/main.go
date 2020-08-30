package main

import (
	"fmt"
	"golang.org/x/text/number"
	"io"
	"io/ioutil"
	"os"
	"strings"
	//"path/filepath"
	//"strings"
)

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		cnt := strings.Count(path, "\\")
		if file.IsDir() {
			fmt.Printf("%", file.Name())
			dirTree(out, path+file.Name()+"\\", printFiles)
		} else {
			if printFiles {
				fmt.Println(file.Name(), printFiles)
			}
		}
	}
	return
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	//fmt.Println(len(os.Args) == 3 && os.Args[2] == "-f")
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
