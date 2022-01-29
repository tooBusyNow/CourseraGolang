package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type fsCompInfo struct {
	sysPath string
	depth   int
	isDir   bool
	size    int64

	parent *fsCompInfo
}

type ByCase []fsCompInfo

func (fsComp ByCase) Len() int           { return len(fsComp) }
func (fsComp ByCase) Swap(i, j int)      { fsComp[i], fsComp[j] = fsComp[j], fsComp[i] }
func (fsComp ByCase) Less(i, j int) bool { return fsComp[i].sysPath < fsComp[j].sysPath }

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

func dirTree(out io.Writer, path string, printFiles bool) error {

	fsComps := getAndSortComponents(path, printFiles)
	for _, comp := range fsComps {
		fmt.Println(comp.sysPath, " ", comp.parent)
	}
	return nil
}

func getAndSortComponents(path string, printFiles bool) []fsCompInfo {

	var fsComps []fsCompInfo
	var root *fsCompInfo
	var dirParent *fsCompInfo

	if err := filepath.Walk(path,

		func(path string, info os.FileInfo, err error) error {
			if !printFiles && !info.IsDir() {
				return nil
			}

			contDepth := strings.Count(path, "\\")
			if dirParent != nil && contDepth <= dirParent.depth {
				dirParent = root
			}

			node := fsCompInfo{
				sysPath: path,
				depth:   strings.Count(path, "\\"),
				isDir:   info.IsDir(),
				size:    info.Size(),
				parent:  dirParent,
			}

			if node.parent == nil {
				root = &node
			}

			fsComps = append(fsComps, node)
			if info.IsDir() {
				dirParent = &node
			}

			return nil
		}); err != nil {
		panic("error occured while scaning dir")
	}

	sort.Sort(ByCase(fsComps))
	return fsComps
}
