package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	return scanPath("", out, path, printFiles)
}

func scanPath(preStr string, output io.Writer, currDir string, printFiles bool) error {
	dirObj, err := os.Open(currDir)
	if err != nil {
		return err
	}
	defer dirObj.Close()
	dirName := dirObj.Name()
	dirItems, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}
	var fiMap map[string]os.FileInfo = map[string]os.FileInfo{}
	var fNames []string
	for _, item := range dirItems {
		if !printFiles && !item.IsDir() {
			continue
		}
		fNames = append(fNames, item.Name())
		fiMap[item.Name()] = item
	}
	sort.Strings(fNames)
	dirItems = make([]os.FileInfo, len(fNames))
	for i, strName := range fNames {
		dirItems[i] = fiMap[strName]
	}
	itemsCount := len(dirItems)
	for i, item := range dirItems {
		if item.IsDir() {
			var fullStr string
			if i+1 < itemsCount {
				fmt.Fprintf(output, preStr+"├───"+"%s\n", item.Name())
				fullStr = preStr + "│\t"
			} else {
				fmt.Fprintf(output, preStr+"└───"+"%s\n", item.Name())
				fullStr = preStr + "\t"
			}
			newDir := filepath.Join(currDir, item.Name())
			err = scanPath(fullStr, output, newDir, printFiles)
			if err != nil {
				fmt.Fprintf(output, preStr+"    ["+"%s]\n", err.Error())
			}
		} else if printFiles {
			if item.Size() > 0 {
				if i+1 < itemsCount {
					fmt.Fprintf(output, preStr+"├───%s (%vb)\n", item.Name(), item.Size())
				} else {
					fmt.Fprintf(output, preStr+"└───%s (%vb)\n", item.Name(), item.Size())
				}
			} else {
				if i+1 < itemsCount {
					fmt.Fprintf(output, preStr+"├───%s (empty)\n", item.Name())
				} else {
					fmt.Fprintf(output, preStr+"└───%s (empty)\n", item.Name())
				}
			}
		}
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		log.Fatalf("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
