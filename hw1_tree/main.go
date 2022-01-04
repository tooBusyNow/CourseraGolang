package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type fileInfo struct {
	filePath string
	depth    int
}

type ByCase []fileInfo

func (files ByCase) Len() int           { return len(files) }
func (files ByCase) Swap(i, j int)      { files[i], files[j] = files[j], files[i] }
func (files ByCase) Less(i, j int) bool { return files[i].filePath < files[j].filePath }

func main() {

	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	var files []fileInfo

	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !printFiles && !info.IsDir() {
			return nil
		}

		files = append(files, fileInfo{
			filePath: path,
			depth:    strings.Count(path, "\\"),
		})
		return nil
	}); err != nil {
		panic("error occured while scaning dir")
	}

	sort.Sort(ByCase(files))

	for _, file := range files {
		fmt.Println(file.filePath, file.depth)
	}

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return nil
}
