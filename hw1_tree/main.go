package main

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var PREFIX_PIPE = "│\t"
var ELBOW = "└───"
var TEE = "├───"
var SPACE_PREFIX = "\t"

func main() {

	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, root string, printFiles bool) error {
	var pipesCounter, spaceCounter int = 0, 0
	recursiveTree(printFiles, out, root, pipesCounter, spaceCounter)
	return nil
}

func getStringSize(fsComp os.FileInfo) string {

	var size string

	if fsComp.Size() == 0 {
		size = " (empty)"
	} else {
		var withContent = []string{" (", strconv.Itoa(int(fsComp.Size())), "b)"}
		size = strings.Join(withContent, "")
	}
	return size
}

func excludeFiles(levelComps []fs.FileInfo) []fs.FileInfo {
	var result []fs.FileInfo
	for _, comp := range levelComps {
		if comp.IsDir() {
			result = append(result, comp)
		}
	}
	return result
}

func recursiveTree(flag bool, out io.Writer, path string, pipesCounter int, spacesCounter int) error {

	levelComps, err := ioutil.ReadDir(path)
	if !flag {
		levelComps = excludeFiles(levelComps)
	}

	levelCount := len(levelComps)
	size := ""

	connector := strings.Repeat(PREFIX_PIPE, int(pipesCounter)) +
		strings.Repeat(SPACE_PREFIX, int(spacesCounter))

	for idx, fsComp := range levelComps {

		var outputArr []string

		if !fsComp.IsDir() {
			size = getStringSize(fsComp)
		}

		if idx == levelCount-1 {
			pipesCounter -= 1
			spacesCounter += 1
			outputArr = []string{connector, ELBOW + fsComp.Name(), size, "\n"}
		} else {
			outputArr = []string{connector, TEE + fsComp.Name(), size, "\n"}
		}

		temp := strings.Join(outputArr, "")
		out.Write([]byte(temp))

		if fsComp.IsDir() {
			recursiveTree(flag, out, path+string(os.PathSeparator)+fsComp.Name(), pipesCounter+1, spacesCounter)
		}
		size = ""
	}
	return err
}
